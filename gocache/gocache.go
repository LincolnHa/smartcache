package gocache

import (
	"lincoln/smartcache/lru"
	"net/http"
)

type GoCache struct {
	commNodes  nodeHttp  //管理跟其他GoCache的通信
	innerCache lru.Cache //缓存对象（用Lru管理）
}

//NewCache 初始化GoCache
func NewCache(nodeaddrs []string) GoCache {
	smartCache := GoCache{
		commNodes: nodeHttp{nodeAddrs: nodeaddrs},
	}

	//将各个cache 节点 hash分布
	smartCache.commNodes.NodeHash()

	//缓存对象开始用lru算法监控剔除最近最少使用的缓存对象
	smartCache.innerCache.Start()

	return smartCache
}

//ServeHTTP 实现http的接口
func (cache GoCache) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

}
