package main

import (
	"fmt"
	"sync"
)

// stage 1: gen
func gen(nums ...int) <-chan int {
	out := make(chan int, len(nums))
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

// stage 2: sq
func sq(in <- chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

// stage 3: print
func printSqs() {
	// set up pipeline
	// gen func provides the outbound channel
	c1 := gen(2,4,6,8,10)
	c2 := gen(3,5,7,9,11)

	for n := range merge(sq(c1), sq(c2)) {
		fmt.Println(n)
	}
}

func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	wg.Add(len(cs))
	out := make(chan int, 1)

	// new gorountine for for input channel to copy val to output chan
	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	for _, c := range cs {
		go output(c)
	}

	// separate go routine which closes out channel after all input channels are closed
	// synchronized using WaitGroup
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// main function
func main()  {
	printSqs()
}