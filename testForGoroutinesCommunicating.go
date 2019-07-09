package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	shutlist := make(map[string]bool)
	var index int
	cc := make(chan int, 1)
	hh := make(chan string, 1)
	shut1 := make(chan bool, 1)
	shut1 <- true
	shut2 := make(chan bool, 1)
	shut2 <- true
	shutlist["q1"] = true
	shutlist["q2"] = true
	shutlist["q"] = true
	switch len(os.Args) {
	case 1:
		fmt.Printf("os.Args: %v\n", os.Args)
		fmt.Printf("no input\n")
		cc <- 1
		hh <- "hello!"
		go defaultGO(cc, shut1)
		go alwaysGO(hh, shut2)
	case 2:
		fmt.Printf("os.Args: %v\n", os.Args)
		fmt.Printf("one input\n")
	case 3:
		fmt.Printf("os.Args: %v\n", os.Args)
		fmt.Printf("two inputs\n")
	}
	for {
		var ccc string
		_, _ = fmt.Scanln(&ccc)
		if ccc == "1" {
			if shutlist["q1"] == true {
				cc <- index
				shut1 <- true
			} else {
				fmt.Println("q1 goroutine is already shut!!!")
			}
		} else if ccc == "q" {
			fmt.Println("all shut down!")
			return
		} else if ccc == "q1" {
			cc <- index
			shut1 <- false
			shutlist["q1"] = false
		} else if ccc == "q2" {
			hh <- ccc + strconv.Itoa(index)
			shut2 <- false
			shutlist["q2"] = false
		} else {
			if shutlist["q2"] == true {
				hh <- ccc + strconv.Itoa(index)
				shut2 <- true
			} else {
				fmt.Println("q2 goroutine is already shut!!!")
			}
		}
		index++
	}
}

func defaultGO(cc chan int, shut chan bool) {
	var getIN int
	var index int
	var isShut bool
	for {
		getIN = <-cc
		fmt.Printf("index %d, getIN %d\n", index, getIN)
		index++
		isShut = <-shut
		if isShut == false {
			fmt.Println("shut goroutine1 now")
			return
		}
	}

}

func alwaysGO(hh chan string, shut chan bool) {
	var hhhIN string
	var index int
	var isShut bool
	for {
		hhhIN = <-hh
		fmt.Printf("index: %d, hhhIN %v\n", index, hhhIN)
		index++
		isShut = <-shut
		if isShut == false {
			fmt.Println("shut goroutine2 now")
			return
		}
	}
}
