package tools

import (
	"strings"
	"sort"
	"encoding/json"
	//"fmt"
)

func InArrayStr(value string, strArr []string) bool{

	if value == "" || len(value) == 0 {
		return false
	}
	if len(strArr) == 0 {
		return false
	}

	for _, v := range strArr {
		if v == value {
			return true
		}
	}
	return false
}

func ArrayRemoveIndexForStr(arr []string, index int) []string {
	if len(arr) == 0 || index < 0 || index >= len(arr) {
		return arr
	}

	return append(arr[:index], arr[index+1:]...)
}

func ArrayRemoveValForStr(arr []string, val string) []string {
	if len(arr) == 0 || val == "" {
		return arr
	}

	index := -1
	for i, v := range arr {
		if v == val {
			index = i
			break
		}
	}

	return ArrayRemoveIndexForStr(arr, index)
}

func splitUrlParamsStr(urlParamsStr string) map[string]string {

	paramsMap := make(map[string]string)
	if len(urlParamsStr) == 0 {
		return paramsMap
	}

	urlParamsArr := strings.Split(urlParamsStr, "&")
	for _, v := range urlParamsArr {
		itemArr := strings.Split(v, "=")
		paramsMap[itemArr[0]] = itemArr[1]
	}
	return paramsMap
}

func SortJsonKey(jsonStr string, allowKeys []string) string {

	var strMap = make(map[string]interface{})
	var sortMap = make(map[string]interface{})
	var keyArr []string
	err := json.Unmarshal([]byte(jsonStr), &strMap)
	if err != nil {
		return ""
	}

	for k, _ := range strMap {

		if InArrayStr(k, allowKeys) {
			keyArr = append(keyArr, k)
		}
	}
	sort.Strings(keyArr)
	for _, v := range keyArr {
		sortMap[v] = strMap[v]
	}

	sortStr, err := json.Marshal(sortMap)
	if err != nil {
		return ""
	}
	return string(sortStr)
}

func SortUrlParamsStr(urlParamsStr string, allowKeys []string) string {

	if len(urlParamsStr) == 0 {
		return urlParamsStr
	}
	var allKeys []string
	paramsMap := make(map[string]string)
	urlParamsArr := strings.Split(urlParamsStr, "&")
	for _, v := range urlParamsArr {
		if v == "" {
			continue
		}
		itemArr := strings.Split(v, "=")
		if len(itemArr) == 0 || itemArr[0] == "" {
			continue
		}

		if InArrayStr(itemArr[0], allowKeys) {
			allKeys = append(allKeys, itemArr[0])
			paramsMap[itemArr[0]] = itemArr[1]
		}
	}
	if len(allKeys) == 0 {
		return ""
	}
	sort.Strings(allKeys)
	newStr := ""
	for _, k := range allKeys {
		newStr += k + "=" + paramsMap[k] + "&"
	}

	return string(newStr[:len(newStr)-1])
}