package models

import "net/http"

//@Title  返回数据结构处理
//@Author cbping
type RespInfo struct {
	F_responseNo  int         `required:"true" description:"响应码"`
	F_responseMsg string      `description:"响应码描述"`
	F_data        interface{} `description:"响应内容"`
}

//响应码
//@Description
const (
	//commom
	RESP_OK          = 10000
	RESP_ERR         = 10001
	RESP_PARAM_ERR   = 10002
	RESP_TOKEN_ERR   = 10003
	RESP_NO_ACCESS   = 10004
	RESP_APP_NOT_ON  = 10005
	RESP_UNKNOWN_ERR = 10006
)

//响应码描述
var respMsg map[int]string = map[int]string{
	RESP_OK:          "操作成功",
	RESP_ERR:         "操作失败",
	RESP_PARAM_ERR:   "参数错误",
	RESP_TOKEN_ERR:   "签名认证错误",
	RESP_NO_ACCESS:   "对不起，您没有此操作权限",
	RESP_APP_NOT_ON:  "暂时未提供服务",
	RESP_UNKNOWN_ERR: "未知错误",

	http.StatusBadRequest:          "参数有误（Bad Request）",
	http.StatusUnauthorized:        "用户认证不成功",
	http.StatusForbidden:           "拒绝访问",
	http.StatusNotFound:            "资源不存在",
	http.StatusMethodNotAllowed:    "此方法未允许访问",
	http.StatusRequestTimeout:      "请求远程服务器超时",
	http.StatusInternalServerError: "服务器内部错误",
	http.StatusServiceUnavailable:  "服务不可用",
}

//
func NewRespInfo(code int, msg string, data interface{}) (respInfo *RespInfo) {
	if "" == msg {
		msg = GetMsgWithCode(code)
	}
	respInfo = &RespInfo{code, msg, data}
	return
}

//
func NewRespInfoWithCode(code int, data interface{}) (respInfo *RespInfo) {
	respInfo = NewRespInfo(code, "", data)
	return
}

func GetMsgWithCode(code int) string {
	val, ok := respMsg[code]
	if !ok {
		val = ""
	}
	return val
}

// @Title 转化成map
// @Description
// @Return map[string]interface{}
func (ri *RespInfo) ToStringMap() (res map[string]interface{}) {
	res = make(map[string]interface{})
	res["F_responseNo"] = ri.F_responseNo
	res["F_responseMsg"] = ri.F_responseMsg
	if nil != ri.F_data {
		res["F_data"] = ri.F_data
	}
	return
}

// @Title SetData
// @Description
func (ri *RespInfo) SetData(data interface{}) {
	ri.F_data = data
}
