package version

import "testing"

func TestVersion(t *testing.T) {
	if VERSION == "" {
		t.Fatalf("version.VERSION is empty")
	}

	if MAXLEN != 16 {
		t.Fatalf("version.MAXLEN must be 16 and can not be changed")
	}
}
