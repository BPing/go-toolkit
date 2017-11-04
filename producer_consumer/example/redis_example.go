package example

import (
	"encoding/json"
	"fmt"
	"github.com/BPing/go-toolkit/cache/mredis"
	. "github.com/BPing/go-toolkit/producer_consumer"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

type Message struct {
	Key string
}

func (msg *Message) Id() string {
	return msg.Key
}

func NewMessage(id string) *Message {
	return &Message{id}
}

var containerRedis *ContainerRedis

func init() {

	dialFunc := func(config mredis.PoolConfig) (c redis.Conn, err error) {
		c, err = redis.Dial(
			config.Network,
			config.Address,
			redis.DialDatabase(config.DbNum),
			redis.DialPassword(config.Password))
		return
	}

	MRedis, irErr := mredis.NewRedisPool(mredis.PoolConfig{
		Address:  "127.0.0.1:6379",
		Network:  "tcp",
		Password: "",
		MaxIdle:  10,
		DbNum:    0,
		Df:       dialFunc})
	if irErr != nil {
		fmt.Println("init redis error:", irErr)
	} else {
		MRedis.SetPowerOn(true)
	}
	containerRedis, _ = NewContainerCachePC(MRedis, func(msg IMessage) {
		//time.Sleep(time.Millisecond*100)
		fmt.Println("c----------------")
		if msg == nil {
			return
		}
		fmt.Println(msg.Id())
	}, func(msgByte []byte) (IMessage, error) {
		msg := &Message{}
		err := json.Unmarshal(msgByte, msg)
		//fmt.Println("Unmarshal",string(msgByte),err)
		return msg, err
	}, func(msg IMessage) ([]byte, error) {
		msgByte, err := json.Marshal(msg.(*Message))
		//fmt.Println("Marshal",string(msgByte),err)
		return msgByte, err
	})
}

func RedisExample() {
	containerRedis.ReadTimeout = 1
	containerRedis.MsgLen = 10
	containerRedis.Record = func(tag string, msg interface{}) {
		fmt.Println(tag, msg)
	}
	containerRedis.Consume()

	go func() {
		for i := 0; i < 50; i++ {
			msg := NewMessage("goone-" + strconv.Itoa(i))
			containerRedis.Produce(msg)
		}
	}()

	go func() {
		for i := 0; i < 50; i++ {
			msg := NewMessage("gotwo-" + strconv.Itoa(i))
			containerRedis.Produce(msg)
			time.Sleep(time.Millisecond * 20)
		}
	}()

	go func() {
		for i := 0; i < 50; i++ {
			msg := NewMessage("gothree-" + strconv.Itoa(i))
			containerRedis.Produce(msg)
			time.Sleep(time.Millisecond * 100)
			fmt.Println(containerRedis.NumGoroutine())
		}
	}()

	time.Sleep(time.Second * 20)
}
