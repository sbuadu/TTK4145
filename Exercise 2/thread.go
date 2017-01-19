// Go 1.2
// go run thread.go

package main

import (
    . "fmt"
    "runtime"
    "time"
)

func someGoroutine1(ch chan int) {
	for j := 0; j < 1000000; j++{
		i := <-ch
		i += 1
		ch <- i
	}
}

func someGoroutine2(ch chan int) {
	for j := 0; j < 1000000; j++{
		i := <-ch
		i -= 1
		ch <- i
	}
}
func main() {
	ch := make(chan int, 1)
	ch <- 0
	runtime.GOMAXPROCS(runtime.NumCPU())    // I guess this is a hint to what GOMAXPROCS does...
                                            // Try doing the exercise both with and without it!
    	go someGoroutine1(ch)                      // This spawns someGoroutine() as a goroutine
	go someGoroutine2(ch)
    // We have no way to wait for the completion of a goroutine (without additional syncronization of some sort)
    // We'll come back to using channels in Exercise 2. For now: Sleep.
    time.Sleep(1000*time.Millisecond)
    Println(<-ch)
}
