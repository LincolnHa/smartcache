package main

import (
	"fmt"
	"lincoln/smartcache/gocache"
	"net/http"
)

func main() {

	//各个cache 节点1
	nodeaddrs := []string{
		"http://192.168.1.102:8001",
		"http://192.168.1.102:8002",
	}

	//将cache 节点 通过一致性hash，分散缓存key_value
	thisNode := "http://192.168.1.102:8002"
	smartCache := gocache.New(nodeaddrs, thisNode, "goChache")

	//启动节点，监听
	err := http.ListenAndServe(":8002", smartCache)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
}
