package Gache

/*
	Q：为什么要实现这两个接口
	A：试想，有N个节点，多线程取值，此时是不是每个远程节点都需要
		分配一个获取方法（Getter),PeerPicker接口是根据key返回
		相应的Getter，而PeerGetter接口则是提供了自定义接口
		本文以http为例，peerGetter就是httpGetter，如果用的不是
		http协议，更甚至是用数据库等，用户就可以自定义Getter了
*/
/*
	PeerPicker根据传入的对应的key
	选择相应节点的PeerGetter
*/
type PeerPicker interface {
	PickPeer(key string)(peer PeerGetter,ok bool)
}

/*
	PeerGetter对于Http客户端，通过它得Get函数，
	从group找到对于的缓存值
	所以写完PeerGetter抽象接口后需要去http.go
	中实现PeerGetter接口
*/
type PeerGetter interface {
	Get(group string,key string)([]byte,error)
}