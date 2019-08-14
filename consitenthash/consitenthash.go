package consitenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

var (
	repeatCount = 10
)

//Hash 一致性hash对象
type Hash struct {
	nodes       []int          //hash后的节点位置(包括虚拟节点)
	repeatCount int            //每个真实节点拓展出来到的虚拟节点个数
	addrs       map[int]string //虚拟节点对应的真实地址
}

//将节点 hash
func (hash Hash) StartHash(addrs []string) {

	//将各个addrs虚拟为 hash环的某个节点位置
	for _, addr := range addrs {
		//将地址  hash为unit32,作为节点位置
		for i := 0; i < repeatCount; i++ {
			node := int(crc32.ChecksumIEEE([]byte(strconv.Itoa(i) + addr)))
			hash.nodes = append(hash.nodes, node)
			hash.addrs[node] = addr //设置节点位置对应的真实地址
		}

	}

	sort.Ints(hash.nodes)
}

//根据Key, 计算Hash， 获取真实的addr
func (hash Hash) GetNode(key string) string {
	node := int(crc32.ChecksumIEEE([]byte(key)))
	idx := sort.Search(len(hash.nodes), func(i int) bool { return hash.nodes[i] >= node })

	if idx == len(hash.nodes) {
		idx = 0
	}

	addr := hash.addrs[idx]

	return addr
}
