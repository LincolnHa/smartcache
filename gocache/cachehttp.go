package gocache

import (
	"lincoln/smartcache/consitenthash"
)

type nodeHttp struct {
	nodeAddrs []string           //其他gocache节点地址
	selfAddr  string             //本地地址
	consiHash consitenthash.Hash //一致性Hash对象
}

//NodeHash 将各个节点hash
func (nHttp nodeHttp) Hash() {
	nHttp.consiHash.StartHash(nHttp.nodeAddrs)
}

func (nHttp nodeHttp) GetAddr(key string) (string, bool) {
	//根据key和hash算法 获取该Key所在的节点 链接
	currAddr := nHttp.consiHash.GetNode(key)

	//该地址是本地地址
	if currAddr == nHttp.selfAddr {
		return currAddr, true
	}

	return currAddr, false
}
