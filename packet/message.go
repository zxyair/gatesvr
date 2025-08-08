package packet

type Message struct {
	Seq        int32  // 序列号
	Route      int32  // 路由ID
	IsCritical bool   // 是否关键消息
	Buffer     []byte // 消息内容
}
type Notification struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
