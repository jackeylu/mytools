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

	"github.com/atotto/clipboard"
	"github.com/jackeylu/mytools/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// the filename of the excelFile
	excelFile string
	// finding by the key, may be the student's name or No.
	key string
)

type Student struct {
	Name  string
	No    string
	Class string
	Grade string
}

func (s Student) String() string {
	return fmt.Sprintf("Tag: %s-%s, Class: %s, Grade: %s", s.Name, s.No, s.Class, s.Grade)
}

// studentCmd represents the student command
var studentCmd = &cobra.Command{
	Use:   "student",
	Short: "find the class name of given student with name or student no.",
	Long:  `Find the class name of given student by given dataset.`,
	Run: func(cmd *cobra.Command, args []string) {
		findStudent(excelFile, key)
	},
}

func init() {
	rootCmd.AddCommand(studentCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// studentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// studentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	studentCmd.Flags().StringVarP(&excelFile, "dataset", "d", "", "the dataset file")
	studentCmd.Flags().StringVarP(&key, "key", "k", "", "the key of the student")
}

func findStudent(excelFile, key string) {
	// check if the dataset file exists
	if excelFile == "" {
		excelFile = viper.GetString("lab.all-student")
	}
	if excelFile == "" {
		panic("dataset is empty")
	}
	if key == "" {
		panic("key is empty")
	}
	findStudentByKey(excelFile, key)
}

func findStudentByKey(excelFile, key string) {
	fmt.Fprintln(os.Stderr, "Using namelist file:", excelFile)
	// read the csvfile into slice of Student
	lines := make([]Student, 0)
	util.ReadExcelFile(excelFile, 4, func(line []string) {
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
	}, true)
	// find the student by key
	if !findStudentByKeyInSlice(lines, key) {
		fmt.Printf("Can not find any student with keyword: %s\n", key)
		return
	}
}

func findStudentByKeyInSlice(lines []Student, key string) bool {
	var found bool = false
	for _, line := range lines {
		if line.No == key || line.Name == key {
			found = true
			fmt.Printf("student %s found with result: %v\n", line.Name, line)
			err := clipboard.WriteAll(fmt.Sprintf("%s-%s", line.Name, line.No))
			if err != nil {
				fmt.Println("failed to copy the tag to clipboard:", err)
			} else {
				fmt.Println("The tag is copied to clipboard.")
			}
		}
	}
	return found
}
