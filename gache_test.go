package Gache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)
var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGetterFunc_Get(t *testing.T) {
	var f Getter  = GetterFunc(func(key string) ([]byte,error) {
		return []byte(key),nil
	})

	expect := []byte("key")
	if v,_ := f.Get("key");!reflect.DeepEqual(v,expect){
		t.Error("callback failed")
	}
}

/*
Group Test
*/

/*
在这个测试用例中，我们主要测试了 2 种情况
1）在缓存为空的情况下，能够通过回调函数获取到源数据。
2）在缓存已经存在的情况下，是否直接从缓存中获取，为了实现这一点，
使用 loadCounts 统计某个键调用回调函数的次数，如果次数大于1，则表示调用了多次回调函数，没有缓存。*/
func TestGet(t *testing.T) {
	loadCounts:= make(map[string]int,len(db))
	gs := NewGroup("scores",2<<10,GetterFunc(
		func(key string)([]byte,error) {
			log.Println("[SlowDB] search key",key)
			if v,ok := db[key];ok{
				if _,ok:= loadCounts[key];!ok{
					loadCounts[key] = 0;
				}
				loadCounts[key]+=1
				return []byte(v),nil
			}
			return nil, fmt.Errorf("%s not exist",key)
		}))
	for k,v := range db{
		if view , err := gs.Get(k);err!=nil|| view.toString()!=v{
			t.Fatal("failed to get value of TOM")
		}
		if _,err:=gs.Get(k);err!=nil|| loadCounts[k]>1 {
			t.Fatalf("cache %s miss",k)
		}

		if view,err := gs.Get("unkown");err == nil{
			t.Fatalf("the value of unkown should be empty,but %s got",view)
		}
	}
}





