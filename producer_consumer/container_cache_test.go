package producerConsumer

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"testing"
	"time"
)

type TestRedis struct {
	msgList *list.List
}

func (tr *TestRedis) BLPop(key string, timeout int64) (map[string]string, error) {
	//fmt.Println("BLPop-----------------------")
	resp := make(map[string]string)
	ele := tr.msgList.Front()
	if ele != nil {
		tr.msgList.Remove(ele)
		resp[key] = ele.Value.(string)
		return resp, nil
	}
	time.Sleep(time.Second * time.Duration(timeout))
	return nil, errors.New("list is empty")
}

func (tr *TestRedis) RPush(key string, values ...interface{}) (int64, error) {
	//fmt.Println("RPush-----------------------")
	switch value := values[0].(type) {
	case []byte:
		tr.msgList.PushBack(string(value))
	case string:
		tr.msgList.PushBack(value)

	default:
		return 0, errors.New("type error ([]byte or string)")
	}
	return 0, nil
}

func (tr *TestRedis) LLen(key string) (int64, error) {
	//fmt.Println("LLen-----------------------")
	return int64(tr.msgList.Len()), nil
}

func TestContainerRedis(t *testing.T) {

	redisInstance := &TestRedis{msgList: list.New()}
	consumeFunc := func(msg IMessage) {
		if msg == nil {
			return
		}
		fmt.Println("c----------------")
		fmt.Println(msg.Id())
	}

	Unmarshal := func(msgByte []byte) (IMessage, error) {
		msg := &Message{}
		err := json.Unmarshal(msgByte, msg)
		//fmt.Println("Unmarshal",string(msgByte),err)
		return msg, err
	}

	Marshal := func(msg IMessage) ([]byte, error) {
		msgByte, err := json.Marshal(msg.(*Message))
		//fmt.Println("Marshal",string(msgByte),err)
		return msgByte, err
	}

	container, err := NewContainer(Config{
		Type:   CacheType,
		MsgLen: 10,
		CacheInstance:redisInstance,
		ConsumeFunc:consumeFunc,
		Unmarshal:Unmarshal,
		Marshal:Marshal,
		AssistIdleKeepAlive:1,
	})

	fmt.Println(err)
	container.Consume()

	go func() {
		for i := 0; i < 50; i++ {
			msg := NewMessage("goone-" + strconv.Itoa(i))
			container.Produce(msg)
		}
	}()

	go func() {
		for i := 0; i < 50; i++ {
			msg := NewMessage("gotwo-" + strconv.Itoa(i))
			container.Produce(msg)
			time.Sleep(time.Millisecond * 20)
		}
	}()

	go func() {
		for i := 0; i < 50; i++ {
			msg := NewMessage("gothree-" + strconv.Itoa(i))
			container.Produce(msg)
			time.Sleep(time.Millisecond * 100)
			fmt.Println(container.NumGoroutine())
		}
	}()

	time.Sleep(time.Second * 6)
}

func TestContainerRedisErr(t *testing.T) {
	redisInstance := &TestRedis{msgList: list.New()}
	consumeFunc := func(msg IMessage) {
		fmt.Println("消费：", msg.Id(), "协程数目：", runtime.NumGoroutine())
	}

	Unmarshal := func(msgByte []byte) (IMessage, error) {
		msg := &Message{}
		err := json.Unmarshal(msgByte, msg)
		//fmt.Println("Unmarshal",string(msgByte),err)
		return msg, err
	}

	Marshal := func(msg IMessage) ([]byte, error) {
		msgByte, err := json.Marshal(msg.(*Message))
		//fmt.Println("Marshal",string(msgByte),err)
		return msgByte, err
	}

	_, err := NewContainer(Config{
		Type:   CacheType,
		MsgLen: 0,
		CacheInstance:redisInstance,
	})

	//_, err = NewContainerPC(0, consumeFunc)
	if err != ErrConsumeFuncNil {
		t.Fatal(ErrConsumeFuncNil)
	}

	_, err = NewContainerCachePC(nil, consumeFunc,Unmarshal,Marshal)
	if err != ErrCacheInstanceNil {
		t.Fatal(ErrCacheInstanceNil)
	}

	_, err = NewContainerCachePC(redisInstance, consumeFunc,nil,Marshal)

	if err != ErrUnmarshalNil {
		t.Fatal(ErrUnmarshalNil)
	}

	_, err = NewContainerCachePC(redisInstance, consumeFunc,Unmarshal,nil)

	if err != ErrMarshalNil {
		t.Fatal(ErrMarshalNil)
	}
}
