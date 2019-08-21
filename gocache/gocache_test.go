package gocache_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"lincoln/smartcache/cachebyte"
	"lincoln/smartcache/gocache"
	"math/rand"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestSetValue(t *testing.T) {

	//post 数据
	reqbody := gocache.BridgeData_Set{
		Value: cachebyte.CacheByte{
			Raws: []byte("hi"),
		},
		Expire: time.Hour * 24 * 365 * 100,
	}

	//http请求
	data, err := json.Marshal(reqbody)
	if err != nil {
		t.Logf("(1)nodeHttp Set Err:%s", err.Error())
		return
	}

	reqBodyBuffer := bytes.NewBuffer(data)
	sendUrl := "http://192.168.1.102:8002/goChache/Set/hello1"

	res, err := http.Post(sendUrl, "application/json;charset=utf-8", reqBodyBuffer)
	if err != nil {
		t.Logf("Err:%s", err.Error())
		return
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
		t.Logf("(2)nodeHttp Set Err:%s", err.Error())
		return
	}

	if result.RetCode == 1 {
		t.Logf("Success")
		return
	}

	return
}

func TestGetValue(t *testing.T) {
	sendUrl := "http://192.168.1.102:8002/goChache/Get/hello"

	//http请求
	req, _ := http.NewRequest("GET", sendUrl, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("(1)nodeHttp Get Err:%s", err.Error())
		return
	}
	defer res.Body.Close()

	//解析
	body, _ := ioutil.ReadAll(res.Body)
	result := gocache.BridgeData_Get{}

	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		t.Errorf("(2)nodeHttp Get Err:%s", err.Error())
		return
	}

	t.Logf("Get Success, value:%s", string(result.Value.Raws))
}

func BenchmarkPostValue(b *testing.B) {
	for i := 0; i < b.N; i++ {

		//post 数据
		reqbody := gocache.BridgeData_Set{
			Value: cachebyte.CacheByte{
				Raws: []byte("hi"),
			},
			Expire: time.Hour * 24 * 365 * 100,
		}

		//http请求
		data, err := json.Marshal(reqbody)
		if err != nil {
			b.Logf("(1)nodeHttp Set Err:%s", err.Error())
			return
		}

		reqBodyBuffer := bytes.NewBuffer(data)
		sendUrl := "http://192.168.1.102:8002/goChache/Set/hello" + strconv.Itoa(rand.Int())

		res, err := http.Post(sendUrl, "application/json;charset=utf-8", reqBodyBuffer)
		if err != nil {
			b.Logf("Err:%s", err.Error())
			return
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
			b.Logf("(2)nodeHttp Set Err:%s", err.Error())
			return
		}

		if result.RetCode == 1 {
			b.Logf("Success")
			return
		}

		return
	}

}
