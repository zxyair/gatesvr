package xstring_test

import (
	"gatesvr/utils/xstring"
	"testing"
)

func Test_PaddingPrefix(t *testing.T) {
	t.Log(xstring.PaddingPrefix("1", "0", 3))
	t.Log(xstring.PaddingPrefix("001", "0", 3))
	t.Log(xstring.PaddingPrefix("0001", "0", 3))
	t.Log(xstring.PaddingPrefix("1", "00", 3))
}
