package main

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Worker stopper: ", ctx.Err())
			return
		default:
			fmt.Println("Working ...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	for i := 1; i <= 3; i++ {
		go worker(ctx)
	}

	time.Sleep(2 * time.Second)

	cancel()

	time.Sleep(1 * time.Second)

}
