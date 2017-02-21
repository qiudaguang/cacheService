package main

import (
	"fmt"
	"net/http"
	"./proxy"
	"./config"
	"os"
	"./cacheService"
	"./manage"
)

func startProxy(w http.ResponseWriter, r *http.Request) {

	if proxy.IsStaticRes(r.URL.Path) {
		proxy.StaticService(w, r)
	} else {
		apiHost := conf.String("api::host")
		apiPort := conf.String("api::port")
		if apiPort != "80" {
			apiHost += apiPort
		}
		proxy.ReverseProxy(w, r, "http://" + apiHost)
	}
}

func StartWebService() {

	cacheService.InitCacheConfig()
	cacheService.BootRegularlyUpdate()

	http.HandleFunc("/", startProxy)
	http.HandleFunc("/manage/", manage.LocateRequest)
	err := http.ListenAndServe(":"+conf.String("server::port"), nil)
	if err != nil {
		fmt.Println("err:", err)
		os.Exit(2)
	}
}

var conf  = config.GetConfig()

func main() {

	StartWebService()
}
