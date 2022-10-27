package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
)

// 读入数据文件
func readData(fileName string, threads int) []string {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	fileSize := stat.Size()                // 获取文件大小
	chunkSize := fileSize / int64(threads) // 计算分块大小
	chunkSize = (chunkSize - 1) / 4 * 4    // 确保为数据类型大小的整数倍
	//fmt.Println(fileSize, chunkSize)

	var offset int64
	tmpFile := make([]string, threads)
	// 将每一段文件分别进行排序
	for i := 0; i < threads; i++ {
		if i == threads-1 {
			chunkSize = fileSize - offset
		}
		//fmt.Println(offset)
		in, count := loadFile(fileName, offset, chunkSize)
		offset += chunkSize
		tmpFile[i] = internalSort(i, in, count)
	}
	return tmpFile
}

// 将文件装载到内存中，并通过通道发送
func loadFile(fileName string, offset int64, chunkSize int64) (chan int32, int) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	if offset > 0 {
		_, err := file.Seek(offset, 0)
		if err != nil {
			panic(err)
		}
	}
	if chunkSize == -1 {
		stat, err := file.Stat()
		if err != nil {
			panic(err)
		}
		chunkSize = stat.Size()
	}
	out := make(chan int32)
	var count int
	var totSize int64
	go func() {
		var m int32
		for {
			data := readNextBytes(file, 4)
			buf := bytes.NewBuffer(data)
			err = binary.Read(buf, binary.LittleEndian, &m)
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
				return
			}
			out <- m
			count++
			totSize += int64(binary.Size(data))
			if totSize >= chunkSize {
				break
			}
		}
		close(out)
	}()
	return out, count
}

func readNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	if err != nil {
		panic(err)
		return nil
	}
	return bytes
}

func readTmpFile(fileName string) <-chan int32 {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
		return nil
	}
	out := make(chan int32, 1024)
	go func() {
		var m int32
		for {
			data := readNextBytes(file, 4)
			buf := bytes.NewBuffer(data)
			err = binary.Read(buf, binary.LittleEndian, &m)
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
				return
			}
			out <- m
		}
		close(out)
	}()
	return out
}
