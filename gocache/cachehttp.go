package gocache

import (
	"lincoln/smartcache/consitenthash"
)

type nodeHttp struct {
	nodeAddrs []string           //其他gocache节点地址
	consiHash consitenthash.Hash //一致性Hash对象
}

//NodeHash 将各个节点hash
func (nHttp nodeHttp) NodeHash() {
	nHttp.consiHash.StartHash(nHttp.nodeAddrs)
}
