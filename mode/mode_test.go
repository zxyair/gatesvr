package mode_test

import (
	"flag"
	"testing"

	"gatesvr/mode"
)

func TestGetMode(t *testing.T) {
	flag.Parse()

	t.Log(mode.GetMode())
}
