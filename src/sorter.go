package main

import (
	"fmt"
	"sort"
	"strconv"
)

func internalSort(id int, in chan int32, count int) string {
	data := make([]int32, count)
	for num := range in {
		data = append(data, num)
	}
	sort.Slice(data, func(i, j int) bool { return data[i] < data[j] })

	out := make(chan int32)
	go func() {
		for _, num := range data {
			out <- num
		}
		close(out)
	}()
	
	fileName := "data.tmp" + strconv.Itoa(id)
	fmt.Printf("%s: created. \n", fileName)
	writeData(fileName, out)
	fmt.Printf("%s: sorted. \n", fileName)
	return fileName
}
