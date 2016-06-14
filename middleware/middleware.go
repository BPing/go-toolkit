package middleware

import "net/http"

// 把中间件和原本的路由处理器封装在一起， 先执行中间件，如果中间件没有提前结束请求， 最终会把执行权归还给原本的路由处理器。
// 中间件允许注册多个，执行顺序和注册顺序一致。 其实原本的路由处理器也可以看做一个中间件了，不过，它是放在最后一个执行位置上（除了末尾的空中间件）。
// 源自开源项目：https://github.com/urfave/negroni


//中间件接口
type MiddleWare interface {
	ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type MiddleWareFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (h MiddleWareFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h(rw, r, next)
}

//中间件列表结构
//适配器模式
type middlewareHandler struct {
	handler MiddleWare
	next    *middlewareHandler
}

func (m *middlewareHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if (nil == m.handler) {
		return
	}
	m.handler.ServeHTTP(rw, r, m.next.ServeHTTP)
}

//包装标准库的http.Handler
func Wrap(handler http.Handler) MiddleWare {
	return MiddleWareFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		handler.ServeHTTP(rw, r)
		next(rw, r)
	})
}

//Cbping是一堆中间件处理程序管理器，
//可以当作http.handler被调用
//通过RegisterMiddlewareHandler注册中间件
type Cbping struct {
	//链表头
	//由中间件和路由处理器组建而成
	//路由处理器处于链表末端（除了末尾的空中间件）
	middlewareHead middlewareHandler
	//中间件数组
	middlewares    []MiddleWare

	//路由处理器
	//原始路由处理器
	mux            http.Handler
}

func (c *Cbping) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	//解析参数
	r.ParseForm()
	//中间件处理
	c.middlewareHead.ServeHTTP(rw, r)
}

//引导初始
func (c *Cbping)Bootstrap() {
	if (nil == c.mux) {
		c.mux = http.DefaultServeMux
	}
	c.middlewares = append(c.middlewares, Wrap(c.mux))
	c.middlewareHead = build(c.middlewares)
}

//运行
func (c *Cbping)Run(addr string) {
	c.Bootstrap()
	http.ListenAndServe(addr, c)
}

func (c *Cbping)RegisterMiddlewareHandlFunc(handlers... func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)) {
	for _, handler := range handlers {
		c.RegisterMiddleWare(MiddleWareFunc(handler))
	}

}

//注册中间件
//中间件执行顺序和注册顺序一致
func (c *Cbping)RegisterMiddleWare(handler MiddleWare) {
	c.middlewares = append(c.middlewares, handler)
}

//注册原本路由处理器
func (c *Cbping)MuxHandler(muxHandler http.Handler) {
	c.mux = muxHandler
}

func New() *Cbping {
	return &Cbping{}
}

//递归构建执行链表
func build(handlers []MiddleWare) middlewareHandler {
	var next middlewareHandler

	if len(handlers) == 0 {
		return voidMiddlewareHandler()
	} else if len(handlers) > 1 {
		next = build(handlers[1:])
	} else {
		next = voidMiddlewareHandler()
	}

	return middlewareHandler{handlers[0], &next}
}

func voidMiddlewareHandler() middlewareHandler {
	return middlewareHandler{
		MiddleWareFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {}),
		&middlewareHandler{},
	}
}