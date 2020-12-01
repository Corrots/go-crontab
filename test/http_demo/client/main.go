package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"log"
)

func main() {
	// 创建连接池
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // 连接超时时间
			KeepAlive: 30 * time.Second, // 长连接超时时间
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second, // TLS握手超时时间
		MaxIdleConns:          100,              // 最大空闲连接
		IdleConnTimeout:       90 * time.Second, // 空闲超时时间
		ExpectContinueTimeout: 1 * time.Second,  // "Expect: 100-continue" header的超时时间
	}
	url := "http://127.0.0.1:8080/abc"
	// 创建客户端
	client := http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // 请求超时时间
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalf("[%s] get err: %v\n", url, err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Printf("[%s] response: %s\n", url, b)
}
