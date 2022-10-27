package main

import (
	"fmt"
	"math/rand"
	"time"
)

// 生成数据文件
func generateData(fileName string, n int) {
	rand.Seed(time.Now().Unix())

	fmt.Printf("Random data: generating... \n")

	// 生成随机数
	p := randomSource(n)
	writeData(fileName, p)

	fmt.Printf("Random data: generated. \n")
}

// 生成随机数，并返回一个channel
func randomSource(count int) <-chan int32 {
	out := make(chan int32)
	go func() {
		for i := 0; i < count; i++ {
			v := rand.Int31()%100 + 1
			out <- v
		}
		close(out)
	}()
	return out
}
