package cacheService

import (
	"../redisConn"
	"github.com/garyburd/redigo/redis"
	"errors"
	"../sendRequest"
	"strings"
	"fmt"
	"../customType"
	"../config"
	"../tools"
	"sync"
	"strconv"
	"net/url"
)

const API_CACHE_KEY_LIST = "API_CACHE_KEY_LIST"
const API_ITEM_CACHE = "API_ITEM_CACHE"
const API_VISITS_COUNT = "API_VISITS_COUNT"

//var redisVar, _ = redisConn.GetRedis()
var pool = redisConn.NewPool()
var conf  = config.GetConfig()
var once sync.Once
var once2 sync.Once
var groupList []customType.ApiConfType
var CacheConfType customType.CacheConfType
var lock sync.Mutex

func ClearCache() {

	redisC := pool.Get()
	defer redisC.Close()

	groupList := GetGroupConfList()
	for _, groupItem := range groupList {

		itemList := GetGroupItem(groupItem.UrlPath)
		for _, itemKey := range itemList {
			redisC.Do("DEL", itemKey)
		}
		redisC.Do("DEL", API_CACHE_KEY_LIST + ":" + groupItem.UrlPath)
		redisC.Do("DEL", API_VISITS_COUNT + ":" + groupItem.UrlPath)
	}
}

func InitCacheConfig() error{

	ClearCache()

	InitCacheServiceConfig()
	InitGroupConfList()
	return nil
}

func InitCacheServiceConfig() error{
	once2.Do(func() {
		expireTime := conf.String("cache_service::expiretime")
		expireTimeInt, err := strconv.Atoi(expireTime); if err == nil {
			CacheConfType.ExpireTime = expireTimeInt
		}
	})
	return nil
}

func InitGroupConfList() {

	once.Do(func() {
		//var err error
		groupListNew, err := config.GetApiConfig(); if err == nil {
			groupList = groupListNew
		}
	})
}

func BuildCacheKey(groupName string, paramsStr string, reqHead map[string][]string, groupConfig customType.ApiConfType) (string, error) {

	if paramsStr == "" || len(groupConfig.ParamsArr) == 0{
		key := API_ITEM_CACHE + ":" + groupName
		return url.QueryEscape(key), nil
	}

	contentType := ""
	if len(reqHead["Content-Type"]) > 0 {
		contentType = reqHead["Content-Type"][0]
	}

	switch contentType {
		case "application/json":
			paramsStr = tools.SortJsonKey(paramsStr, groupConfig.ParamsArr)
		default://默认为get url
			paramsStr = tools.SortUrlParamsStr(paramsStr, groupConfig.ParamsArr)
	}

	/*paramsStr = strings.Replace(paramsStr, "\n", "", -1)
	paramsStr = strings.Replace(paramsStr, "\r", "", -1)*/
	if paramsStr == "" {
		key := API_ITEM_CACHE + ":" + groupName
		return url.QueryEscape(key), nil
	}
	key := API_ITEM_CACHE + ":" + groupName + "?" + paramsStr
	return url.QueryEscape(key), nil
}

func LookupCache(groupName string, url string, paramsStr string, method string, reqHead map[string][]string) (string, error){

	if tools.InArrayStr(groupName, GetGroupNameList()) == false {
		return "", errors.New("not found group interface")
	}
	groupConfig, err := GetGroupConfig(groupName)
	if err != nil {
		return "", err
	}

	go func(groupName string) {
		IncrementVisitsCount(groupName)
		isUpdate := CheckUpdateByVisistCount(groupConfig)
		if isUpdate == true {
			UpdateCacheForApiGroup(groupName)
		}
	}(groupName)

	cacheKey, _ := BuildCacheKey(groupName, paramsStr, reqHead, groupConfig)
	respType, err := GetApiItemCache(cacheKey); if err != nil {
		return "", err
	}
	if len(respType.Resp) > 0 {
		//fmt.Println("cache: fit")
		return respType.Resp, nil
	}
	fmt.Println("cache: miss")

	respStr, respErr := UpdateCacheForApi(groupName, url, paramsStr, method, cacheKey, groupConfig)
	return respStr, respErr
}

