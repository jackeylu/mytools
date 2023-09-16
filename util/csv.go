package util

import (
	"encoding/csv"
	"fmt"
	"os"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func ReadCsvFile(csvfile string, columnsSize int, f func([]string)) {
	if !fileExists(csvfile) {
		panic(fmt.Sprintf("file [%s] not found", csvfile))
	}
	file, err := os.Open(csvfile)
	if err != nil {
		panic(fmt.Errorf("error opening namelist file:%v", err))
	}

	defer file.Close()
	// 创建CSV阅读器
	reader := csv.NewReader(file)

	// 读取CSV文件的第一行，作为列名
	columns, err := reader.Read()
	if err != nil {
		panic(fmt.Errorf("Error reading header: %v", err))
	}
	if columns == nil || len(columns) != columnsSize {
		panic(fmt.Errorf("error reading header, expected %d columns but not matched", columnsSize))
	}
	// 读取文件中的每一行
	lines, err := reader.ReadAll()
	if err != nil {
		panic(fmt.Errorf("error reading namelist file:%v", err))
	}
	for _, line := range lines {
		f(line)
	}
}
