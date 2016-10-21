package producerConsumer

// 容器接口
// 实现基于生产/消费模式
//
//   1、Produce(msg interface{}) 生产信息，把消息放入消息列表中。
//   2、Consume() 消费消息。
type IContainer interface {
	// 生产消息
	Produce(msg IMessage) error
	// 消费消息
	Consume() error
	// 消费消息的协程数目
	NumGoroutine() (master, assistActive int64)
}


// 消息接口
type IMessage interface {
	// 标识
	// 此字段不能为空，
	// 否则会被当做无效数据抛弃。
	// 容器不对此标识的唯一性感兴趣。
	// 用户可以自行确保此标识的唯一性
	Id() string
}