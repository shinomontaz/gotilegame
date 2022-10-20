package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println(runtime.GOMAXPROCS(1))
	for i := 0; i < 3; i++ {
		go func() {
			var x int
			for {
				x++
			}
		}()
	}

	fmt.Println("ok")
}
