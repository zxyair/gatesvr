package xos_test

import (
	"gatesvr/utils/xos"
	"testing"
)

func TestCreate(t *testing.T) {
	_, err := xos.Create("./pprof/server/cpu_profile")
	if err != nil {
		t.Fatal(err)
	}

}
