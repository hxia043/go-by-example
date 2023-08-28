package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func handleGracefulShutdown(sum *int) {
	fmt.Println("clean app...")
	*sum = 0
	fmt.Printf("reset sleep start time to %d second\n", *sum)
	fmt.Println("clean app finished.")
}

func main() {
	shutdownsignal := make(chan os.Signal)
	go signal.Notify(shutdownsignal, syscall.SIGTERM, syscall.SIGKILL)

	gracefulshutdownsignal := make(chan os.Signal)
	go signal.Notify(gracefulshutdownsignal, syscall.SIGINT)

	// main logic of app
	sum := 0
	go func(sum int) {
		for {
			sum++
			time.Sleep(time.Second)
			fmt.Printf("sleep %d second...\n", sum)
		}
	}(sum)

	select {
	case s := <-shutdownsignal:
		fmt.Println("shutdown with signal: ", s)
	case gs := <-gracefulshutdownsignal:
		fmt.Println("graceful shutdown with signal: ", gs)
		handleGracefulShutdown(&sum)
	}
}
