package peterstd

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var (
	httpClient = &http.Client{
		Timeout: 4 * time.Second,
	}
)

func ContentHttpRequest(method string, url string,
	requestHead map[string]string, query url.Values) ([]byte, error) {
	//
	data, err := CommonHttpRequest(httpClient, method, url, requestHead, query)
	if err != nil {
		data, err = CommonHttpRequest(httpClient, method, url, requestHead, query)
		if err != nil {
			Errorln("FlashHttpRequest.CommonHttpRequest Do Error", url, err)
			return data, err
		}
	}
	return data, err
}

func CommonHttpRequest(httpClient *http.Client, method string, url string,
	requestHead map[string]string, query url.Values) ([]byte, error) {
	request, err := http.NewRequest(method, url, bytes.NewBufferString(query.Encode()))
	if err != nil {
		return nil, err
	}
	if requestHead != nil {
		for _key, _value := range requestHead {
			request.Header.Set(_key, _value)
		}
	}
	response, err := httpClient.Do(request)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(response.Body)
	return data, err
}
