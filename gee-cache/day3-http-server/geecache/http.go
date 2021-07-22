package geecache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// 提供被其他节点访问的能力 （基于http）

const defaultBasePath = "/_geecache/"

// 承载节点间HTTP通信的核心结构(包含服务端和客户端)
type HTTPPool struct {
	self		string		// 记录自己的地址，包括主机名/IP和端口
	basePath	string		// 作为节点间通讯地址的前缀，默认是 /_geecache/,
	// 那么http://example.com/_geecache/ 开头的请求，就用于节点间的访问，
	// 因为一个主机上还可能承载其他的服务，加一段Path是一个好习惯，比如，大部分网站的api接口，一般以 /api 作为前缀
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