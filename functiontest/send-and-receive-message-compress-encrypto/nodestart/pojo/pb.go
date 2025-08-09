package pojo

// 请求
type GreetReq struct {
	Message string `json:"message"`
}
type AuthuritionReq struct {
	Message int64 `json:"message"`
}
type AuthuritionRes struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 响应
type GreetRes struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
