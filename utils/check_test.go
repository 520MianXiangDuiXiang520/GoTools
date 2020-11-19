package utils

import (
	"fmt"
	"testing"
)

type TestStruct struct {
	Name string `check:"not null; len:[0, 12];"`
	Age  int    `json:"age" check:"not null; size: [1, 150]"`
	Q    []int  `check:"len: [1, 3]"`
}

func TestCheckRequest(t *testing.T) {
	ok := CheckRequest(&TestStruct{
		Name: "name",
		Age:  10,
	})
	if ok {
		t.Error("fail")
	}
}

func TestFindNum(t *testing.T) {
	fmt.Println(findNum("len: [0, 10]"))
}
