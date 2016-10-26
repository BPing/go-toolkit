package producerConsumer

const (
	MasterRunner = true
	AssistRunner = false

	ErrTag   = "PC-ErrTag"
	DebugTag = "PC-DebugTag"
	PanicTag = "PC-PanicTag"

	CacheType=ContainerType("cache")
	ChannelType=ContainerType("channel")
)

type ContainerType string

// 基本属性
type CBaseInfo struct {

	// 活跃协助消费协程数目
	assistActiveNum int64

	// 主消费协程数目
	masterNum int64

	// 记录内部日志信息
	Record func(tag string, msg interface{})
}

func (cbi *CBaseInfo) record(tag string, msg interface{}) {
	if cbi.Record != nil && msg != nil {
		cbi.Record(tag, msg)
	}
}

// 初始配置项
type Config struct {

	// * channel型(ChannelType)：基于缓冲channel队列实现的。
	// * cache型(CacheType)：基于redis型队列实现的。
	Type ContainerType

	// 消费信息的函数
	// 信息体最终落到此函数处理
	// 由用户自定义函数实体内容
	ConsumeFunc func(IMessage)

	// 消息队列长度
	//  如果为channel型，此变量为int类型。请自行处理不一致。
	MsgLen int64

	// 空闲存活时间（针对AssistRunner）,单位为秒（s）
	//   如果是redis型的，此值等同于ReadTimeout。
	AssistIdleKeepAlive int64

	// 记录内部日志信息
	Record func(tag string, msg interface{})

	// redis型
	// ---------------------------------------------------------

	// redis 队列名字（唯一标识）
	//  针对redis型
	CacheListKey string

	// 把redis队列中的字符串信息解析成相应的信息结构体
	//  针对redis型
	Unmarshal func([]byte) (IMessage, error)

	// 把信息结构体序列化成字符串，以便保存到redis队列中
	//  针对redis型
	Marshal func(IMessage) ([]byte, error)

	// redis 操作实例 实现接口IRedis
	//  针对redis型
	CacheInstance ICache
}