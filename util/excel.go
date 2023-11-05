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

// AppendExcelFileByFunction 覆盖写数据到excel文件中
func AppendExcelFileByFunction(excelFile string, columns []string, f func() [][]string) {
	WriteOrAppendExcelFile(excelFile, columns, f(), true)
}

// WriteExcelFile 覆盖写数据到excel文件中
func WriteExcelFile(excelFile string, headers []string, rows [][]string) error {
	return WriteOrAppendExcelFile(excelFile, headers, rows, false)
}

func WriteOrAppendExcelFile(excelFile string, headers []string, rows [][]string, append bool) error {
	if len(rows) == 0 {
		fmt.Print("No data to handle.")
		return nil
	}
	if len(headers) != len(rows[0]) {
		return errors.New("columns size not match data size")
	}

	file := excelize.NewFile()
	defer file.Close()

	rowNum := -1
	if append && fileExists(excelFile) {
		ReadExcelFile(excelFile, func(i int, row []string) error {
			cell, err := excelize.CoordinatesToCellName(1, i+1)
			if err != nil {
				fmt.Println(err)
				return err
			}
			file.SetSheetRow("Sheet1", cell, &row)
			rowNum = i
			return nil
		}, false)
	}
	if rowNum == -1 {
		// write header
		cell, err := excelize.CoordinatesToCellName(1, 1)
		if err != nil {
			fmt.Println(err)
			return err
		}
		file.SetSheetRow("Sheet1", cell, &headers)
		rowNum += 1
	}
	for idx, row := range rows {
		cell, err := excelize.CoordinatesToCellName(1, rowNum+idx+2)
		if err != nil {
			fmt.Println(err)
			return err
		}
		file.SetSheetRow("Sheet1", cell, &row)
	}

	file.SaveAs(excelFile)
	return nil
}
