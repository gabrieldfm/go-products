package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gabrieldfm/go-products/internal/db"
)

var (
	wg   sync.WaitGroup
	urls []string = []string{"https://www.linkedin.com/feed/", "https://www.linkedin.com/notifications/?filter=all"}
)

type VisitedLink struct {
	Link        string    `bson: "link"`
	VisitedDate time.Time `bson: "visiteddate"`
}

func periodicTask(url string) {
	fmt.Println("Request url " + url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Request failed " + err.Error())
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error when sending request: ", err.Error())

	}
	defer resp.Body.Close()

	visitedLink := VisitedLink{
		Link:        url,
		VisitedDate: time.Now(),
	}

	db.Insert("links", visitedLink)

	fmt.Println("Executing worker:", time.Now().In(time.FixedZone("America/Sao_Paulo", -3*60*60)).Format(time.RFC3339))
	fmt.Println(resp.StatusCode)
	time.Sleep(1 * time.Second)

}

// workerContinuousExecution goroutine for periodic task
func workerContinuousExecution(ctx context.Context) {
	defer wg.Done()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	fmt.Println("Worker started.")

	periodicTask(urls[0])

	for {
		select {
		case <-ticker.C:
			for _, v := range urls {
				periodicTask(v)
			}
		case <-ctx.Done():
			fmt.Println("Worker execution, reciving signal to shutdown.")
			return
		}
	}
}

func main() {
	fmt.Println("Start worker...")
	wg.Add(1)

	// Create context to control worker life cycle
	ctx, cancel := context.WithCancel(context.Background())

	go workerContinuousExecution(ctx)

	// Interruption signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for interruption signals
	sig := <-sigChan
	fmt.Printf("Shutdown signal received: %v\n", sig)
	cancel()

	fmt.Println("Waiting worker execution...")
	wg.Wait()

	fmt.Println("Service finish.")
}
