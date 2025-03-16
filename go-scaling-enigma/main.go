package main

import (
	"fmt"
	"scaling-enigma/go-scaling-enigma/main.go/server"
	"scaling-enigma/go-scaling-enigma/main.go/tui"
	"sync"
)

func main() {
	fmt.Println("Lets Go!")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.ServeGin()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		tui.StartCLI()
	}()

	wg.Wait()

}
