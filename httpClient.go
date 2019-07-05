// Copyright 2014-2019 Tyring <wxytom@163.com>.
// Licensed under the MIT license.
// Powerful and easy to use http client

package httpclient

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	POST_DATA_TYPE_JSON = 1
	POST_DATA_TYPE_FORM = 2
)

type HttpClient struct {
	postData 	  map[string]interface{} // save post data
	postContents  []byte				 // Custom send data
	headers       map[string]string		 // http headers
	timeOut       time.Duration
	postDataType  int
	proxy         string
	useGZip       bool
	receiveCookie string
}

func NewHttpClient() *HttpClient {
	httpClient := HttpClient{}
	httpClient.timeOut = time.Second * 30
	httpClient.postDataType = POST_DATA_TYPE_FORM
	return &httpClient
}

// Set contents type
func (httpClient *HttpClient) SetPostDataType(dataType int) {
	httpClient.postDataType = dataType
}

func (httpClient *HttpClient) SetPostData(postData interface{}) {
	httpClient.postContents, _ = json.Marshal(postData)
}

// add form or json 
func (httpClient *HttpClient) AddPostData(key string, value interface{}) {
	if httpClient.postData == nil {
		httpClient.postData = make(map[string]interface{})
	}
	httpClient.postData[key] = value
}

// Set the proxy host:port or http://host:port
// example 127.0.0.1:8888
func (httpClient *HttpClient) SetProxy(proxy string) {
	if strings.Contains(proxy, "://") {
		httpClient.proxy = proxy
	} else {
		httpClient.proxy = "http://" + proxy
	}
}

// like uid=9; unick="ben"
func (httpClient *HttpClient) SetCookie(cookie string) {
	httpClient.AddHeader("Cookie", cookie)
}

// return cookies like uid=9; unick="ben"
func (httpClient *HttpClient) GetCookie() string {
	return httpClient.receiveCookie
}

// add http header 
func (httpClient *HttpClient) AddHeader(key, value string) {
	if httpClient.headers == nil {
		httpClient.headers = make(map[string]string)
	}
	httpClient.headers[key] = value
}

// enable or disable gzip
func (httpClient *HttpClient) EncodingGZip(bUse bool) {
	httpClient.useGZip = bUse
}

// set referer url
func (httpClient *HttpClient) SetReferer(refUrl string) {
	httpClient.AddHeader("Referer", refUrl)
}

// The Post request
func (httpClient *HttpClient) Post(link string) (string, error) {
	if httpClient.postContents == nil || len(httpClient.postContents) == 0 {
		httpClient.postContents = httpClient.getPostData()
	}
	return httpClient.do("POST", link, httpClient.postContents)
}

// The GET request
func (httpClient *HttpClient) Get(link string) (string, error) {
	return httpClient.do("GET", link, nil)
}


func (httpClient *HttpClient) getPostData() []byte {
	if httpClient.postData == nil || len(httpClient.postData) == 0 {
		return []byte("")
	}

	if httpClient.postDataType == POST_DATA_TYPE_JSON {
		data, _ := json.Marshal(httpClient.postData)

		// clean postdata
		httpClient.postData = nil
		return data
	} else {
		var data string
		for key, value := range httpClient.postData {
			separate := "&"
			if len(data) == 0 {
				separate = ""
			}
			data += fmt.Sprintf("%s%s=%v", separate, key, value)
		}
		// clean postdata
		httpClient.postData = nil
		return []byte(data)
	}
}

func (httpClient *HttpClient) setHeaders(request *http.Request) {
	for k, v := range httpClient.headers {
		request.Header.Set(k, v)
	}
}

func (httpClient *HttpClient) do(method string, link string, data []byte) (string, error) {

	var request *http.Request
	var err error
	if data != nil && len(data) != 0 {
		// 大于1024字节使用gzip压缩
		if httpClient.useGZip {
			var zBuf bytes.Buffer
			zipWrite := gzip.NewWriter(&zBuf)
			defer zipWrite.Close()
			if _, err = zipWrite.Write(data); err != nil {
				fmt.Println("-----gzip is faild,err:", err)
			}
			request, err = http.NewRequest(method, link, &zBuf)
			request.Header.Add("Content-Encoding", "gzip")
			//request.Header.Add("Accept-Encoding", "gzip")
		} else {
			request, err = http.NewRequest(method, link, bytes.NewReader(data))
		}
	} else {
		request, err = http.NewRequest(method, link, nil)
	}

	// clean postdata
	httpClient.postContents = nil

	if err != nil {
		return "", err
	} else {
		var transport *http.Transport = nil
		if httpClient.proxy != "" {
			URL := url.URL{}
			urlProxy, _ := URL.Parse(httpClient.proxy)
			transport = &http.Transport{
				Proxy: http.ProxyURL(urlProxy),
			}
		} else {
			transport = &http.Transport{}
		}

		netClient := &http.Client{
			Timeout:   httpClient.timeOut,
			Transport: transport,
		}

		// set header
		httpClient.setHeaders(request)

		if response, err := netClient.Do(request); err != nil {
			return "", err
		} else {
			defer response.Body.Close()
			// save recevie cookie
			for _, v := range response.Cookies() {
				separate := "; "
				if httpClient.receiveCookie == "" {
					separate = ""
				}
				httpClient.receiveCookie += fmt.Sprintf("%s%s=%s", separate, v.Name, v.Value)
			}

			if data, err := ioutil.ReadAll(response.Body); err == nil {
				// gzip decompress
				if strings.Contains(response.Header.Get("Accept-Encoding"), "gzip") {
					gzipReader, err := gzip.NewReader(bytes.NewReader(data))
					defer gzipReader.Close()
					if err != nil {
						return string(data), nil
					}
					if unBody, err := ioutil.ReadAll(gzipReader); err != nil {
						return string(data), nil
					} else {
						return string(unBody), nil
					}
				}
				return string(data), err
			} else {
				return "", err
			}
		}
	}
}
