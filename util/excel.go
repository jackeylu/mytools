package util

import (
	"errors"
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

func ReadExcelFile(excelFile string, f func(int, []string) error, ignoreHeader bool) error {
	if !fileExists(excelFile) {
		panic(fmt.Sprintf("file [%s] not found", excelFile))
	}
	file, err := excelize.OpenFile(excelFile)
	if err != nil {
		return fmt.Errorf("error opening namelist file:%v", err)
	}

	defer file.Close()

	// 获取第一张表
	rows, err := file.GetRows(file.GetSheetName(0))
	if err != nil {
		return fmt.Errorf("error getting rows from first sheet:%v", err)
	}
	if ignoreHeader {
		rows = rows[1:]
	}
	for i, row := range rows {
		if err = f(i, row); err != nil {
			return fmt.Errorf("error on handle %d row [%v] with %v", i, row, err)
		}
	}
	return nil
}

func WriteExcelFileByFunction(excelFile string, columns []string, f func() [][]string) {
	WriteExcelFile(excelFile, columns, f())
}

func WriteExcelFile(excelFile string, columns []string, data [][]string) error {
	if len(data) == 0 {
		fmt.Print("No data to handle.")
		return nil
	}
	if len(columns) != len(data[0]) {
		return errors.New("columns size not match data size")
	}

	file := excelize.NewFile()
	defer file.Close()

	// write header
	cell, err := excelize.CoordinatesToCellName(1, 1)
	if err != nil {
		fmt.Println(err)
		return err
	}
	file.SetSheetRow("Sheet1", cell, &columns)

	for idx, row := range data {
		cell, err := excelize.CoordinatesToCellName(1, idx+2)
		if err != nil {
			fmt.Println(err)
			return err
		}
		file.SetSheetRow("Sheet1", cell, &row)
	}

	file.SaveAs(excelFile)
	return nil
}
