package gocache

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lincoln/smartcache/cachebyte"
	"lincoln/smartcache/consitenthash"
	"net/http"
	"time"
)

type nodeHttp struct {
	nodeAddrs []string           //其他gocache节点地址
	selfAddr  string             //本地地址
	consiHash consitenthash.Hash //一致性Hash对象
}

//NodeHash 将各个节点hash
func (nHttp *nodeHttp) HashAddr() {
	nHttp.consiHash.New()

	nHttp.consiHash.StartHash(nHttp.nodeAddrs)
}

func (nHttp *nodeHttp) GetAddr(key string) (string, bool) {
	//根据key和hash算法 获取该Key所在的节点 链接
	currAddr := nHttp.consiHash.GetNode(key)

	//该地址是本地地址
	if currAddr == nHttp.selfAddr {
		return currAddr, true
	}

	return currAddr, false
}

func (nHttp *nodeHttp) Get(sendUrl string) (*cachebyte.CacheByte, bool) {
	fmt.Printf("Get: sendUrl:%s\r\n", sendUrl)

	//http请求
	req, _ := http.NewRequest("GET", sendUrl, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("(1)nodeHttp Get Err:%s\r\n", err.Error())
		return nil, false
	}
	defer res.Body.Close()

	//解析
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Printf("body:%s\r\n", string(body))
	result := BridgeData_Get{}

	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		fmt.Printf("(2)nodeHttp Get Err:%s\r\n", err.Error())
		return nil, false
	}

	return &result.Value, true
}

func (nHttp *nodeHttp) Set(sendUrl string, key string, sendvalue cachebyte.CacheByte, expires time.Duration) bool {
	fmt.Printf("Set: sendUrl:%s\r\n", sendUrl)

	//post 数据
	reqbody := BridgeData_Set{
		Value:  sendvalue,
		Expire: time.Hour * 24 * 365 * 100,
	}

	//http请求
	data, err := json.Marshal(reqbody)
	if err != nil {
		fmt.Printf("(1)nodeHttp Set Err:%s\r\n", err.Error())
		return false
	}

	reqBodyBuffer := bytes.NewBuffer(data)
	res, err := http.Post(sendUrl, "application/json;charset=utf-8", reqBodyBuffer)
	if err != nil {
		fmt.Printf("(2)nodeHttp Set Err:%v\r\n", err)
		return false
	}

	defer res.Body.Close()

	//解析
	resbody, _ := ioutil.ReadAll(res.Body)
	result := struct {
		Key     string
		Method  string
		RetCode int
		Msg     string
	}{}

	err = json.Unmarshal([]byte(resbody), &result)
	if err != nil {
		fmt.Printf("(3)nodeHttp Set Err:%s\r\n", err.Error())
		return false
	}

	if result.RetCode == 1 {
		return true
	}

	return false
}
