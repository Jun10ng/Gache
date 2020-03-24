package Gache

import (
	"Gache/lru"
	"sync"
)
/*
并发控制

封装lru，封装get和add方法，并添加互斥锁mu
*/
type cache struct {
	mu sync.Mutex
	lru *lru.Cache
	cacheBytes int64
}

/*
延迟初始化
在add方法中 ， 判断c.lru是否为nil，如果等于nil再创建
这种方法称为延迟初始化，一个对象的延迟初始化意味着该对象的
创建将会延迟至第一次使用该对象时。
*/
func (c *cache) add(key string, value ByteView){
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil{
		c.lru = lru.New(c.cacheBytes,nil)
	}
	c.lru.Add(key,value)
}

func (c *cache) get(key string) (value ByteView,ok bool){
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru==nil {
		return
	}
	if v,ok:=c.lru.Get(key);ok{
		return v.(ByteView),ok
	}
	return
}