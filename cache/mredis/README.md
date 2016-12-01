# redis
    redis常用命令操作封装。根据具体情况而定
    

# 快速开始
```go

dialFunc := func(config PoolConfig) (c redis.Conn, err error) {
c, err = redis.Dial(config.Network, config.Address)
if err != nil {
    return nil, err
}

if config.password != "" {
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

 config=PoolConfig{
        Network  :"tcp",
      Address  :   "127.0.0.1:6379",
      MaxIdle  :10,
      Password :"123456",
      DbNum    :0,
      Df       :dialFunc,
 }

redis, _ := NewRedisPool(config)

```

# API

* `字符串类型相关命令操作`

`GET`
```go
 redis.Get("cacheKey") 
```

`SET`
```go
redis.Set("cacheKey","value") 
```

* `哈希(Hash)类型相关命令操作`
 Redis hash 是一个string类型的field和value的映射表，hash特别适合用于存储对象。
 Redis 中每个 hash 可以存储 232 - 1 键值对（40多亿）。
  
`HGET`
```go
redis.HGet("cacheKey","field")
```

# 依赖
  
* go get github.com/garyburd/redigo/redis

# redis

* [redis 命令](#http://www.redis.net.cn/order/)
