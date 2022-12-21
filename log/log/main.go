package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func main() {
	l := log.Default()

	f, _, _, ok := runtime.Caller(1)
	if !ok {
		println("functionName: runtime.Caller: failed")
		os.Exit(1)
	}

	format := "print info %s"
	a := "hello, world"
	funcShortName := filepath.Base(runtime.FuncForPC(f).Name())
	info := fmt.Sprintf(format, a)

	l.Printf("[%v] %v():: %v", "Info", funcShortName, info)
}
