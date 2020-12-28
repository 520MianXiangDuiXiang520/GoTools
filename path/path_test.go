package path

import (
	"testing"
)

func TestIsAbs(t *testing.T) {
	if !IsAbs("/usr/bin") {
		t.Error()
	}
	if !IsAbs("E:\\GinTools\\path\\path_test.go") {
		t.Error()
	}
	if IsAbs("./src") {
		t.Error()
	}
}
