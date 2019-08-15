package gocache

import (
	"fmt"
	"lincoln/smartcache/lru"
	"net/http"
	"strings"
)

type GoCache struct {
	baseName   string    //Http请求时的链接信息
	cacheNodes nodeHttp  //管理跟其他GoCache的通信
	innerCache lru.Cache //缓存对象（用Lru管理）
}

//New 初始化GoCache
func New(nodeaddrs []string, localAddr string, baseNameInURL string) GoCache {
	smartCache := GoCache{
		baseName:   baseNameInURL,
		cacheNodes: nodeHttp{nodeAddrs: nodeaddrs, selfAddr: localAddr},
	}

	//将各个cache 节点 hash分布
	smartCache.cacheNodes.Hash()

	return smartCache
}

//ServeHTTP 实现http的接口
func (cache GoCache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, cache.baseName) {
		http.NotFound(w, r)
		return
	}

	splits := strings.SplitN(r.URL.Path, "/", 2)
	if len(splits) != 2 {
		http.Error(w, "error url", http.StatusForbidden)
		return
	}

	//获取请求的key
	key := splits[1]

	//获取key对应的value
	value, ok := cache.GetValue(key)
	if !ok {
		http.Error(w, "error url", http.StatusNotFound)
		return
	}

	//返回
	response := fmt.Sprintf("{\"key\":%s,\"value\":%s}", key, value)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(response))
}

//GetValue 根据key获取value(可能在本地或从其他节点获取)
func (cache GoCache) GetValue(key string) (string, bool) {
	//获取节点通信地址
	addr, isLocal := cache.cacheNodes.GetAddr(key)

	//Key在当前节点
	if isLocal {

		return addr, true
	}

	//Key在其他节点，通过http获取
	return addr, true
}
