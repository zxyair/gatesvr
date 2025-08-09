package logics

import "gatesvr/cluster/client"

// 强制断开连接
func ForceCloseConn(conn *client.Conn) {
	conn.Close(true)
}

// 优雅断开连接
func GreceCloseConn(conn *client.Conn) {
	conn.Close(false)
}
