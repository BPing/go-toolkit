package mredis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"testing"
	"time"
)

var redisPool *RedisPool
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
		Network:  "tcp",
		Address:  "127.0.0.1:6379",
		MaxIdle:  10,
		Password: "",
		DbNum:    0,
		Df:       dialFunc,
	}

	redisPool, errP = NewRedisPool(config)

	fmt.Println(errP)
}

func TestRedisPoolString(t *testing.T) {
	if nil == redisPool {
		t.Fatal("redisPool is nil")
	}
	testKey := "testKey"
	testValStr := "testVal"
	testValNum := int64(12)
	redisPool.Set(testKey, testValStr)
	str, _ := redisPool.Get(testKey)
	if str != testValStr {
		t.Fatal("set the val fail")
	}

	redisPool.Del(testKey)

	str, _ = redisPool.Get(testKey)
	if str == testValStr {
		t.Fatal("Del fail")
	}

	redisPool.Set(testKey, fmt.Sprintf("%d", testValNum))

	num, err := redisPool.Decr(testKey)
	if num != testValNum-1 || err != nil {
		t.Fatal("Decr fail")
	}

	num, err = redisPool.Incr(testKey)
	if num != testValNum || err != nil {
		t.Fatal("Incr fail")
	}

	redisPool.Del(testKey)

	lenAppend, err := redisPool.Append(testKey, testValStr)
	if err != nil || lenAppend != len([]byte(testValStr)) {
		t.Fatal("Append fail")
	}

	//expire
	_, err = redisPool.Expire(testKey, 1)
	if nil != err {
		t.Fatal("Expire fail" + err.Error())
	}
	time.Sleep(time.Millisecond * 1100)
	str, _ = redisPool.Get(testKey)
	if str == testValStr {
		t.Fatal("Expire fail")
	}

	redisPool.Set(testKey, testValStr)

	_, err = redisPool.PExpire(testKey, 1000)
	if nil != err {
		t.Fatal("PExpire fail" + err.Error())
	}
	time.Sleep(time.Millisecond * 1100)
	str, _ = redisPool.Get(testKey)
	if str == testValStr {
		t.Fatal("PExpire fail")
	}

	redisPool.Del(testKey)
	err = redisPool.SetEx(testKey, testValStr, 1)
	if nil != err {
		t.Fatal("SetEx fail" + err.Error())
	}
	time.Sleep(time.Millisecond * 1100)
	str, err = redisPool.Get(testKey)
	fmt.Println(err)
	if str == testValStr {
		t.Fatal("SetEx fail", err)
	}

	redisPool.Del(testKey)
}

func TestRedisPoolStringJson(t *testing.T) {
	if nil == redisPool {
		t.Fatal("redisPool is nil")
	}

	testKey := "testKeyJson"
	testStrVal := "testStrVal"
	type testStruct struct {
		Test1 string
		Test2 bool
		Test3 int
	}
	testStructVal := testStruct{"hello", true, 99}

	err := redisPool.SetJson(testKey, testStrVal)
	if nil != err {
		t.Fatal("SetJson fail:" + err.Error())
	}

	var testReply string
	err = redisPool.GetJson(testKey, &testReply)
	if nil != err {
		t.Fatal("GetJson fail:" + err.Error())
	}
	fmt.Println("testReply:", testReply)

	err = redisPool.SetExJson(testKey, testStructVal, 1)
	if nil != err {
		t.Fatal("SetExJson fail:" + err.Error())
	}

	var testReplyStruct testStruct
	err = redisPool.GetJson(testKey, &testReplyStruct)
	if nil != err {
		t.Fatal("SetExJson fail:" + err.Error())
	}
	fmt.Println("testReplyStruct:", testReplyStruct)

	time.Sleep(time.Millisecond * 1100)

	err = redisPool.GetJson(testKey, &testReplyStruct)
	if nil != err && ErrNil != err {
		t.Fatal("SetExJson fail:" + err.Error())
	}
	fmt.Println("testReplyStruct:", testReplyStruct)

}

