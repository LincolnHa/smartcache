package lru

import (
	"container/list"
	"sync"
)

//Cache 实现了Lru算法的 缓存结构
type Cache struct {
	mu      sync.Mutex             //锁
	list    list.List              //按访问时间排序的key的缓存列表
	data    map[string]interface{} //缓存的 key-value
	maxSize int                    //最大的可分配的内存大小
}

//Add 添加缓存对象
func (cache *Cache) Add() {

}

//RemoveOldest 移除最近最少使用的对象
func (cache *Cache) RemoveOldest() {

}
