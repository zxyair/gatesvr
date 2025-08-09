package mode_test

import (
	"flag"
	"testing"
)

func TestGetMode(t *testing.T) {
	flag.Parse()

	t.Log(GetMode())
}
