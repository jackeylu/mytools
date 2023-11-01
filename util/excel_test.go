package util

import (
	"fmt"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestWriteExcelFile(t *testing.T) {
	excelFile := "test.xlsx"

	// Test columns size not match data size
	columns2 := []string{"Name", "Age"}
	data2 := [][]string{
		{"John Doe", "25\n36"},
	}
	err := WriteExcelFile(excelFile, columns2, data2)
	if err != nil {
		t.Error(err)
	}

	columns := []string{"Name", "Age", "Email"}
	// Test empty data
	err = WriteExcelFile(excelFile, columns, [][]string{})
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Test valid data
	data := [][]string{
		{"John Doe", "25\n30\nabc", "john@example.com"},
		{"Jane Smith", "30", "jane@example.com"},
	}
	err = WriteExcelFile(excelFile, columns, data)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Test file content
	content := [][]string{}
	ReadExcelFile(excelFile, len(columns), func(s []string) {
		content = append(content, s)
	}, false)

	expectedContent := [][]string{
		{"Name", "Age", "Email"},
		{"John Doe", "25\n30\nabc", "john@example.com"},
		{"Jane Smith", "30", "jane@example.com"},
	}

	if !checkContent(content, expectedContent) {
		t.Errorf("Expected file content to be %v, but got %v", expectedContent, content)
	}
}

func checkContent(content, expectedContent [][]string) bool {
	if len(content) != len(expectedContent) {
		return false
	}

	for i, rowContent := range content {
		if len(rowContent) != len(expectedContent[i]) {
			return false
		}

		for j, cellContent := range rowContent {
			if cellContent != expectedContent[i][j] {
				return false
			}
		}
	}

	return true
}

func TestExcelExample(t *testing.T) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	for idx, row := range [][]interface{}{
		{nil, "Apple", "Orange", "Pear"}, {"Small", 2, 3, 3},
		{"Normal", 5, 2, 4}, {"Large", 6, 7, 8},
	} {
		cell, err := excelize.CoordinatesToCellName(1, idx+1)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetSheetRow("Sheet1", cell, &row)
	}
	if err := f.AddChart("Sheet1", "E1", &excelize.Chart{
		Type: excelize.Col3DClustered,
		Series: []excelize.ChartSeries{
			{
				Name:       "Sheet1!$A$2",
				Categories: "Sheet1!$B$1:$D$1",
				Values:     "Sheet1!$B$2:$D$2",
			},
			{
				Name:       "Sheet1!$A$3",
				Categories: "Sheet1!$B$1:$D$1",
				Values:     "Sheet1!$B$3:$D$3",
			},
			{
				Name:       "Sheet1!$A$4",
				Categories: "Sheet1!$B$1:$D$1",
				Values:     "Sheet1!$B$4:$D$4",
			},
		},
		Title: []excelize.RichTextRun{
			{
				Text: "Fruit 3D Clustered Column Chart",
			},
		},
		Legend: excelize.ChartLegend{
			ShowLegendKey: false,
		},
		PlotArea: excelize.ChartPlotArea{
			ShowBubbleSize:  true,
			ShowCatName:     false,
			ShowLeaderLines: false,
			ShowPercent:     true,
			ShowSerName:     true,
			ShowVal:         true,
		},
	}); err != nil {
		fmt.Println(err)
		return
	}
	// Save spreadsheet by the given path.
	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}
