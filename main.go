package main

import (
	"lincoln/smartcache/gocache"
	"net/http"
)

func main() {

	//各个cache 节点
	nodeaddrs := []string{
		"http://192.168.1.105:8001",
		"http://192.168.1.105:8002",
		"http://192.168.1.105:8003",
	}

	//将cache 节点 通过一致性hash，分散缓存key_value
	smartCache := gocache.NewCache(nodeaddrs)

	//启动节点，监听
	http.ListenAndServe(":8001", smartCache)
}
