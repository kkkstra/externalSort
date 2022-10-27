package main

import (
	"flag"
)

//type Data struct {
//	Num int32
//}

var (
	fileName = flag.String("f", "data", "the filename of the data file")
	threads  = flag.Int("t", 16, "the number of the threads")
	count    = flag.Int("n", 10000000, "the number of the numbers to be generated")
	genData  = flag.Bool("g", false, "whether generate random data")
)

func main() {
	flag.Parse()
	if *genData {
		generateData(*fileName, *count) // 随机生成数据
	}
	tmpFile := readData(*fileName, *threads) // 分块进行排序
	mergeData2(tmpFile, *threads)            // 进行多路合并
}
