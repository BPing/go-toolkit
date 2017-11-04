package producerConsumer

import (
	"errors"
	"time"
)

// 新建
func NewContainer(config Config) (IContainer, error) {

	switch config.Type {
	case ChannelType:
		container, err := NewContainerPC(int(config.MsgLen), config.ConsumeFunc)
		if err != nil {
			return nil, err
		}
		container.SetAssistIdleKeepAlive(time.Duration(config.AssistIdleKeepAlive) * time.Second)
		container.Record = config.Record
		return container, nil
	case CacheType:
		container, err := NewContainerCachePC(config.CacheInstance, config.ConsumeFunc, config.Unmarshal, config.Marshal)
		if err != nil {
			return nil, err
		}
		container.Record = config.Record
		container.MsgLen = config.MsgLen
		container.ReadTimeout = config.AssistIdleKeepAlive
		return container, nil

	default:
		return nil, errors.New("invalid type")
	}
	return nil, errors.New("fail to create container")
}
