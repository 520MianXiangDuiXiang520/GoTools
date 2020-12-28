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

func TestCheck2(t *testing.T) {
	s := &TestStruct{
		Name: "12345",
		Age:  1234,
		Q:    []int{1, 2},
	}
	s.FName = "123"
	ok := Check(s)
	if ok {
		t.Errorf("")
	}
}

func TestFindNum(t *testing.T) {
	fmt.Println(findNum("len: [0, 10]"))
}

func TestCheckRequest2(t *testing.T) {
	s := "111"
	CheckRequest(&s)
}

type numTest struct {
	NumInt    int    `check:"more: 10"`
	NumUint   uint   `check:"less: 20"`
	NumInt64  int64  `check:"equal: 32"`
	NumInt32  int32  `check:"equal: 32"`
	NumUint32 uint32 `check:"equal: 32"`
}

// 2020/12/27 测试数值大小检查
func TestCheck(t *testing.T) {
	n := numTest{
		NumInt:    15,
		NumUint:   uint(10),
		NumInt32:  int32(32),
		NumInt64:  int64(32),
		NumUint32: uint32(32),
	}
	if !Check(n) {
		t.Error()
	}
}
