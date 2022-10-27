package main

import "fmt"

// 二路归并
func mergeData2(tmpFile []string, threads int) {
	fmt.Printf("")
	ins := make([]<-chan int32, threads)
	for i, fileName := range tmpFile {
		ins[i], _ = loadFile(fileName, 0, -1)
	}
	res := divideData(ins...)
	writeData("result", res)
}

// 递归实现二路归并
func divideData(in ...<-chan int32) <-chan int32 {
	l := len(in)
	if l == 1 {
		return in[0]
	}
	mid := l / 2
	return mergeData(divideData(in[:mid]...), divideData(in[mid:]...))
}

// 接受两个chanel，合并后返回一个channel
func mergeData(in1, in2 <-chan int32) <-chan int32 {
	out := make(chan int32, 1024)
	go func() {
		v1, ok1 := <-in1
		v2, ok2 := <-in2
		for ok1 || ok2 {
			if !ok2 || (ok1 && v1 <= v2) {
				out <- v1
				v1, ok1 = <-in1
			} else {
				out <- v2
				v2, ok2 = <-in2
			}
		}
		close(out)
	}()
	return out
}
