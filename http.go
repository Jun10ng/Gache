package Gache

import (
	"Gache/consistent"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const(
	defaultBasePath = "/_Gache"
	/*
		节点数，后续代码添加
	*/
	defaultReplicas = 50
)



//线程池
/*
	self 用来记录自己的地址，包括主机名/IP和端口
	basePath 作为节点间通讯地址的前缀，默认是/_Gache/
	以http://XXX.com/_Gache/开头的请求，就用于节点间的访问

*/
type HTTPPool struct {
	self 	 string
	basePath string
	/*
		后续代码添加
	*/
	/*
		监视节点和httpGetters
	*/
	mu 		    sync.Mutex
	/*
		一致性哈希算法的Map，见readMe中的一致性哈希章节的示意图
		用来根据具体的key选择节点
	*/
	peers       *consistent.Map
	/*
		key值为：http://xxxxx
		每个远程连接节点对应一个httpGetter
		因为httpGetter与远程节点的地址 baseURL有关
	*/
	httpGetters map[string]*httpGetter
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

/*
	httpGetter 实现 peer.go中的peerGetter接口
	baseURL 表示将要访问的远程节点的地址
	例如http://example.com/gache/
	使用http.Get()方式获取返回值，并转换为[]bytes类型
*/
type httpGetter struct{
	baseURL string
}

func (hg *httpGetter)Get(group string,key string)([]byte,error){
	var b  strings.Builder
	b.WriteString(hg.baseURL)
	/*
		url.QueryEscape用于转义url字段
		可以被安全地防止
	*/
	b.WriteString(url.QueryEscape(group))
	b.WriteString(url.QueryEscape(key))
	s := b.String()

	res,err := http.Get(s)
	if err != nil {
		return nil,err
	}
	defer res.Body.Close()

	if res.StatusCode!= http.StatusOK {
		return nil,fmt.Errorf("server returned : %v",res.Status)
	}

	/*
		bytes是字节类型的数组
		具体设计见cache.go和lru.go
	*/
	bytes,err := ioutil.ReadAll(res.Body)

	/*
		异常
	*/
	if err != nil{
		return nil, fmt.Errorf("reading response body: %v",err)
	}

	return bytes,nil
}

var _ PeerGetter = (*httpGetter)(nil)



/*
	实现PeerPicker接口
*/

/*
	Set()实例化了一致性哈希算法
	并添加了传入的节点
*/
func (p *HTTPPool)Set(peers ...string){
	p.mu.Lock()
	defer p.mu.Unlock()
	/*
		实例化
	*/
	p.peers = consistent.New(defaultReplicas,nil)
	/*
		添加传入节点
	*/
	p.peers.Add(peers...)
	/*
		为每个节点配置Getter
	*/
	p.httpGetters = make(map[string]*httpGetter,len(peers))

	/*
		peer如http://xxxx.com/
		basePath如 _Gache/
		合起来就是一个节点的访问路径
	*/
	for _,peer := range peers{
		p.httpGetters[peer] = &httpGetter{
			baseURL:peer+p.basePath,
		}

	}
}
/*
	PickerPeer() 包装了一致性哈希算法的 Get() 方法，
	根据具体的 key，选择节点，返回节点对应的 HTTP 客户端。
	注意，这只是返回一个节点，而不是返回缓存的value
	返回缓存的实现在gache.go的load()中
*/
func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

var _ PeerPicker = (*HTTPPool)(nil)

















