package ErrNum

import (
	"sync"
)

const ERR_OK = 0
const ERR_UNKNOW = 10
const ERR_NOT_SUP_METHOD = 11
const ERR_PARAMS_ERR = 12

const ERR_UP_C_G_ERR = 10001
const ERR_DEL_C_G_ERR = 10002
const ERR_G_EXIST = 10010

type ErrNum struct {
	Name string
}

var en map[int] string
var once sync.Once

func GetInstance() map[int] string {
	once.Do(func() {
		en = make(map[int] string)
		en[ERR_OK] = ""
		en[ERR_UNKNOW] = "未知错误"
		en[ERR_NOT_SUP_METHOD] = "不支持的请求方法"
		en[ERR_PARAMS_ERR] = "参数有误"
		en[ERR_UP_C_G_ERR] = "更新缓存分组失败"
		en[ERR_DEL_C_G_ERR] = "删除缓存分组失败"
		en[ERR_G_EXIST] = "接口分组已经存在"
	})
	return en
}

func GetErrMsg(errcode int) string{
	en := GetInstance()
	return en[errcode]
}