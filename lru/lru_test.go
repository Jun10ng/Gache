package lru

import (
	"testing"
)
type String string
func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	t.Log("=======TEST GET=====")
	lru := New(int64(0), nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
	t.Log("========TEST GET OVER==========")
}

func TestCache_RemoveOldest(t *testing.T) {
	t.Log("=======TEST REMOVE_OLDEST=======")
	k1,k2,k3 := "key1","key2","key3"
	v1,v2,v3 := "value1","value2","value3"
	cap:=len(k1+k2+v1+v2)
	lru := New(int64(cap),nil)
	lru.Add(k1,String(v1))
	lru.Add(k2,String(v2))
	lru.Add(k3,String(v3))

	if _,ok := lru.Get("key1");ok {
		t.Log("Remove oldest element failed")
	}else{
		v,_ := lru.Get("key2")
		t.Logf("key2 ' value is %s\n",v)
	}
	t.Log("=======TEST REMOVE_OLDEST OVER=======")
}