package client

import "time"

const (
	// 默认失败尝试最大次数
	defaultMaxBadRetryCount = 2

	// 默认慢请求时间
	defaultSlowReqLong = 5 * time.Second
)
