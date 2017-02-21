package customType

import (
	"encoding/json"
)

type ReturnType struct {
	Errcode int
	Errmsg string
	Data interface{}
}

func InitReturnType() {

}

func (result *ReturnType) ToJsonStr() string{

	resStr, _ := json.Marshal(result)
	return string(resStr)
}
