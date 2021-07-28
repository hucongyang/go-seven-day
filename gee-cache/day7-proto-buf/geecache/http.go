package geecache

import (
	"fmt"
	"geecache/consistenthash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// 提供被其他节点访问的能力 （基于http）

// 添加节点选择的功能
const (
	defaultBasePath = "/_geecache/"
	defaultReplicas = 50
)

// 承载节点间HTTP通信的核心结构(包含服务端和客户端)
type HTTPPool struct {
	self		string		// 记录自己的地址，包括主机名/IP和端口
	basePath	string		// 作为节点间通讯地址的前缀，默认是 /_geecache/,
	// 那么http://example.com/_geecache/ 开头的请求，就用于节点间的访问，
	// 因为一个主机上还可能承载其他的服务，加一段Path是一个好习惯，比如，大部分网站的api接口，一般以 /api 作为前缀

	mu 			sync.Mutex
	peers 		*consistenthash.Map      // 类型是一致性哈希算法的 Map，用来根据具体的key选择节点
	httpGetters	map[string]*httpGetter   // 映射远程节点与对应的httpGetter，每一个远程节点对应一个 httpGetters
										 // 因为httpGetter与远程节点的地址baseURL有关
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self: 		self,
		basePath: 	defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 判断路由访问路径是否是 basePath，不是返回错误
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	// 约定访问路径格式为 /<basepath>/<groupname>/<key>，通过groupname得到group实例，再使用group.Get(key)获取缓存数据
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: " + groupName, http.StatusNotFound)
		return
	}
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}

// 实现客户端功能
type httpGetter struct {
	baseUrl		string  // 将要访问的远程节点的地址 如：http://example.com/_geecache/
}
// 使用http.Get()方式获取返回值，并转换为 []bytes 类型
func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseUrl,
		url.QueryEscape(group),
		url.QueryEscape(key),
		)
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}
	return bytes, nil
}

var _ PeerGetter = (*httpGetter)(nil)

// 实例化了一致性哈希算法，并且添加了传入的节点
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)    // 添加传入的节点
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	// 为每一个节点创建了一个HTTP客户端的httpGetter
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseUrl: peer + p.basePath}
	}
}
// 包装了一致性哈希算法的 Get() 方法，根据具体的key，选择节点，返回节点对应的 HTTP 客户端。
func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
  	p.mu.Lock()
  	defer p.mu.Unlock()
  	if peer := p.peers.Get(key); peer != "" && peer != p.self {
	 	p.Log("Pick peer %s", peer)
	 	return p.httpGetters[peer], true
  	}
	return nil, false
}





