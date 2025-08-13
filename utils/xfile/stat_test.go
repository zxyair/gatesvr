package xfile_test

import (
	"gatesvr/utils/xfile"
	"testing"
)

func TestStat(t *testing.T) {
	fi, err := xfile.Stat("a.txt")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(fi.CreateTime())
}
