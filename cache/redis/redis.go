// Copyright 2016  cbping. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// redis 封装包
// 主要引用 github.com/garyburd/redigo/redis
//
//  一、自定义Pool
//       不用github.com/garyburd/redigo/redis包中原有的pool功能。而已重新实现pool功能。
//       新的pool主要参考github.com/fzzy/radix中的pool的功能，通过golang的channel来实现，
//       非常漂亮；同时，pool的Put看起来非常清晰；总体实现比较优雅。
//
// 例子：
//   dialFunc := func(config PoolConfig) (c redis.Conn, err error) {
//	c, err = redis.Dial(config.Network, config.Address)
//	if err != nil {
//		return nil, err
//	}
//
//	if config.password != "" {
//		if _, err := c.Do("AUTH", config.Password); err != nil {
//			c.Close()
//			return nil, err
//		}
//	}
//
//	_, selecterr := c.Do("SELECT", config.DbNum)
//	if selecterr != nil {
//		c.Close()
//		return nil, selecterr
//	}
//	return
//     }
//
//     config=PoolConfig{
//            Network  :"tcp",
//	      Address  :   "127.0.0.1:6379",
//	      MaxIdle  :10,
//	      Password :"123456",
//	      DbNum    :0,
//	      Df       :dialFunc,
//     }
//
// 	redis, _ := NewRedisPool(config)
//
package redis

import (
	"encoding/json"
	"errors"
	"github.com/garyburd/redigo/redis"
)

var (
	ErrNil     = redis.ErrNil
	ErrPowerOn = errors.New("RedisPool:turn on first")
)

// 内部使用了池功能
type RedisPool struct {
	pool    *Pool
	powerOn bool
	//record func(tag, msg string)
}

// 新建
func NewRedisPool(config PoolConfig) (*RedisPool, error) {
	pool, err := NewPool(config)
	if nil != err {
		return nil, err
	}
	return &RedisPool{pool, true}, nil
}

func (rp *RedisPool) SetPowerOn(powerOn bool) *RedisPool {
	rp.powerOn = powerOn
	return rp
}

//
//// 记录信息
//func (rp *RedisPool) log(tag, msg string) {
//	if nil != rp.record {
//		rp.record(tag, msg)
//	}
//}

// 关闭
// 清空内部连接池
func (rp *RedisPool) Close() error {
	if nil != rp.pool {
		rp.pool.Empty()
	}
	return nil
}

// 从池中获取空闲（或者新建）连接处理redis命令操作
// 调用有效连接redis.Conn的Do方法
// 此方法处理连接池连接的取出和放回等额外的相关工作
func (rp *RedisPool) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if !rp.powerOn {
		return nil, ErrPowerOn
	}
	if nil == rp.pool {
		return nil, errors.New("RedisPool:please create a redis pool first")
	}
	conn, err := rp.pool.Get()
	if nil != err {
		return
	}
	reply, err = conn.Do(commandName, args...)
	rp.pool.Put(conn)
	return
}

// 调用内部pool.Get
func (rp *RedisPool) GetConn() (redis.Conn, error) {
	if nil != rp.pool {
		return rp.pool.Get()
	}
	return nil, errors.New("the pool is nil")
}

// 调用内部pool.Put
func (rp *RedisPool) PutConn(conn redis.Conn) error {
	if nil != rp.pool {
		return rp.pool.Put(conn)
	}
	return errors.New("the pool is nil")
}

//----------------------------------------------------------------------------------------------------------------------
// 常用操作
//----------------------------------------------------------------------------------------------------------------------

// 字符串类型相关命令操作
//----------------------------------------------------------------------------------------------------------------------

const (
	okResp = "OK"
)

// 获取指定 key 的值。（string（字符串））
// 如果key不存在或者有异常则返回空字符串
func (rp *RedisPool) Get(key string) (string, error) {
	return redis.String(rp.Do("GET", key))
}

// 设置指定 key 的值
// 注意：一个键最大能存储512MB。
func (rp *RedisPool) Set(key, val string) error {
	v, err := redis.String(rp.Do("SET", key, val))
	if nil == err && v == okResp {
		return nil
	} else if nil == err {
		return errors.New("not ok")
	}
	return err
}

// SETEX key seconds value 将值 value 关联到 key ，
// 并将 key 的过期时间设为 seconds (以秒为单位)。
// @expired 有效时长 (以秒为单位)
func (rp *RedisPool) SetEx(key, val string, expired int64) error {
	v, err := redis.String(rp.Do("SETEX", key, expired, val))
	if nil == err && v == okResp {
		return nil
	} else if nil == err {
		return errors.New("not ok")
	}
	return err
}

