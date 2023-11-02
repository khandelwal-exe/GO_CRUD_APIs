package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("Race Condition")
	wg := &sync.WaitGroup{}
	mut := &sync.Mutex{}

	var score = []int{} // Initialize score as an empty slice

	wg.Add(3)
	go func(Wg *sync.WaitGroup, mut *sync.Mutex) {
		fmt.Println("One R")
		mut.Lock()
		score = append(score, 1)
		mut.Unlock()
		wg.Done()
	}(wg, mut)
	go func(Wg *sync.WaitGroup, mut *sync.Mutex) {
		fmt.Println("Two R")
		mut.Lock()
		score = append(score, 2)
		mut.Unlock()
		wg.Done()
	}(wg, mut)
	go func(Wg *sync.WaitGroup, mut *sync.Mutex) {
		fmt.Println("Three R")
		mut.Lock()
		score = append(score, 3)
		mut.Unlock()
		wg.Done()
	}(wg, mut)
	wg.Wait()
	fmt.Println(score)
}
