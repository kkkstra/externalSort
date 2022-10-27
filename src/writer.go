package main

import (
	"bytes"
	"encoding/binary"
	"os"
)

// 将通道in的数据写入文件filename中
func writeData(fileName string, in <-chan int32) {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
		return
	}
	defer file.Close()
	for v := range in {
		//fmt.Println(v)
		var buf bytes.Buffer
		binary.Write(&buf, binary.LittleEndian, v)
		_, err = file.Write(buf.Bytes())
		if err != nil {
			panic(err)
			return
		}
	}
}
