package main

import (
	"sync"
)

var wg sync.WaitGroup

//使用携程和管道打印字符串
func PrintMain() {
	p := "abcdefg"
	c := make(chan byte, 1000)
	for _, v := range []byte(p) {
		println(v)
		c <- v
	}
	go func() {
		for {
			printString, ok := <-c
			if !ok {
				close(c)
				break
			}
			println(string(printString))

		}
	}()
	wg.Done()
}

func main() {
	wg.Add(1)
	PrintMain()
	wg.Wait()
}
