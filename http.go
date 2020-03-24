package Gache

import (

	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath  = "/_Gache/"

//线程池
/*
	self 用来记录自己的地址，包括主机名/IP和端口
	basePath 作为节点间通讯地址的前缀，默认是/_Gache/
	以http://XXX.com/_Gache/开头的请求，就用于节点间的访问

*/
type HTTPPool struct {
	self 	 string
	basePath string
}

func NewHTTPPool(self string) *HTTPPool  {
	return &HTTPPool{
		self:	self,
		basePath:	defaultBasePath,
	}
}

// Core Code

/*
ServeHTTP 的实现逻辑是比较简单的，首先判断访问路径的前缀是否是 basePath，不是返回错误。
我们约定访问路径格式为 /<basepath>/<groupname>/<key>，
通过 groupname 得到 group 实例，再使用 group.Get(key) 获取缓存数据。
最终使用 w.Write() 将缓存值作为 httpResponse 的 body 返回。
*/
func (p *HTTPPool) Log(format string,v ...interface{}){
	log.Printf("[Server %s] %s",p.self,fmt.Sprintf(format,v...))
}

func (p *HTTPPool) ServeHTTP (w http.ResponseWriter,r *http.Request){
	if !strings.HasPrefix(r.URL.Path,p.basePath){
		panic("HTTPPool serving unexpected path:"+r.URL.Path)
	}
	p.Log("%s %s",r.Method,r.URL.Path)
	parts:=strings.SplitN(r.URL.Path[len(p.basePath):],"/",2)
	if len(parts)!= 2{
		http.Error(w,"bad request",http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key:=parts[1]

	group := GetGroup(groupName)

	if group == nil{
		http.Error(w,"no such group:" + groupName,http.StatusNotFound)
		return
	}

	view,err := group.Get(key)
	if err!=nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}

