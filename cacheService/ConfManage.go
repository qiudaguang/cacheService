package cacheService

import (
	"errors"
	"../customType"
	"../config"
)

func AddGroup(apiType customType.ApiConfType) error{

	existIndex := -1
	for index, groupType := range groupList {
		if groupType.UrlPath == apiType.UrlPath {
			existIndex = index
			groupList[index] = apiType
			break
		}
	}
	if existIndex != -1 {
		return errors.New("接口分组已存在")
	}
	groupList = append(groupList, apiType)
	err := config.UpdateApiConfig(groupList); if err != nil {
		return err
	}
	return nil
}

func EditGroup(apiType customType.ApiConfType) error{

	existIndex := -1
	for index, groupType := range groupList {
		if groupType.UrlPath == apiType.UrlPath {
			existIndex = index
			groupList[index] = apiType
			break
		}
	}
	if existIndex == -1 {
		return errors.New("接口分组不存在")
	}
	err := config.UpdateApiConfig(groupList); if err != nil {
		return err
	}
	return nil
}

func DelGroup(groupName string) error{

	existIndex := -1
	for index, groupType := range groupList {
		if groupType.UrlPath == groupName {
			existIndex = index
			break
		}
	}
	if existIndex == -1 {
		return errors.New("接口分组不存在")
	}
	groupList = append(groupList[:existIndex], groupList[existIndex+1:]...)
	err := config.UpdateApiConfig(groupList); if err != nil {
		return err
	}

	redisC := pool.Get()
	defer redisC.Close()

	itemList := GetGroupItem(groupName)
	for _, item := range itemList {
		redisC.Do("DEL", API_ITEM_CACHE + ":" + item)
	}

	redisC.Do("DEL", API_VISITS_COUNT + ":" + groupName)
	return nil
}
