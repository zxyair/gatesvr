package filter

import "net"

func BlackListCheck(addr net.Addr) bool {
	return false
}
