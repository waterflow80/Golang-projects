package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	userName := fetchUser()
	respch := make(chan any, 2) // The size doesn't matter 
	wg := &sync.WaitGroup{}

	wg.Add(2)
	go fetchUserLikes(userName, respch, wg)
	go fetchUserMatch(userName, respch, wg)

	wg.Wait() // block until 2 wg.Done()
	close(respch)

	for 	resp := range respch {
		fmt.Println("resp ", resp)
	}

	fmt.Println("took ", time.Since(start))

}

func fetchUser() string {
	time.Sleep(time.Microsecond * 100)	

	return "BOB"
}

func fetchUserLikes(userName string, respch chan any, wg *sync.WaitGroup) {
	time.Sleep(time.Microsecond * 150)

	respch <- 11
	wg.Done()
}

func fetchUserMatch(userName string, respch chan any, wg *sync.WaitGroup) {
	time.Sleep(time.Microsecond * 100)

	respch <- "ANNA"
	wg.Done()
}