package Gache

import (
	"fmt"
	"log"
	"sync"
)

/*
负责与外部交互
控制缓存存储和获取的主流程
*/

//回调函数
//当缓存不存在，从其他数据源获取，并添加到缓存
type Getter interface {
	Get(key string)([]byte,error)
}

/*
定义一个函数类型 F，并且实现接口 A 的方法，
然后在这个方法中调用自己。
这是 Go 语言中将其他函数转换为接口 A 的常用技巧
（参数返回值定义与 F 一致）
*/

type GetterFunc func(key string)([]byte,error)

//
func (f GetterFunc) Get(key string)([]byte,error){
	return f(key)
}

/*
Groups
缓存的集合
Groups =Group的集合= n个缓存表
注意区分Groups和Group
*/

type Group struct {
	name string
	getter Getter
	mainCache cache
}

var (
	mu  sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string,cacheBytes int64, getter Getter) *Group{
	//getter是一个函数类型的参数
	if getter == nil{
		panic("nil Getter")
	}

	mu.Lock()
	defer mu.Unlock()
	g:= &Group{
		name: name,
		getter: getter,
		mainCache:cache{cacheBytes:cacheBytes},
	}
	groups[name] = g
	return g
}

//返回相应的group
func GetGroup(name string) *Group{
	mu.RLock()
	g:=groups[name]
	mu.RUnlock()
	return g
}

//get
func (g *Group)Get(key string) (ByteView,error){
	if key==""{
		return ByteView{},fmt.Errorf("key is required")
	}

	if v,ok:=g.mainCache.get(key);ok{
		log.Println("[Gache] hit")
		return v,nil;
	}
	//从其他数据源获取
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}
//从其他数据源获取
func (g *Group) getLocally(key string)(ByteView,error) {
	bytes,err := g.getter.Get(key)
	if err!=nil{
		return ByteView{},err
	}
	value := ByteView{b:cloneBytes(bytes)}
	//将远程数据添加到当前缓存
	g.popularCache(key,value)
	return value,nil
}

func (g *Group) popularCache(key string,value ByteView) {
	g.mainCache.add(key,value)
}








