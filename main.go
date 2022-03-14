package main

import (
	"fmt"
	"health/checking"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var links = []string{
	"https://mail.ru/",
	"https://yandex.ru/",
}

func main() {
	outputChan := make(chan checking.Result)
	workerPool := checking.Run(10, time.Second*10, outputChan)

	workerPool.Start()

	go jobs(workerPool)
	go results(outputChan)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	workerPool.Stop()
}

func results(results chan checking.Result) {
	go func() {
		for r := range results {
			fmt.Println(r.GetInfo())
		}
	}()
}

func jobs(wp *checking.Pool) {
	for {
		for _, link := range links {
			wp.Push(checking.Job{Link: link})
		}
		time.Sleep(time.Second * 10)
	}
}
