package path_tools

import (
	"path"
)

func IsAbs(p string) bool {
	if path.IsAbs(p) {
		return true
	}
	if len(p) > 3 && (p[0] >= 'A' && p[0] <= 'Z') && p[1] == ':' && (p[2] == '/' || p[2] == '\\') {
		return true
	}
	return false
}
