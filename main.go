package main

import (
	"context"
	"fmt"
	"github.com/seanlee0923/plscli"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	cfg := plscli.Config("http://172.16.66.128", "1923", "test-deploy")
	server := plscli.NewClient(cfg)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.RunWithContext(ctx); err != nil {
			fmt.Println(err)
		}
	}()

	// 종료신호 받기
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go testInterval(server)

	<-sigCh
	fmt.Println("Shutting down...")

	cancel()
	wg.Wait()
}

func testInterval(c *plscli.PlsClient) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			leader, err := c.IsLeader()

			fmt.Println(leader, err)
		}
	}
}
