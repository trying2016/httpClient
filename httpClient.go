// Copyright 2014-2019 Tyring <wxytom@163.com>.
// Licensed under the MIT license.
// Powerful and easy to use http client

package httpClient

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
	postData      map[string]interface{} //
	postContents  []byte
	headers       map[string]string
	timeOut       time.Duration
	postDataType  int
	proxy         string
	useGZip       bool
	receiveCookie string
}

func NewHttpClient() *HttpClient {
	hClient := HttpClient{}
	hClient.timeOut = time.Second * 30
	hClient.postDataType = POST_DATA_TYPE_FORM
	return &hClient
}

// Set contents type
func (hClient *HttpClient) SetPostDataType(dataType int) {
	hClient.postDataType = dataType
}

func (hClient *HttpClient) SetPostData(postData interface{}) {
	hClient.postContents, _ = json.Marshal(postData)
}

// add
func (hClient *HttpClient) AddFormData(key string, value interface{}) {
	if hClient.postData == nil {
		hClient.postData = make(map[string]interface{})
	}
	hClient.postData[key] = value
}

// Set the proxy host:port or http://host:port
// example 127.0.0.1:8888
func (hClient *HttpClient) SetProxy(proxy string) {
	if strings.Contains(proxy, "://") {
		hClient.proxy = proxy
	} else {
		hClient.proxy = "http://" + proxy
	}
}

func (hClient *HttpClient) SetCookie(cookie string) {
	hClient.AddHeader("Cookie", cookie)
}

func (hClient *HttpClient) GetCookie() string {
	return hClient.receiveCookie
}

func (hClient *HttpClient) AddHeader(key, value string) {
	if hClient.headers == nil {
		hClient.headers = make(map[string]string)
	}
	hClient.headers[key] = value
}

//
func (hClient *HttpClient) EncodingGZip(bUse bool) {
	hClient.useGZip = bUse
}

// Post
func (hClient *HttpClient) Post(link string) (string, error) {
	if hClient.postContents == nil || len(hClient.postContents) == 0 {
		hClient.postContents = hClient.getPostData()
	}
	return hClient.do("POST", link, hClient.postContents)
}

func (hClient *HttpClient) Get(link string) (string, error) {
	httpClient.postDataType = POST_DATA_TYPE_FORM
	formData := httpClient.getPostData()
	if len(formData) > 0 {
		formStr := string(formData)
		if strings.Contains(link, "?") {
			link += formStr
		}else{
			link += "?" + formStr
		}
	}
	return hClient.do("GET", link, nil)
}

func (hClient *HttpClient) getPostData() []byte {
	if hClient.postData == nil || len(hClient.postData) == 0 {
		return []byte("")
	}

	if hClient.postDataType == POST_DATA_TYPE_JSON {
		data, _ := json.Marshal(hClient.postData)

		// clean postdata
		hClient.postData = nil
		return data
	} else {
		var data string
		for key, value := range hClient.postData {
			separate := "&"
			if len(data) == 0 {
				separate = ""
			}
			data += fmt.Sprintf("%s%s=%v", separate, key, value)
		}
		// clean postdata
		hClient.postData = nil
		return []byte(data)
	}
}

func (hClient *HttpClient) SetReferer(refUrl string) {
	hClient.AddHeader("Referer", refUrl)
}

func (hClient *HttpClient) setHeaders(request *http.Request) {
	for k, v := range hClient.headers {
		request.Header.Set(k, v)
	}
}

func (hClient *HttpClient) do(method string, link string, data []byte) (string, error) {
	var request *http.Request
	var err error
	if data != nil && len(data) != 0 {
		// 大于1024字节使用gzip压缩
		if hClient.useGZip {
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
	hClient.postContents = nil

	if err != nil {
		return "", err
	} else {
		var transport *http.Transport = nil
		if hClient.proxy != "" {
			URL := url.URL{}
			urlProxy, _ := URL.Parse(hClient.proxy)
			transport = &http.Transport{
				Proxy: http.ProxyURL(urlProxy),
			}
		} else {
			transport = &http.Transport{}
		}

		netClient := &http.Client{
			Timeout:   hClient.timeOut,
			Transport: transport,
		}

		// set header
		hClient.setHeaders(request)

		if response, err := netClient.Do(request); err != nil {
			return "", err
		} else {
			defer response.Body.Close()
			// save recevie cookie
			for _, v := range response.Cookies() {
				separate := "; "
				if hClient.receiveCookie == "" {
					separate = ""
				}
				hClient.receiveCookie += fmt.Sprintf("%s%s=%s", separate, v.Name, v.Value)
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
