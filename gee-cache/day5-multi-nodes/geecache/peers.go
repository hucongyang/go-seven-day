package geecache

type PeerPicker interface {
	// 根据传入的key选择相应节点PeerGetter
	PickPeer(key string) (peer PeerGetter, ok bool)
}
// 节点接口
type PeerGetter interface {
	// 从对应group查找缓存值, PeerGetter 对应流程中的HTTP客户端
	Get(group string, key string) ([]byte, error)
}