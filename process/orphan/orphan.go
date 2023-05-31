package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	name := "/bin/bash"
	argv := []string{name, "-c", `sleep 100 && echo "child: exit"`}
	attr := &os.ProcAttr{
		Dir:   "/tmp",                                     // workdir of new processor
		Env:   os.Environ(),                               // env list of new processor
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}, // stream output for new process
	}

	// create sub processor
	pro, err := os.StartProcess(name, argv, attr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("parent pid is: %d,chile pid is: %d", os.Getpid(), pro.Pid)

	// sleep 30s for main processor
	time.Sleep(30 * time.Second)
	fmt.Println("parent out!!")
}
