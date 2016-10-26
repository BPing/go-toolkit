package producerConsumer

// 容器接口
//   生产/消费模式
//
//    通过调用`Consume()`可以产生一个主要消费协程。主协程将一直存在，在没有消息体处理的时候进入阻塞等待。
// 可以通过调用`Consume()`的次数来控制产生主协程的数目。
//    当消息体队列的已满，则会产生协助协程消费消息体。协助协程在消息体猛涨时候出现，在没有消息体处理的时候
// 阻塞等待一定时间后将被销毁。
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

//// 消息队列
//type IMessageList interface {
//        // 移出并获取列表的第一个消息
//	// 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止
//	BPop(timeout int64)(IMessage, error)
//
//	// 在列表尾部中添加一个消息
//	Push(IMessage) (int64, error)
//
//	// 获取列表长度
//	LLen() (int64, error)
//}

// 缓存Cache接口
type ICache interface {
	// BLPOP key1 timeout(秒)
	// 移出并获取列表的第一个元素，
	// 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止。
	BLPop(key string,timeout int64)(map[string]string, error)

	// 在列表尾部中添加一个或多个值
        RPush(key string,values ... interface{}) (int64, error)

	// 获取列表长度
	LLen(key string) (int64, error)
}