package main

import (
	"bytes"
	"encoding/binary"
	"os"
)

// 将通道in的数据写入文件filename中
func writeData(fileName string, in <-chan int32, flag bool) {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
		return
	}
	defer file.Close()
	for v := range in {
		var buf bytes.Buffer
		if flag {
			//fmt.Printf("%d ", v)
		}
		binary.Write(&buf, binary.LittleEndian, v)
		_, err = file.Write(buf.Bytes())
		if err != nil {
			panic(err)
			return
		}
	}
}
