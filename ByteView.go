package Gache

//支持并发读写
//ByteView 表示缓存值 只读
//选择byte类型是为了能够支持任意数据类型的存储，
//例如字符串，图片
//b是只读的，使用ByteSlice方法返回一个拷贝
//防止缓存值被外部程序修改
type ByteView struct{
	b []byte
}

//实现了 Value 接口
func (v ByteView) Len() int  {
	return len(v.b)
}

//返回切片的副本
func (v ByteView) ByteSlice() []byte  {
	return cloneBytes(v.b)
}

//toString
func (v ByteView) toString() string  {
	return string(v.b)
}

func cloneBytes(bytes []byte) []byte {
	c:= make([]byte,len(bytes))
	copy(c,bytes)
	return c
}