func TestRedisPoolMap(t *testing.T) {
	if nil == redisPool {
		t.Fatal("redisPool is nil")
	}

	testKey := "testKeyMap"
	testfield1 := "testfield1"
	testfield2 := "testfield2"
	testfield3 := "testfield3"
	testVal1 := "testVal1"
	testVal2 := "testVal2"
	testVal3 := "testVal3"

	err := redisPool.HMSet(testKey, testfield1, testVal1, testfield2, testVal2)
	if nil != err {
		t.Fatal("HMSet fail" + err.Error())
	}

	row, err := redisPool.HSet(testKey, testfield3, testVal3)
	if nil != err {
		t.Fatal("HSet fail" + err.Error())
	}
	fmt.Println("row:", row)

	str, err := redisPool.HGet(testKey, testfield3)
	if nil != err || str != testVal3 {
		t.Fatal("HGet fail" + err.Error())
	}

	strMap, err := redisPool.HGetAll(testKey)
	if nil != err || strMap[testfield3] != testVal3 {
		t.Fatal("HGetAll fail" + err.Error())
	}
	fmt.Println(strMap)

	lenMap, err := redisPool.HLen(testKey)
	if nil != err || lenMap != 3 {
		t.Fatal("HLen fail" + err.Error())
	}

	row, err = redisPool.HDel(testKey, testfield3)
	if nil != err || row != 1 {
		t.Fatal("HDel fail", err)
	}
}

func TestRedisPoolList(t *testing.T) {
	testListKey := "testListKey"
	testNotExistListKey := "testNotExistListKey"
	testListVal1 := "testListVal1"
	testListVal2 := "testListVal2"
	testListVal3 := "testListVal3"

	// 清理
	_, err := redisPool.LRem(testListKey, 0, testListVal1)
	if nil != err {
		t.Fatal("LTrim fail", err)
	}
	_, err = redisPool.LRem(testListKey, 0, testListVal2)
	_, err = redisPool.LRem(testListKey, 0, testListVal3)

	intE, err := redisPool.LPush(testListKey, testListVal1, testListVal2, testListVal3)
	if nil != err {
		t.Fatal("LPush", err)
	}

	strVal, err := redisPool.LPop(testListKey)
	if nil != err || strVal != testListVal3 {
		t.Fatal("LPop fail", err)
	}
	fmt.Println(strVal)

	err = redisPool.LTrim(testListKey, 0, 0)
	if nil != err {
		t.Fatal("LTrim fail", err)
	}

	// 清理
	_, err = redisPool.LRem(testListKey, 0, testListVal2)

	intE, err = redisPool.RPush(testListKey, testListVal1, testListVal2, testListVal3)
	if nil != err || intE != 3 {
		t.Fatal("RPush", err)
	}

	strVal, err = redisPool.RPop(testListKey)
	if nil != err || strVal != testListVal3 {
		t.Fatal("RPop fail", err)
	}
	fmt.Println(strVal)

	strVal, err = redisPool.LPop(testListKey)
	if nil != err || strVal != testListVal1 {
		t.Fatal("RPush fail", err)
	}
	fmt.Println(strVal)

	// pushx
	intE, err = redisPool.RPushX(testNotExistListKey, testListVal1)
	if nil != err || (nil == err && intE != 0) {
		t.Fatal("RPushX fail", err)
	}

	intE, err = redisPool.LPushX(testNotExistListKey, testListVal1)
	if nil != err || (nil == err && intE != 0) {
		t.Fatal("LPushX fail", err, intE)
	}

	// LRange
	strs, err := redisPool.LRange(testListKey, 0, 10)
	if nil != err {
		t.Fatal("LRange fail", err)
	}
	fmt.Println(strs)

	// LLen
	len, err := redisPool.LLen(testListKey)
	if nil != err || len != 1 {
		t.Fatal("LLen fail", err)
	}
	fmt.Println(strs)

	// block
	strMap, err := redisPool.BLPop(testListKey, 1)
	if nil != err || strMap[testListKey] != testListVal2 {
		t.Fatal("BLPop fail", err)
	}
	fmt.Println(strMap)

	strMap, err = redisPool.BLPop(testListKey, 1)
	if nil == err {
		t.Fatal("BLPop fail non")
	}
	fmt.Println(strMap, err)

	intE, err = redisPool.RPush(testListKey, testListVal1)
	if nil != err {
		t.Fatal("BRPop RPush", err)
	}

	strMap, err = redisPool.BRPop(testListKey, 1)
	if nil != err || strMap[testListKey] != testListVal1 {
		t.Fatal("BRPop fail", err)
	}
	fmt.Println(strMap, time.Now().Unix())

	strMap, err = redisPool.BRPop(testListKey, 1)
	if nil == err {
		t.Fatal("BRPop fail non")
	}
	fmt.Println(strMap, err, time.Now().Unix())

}

func TestSync(t *testing.T) {

}
