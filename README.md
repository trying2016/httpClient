Easy HTTP client for golang

Installation:
go get github.com/trying2016/httpClient

Quick Start:
package main

import (
	"github.com/trying2016/httpClient"
	"fmt"
)

func main(){
	hClient := httpClient.NewHttpClient()
	fmt.Printf(hClient.Get("http://www.jd.com"))
}

Sending Request:

    // get : http://www.jd.com?uid=1&name=kk
    hClient := httpClient.NewHttpClient()
    hClient.AddFormData("uid",1)
    hClient.AddFormData("name","kk")
    fmt.Printf(hClient.Get("http://www.jd.com"))

    // post : body uid=1&name=kk
    hClient := httpClient.NewHttpClient()
    hClient.AddFormData("uid",1)
    hClient.AddFormData("name","kk")
    fmt.Printf(hClient.Post("http://www.jd.com"))

    // post : send the body as json
    hClient := httpClient.NewHttpClient()
    hClient.SetPostDataType(httpClient.POST_DATA_TYPE_JSON)
    hClient.AddFormData("uid",1)
    hClient.AddFormData("name","kk")
    fmt.Printf(hClient.Post("http://www.jd.com"))

    // set send cookie and get receive cookie
    hClient := httpClient.NewHttpClient()
    hClient.SetCookie("uid=9; unick=\"ben\"")
    // post
    fmt.Printf(hClient.Post("http://www.jd.com"))
    // or get
    fmt.Printf(hClient.Get("http://www.jd.com"))
    // get receive cookie
    fmt.Printf(hClient.GetCookie())

    // The custom head
    hClient := httpClient.NewHttpClient()
    hClient.AddHeader("user-agent",	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.1 Safari/605.1.15")
    fmt.Printf(hClient.Get("http://www.jd.com"))

    // encoding gzip 
    hClient := httpClient.NewHttpClient()
    hClient.EncodingGZip(true)
    // set post data
    // hClient.SetPostData(...)
    fmt.Printf(hClient.Post("http://www.jd.com"))

    // set proxy
    hClient := httpClient.NewHttpClient()
    hClient.SetProxy("127.0.0.1:8888")
    fmt.Printf(hClient.Get("http://www.jd.com"))