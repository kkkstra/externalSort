package main

import (
	"fmt"
	"math"
)

// 二路归并
func mergeTmpFiles2(tmpFile []string, threads int) {
	fmt.Printf("result: binary merging...\n")
	ins := make([]<-chan int32, threads)
	for i, fileName := range tmpFile {
		ins[i], _ = loadFile(fileName, 0, -1)
	}
	res := divideData(ins...)
	writeData("result", res, true)
	fmt.Printf("result: merged\n")
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

const (
	MAXKEY int32 = math.MaxInt32
	MINKEY int32 = math.MinInt32
)

// K路归并
func mergeTmpFiles(tmpFile []string, threads int) {
	fmt.Printf("result: merging...\n")
	ins := make([]<-chan int32, threads)
	for i, fileName := range tmpFile {
		ins[i], _ = loadFile(fileName, 0, -1)
	}

	// 初始化各个管道的值
	val := make([]int32, threads+1)
	for i := 0; i < threads; i++ {
		val[i] = input(i, ins, threads)
		//fmt.Printf("%d %d\n", i, val[i])
	}

	loserTree := createLoserTree(val, threads) // 创建败者树
	//for i := 0; i < threads; i++ {
	//	fmt.Printf("%d %d\n", i, loserTree[i])
	//}

	out := make(chan int32, 1024)
	go func() {
		for val[loserTree[0]] != MAXKEY {
			//fmt.Printf("\n!%d\n", val[loserTree[0]])
			out <- val[loserTree[0]]
			val[loserTree[0]] = input(loserTree[0], ins, threads)
			modifyLoserTree(loserTree, val, loserTree[0], threads)
		}
		close(out)
	}()
	writeData("result", out, true)
	fmt.Printf("result: merged\n")
}

// 从对应管道中读入数据
func input(id int, in []<-chan int32, threads int) int32 {
	if id == threads {
		return MAXKEY
	}
	val, ok := <-in[id]
	if ok {
		return val
	} else {
		return MAXKEY
	}
}

// 创建败者树
func createLoserTree(val []int32, threads int) []int {
	loserTree := make([]int, threads+1)
	val[threads] = MINKEY
	for i := 0; i < threads; i++ {
		loserTree[i] = threads
	}
	for i := threads - 1; i >= 0; i-- {
		modifyLoserTree(loserTree, val, i, threads)
	}
	return loserTree
}

func modifyLoserTree(loserTree []int, val []int32, s int, threads int) {
	t := (s + threads) / 2
	for t > 0 {
		if val[s] > val[loserTree[t]] {
			s, loserTree[t] = loserTree[t], s
		}
		t /= 2
	}
	loserTree[0] = s
}
