package array

import (
	"testing"
)

func TestGetVersion(t *testing.T) {
	a, _ := NewU32(nil, nil)
	if a.GetVersion() != "a32" {
		t.Errorf("version should be a32")
	}
}
