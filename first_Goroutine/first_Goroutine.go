package main

import (
	"fmt"
	"sync"
	"time"
)

// We pass a pointer to the wg so that it can decrement the wait counter
// We have to do this to al the functions that we want to run as goroutines,
// otherwise the main function will end before the goroutines finish their work

func firstgoroutine(pais string, wg *sync.WaitGroup) {

	// This will decrement the wait counter when the function finishes
	// (whether it's due to an error or after a normal execution)
	defer wg.Done()

	for i := 0; i < 3; i++ {

		fmt.Printf("Pais: %s \n", pais)
		time.Sleep(500 * time.Millisecond) // Sleep for 500 milliseconds to simulate work
	}

}

func main() {

	var wg sync.WaitGroup
	var pais string = "Argentina"

	wg.Add(3)

	go firstgoroutine(pais, &wg) // Start the goroutine and pass the wait group
	pais = "Colombia"
	go firstgoroutine(pais, &wg)
	pais = "Venezuela"
	go firstgoroutine(pais, &wg)

	wg.Wait() // wait until the waitgroup is done

	fmt.Printf("Te waitgroup is done")
}
