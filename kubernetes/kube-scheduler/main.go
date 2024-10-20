package main

import (
	"fmt"
	"math/rand"
	"os"
	"prebind/prebind"
	"time"

	"k8s.io/kubernetes/cmd/kube-scheduler/app"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	command := app.NewSchedulerCommand(
		app.WithPlugin(prebind.Name, prebind.New),
	)

	if err := command.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

}
