package main

import (
	"fmt"
	"time"
)

func squared(a float32, channel chan float32) {

	time.Sleep(500 * time.Millisecond)
	channel <- (a * a)
}

func main() {

	var x float32 = 5
	var y float32 = 0
	channel := make(chan float32)

	fmt.Printf("Go channel example \n")

	go squared(x, channel)

	y = <-channel

	fmt.Printf("El cuadrado de %f es: %f \n", x, y)

}
