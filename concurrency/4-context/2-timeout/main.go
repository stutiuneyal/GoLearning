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

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go worker(ctx)

	time.Sleep(5 * time.Second)

}
