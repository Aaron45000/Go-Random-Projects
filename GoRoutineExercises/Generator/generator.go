package main

import (
	"fmt"
	"time"
)

func generator(channel chan int) {

	for i := 1; i <= 5; i++ {

		time.Sleep(500 * time.Millisecond)
		channel <- (2 * i)
	}

	close(channel)

	fmt.Printf("The channel was closed \n")

}

func main() {

	channel := make(chan int)

	fmt.Printf("Generator function \n")

	go generator(channel)

	// Like and iterator for the values sent on the channel
	for numero := range channel {

		fmt.Printf("current even number: %d \n", numero)
	}

	fmt.Printf("Generator close the channel \n")

}
