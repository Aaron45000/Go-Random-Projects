package main

import "fmt"

func main() {

	var numero int = 15

	fmt.Printf("Hola, mundo %d\n", numero)

	for i := 0; i <= numero; i++ {
		if i%5 == 0 {

			fmt.Printf("%d es multiplo de 5 \n", i)

		} else {

			fmt.Printf("%d\n", i)

		}

	}

}
