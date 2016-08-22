package redis

import (
	"testing"
	"github.com/garyburd/redigo/redis"
	"fmt"
	"time"
)

var reidsPool *RedisPool
var errP error

func init() {

	dialFunc := func(config PoolConfig) (c redis.Conn, err error) {
		c, err = redis.Dial(config.Network, config.Address)
		if err != nil {
			return nil, err
		}

		if config.Password != "" {
			if _, err := c.Do("AUTH", config.Password); err != nil {
				c.Close()
				return nil, err
			}
		}

		_, selecterr := c.Do("SELECT", config.DbNum)
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}

	config := PoolConfig{
		Network  :"tcp",
		Address  : "127.0.0.1:6379",
		MaxIdle  :10,
		Password :"",
		DbNum    :0,
		Df       :dialFunc,
	}

	reidsPool, errP = NewRedisPool(config)

	fmt.Println(errP)
}

func TestRedisPoolString(t *testing.T) {
	if nil == reidsPool {
		t.Fatal("reidsPool is nil")
	}
	testKey := "testKey"
	testValStr := "testVal"
	testValNum := int64(12)
	reidsPool.Set(testKey, testValStr)
	str, _ := reidsPool.Get(testKey)
	if str != testValStr {
		t.Fatal("set the val fail")
	}

	reidsPool.Del(testKey)

	str, _ = reidsPool.Get(testKey)
	if str == testValStr {
		t.Fatal("Del fail")
	}

	reidsPool.Set(testKey, fmt.Sprintf("%d", testValNum))

	num, err := reidsPool.Decr(testKey)
	if num != testValNum - 1 || err != nil {
		t.Fatal("Decr fail")
	}

	num, err = reidsPool.Incr(testKey)
	if num != testValNum || err != nil {
		t.Fatal("Incr fail")
	}

	reidsPool.Del(testKey)

	lenAppend, err := reidsPool.Append(testKey, testValStr)
	if err != nil || lenAppend != len([]byte(testValStr)) {
		t.Fatal("Append fail")
	}

	//expire
	_, err = reidsPool.Expire(testKey, 1)
	if nil != err {
		t.Fatal("Expire fail" + err.Error())
	}
	time.Sleep(time.Millisecond * 1100)
	str, _ = reidsPool.Get(testKey)
	if str == testValStr {
		t.Fatal("Expire fail")
	}

	reidsPool.Set(testKey, testValStr)

	_, err = reidsPool.PExpire(testKey, 1000)
	if nil != err {
		t.Fatal("PExpire fail" + err.Error())
	}
	time.Sleep(time.Millisecond * 1100)
	str, _ = reidsPool.Get(testKey)
	if str == testValStr {
		t.Fatal("PExpire fail")
	}

	reidsPool.Del(testKey)
	err = reidsPool.SetEx(testKey, testValStr, 1)
	if nil != err {
		t.Fatal("SetEx fail" + err.Error())
	}
	time.Sleep(time.Millisecond * 1100)
	str, _ = reidsPool.Get(testKey)
	if str == testValStr {
		t.Fatal("SetEx fail")
	}

	reidsPool.Del(testKey)
}

func TestRedisPoolMap(t *testing.T) {
	if nil == reidsPool {
		t.Fatal("reidsPool is nil")
	}

	testKey := "testKeyMap"
	testfield1 := "testfield1"
	testfield2 := "testfield2"
	testfield3 := "testfield3"
	testVal1 := "testVal1"
	testVal2 := "testVal2"
	testVal3 := "testVal3"

	err := reidsPool.HMSet(testKey, testfield1, testVal1, testfield2, testVal2)
	if nil != err {
		t.Fatal("HMSet fail" + err.Error())
	}

	row, err := reidsPool.HSet(testKey, testfield3, testVal3)
	if nil != err {
		t.Fatal("HSet fail" + err.Error())
	}
	fmt.Println("row:", row)

	str, err := reidsPool.HGet(testKey, testfield3)
	if nil != err || str != testVal3 {
		t.Fatal("HGet fail" + err.Error())
	}

	strMap, err := reidsPool.HGetAll(testKey)
	if nil != err || strMap[testfield3] != testVal3 {
		t.Fatal("HGetAll fail" + err.Error())
	}
	fmt.Println(strMap)

	lenMap, err := reidsPool.HLen(testKey)
	if nil != err || lenMap != 3 {
		t.Fatal("HLen fail" + err.Error())
	}

	row, err = reidsPool.HDel(testKey, testfield3)
	if nil != err || row != 1 {
		t.Fatal("HDel fail" + err.Error())
	}
}

func TestSync(t *testing.T) {

}