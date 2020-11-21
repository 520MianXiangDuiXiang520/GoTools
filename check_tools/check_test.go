package check_tools

import (
	"fmt"
	"testing"
)

type Father struct {
	FName string `check:"not null"`
}

type TestStruct struct {
	Father
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

func TestCheckRequest2(t *testing.T) {
	s := "111"
	CheckRequest(&s)
}
