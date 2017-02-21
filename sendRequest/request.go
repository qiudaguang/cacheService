package sendRequest

import (
	"net/http"
	"strings"
	"io/ioutil"
	"errors"
)

func SendHttpGet(url string) (string, error) {

	resp, respErr := http.Get(url)
	if respErr != nil {
		return "", respErr
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return string(body), errors.New(resp.Status)
	}
	return string(body), nil
}

func SendHttpPost(url string, body string) (string , error){

	resp, err := http.Post(url, "application/json", strings.NewReader(body))
	if err != nil {
		return "", err
	}

	body2, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return string(body2), errors.New(resp.Status)
	}
	return string(body2), nil
}
