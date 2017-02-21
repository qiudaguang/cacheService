package manage

import (
	"net/http"
	"../cacheService"
	"../customType"
	"fmt"
	"../tools"
	"../ErrNum"
	"../redisConn"
	"strconv"
	"github.com/garyburd/redigo/redis"
)

type ViewType struct {

}

var pool = redisConn.NewPool()

func (v *ViewType) NotFound(w http.ResponseWriter, r *http.Request) {

	var result customType.ReturnType
	result.Errcode = 404
	result.Errmsg = "未找到页面"

	ReturnResult(w, result)
}

func (v *ViewType) Index(w http.ResponseWriter, r *http.Request) customType.ReturnType{

	var result customType.ReturnType
	result.Data = "hello"
	return result
}

func (v *ViewType) GroupList(w http.ResponseWriter, r *http.Request) {

	type returnConfType struct {
		customType.ApiConfType
		VisitsCount int
	}

	var result customType.ReturnType
	var typeList []returnConfType
	groupList := cacheService.GetGroupNameList()
	for _, groupname := range groupList {
		var itemType returnConfType
		itemConfig, err := cacheService.GetGroupConfig(groupname); if err != nil {
			continue
		}
		itemType.ApiConfType = itemConfig
		itemType.VisitsCount = cacheService.GetVisitsCount(groupname)
		typeList = append(typeList, itemType)
	}

	result.Data = typeList
	ReturnResult(w, result)
}

func (v *ViewType) UpdateCacheForApiGroup(w http.ResponseWriter, r *http.Request) {

	result := customType.ReturnType{}
	if r.Method != "POST" {
		result.Errcode = ErrNum.ERR_NOT_SUP_METHOD
		ReturnResult(w, result)
		return
	}
	r.ParseForm()
	groupname :=  r.PostFormValue("groupname")

	err := cacheService.UpdateCacheForApiGroup(groupname)
	if err != nil {
		result.Errcode = ErrNum.ERR_UP_C_G_ERR
		result.Errmsg = err.Error()
	}

	ReturnResult(w, result)
}

func (v *ViewType) UpdateCacheForApiKey(w http.ResponseWriter, r *http.Request) {

	result := customType.ReturnType{}
	if r.Method != "POST" {
		result.Errcode = ErrNum.ERR_NOT_SUP_METHOD
		ReturnResult(w, result)
		return
	}
	r.ParseForm()
	groupName :=  r.PostFormValue("groupName")
	cacheKey :=  r.PostFormValue("cacheKey")
	if groupName == "" || cacheKey == "" {
		result.Errcode = ErrNum.ERR_PARAMS_ERR
		ReturnResult(w, result)
		return
	}

	itemconfig, _ := cacheService.GetGroupConfig(groupName)
	respType, _ := cacheService.GetApiItemCache(cacheKey)
	cacheService.UpdateCacheForApi(groupName, respType.ApiUrl, respType.ParamsStr, respType.Method, cacheKey, itemconfig)

	ReturnResult(w, result)
}

//某个分组下的缓存列表
func (v *ViewType) GetCacheItemList(w http.ResponseWriter, r *http.Request) {

	result := customType.ReturnType{}
	if r.Method != "POST" {
		result.Errcode = ErrNum.ERR_NOT_SUP_METHOD
		ReturnResult(w, result)
		return
	}

	r.ParseForm()
	groupname :=  r.PostFormValue("groupname")

	type rt2 struct {
		customType.RespCacheType
		Ttl	int
		CacheKey string
	}

	redisC := pool.Get()
	defer redisC.Close()

	var itemTypeList []rt2
	itemList := cacheService.GetGroupItem(groupname)
	for _, itemKey := range itemList {
		respType, _ := cacheService.GetApiItemCache(itemKey)
		if (respType == customType.RespCacheType{}) {
			//为空已过期，将key移除list
			cacheService.RemoveFromGroupItem(groupname, itemKey)
			continue
		}

		var itemRt2 rt2
		itemRt2.CacheKey = itemKey
		itemRt2.RespCacheType = respType
		itemRt2.Ttl, _ = redis.Int(redisC.Do("TTL", itemKey))

		itemTypeList = append(itemTypeList, itemRt2)
	}

	result.Data = itemTypeList

	ReturnResult(w, result)
}

