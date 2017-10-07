[尝试使用新的版本](https://github.com/BPing/go-toolkit/tree/master/http-client)

# Curl
    curl请求。
    
     
* Client：封装http.Client.修饰http.Client。处理请求过程加入一些必要的自定义处理。如：慢请求记录
* Request：接口类型。构建并返回http.Request
* Response：封装http.Response.修饰http.Response。集成一些常用的处理响应内容方法。如：`ToJson()` 返回json格式内容

# TODO
1、加上熔断模式机制？
2、请求失败重试？

