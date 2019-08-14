package gocache

import (
	"lincoln/smartcache/lru"
	"net/http"
)

type GoCache struct {
	cacheNodes nodeHttp  //管理跟其他GoCache的通信
	innerCache lru.Cache //缓存对象（用Lru管理）
}

//New 初始化GoCache
func New(nodeaddrs []string) GoCache {
	smartCache := GoCache{
		cacheNodes: nodeHttp{nodeAddrs: nodeaddrs},
	}

	//将各个cache 节点 hash分布
	smartCache.cacheNodes.Hash()

	return smartCache
}

//ServeHTTP 实现http的接口
func (cache GoCache) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

}
