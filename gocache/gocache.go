package gocache

import (
	"encoding/json"
	"fmt"
	"lincoln/smartcache/cachebyte"
	"lincoln/smartcache/lru"
	"net/http"
	"strings"
	"time"
)

//GoCache 缓存对象
type GoCache struct {
	baseName   string    //Http请求时的链接信息
	cacheNodes nodeHttp  //管理跟其他GoCache的通信
	innerCache lru.Cache //缓存对象（用Lru管理）
}

type BridgeData_Get struct {
	Key    string
	Method string
	Value  cachebyte.CacheByte
}

type BridgeData_Set struct {
	Value  cachebyte.CacheByte
	Expire time.Duration
}

//New 初始化GoCache(把节点Addr hash分布)
func New(nodeaddrs []string, localAddr string, baseNameInURL string) *GoCache {
	smartCache := GoCache{
		baseName:   baseNameInURL,
		cacheNodes: nodeHttp{nodeAddrs: nodeaddrs, selfAddr: localAddr},
	}

	//初始化innerCache
	smartCache.innerCache.New()

	//将各个cache 节点 hash分布
	smartCache.cacheNodes.HashAddr()

	return &smartCache
}

//ServeHTTP 处理其他节点请求过来的数据
func (cache *GoCache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//获取请求参数
	splits := strings.SplitN(r.URL.Path, "/", 4)
	if len(splits) != 4 {
		http.NotFound(w, r)
		return
	}

	baseName := splits[1]
	method := splits[2]
	key := splits[3]

	//请求错误
	if baseName != cache.baseName {
		http.NotFound(w, r)
		return
	}

	//要获取key对应的value
	if method == "Get" {
		value, ok := cache.GetValue(key)
		if !ok {
			http.NotFound(w, r)
			return
		}

		data := BridgeData_Get{
			Key:    key,
			Method: method,
			Value:  *value,
		}

		resData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, "(2)error url", http.StatusNotFound)
			return
		}

		//响应结构为 BridgeData_Get的数据
		w.Header().Set("Content-Type", "application/json")
		w.Write(resData)
		return
	}

	//存储key
	if method == "Set" {
		//post 数据
		reqbody := BridgeData_Set{}
		err := json.NewDecoder(r.Body).Decode(&reqbody)
		if err != nil {
			http.Error(w, "error url", http.StatusNotFound)
			return
		}

		retCode := 0
		if ok := cache.SetValue(key, reqbody); ok {
			retCode = 1
		}

		//返回结果
		response := fmt.Sprintf("{\"Key\":\"%s\",\"Method\":\"%s\",\"RetCode\":%d,\"Msg\":\"%s\"}", key, method, retCode, "ok")

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
		return
	}

}

//GetValue 根据key获取value(可能在本地或从其他节点获取)
func (cache *GoCache) GetValue(key string) (*cachebyte.CacheByte, bool) {
	var obj *cachebyte.CacheByte
	var ok bool

	//获取节点通信地址(ip:port)
	addr, isLocal := cache.cacheNodes.GetAddr(key)

	if isLocal {
		//Key在当前节点
		obj, ok = cache.innerCache.Get(key)
	} else {
		//Key在其他节点，通过http获取
		url := addr + "/" + cache.baseName + "/Get" + "/" + key
		obj, ok = cache.cacheNodes.Get(url)
	}

	//有值
	if ok && obj != nil {
		return obj, true
	}

	return nil, false
}

//SetValue 缓存Value
func (cache *GoCache) SetValue(key string, data BridgeData_Set) bool {
	var ok bool

	//获取节点通信地址
	addr, isLocal := cache.cacheNodes.GetAddr(key)

	if isLocal {
		//Key在当前节点
		ok = cache.innerCache.Set(key, data.Value, data.Expire)
	} else {
		//Key在其他节点，直接通过http获取
		url := addr + "/" + cache.baseName + "/Set" + "/" + key
		ok = cache.cacheNodes.Set(url, key, data.Value, data.Expire)
	}

	return ok
}
