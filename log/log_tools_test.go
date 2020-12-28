package utils

import (
	"fmt"
	"strconv"
	"testing"
)

func TestCatchException(t *testing.T) {
	s := []string{
		"1", "-", "3", "4",
	}
	for _, v := range s {
		_, e := strconv.Atoi(v)
		ExceptionLog(e, fmt.Sprintf("test fail id = %v", 10))
	}
}
