package util

import (
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func ReadExcelFile(excelFile string, columnsSize int, f func([]string)) {
	if !fileExists(excelFile) {
		panic(fmt.Sprintf("file [%s] not found", excelFile))
	}
	file, err := excelize.OpenFile(excelFile)
	if err != nil {
		panic(fmt.Errorf("error opening namelist file:%v", err))
	}

	defer file.Close()

	// 获取第一张表
	rows, err := file.GetRows(file.GetSheetName(0))
	if err != nil {
		panic(fmt.Errorf("error getting rows from first sheet:%v", err))
	}

	for _, row := range rows {
		f(row)
	}
}
