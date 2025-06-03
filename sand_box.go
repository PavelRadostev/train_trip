package main

import (
	"fmt"
	"time"
)

func greet(c chan string) {
	for {
		name := <-c
		fmt.Println("Hello " + name + "!")
		time.Sleep(1000 * time.Microsecond)
		c <- "Name received: " + name // Send a message back to the channel
	}
}

func main() {
	fmt.Println("main() started")
	c := make(chan string)

	go greet(c)

	c <- "Alice"

	resp := <-c
	c <- "Mike"
	fmt.Println(resp + " has been processed by the goroutine")
	fmt.Println("main() waiting for goroutine to finish")
	time.Sleep(1000 * time.Microsecond) // Give goroutine time to finish
	close(c)                            // Close the channel to signal no more data will be sent
	fmt.Println("main() stopped")
}
