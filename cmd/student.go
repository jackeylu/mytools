/*
Copyright © 2023 Lyu Lin <lvlin@whu.edu.cn>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/jackeylu/mytools/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// the filename of the excelFile
	excelFile string
	// finding by the keys, may be the student's name or No.
	keys []string
	// showing the found result with name first or NO. first?
	reverse bool
)

type Student struct {
	Name  string
	No    string
	Class string
	Grade string
}

func (s Student) String() string {
	return fmt.Sprintf("Name: %s, NO: %s, Class: %s, Grade: %s", s.Name, s.No, s.Class, s.Grade)
}

// studentCmd represents the student command
var studentCmd = &cobra.Command{
	Use:   "student",
	Short: "根据学生的完整姓名或学号进行查询.",
	Long:  `在指定的excel文件中查找学生信息，如果找到将会输出学生信息，包括姓名、学号、班级、年级.`,
	Run: func(cmd *cobra.Command, args []string) {
		findStudent(excelFile, keys)
	},
}

func init() {
	rootCmd.AddCommand(studentCmd)

	studentCmd.Flags().StringVarP(&excelFile, "dataset", "d", "", "the dataset file")
	studentCmd.Flags().StringSliceVarP(&keys, "keys", "k", []string{}, "the key text of students")
	studentCmd.Flags().BoolVarP(&reverse, "reverse", "r", false, "student no first or name first?")
}

func findStudent(excelFile string, keys []string) {
	// check if the dataset file exists
	if excelFile == "" {
		excelFile = viper.GetString("lab.all-student")
	}
	if excelFile == "" {
		panic("dataset is empty")
	}
	if len(keys) == 0 {
		panic("keys is empty")
	}
	fmt.Println(keys)
	findStudentByKeys(excelFile, keys)
}

func findStudentByKeys(excelFile string, keys []string) {
	fmt.Fprintln(os.Stderr, "Using namelist file:", excelFile)
	// read the csvfile into slice of Student
	lines := make([]Student, 0)
	util.ReadExcelFile(excelFile, func(_ int, line []string) error {
		if len(line) != 4 {
			panic(fmt.Errorf("error reading namelist fields:%v, expected four columns but not matched",
				line))
		}
		lines = append(lines, Student{
			Name:  line[0],
			No:    line[1],
			Class: line[2],
			Grade: line[3],
		})
		return nil
	}, true)
	// find the student by key
	if !findStudentByKeyInSlice(lines, strings.TrimSpace(keys[len(keys)-1])) {
		fmt.Printf("Can not find any student with keyword: %s\n", keys[len(keys)-1])
	}
}

func findStudentByKeyInSlice(lines []Student, key string) bool {
	var found bool = false
	for _, line := range lines {
		if line.No == key || line.Name == key {
			found = true
			fmt.Printf("student %s found with result: %v\n", line.Name, line)
			var err error
			var tag string
			if !reverse {
				tag = fmt.Sprintf("%s-%s", line.Name, line.No)
			} else {
				tag = fmt.Sprintf("%s-%s", line.No, line.Name)
			}
			err = clipboard.WriteAll(tag)
			if err != nil {
				fmt.Println("failed to copy the tag to clipboard:", err)
			} else {
				fmt.Printf("The tag %s is copied to clipboard.\n", tag)
			}
		}
	}
	return found
}