func UpdateAllCache() {
	groupList := GetGroupConfList()
	for _, group := range groupList {

		UpdateCacheForApiGroup(group.UrlPath)
	}
}

func UpdateCacheForApiGroup(groupName string) error{

	itemList := GetGroupItem(groupName)
	if len(itemList)  == 0 {
		return errors.New("没有需要更新的接口")
	}

	itemconfig, _ := GetGroupConfig(groupName)
	for _, itemKey := range itemList {

		respType, _ := GetApiItemCache(itemKey)
		if respType.ApiUrl == "" {
			continue
		}
		UpdateCacheForApi(groupName, respType.ApiUrl, respType.ParamsStr, respType.Method, itemKey, itemconfig)
	}
	return nil
}

func UpdateCacheForApi(groupName string, apiUrl string, paramsStr string, method string, cacheKey string, itemconfig customType.ApiConfType) (string, error){

	var respStr string
	var respErr error
	method = strings.ToUpper(method)
	switch method {
		case "POST":
			respStr, respErr = sendRequest.SendHttpPost(apiUrl, paramsStr)
		case "GET":
			respStr, respErr = sendRequest.SendHttpGet(apiUrl)
		default:
			return "", errors.New("invalid method")
	}

	if respErr != nil {
		return respStr, respErr
	}
	var expireTime = 3600
	if itemconfig.ExpireTime > 0 {
		expireTime = itemconfig.ExpireTime
	}else {
		expireTime = CacheConfType.ExpireTime
	}

	var args =  []interface{}{cacheKey}
	args = append(args, "Resp", respStr)
	args = append(args, "ParamsStr", paramsStr)
	args = append(args, "Method", method)
	args = append(args, "ApiUrl", apiUrl)

	redisC := pool.Get()
	defer redisC.Close()
	redisC.Do("HMSET", args...)
	redisC.Do("Expire", cacheKey, expireTime)

	AddUrlToGroupItem(groupName, cacheKey)
	return respStr, nil
}


func GetVisitsCount(groupName string) int{
	redisC := pool.Get()
	defer redisC.Close()
	visistCount, _ := redis.Int(redisC.Do("GET", API_VISITS_COUNT + ":" + groupName))
	return visistCount
}

func IncrementVisitsCount(groupName string) {

	key := API_VISITS_COUNT + ":" + groupName

	redisC := pool.Get()
	defer redisC.Close()
	redisC.Do("INCR", key)
}

func AddUrlToGroupItem(groupName string, cacheKey string) (error) {

	redisC := pool.Get()
	defer redisC.Close()
	_, err := redisC.Do("SADD", API_CACHE_KEY_LIST + ":" + groupName, cacheKey)
	return err
}

func RemoveFromGroupItem(groupName string, cacheKey string) error {

	redisC := pool.Get()
	defer redisC.Close()
	_, err := redisC.Do("SREM", API_CACHE_KEY_LIST + ":" + groupName, cacheKey)
	return err
}

func GetGroupConfList() []customType.ApiConfType {
	return groupList
}

func GetGroupNameList() []string {

	var groupNameList []string
	groupConfList := GetGroupConfList()
	for _, groupConf := range groupConfList {
		groupNameList = append(groupNameList, groupConf.UrlPath)
	}
	return groupNameList
}

func GetGroupItem(groupName string) []string {

	redisC := pool.Get()
	defer redisC.Close()
	list, _ := redis.Strings(redisC.Do("SMEMBERS", API_CACHE_KEY_LIST + ":" + groupName))
	return list
}

func GetApiItemCache(cacheKey string) (customType.RespCacheType, error){

	var cache customType.RespCacheType

	redisC := pool.Get()
	defer redisC.Close()
	cacheVal, err := redis.Values(redisC.Do("HGETALL", cacheKey))

	if err != nil {
		return cache, err
	}

	err2 := redis.ScanStruct(cacheVal, &cache)
	if err2 != nil {
		return cache, err2
	}
	return cache, nil
}

func GetGroupConfig(groupName string) (customType.ApiConfType, error){

	for _, groupType := range groupList {

		if groupType.UrlPath == groupName {
			return groupType, nil
		}
	}
	return customType.ApiConfType{}, errors.New("not found config")
}