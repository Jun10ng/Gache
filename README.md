
##golang 实现 lru

 代码

lru/lru.go

```
package lru

import (
	"container/list"
	//"github.com/astaxie/beego/validation"
)

type lru struct {
	maxSize int64
	curSize int64
	ll *list.List
	cache map[string]*list.Element
}

/*
链表节点中的Value值的结构体
一定要定义一个结构体来存放key和value
因为，在lru中，map的value是指向链表的节点
当我们删除节点的时候，不仅要删除链表中的节点，
还要删除map的值，所以需要从链表中取key，然后反推map，
删除键值对，这样才是完整的删除
*/
type entry struct{
	key string
	value string
	/*
	实际使用的时候，value应该定义为interface，
	或者是自己声明的一个interface
	*/
}

func NewLru(maxSize int64) *lru {
	return &lru{
		maxSize:maxSize,
	}
}

func (lru *lru) get(key string) (s string ,ok bool ){
	if e,ok := lru.cache[key]; ok{
		lru.ll.MoveToFront(e)
		entry := e.Value.(*entry)
		return entry.value,true
	}
	return

}

func (lru *lru) add(key string,value string){
	if e,ok:= lru.cache[key];ok {
		//更新
		lru.ll.MoveToFront(e)
		e.Value.(*entry).value = value
	} else {
		//添加
		newE := &entry{key,value}
		lru.cache[key] = lru.ll.PushFront(newE)
		lru.curSize++
	}
	for lru.maxSize>0 && lru.curSize > lru.maxSize {
		lru.removeOldest()
	}
}

func (lru *lru) removeOldest() {
	//链表
	e:= lru.ll.Back()
	lru.ll.Remove(e)
	//哈希表
	delete(lru.cache,e.Value.(*entry).key)
	lru.curSize--
}




```



### 测试代码

**记得使用go.mod**

lru/lru_test.go

```
package lru

import (
	//"reflect"
	"testing"
)

func TestLru_Get(t *testing.T) {
	lru := New(int64(2))
	lru.Add("key1","value1")
	lru.Add("key2","value2")
	if v1,ok:=lru.Get("key1");ok {
		t.Logf("key1 's value is %s",v1)
	}else { t.Log("missing key1")}

	if v2, ok := lru.Get("key2"); ok {
		t.Logf("key2's value is %s",v2)
	}else { t.Log("missing key2")}
}

func TestLru_RemoveOldest(t *testing.T) {
	lru:=New(int64(2))
	lru.Add("key1","value1")
	lru.Add("key2","value2")
	lru.Add("key3","value3")
	//key1 应该被删除

	if _,ok:= lru.Get("key");!ok{
		t.Log("missing key1")
	}else{
		t.Log("remove oldest failed!")
	}

}
```

## 实现单机并发
cache.go
主要是对lru内的数据结构的方法加锁，
这里有个小点注意下，cache并没有new方法，因为采用的是延迟初始化
在add方法中，判断c.lru是否为nil，如果等于nil再创建
这种方法称为延迟初始化，一个对象的延迟初始化意味着该对象的
创建将会延迟至第一次使用该对象时。
这个方法在redis中很常见，因为能一定程度上提高性能
```
type cache struct {
	mu sync.Mutex
	lru *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView){
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil{
		c.lru = lru.New(c.cacheBytes,nil)
	}
	c.lru.Add(key,value)
}
```
##http调用
TODO

## 一致性哈希
TODO