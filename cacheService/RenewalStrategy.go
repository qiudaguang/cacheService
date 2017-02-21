package cacheService

import (
	"fmt"
	"time"
	"../customType"
	"github.com/garyburd/redigo/redis"
)

func BootRegularlyUpdate() {

	updatecycle, _ := conf.Int64("cache_service::updatecycle")
	if updatecycle <= 0 {
		fmt.Println("BootRegularlyUpdate err, updatecycle 无效")
		return
	}

	ticket := time.NewTicker(time.Duration(updatecycle) * time.Second)
	go func() {
		for t := range ticket.C {
			fmt.Println("BootRegularlyUpdate", t)
			UpdateAllCache()
		}
	}()
}


func CheckUpdateByVisistCount(itemconfig customType.ApiConfType) bool {

	if itemconfig.CheckCount <= 0 {
		return false
	}

	redisC := pool.Get()
	defer redisC.Close()
	visistCount, err := redis.Int(redisC.Do("GET", API_VISITS_COUNT + ":" + itemconfig.UrlPath))
	if err != nil || visistCount == 0 {
		return false
	}

	if visistCount % itemconfig.CheckCount == 0 {
		return true
	}

	return false
}