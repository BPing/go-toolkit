package producerConsumer

import (
	"fmt"
	"runtime"
	"strconv"
	"testing"
	"time"
)

type Message struct {
	Key string
}

func(msg *Message)Id()string{
	return msg.Key
}

func NewMessage(id string)*Message{
    return &Message{id}
}

func TestContainer(t *testing.T) {
	consumeFunc := func(msg IMessage) {
		fmt.Println("消费：", msg.Id(), "协程数目：", runtime.NumGoroutine())
	}

	container, _ := NewContainer(Config{
		Type:ChannelType,
		MsgLen:20,
		ConsumeFunc:consumeFunc,
		AssistIdleKeepAlive:1,
	})
	container.Consume()

	go func() {
		for i := 0; i < 50; i++ {
			msg:= NewMessage("goone-"+strconv.Itoa(i))
			container.Produce(msg)
		}
	}()

	go func() {
		for i := 0; i < 50; i++ {
			msg:= NewMessage("gotwo-"+strconv.Itoa(i))
			container.Produce(msg)
			time.Sleep(time.Millisecond * 20)
		}
	}()

	go func() {
		for i := 0; i < 50; i++ {
			msg:= NewMessage("gothree-"+strconv.Itoa(i))
			container.Produce(msg)
			time.Sleep(time.Millisecond * 100)
		}
	}()

	time.Sleep(time.Second * 3)
}

func TestContainerErr(t *testing.T) {
	consumeFunc := func(msg IMessage) {
		fmt.Println("消费：", msg.Id(), "协程数目：", runtime.NumGoroutine())
	}

	_, err := NewContainer(Config{
		Type:ChannelType,
		MsgLen:0,
		ConsumeFunc:consumeFunc,
	})

	//_, err = NewContainerPC(0, consumeFunc)
	if err != ChanLenErr {
		t.Fatal(err)
	}

	_, err = NewContainerPC(20, nil)

	if err != ConsumeFuncNilErr {
		t.Fatal(ConsumeFuncNilErr)
	}

}