// EXPIRE key seconds 为给定 key 设置过期时间。(以秒为单位)。
func (rp *RedisPool) Expire(key string, expired int64) (int, error) {
	return redis.Int(rp.Do("EXPIRE", key, expired))
}

// PEXPIRE key milliseconds 设置 key 的过期时间亿以毫秒计。
func (rp *RedisPool) PExpire(key string, expired int64) (int, error) {
	return redis.Int(rp.Do("PEXPIRE", key, expired))
}

// 删除指定 key 的值
// @return 返回删除个数。
func (rp *RedisPool) Del(key string) (int, error) {
	return redis.Int(rp.Do("DEL", key))
}

// INCR key 将 key 中储存的数字值增一。
// @return int64 增一之后的数字值
func (rp *RedisPool) Incr(key string) (int64, error) {
	return redis.Int64(rp.Do("INCR", key))
}

// DECR key 将 key 中储存的数字值减一。
// @return int64 减一之后的数字值
func (rp *RedisPool) Decr(key string) (int64, error) {
	return redis.Int64(rp.Do("DECR", key))
}

// APPEND key value 如果 key 已经存在并且是一个字符串，
// APPEND 命令将 value 追加到 key 原来的值的末尾。
// 否则，新建key/value
func (rp *RedisPool) Append(key, val string) (int, error) {
	return redis.Int(rp.Do("APPEND", key, val))
}

//json
//@see Get()
func (rp *RedisPool) GetJson(key string, reply interface{}) (err error) {
	rstr, err := rp.Get(key)
	if nil != err {
		return err
	}
	err = json.Unmarshal([]byte(rstr), reply)
	return
}

//json
//@see Set()
func (rp *RedisPool) SetJson(key string, val interface{}) (err error) {
	dec, err := json.Marshal(val)
	if nil != err {
		return err
	}
	err = rp.Set(key, string(dec))
	return
}

//json
//@see SetEx()
func (rp *RedisPool) SetExJson(key string, val interface{}, expired int64) (err error) {
	dec, err := json.Marshal(val)
	if nil != err {
		return err
	}
	err = rp.SetEx(key, string(dec), expired)
	return
}

// 哈希(Hash)类型相关命令操作
//
// Redis hash 是一个string类型的field和value的映射表，hash特别适合用于存储对象。
// Redis 中每个 hash 可以存储 232 - 1 键值对（40多亿）。
//----------------------------------------------------------------------------------------------------------------------

// HGET key field 获取存储在哈希表中指定字段的值
// 如果key或field不存在或者有异常则返回空字符串
func (rp *RedisPool) HGet(key, field string) (string, error) {
	return redis.String(rp.Do("HGet", key, field))
}

// HSET key field value 将哈希表 key 中的字段 field 的值设为 value
func (rp *RedisPool) HSet(key, field, value string) (int, error) {
	return redis.Int(rp.Do("HSET", key, field, value))
}

// HMSET key field1 value1 [field2 value2 ]
// 同时将多个 field-value (域-值)对设置到哈希表 key 中。
func (rp *RedisPool) HMSet(key string, field_value ...interface{}) error {
	args := append([]interface{}{key}, field_value...)
	v, err := redis.String(rp.Do("HMSET", args...))
	if nil == err && v == okResp {
		return nil
	} else if nil == err {
		return errors.New("not ok")
	}
	return err
}

// HDEL key field2 [field2] 删除一个或多个哈希表字段
func (rp *RedisPool) HDel(key string, field ...interface{}) (int, error) {
	return redis.Int(rp.Do("HDEL", append([]interface{}{key}, field...)...))
}

// HGETALL key 获取在哈希表中指定 key 的所有字段和值
func (rp *RedisPool) HGetAll(key string) (map[string]string, error) {
	v, err := redis.Strings(rp.Do("HGETALL", key))
	if nil == err {
		redisMap := make(map[string]string)
		mapLen := len(v)
		for index := 0; index < mapLen; index += 2 {
			redisMap[v[index]] = v[index+1]
		}
		return redisMap, nil
	}

	return nil, err

}

// HLEN key 获取哈希表中字段的数量
func (rp *RedisPool) HLen(key string) (int, error) {
	return redis.Int(rp.Do("HLEN", key))
}

// 列表类型相关命令操作
//
// Redis列表是简单的字符串列表，按照插入顺序排序。你可以添加一个元素导列表的头部（左边）或者尾部（右边）
// 一个列表最多可以包含 232 - 1 个元素 (4294967295, 每个列表超过40亿个元素)。
//----------------------------------------------------------------------------------------------------------------------
