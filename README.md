# smartcache
当前的缓存库简单的使用了 多节点(http协议) + lru 算法

请求key：http://192.168.1.102:8002/goCache/Get/hello1

goCache: 组名
Get:要执行的动作
hello1: key名

响应
{
    "Key": "hello1",
    "Method": "Get",
    "Value": {
        "Raws": "aGk="  //Raws是 []byte类型
    }
}


设置key: http://192.168.1.102:8002/goChache/Set/hello1
goCache: 组名
Set:要执行的动作
hello1: key名

响应
{
	"Key": "hello1",
	"Method": "Set",
	"RetCode": 0,  //操作结果 1=成功  0=失败
	"Msg": "ok"
}

 