//添加分组
func (v *ViewType) AddApiGroup(w http.ResponseWriter, r *http.Request) {

	result := customType.ReturnType{}

	if r.Method != "POST" {
		result.Errcode = ErrNum.ERR_NOT_SUP_METHOD
		ReturnResult(w, result)
		return
	}

	r.ParseForm()
	var apiType customType.ApiConfType
	apiType.UrlPath = r.PostFormValue("UrlPath")
	apiType.ExpireTime, _ = strconv.Atoi(r.PostFormValue("ExpireTime"))
	apiType.CheckCount, _ = strconv.Atoi(r.PostFormValue("CheckCount"))
	apiType.ParamsArr = r.Form["ParamsArr[]"]
	if len(apiType.ParamsArr) == 0 {
		apiType.ParamsArr = []string{}
	}

	groupList := cacheService.GetGroupNameList()
	exist := tools.InArrayStr(apiType.UrlPath, groupList); if exist == true {
		result.Errcode = ErrNum.ERR_G_EXIST
		ReturnResult(w, result)
		return
	}

	err := cacheService.AddGroup(apiType); if err != nil {
		result.Errmsg = err.Error()
	}

	ReturnResult(w, result)
}

//获得分组配置
func (v *ViewType) GetApiGroupConf(w http.ResponseWriter, r *http.Request) {

	result := customType.ReturnType{}

	if r.Method != "POST" {
		result.Errcode = ErrNum.ERR_NOT_SUP_METHOD
		ReturnResult(w, result)
		return
	}

	r.ParseForm()
	groupname := r.PostFormValue("groupName")
	existIndex := -1
	grouplist := cacheService.GetGroupConfList()
	for index, groupType := range grouplist {
		if groupType.UrlPath == groupname {
			existIndex = index
			break
		}
	}
	if existIndex == -1 {
		result.Errmsg = "接口分组不存在"
		ReturnResult(w, result)
		return
	}
	fmt.Println(existIndex)
	result.Data = grouplist[existIndex]

	ReturnResult(w, result)
}

//修改分组
func (v *ViewType) EditApiGroup(w http.ResponseWriter, r *http.Request) {

	result := customType.ReturnType{}

	if r.Method != "POST" {
		result.Errcode = ErrNum.ERR_NOT_SUP_METHOD
		ReturnResult(w, result)
		return
	}

	r.ParseForm()
	var apiType customType.ApiConfType
	apiType.UrlPath = r.PostFormValue("UrlPath")
	apiType.ExpireTime, _ = strconv.Atoi(r.PostFormValue("ExpireTime"))
	apiType.CheckCount, _ = strconv.Atoi(r.PostFormValue("CheckCount"))
	apiType.ParamsArr = r.Form["ParamsArr[]"]
	if len(apiType.ParamsArr) == 0 {
		apiType.ParamsArr = []string{}
	}

	err := cacheService.EditGroup(apiType); if err != nil {
		result.Errmsg = err.Error()
	}

	fmt.Println("this is here")
	ReturnResult(w, result)
}

//删除分组
func (v *ViewType) DelApiGroup(w http.ResponseWriter, r *http.Request) {

	result := customType.ReturnType{}

	if r.Method != "POST" {
		result.Errcode = ErrNum.ERR_NOT_SUP_METHOD
		ReturnResult(w, result)
		return
	}

	r.ParseForm()
	groupName := r.PostFormValue("groupname")
	err := cacheService.DelGroup(groupName); if err != nil {
		result.Errcode = ErrNum.ERR_DEL_C_G_ERR
		result.Errmsg = err.Error()
	}

	ReturnResult(w, result)
}

func ReturnResult(w http.ResponseWriter, result customType.ReturnType) {

	if result.Errcode != 0 && result.Errmsg == "" {
		result.Errmsg = ErrNum.GetErrMsg(result.Errcode)
	}
	fmt.Fprint(w, result.ToJsonStr())
	return
}