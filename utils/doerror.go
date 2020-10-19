package utils

import (
	"log"
	"runtime"
)

func ExceptionLog(e error, mes string) {
	if e != nil {
		pc, _, line, _ := runtime.Caller(1)
		fName := runtime.FuncForPC(pc).Name()
		log.Printf("[Error] %v:%v  %v", fName, line, mes)
		log.Printf("[Error] %v", e)
	}
}

func LogPlus(msg string) {
	pc, _, line, _ := runtime.Caller(1)
	fName := runtime.FuncForPC(pc).Name()
	log.Printf("[Info] %v:%v  %v", fName, line, msg)
}
