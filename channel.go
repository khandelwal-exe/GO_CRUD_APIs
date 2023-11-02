package main

import (
	"fmt"
	"sync"
)

func worker(id int, results chan int, wg *sync.WaitGroup) {
	defer wg.Done() // Mark the worker as done when it's finished
	result := id * 2
	results <- result
}

func main() {
	results := make(chan int)
	var wg sync.WaitGroup

	// Start the first worker
	wg.Add(1)
	go worker(1, results, &wg)

	// Start the second worker
	wg.Add(1)
	go worker(2, results, &wg)

	// Use a goroutine to collect results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect and print the results
	for result := range results {
		fmt.Printf("Received result: %d\n", result)
	}
}
