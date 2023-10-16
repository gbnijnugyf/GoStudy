package main

import "fmt"

func CalSquare() {
	src := make(chan int)
	//带缓冲的channel规避生产与消费速度不匹配问题
	dest := make(chan int, 3)
	go func() {
		defer close(src)
		for i := 0; i < 10; i++ {
			src <- i
		}
	}()
	go func() {
		defer close(dest)
		for i := range src {
			dest <- i * i
		}
	}()
	for i := range dest {
		fmt.Println(i)
	}
}

func main() {
	CalSquare()
}
