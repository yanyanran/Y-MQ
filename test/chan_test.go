package main

import (
	"fmt"
	"testing"
)

var flagerr int

func Max() bool {
	return false
}

func TestChan(t *testing.T) {
	ch := make(chan int, 5)
	/*ch <- 18*/
	close(ch)
	err := Max()
	if !err {
		flagerr = 1
	}
	// -------------
	if flagerr == 1 {
		return
	}
	x, ok := <-ch
	if ok {
		fmt.Println("ch received: ", x)
	} else {
		fmt.Println("[1]channel closed, data invalid.")
	}

	x, ok = <-ch
	if !ok {
		fmt.Println("[2]channel closed, data invalid.")
	}

	chh := make(chan int, 1) // 0-> deadlock
	fmt.Println(111111111111111111)
	chh <- 18
	fmt.Println(222222222222222222)
	close(chh)
	y, okk := <-chh
	if okk {
		fmt.Println("chh received: ", y)
	} else {
		fmt.Println("[3]channel closed, data invalid. chh received: ", y)
	}

	y, okk = <-chh
	if !okk {
		fmt.Println("[4]channel closed, data invalid. chh received: ", y) // 0
	}
}
