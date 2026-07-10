package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	monitoringDuration = 10 * time.Second
)

type ServiceStatus struct {
	ServiceName string
	Status      string
	CheckedAt   time.Time
}

/*
statusChannel chan ServiceStatus -> two way channel, can send and receive
statusChannel chan<- ServiceStatus -> only send
statusChannel <-chan ServiceStatus -> only receive
*/
func montoringService(ctx context.Context, wg *sync.WaitGroup, serviceName string, interval time.Duration, statusChannel chan<- ServiceStatus) {
	defer wg.Done()

	// Ticker -> triggers at configured interval
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		// stop monitoring when the context is cancelled
		case <-ctx.Done():
			fmt.Printf("%s monitor stopped\n", serviceName)
			return

			// perform a health check whenever the ticker fires
		case checkedAt := <-ticker.C:

			status := ServiceStatus{
				ServiceName: serviceName,
				Status:      "Healthy",
				CheckedAt:   checkedAt,
			}

			// Use another select so the go routine does not remain blocked
			// while sending the cancellation
			select {
			case statusChannel <- status:
			case <-ctx.Done():
				return
			}
		}
	}

}

func displayStatus(status ServiceStatus) {
	fmt.Printf(
		"Service: %-20s | Status: %-10s | Checked At: %s\n",
		status.ServiceName,
		status.Status,
		status.CheckedAt.Format("15:04:05"),
	)
}

func main() {

	// create our timeout context
	ctx, cancel := context.WithTimeout(context.Background(), monitoringDuration)
	defer cancel()

	// create a shared channel
	statusChannel := make(chan ServiceStatus) //-> map[string]int

	// create a WaitGroup
	var wg sync.WaitGroup

	// Details of the service that needs to be monitored, and the frequency of monitoring
	services := map[string]time.Duration{
		"API Service":       1 * time.Second,
		"Database Service":  2 * time.Second,
		"Messaging Service": 3 * time.Second,
	}

	// start one goroutine for each service
	for serviceName, interval := range services {
		wg.Add(1)

		go montoringService(
			ctx,
			&wg,
			serviceName,
			interval,
			statusChannel,
		)
	}

	// wait for all the service monitors to stop before closing the channel
	go func() {
		wg.Wait()
		close(statusChannel)
	}()

	fmt.Println("Concurrent Service Health montoring started")
	fmt.Println("Montoring Duration: ", monitoringDuration)
	fmt.Println()

	// receive results until every monitor has stopped and channel is closed

	// for {
	// 	status, ok := <-statusChannel
	// 	if !ok {
	// 		fmt.Println("Status channel closed")
	// 		return
	// 	}

	// 	displayStatus(status)
	// }

	for status := range statusChannel {
		displayStatus(status)
	}

	fmt.Println()
	fmt.Println("Monitoring duration completed")
	fmt.Println("All monitors stopped")
	fmt.Println("Aplication completed gracefully")

}
