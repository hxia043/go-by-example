package main

//#include <stdio.h>
import "C"

func main() {
	s := C.CString("hello world.")
	C.puts(s)
}
