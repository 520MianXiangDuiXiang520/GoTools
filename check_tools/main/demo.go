package main

import (
	"fmt"
	"github.com/520MianXiangDuiXiang520/GinTools/check_tools"
)

type DemoFather struct {
	FName string `check:"not null"`
}

type Demo struct {
	DemoFather
	Name string `check:"not null; len:[0, 12];"`
	Age  int    `check:"not null; size: [1, 150]"`
	Mail string `check:"not null; email"`
}

func main() {
	req := &Demo{}
	req.FName = "12"
	req.Name = "1"
	req.Age = 10
	req.Mail = "15364968962@163,com"
	if check_tools.CheckRequest(req) {
		fmt.Printf("pass")
	}
}

// 2020/11/21 15:07:39 [Check] 15364968962@163,com is not a email
// 2020/11/21 15:07:39 [check] [Mail] Failed label inspection
