package proxy

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"bytes"
	"../cacheService"
	//"strings"
	"strings"
	"os"
)

func ReverseProxy(w http.ResponseWriter, r *http.Request, host string) {

	url := host + r.URL.String()

	paramsStr := ""
	if r.Method == "POST" {
		r_body, _ := ioutil.ReadAll(r.Body)
		paramsStr = bytes.NewBuffer(r_body).String()
	} else if r.Method == "GET" {
		paramsStr = ""	//需截取url参数
		arr := strings.Split(r.URL.String(), "?")
		if len(arr) > 1 {
			paramsStr = arr[1]
		}
	}

	respStr, err := cacheService.LookupCache(r.URL.Path, url, paramsStr, r.Method, r.Header); if err != nil {
		fmt.Println("lookup err:", err)
		http.Error(w, err.Error(), 500)
	}

	fmt.Fprint(w, string(respStr))
	return
}

func StaticService(w http.ResponseWriter, r *http.Request) {

	staticExt := GetStaticExt()
	ext := ParseUrlExt(r.URL.Path)
	w.Header().Set("content-type",  staticExt[ext])

	fin, err := os.Open("D:\\coding\\go\\test\\cache" + r.URL.Path)

	if err != nil && os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(""))
	}
	defer fin.Close()
	fd, _ := ioutil.ReadAll(fin)

	w.Write(fd)
}

func GetStaticExt() map[string] string{

	mimeTypes := make(map[string] string)
	mimeTypes["css"] = "text/css"
	mimeTypes["js"] = "application/javascript"
	mimeTypes["jpeg"] = "image/jpeg"
	mimeTypes["jpg"] = "image/jpeg"
	mimeTypes["jpe"] = "image/jpeg"
	mimeTypes["png"] = "image/png"
	mimeTypes["gif"] = "image/gif"
	mimeTypes["ico"] = "image/x-icon"

	return mimeTypes
}

func ParseUrlExt(urlPath string) string{

	lastIndex := len(urlPath)
	index := strings.LastIndex(urlPath, "?")
	if index > 0 {
		lastIndex = index
	}
	index = strings.LastIndex(urlPath, "#")
	if index > 0 && index < lastIndex {
		lastIndex = index
	}
	startIndex := strings.LastIndex(urlPath, ".")
	if startIndex == -1 {
		return ""
	}
	request_type := urlPath[startIndex+1:lastIndex]
	return request_type
}

func IsStaticRes(urlPath string) bool{

	urlExt := ParseUrlExt(urlPath)
	if urlExt == "" {
		return false
	}
	staticExt := GetStaticExt()
	contentType := staticExt[urlExt]
	if contentType == "" {
		return false
	}
	return true
}