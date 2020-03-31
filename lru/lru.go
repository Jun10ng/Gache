package lru

import (
	"container/list"
)


type Value interface {

	//返回占用的内存大小
	Len() int
}

type entry struct {
	key string
	value Value
}

type Cache struct {
	//允许使用的最大内存
	maxBytes int64

	//当前已使用的内存
	nbytes int64


	ll *list.List

	cache map[string] *list.Element

	//某条记录被移除时的回调函数，可以是nil
	OnEvicted func(key string, value Value)

}

//查找
func (c *Cache) Get(key string)(value Value, ok bool){
	if e,ok := c.cache[key];ok{
		c.ll.MoveToFront(e);
		kv:=e.Value.(*entry)
		return kv.value,true
	}
	return
}

//增加and修改
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		//回调函数不为空的情况下，使用回调函数
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

