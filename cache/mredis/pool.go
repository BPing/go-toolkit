package mredis

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"time"
)

// 配置信息
type PoolConfig struct {
	Network        string
	Address        string
	MaxIdle        int
	DbNum          int
	Password       string
	IdleTimeout    time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	ConnectTimeout time.Duration
	Df             func(config PoolConfig) (redis.Conn, error)
}

// 池
// 获取一个可用的redis连接。
// 如果没有可用的它将新建一个redis连接
// 主要方法：
//  pool.Get() 从池中获取连接
//  pool.Put() 重新放回池中
type Pool struct {
	config PoolConfig
	pool   chan redis.Conn
	df     func(config PoolConfig) (redis.Conn, error)
}

// 新建池
//
// 例子：
/*         dialFunc := func(config PoolConfig) (c redis.Conn, err error) {
c, err = redis.Dial(config.Network, config.Address)
if err != nil {
	return nil, err
}

if config.Password != "" {
	if _, err := c.Do("AUTH", config.Password); err != nil {
		c.Close()
		return nil, err
	}
}

_, selecterr := c.Do("SELECT", config.DbNum)
if selecterr != nil {
	c.Close()
	return nil, selecterr
}
       return
}

config := PoolConfig{
	Network  :"tcp",
	Address  : "127.0.0.1:6379",
	MaxIdle  :10,
	Password :"123456",
	DbNum    :0,
	Df       :dialFunc,
}

*/
//
// 	p, _ := NewPool(config)
//
func NewPool(config PoolConfig) (*Pool, error) {
	if config.Df == nil {
		return nil, errors.New("dialFunc is nil")
	}
	pool := make([]redis.Conn, 0, config.MaxIdle)
	for i := 0; i < config.MaxIdle; i++ {
		client, err := config.Df(config)
		if err != nil {
			for _, client = range pool {
				client.Close()
			}
			return nil, err
		}
		if client != nil {
			pool = append(pool, &idleConn{client, time.Now()})
		}
	}
	p := Pool{
		config: config,
		pool:   make(chan redis.Conn, config.MaxIdle),
		df:     config.Df,
	}
	for i := range pool {
		p.pool <- pool[i]
	}
	return &p, nil
}

// 空闲连接
type idleConn struct {
	redis.Conn
	t time.Time
}

//type DialFunc func(network, addr string) (redis.Conn, error)

// 获取一个可用的redis连接。如果没有可用的它将新建一个redis连接
func (p *Pool) Get() (redis.Conn, error) {
	select {
	case conn := <-p.pool:
		return conn, nil
	default:
		return p.df(p.config)
	}
}

// 把一个redis连接放回池中。如果池已经满了，将关闭此redis连接。
// 如果此redis连接已经关闭（由于连接失败或其他原因），则不应该再放回池中。
func (p *Pool) Put(conn redis.Conn) error {
	select {
	case p.pool <- conn:
		return nil
	default:
		return conn.Close()
	}
}

// 移除并关闭池中所有的redis连接。
// 假设没有其他的连接等着被放回
// 这个方法有效地关闭和清理池。
func (p *Pool) Empty() {
	var conn redis.Conn
	for {
		select {
		case conn = <-p.pool:
			conn.Close()
		default:
			return
		}
	}
}
