package producerConsumer

import (
	"errors"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	// （秒）
	DefaultReadTimeout = 3
)

var (
	ErrMarshalNil       = errors.New("ContainerRedis: Marshal is nil")
	ErrConsumeFuncNil   = errors.New("ContainerRedis: consumeFunc is nil")
	ErrUnmarshalNil     = errors.New("ContainerRedis: unmarshal is nil")
	ErrCacheInstanceNil = errors.New("ContainerRedis: cacheInstance is nil")
)

//  以Cache(如：redis)队列作为信息队列
type ContainerRedis struct {
	CBaseInfo
	// 消费信息的函数
	// 信息体最终落到此函数处理
	// 由用户自定义函数实体内容
	consumeFunc func(IMessage)

	// 把Cache队列中的字符串信息解析成相应的信息结构体
	Unmarshal func([]byte) (IMessage, error)

	// 把信息结构体序列化成字符串，以便保存到Cache(如：redis)队列中
	Marshal func(IMessage) ([]byte, error)

	// Cache(如：redis)操作实例
	//    实现接口cache
	cacheInstance ICache

	// 读取信息列表超时时间(秒)
	//  如果列表没有元素会阻塞列表直到等待超时
	ReadTimeout int64

	// 信息队列长度。
	// 如果为零，代表程序不限制（则以Cache(如：redis)的队列最大限制为准）
	// 但不保证真正Cache(如：redis)队列长度一定会小于此值。
	// 另外，此值关系到是否有协助协程产生。
	MsgLen int64

	// cache 队列名字（唯一标识）
	CacheListKey string
}

// 新建生产/消费模式容器
func NewContainerCachePC(cacheInstance ICache, consumeFunc func(IMessage), unmarshal func([]byte) (IMessage, error), marshal func(IMessage) ([]byte, error)) (*ContainerRedis, error) {

	if cacheInstance == nil {
		return nil, ErrCacheInstanceNil
	}

	if consumeFunc == nil {
		return nil, ErrConsumeFuncNil
	}
	if unmarshal == nil {
		return nil, ErrUnmarshalNil
	}
	if marshal == nil {
		return nil, ErrMarshalNil
	}
	return &ContainerRedis{
		consumeFunc:   consumeFunc,
		Unmarshal:     unmarshal,
		Marshal:       marshal,
		CacheListKey:  "ContainerCache-" + strconv.FormatInt(time.Now().UnixNano(), 10),
		cacheInstance: cacheInstance,
	}, nil
}

// 消费
func (cr *ContainerRedis) Consume() error {
	if cr.cacheInstance == nil {
		return ErrCacheInstanceNil
	}
	cr.consume(MasterRunner, nil)
	return nil
}

// 生产
func (cr *ContainerRedis) Produce(msg IMessage) error {
	if cr.cacheInstance == nil {
		return ErrCacheInstanceNil
	}
	if cr.MsgLen > 0 {
		if llen, err := cr.cacheInstance.LLen(cr.CacheListKey); err == nil && llen >= cr.MsgLen {
			cr.consume(AssistRunner, msg)
			return nil
		}

	}
	if cr.Marshal == nil {
		return ErrMarshalNil
	}
	msgBytes, err := cr.Marshal(msg)
	if nil == err {
		_, err = cr.cacheInstance.RPush(cr.CacheListKey, msgBytes)
	}
	return err
}

//
func (cr *ContainerRedis) NumGoroutine() (master, assistActive int64) {
	master = atomic.LoadInt64(&cr.masterNum)
	assistActive = atomic.LoadInt64(&cr.assistActiveNum)
	return
}

// 开启消费协程
// @master 是否主要消费协程。
//         主要消费协程一直执行
//         协助协程是在消息过多的时候开启，在没有消息体的时候结束。
// @argMsg 队列已满，放不进去的消息，协助协程消费的第一个消息。
func (cr *ContainerRedis) consume(master bool, argMsg IMessage) {
	if cr.ReadTimeout <= 0 {
		cr.ReadTimeout = DefaultReadTimeout
	}
	if master == MasterRunner {
		go func() {
			defer cr.catch(MasterRunner)
			//delayTime := 125
			for {
				msgMap, err := cr.cacheInstance.BLPop(cr.CacheListKey, cr.ReadTimeout)
				if nil == msgMap {
					cr.record(DebugTag, err)
					//time.Sleep(time.Millisecond * delayTime)
					//delayTime = 2 * delayTime
					//if delayTime > 1000 {
					//	delayTime = 1000
					//}
					continue
				}
				//delayTime=125
				if nil != cr.Unmarshal {
					msg, err := cr.Unmarshal([]byte(msgMap[cr.CacheListKey]))
					if nil == err {
						cr.consumeFunc(msg)
					} else {
						cr.record(ErrTag, err)
					}

				}

			}
		}()
		atomic.AddInt64(&cr.masterNum, 1)

	} else if master == AssistRunner {

		atomic.AddInt64(&cr.assistActiveNum, 1)

		go func() {
			defer cr.catch(AssistRunner)
			//先消费放不进队列的消息
			if nil != cr.consumeFunc && nil != argMsg && argMsg.Id() != "" {
				cr.consumeFunc(argMsg)
			}
			for {
				msgMap, err := cr.cacheInstance.BLPop(cr.CacheListKey, cr.ReadTimeout)
				if nil == msgMap {
					cr.record(DebugTag, err)
					return
				}
				if nil != cr.Unmarshal {
					msg, err := cr.Unmarshal([]byte(msgMap[cr.CacheListKey]))
					if nil == err {
						cr.consumeFunc(msg)
					} else {
						cr.record(ErrTag, err)
					}
				}
			}
		}()
	}
}

// 协程收尾工作，比如：捕捉异常、主要协程恢复
func (cr *ContainerRedis) catch(master bool) {
	if master == MasterRunner {
		atomic.AddInt64(&cr.masterNum, -1)
	} else {
		atomic.AddInt64(&cr.assistActiveNum, -1)
	}
	if err := recover(); err != nil {
		cr.record(PanicTag, err)
		if master == MasterRunner {
			// 开启新的主要协程
			cr.Consume()
		}
	}
}
