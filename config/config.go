package config

import (
	"github.com/lxmgo/config"
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"../customType"
	"sync"
)

var once sync.Once

var staticConf config.ConfigInterface

func GetConfig() config.ConfigInterface{

	once.Do(func() {
		var err error
		staticConf, err = config.NewConfig(os.Args[1] + "/config.ini")
		if err != nil {
			fmt.Println("err:", err)
		}
	})
	return staticConf
}

func GetApiConfig() ([]customType.ApiConfType, error) {

	var apiList []customType.ApiConfType

	configBytes, err := ioutil.ReadFile(os.Args[1] + "/apiConfig.json")
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		return apiList, err
	}

	if err := json.Unmarshal(configBytes, &apiList); err != nil {
		fmt.Println("Unmarshal: ", err.Error())
		return apiList, err
	}
	return apiList, nil
}

func UpdateApiConfig(apiList []customType.ApiConfType) error {

	apiJsonByte, err := json.Marshal(apiList); if err != nil {
		return err
	}
	err = ioutil.WriteFile(os.Args[1] + "/apiConfig.json", apiJsonByte, 0777)
	return err
}