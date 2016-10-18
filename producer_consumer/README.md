# producerConsumer
生产/消费模式

# 快速开始

```go  

    package main
   
    import (
        "fmt"
        "runtime"
        "strconv"
        "time"
        "github.com/BPing/Golib/producer_consumer"
    )
    
    var container *producerConsumer.Container
    
    func init(){
        consumeFunc := func(msg producerConsumer.Message) {
            fmt.Println("消费：", msg.Id, "协程数目：", runtime.NumGoroutine())
        }
    
        container, _ = producerConsumer.NewContainerPC(20, consumeFunc)
        container.Consume()
    }
    
    func main(){
    
        go func() {
            for i := 0; i < 50; i++ {
                msg, _ := producerConsumer.NewMessage("goone-"+strconv.Itoa(i), nil)
                container.Produce(msg)
            }
        }()
    
        go func() {
            for i := 0; i < 50; i++ {
                msg, _ := producerConsumer.NewMessage("gotwo-"+strconv.Itoa(i), nil)
                container.Produce(msg)
                time.Sleep(time.Millisecond * 20)
            }
        }()
    
        time.Sleep(time.Second * 3)
    }

```

# 描述

* 通过调用`Consume()`可以产生一个主要消费协程。主协程将一直存在，在没有消息体处理的时候进入阻塞等待。
  可以通过调用`Consume()`的次数来控制产生主协程的数目。
* 当消息体队列的已满，则会产生协助协程消费消息体。协助协程在消息体猛涨时候出现，在没有消息体处理的时候阻塞等待一定时间后将被销毁。
  协助协程数目不作上限控制。
  
