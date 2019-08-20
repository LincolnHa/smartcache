package cachebyte

//CacheByte 在缓存里的真实存储对象
type CacheByte struct {
	Raws []byte
}

//ToString 将CacheByte 转换为string
func (bytes CacheByte) ToString() string {
	return string(bytes.Raws)
}
