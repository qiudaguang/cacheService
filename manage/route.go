package manage

import (
	"net/http"
	"reflect"
	"strings"
	"fmt"
	"encoding/json"
)

type ControllerType struct {
	Controll string
	Action string
}

func GetDefaultAction() string{
	return "Index"
}

func GetDefaultController() string{
	return "Manage"
}

func BuildContrType(urlPath string) ControllerType {

	contrType := ControllerType{GetDefaultController(), GetDefaultAction()}
	arr := strings.Split(urlPath, "?")
	newStr := strings.FieldsFunc(arr[0], func(s rune) bool{
		if s == '/' {
			return true
		} else {
			return false
		}
	})

	switch len(newStr) {
		case 0:
		case 1:
			contrType.Controll = newStr[0]
		default:
			contrType.Controll = newStr[0]
			contrType.Action = newStr[1]
	}

	return contrType
}

func LocateRequest(w http.ResponseWriter, r *http.Request) {

	contrType := BuildContrType(r.URL.Path)

	viewType := &ViewType{}
	object := reflect.ValueOf(viewType)
	method := object.MethodByName(contrType.Action)
	if method.IsValid() == false {
		method = object.MethodByName("NotFound")
	}

	input := make([]reflect.Value, 2)
	input[0] = reflect.ValueOf(w)
	input[1] = reflect.ValueOf(r)
	values := method.Call(input)

	if len(values) > 0 {
		fmt.Fprint(w, ToJsonStr(values))
	}
}

func ToJsonStr(values []reflect.Value) string{

	resStr, _ := json.Marshal(values)
	return string(resStr)
}