package lru

import (
	"container/list"
	"fmt"
	"lincoln/smartcache/cachebyte"
	"sync"
	"time"
)

//Cache 实现了Lru算法的 缓存结构
type Cache struct {
	mu       sync.Mutex               //锁
	list     list.List                //按访问时间排序的key的缓存列表
	data     map[string]*list.Element //缓存的 key-value
	currSize int                      //当前内存大小
	maxSize  int                      //最大的可分配的内存大小
}

type innerEle struct {
	key   string
	value interface{}
}

func (cache *Cache) New() {
	cache.data = make(map[string]*list.Element)
	cache.maxSize = 4 << 30 //默认缓存4g大小
}

//Get 获取缓存对象
func (cache *Cache) Get(key string) (*cachebyte.CacheByte, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	//获取key, 将key 加到链表头部
	ele, ok := cache.data[key]
	if !ok {
		fmt.Println("not in key")
		return nil, false
	}
	cache.list.MoveToFront(ele)

	//获取value
	ie, ok := ele.Value.(innerEle)
	if !ok {
		fmt.Println("not innerEle")
		return nil, false
	}

	cb, ok := ie.value.(cachebyte.CacheByte)
	return &cb, ok
}

//Set 添加缓存对象(expires 暂不做处理)
func (cache *Cache) Set(key string, val cachebyte.CacheByte, expires time.Duration) bool {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	//查看key是否已存在
	ele, ok := cache.data[key]
	if ok {
		//已存在， 移到最前面
		cache.list.MoveToFront(ele)
		return true
	}

	thisSize := len(val.Raws)

	//先检查是否已达到最大内存限制了
	if thisSize+cache.currSize >= cache.maxSize {
		cache.RemoveOldest()
		cache.currSize = cache.currSize - thisSize
	}

	//将key加到链表头部
	cache.currSize = cache.currSize + thisSize

	ele = cache.list.PushFront(innerEle{
		key:   key,
		value: val,
	})

	cache.data[key] = ele
	fmt.Println("set ok")
	return true
}

//RemoveOldest 移除最近最少使用的对象
func (cache *Cache) RemoveOldest() {
	ele := cache.list.Back()
	cache.list.Remove(ele)
	delete(cache.data, (ele.Value.(innerEle)).key)

	fmt.Printf("移除最近最少使用:%s", (ele.Value.(innerEle)).key)
}
