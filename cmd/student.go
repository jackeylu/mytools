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
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// the filename of the dataset
	dataset string
	// finding by the key, may be the student's name or No.
	key string
)

// studentCmd represents the student command
var studentCmd = &cobra.Command{
	Use:   "student",
	Short: "find the class name of given student with name or student no.",
	Long:  `Find the class name of given student by given dataset.`,
	Run: func(cmd *cobra.Command, args []string) {
		findStudent(dataset, key)
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
	studentCmd.Flags().StringVarP(&dataset, "dataset", "d", "", "the dataset file")
	studentCmd.Flags().StringVarP(&key, "key", "k", "", "the key of the student")
}

func findStudent(dataset, key string) {
	if dataset == "" {
		panic("dataset is empty")
	}
	if key == "" {
		panic("key is empty")
	}
	if !fileExists(dataset) {
		panic("dataset file not found")
	}
	findStudentByKey(dataset, key)
}

func findStudentByKey(dataset, key string) {
	file, err := os.Open(dataset)
	if err != nil {
		panic(fmt.Errorf("error opening namelist file:%v", err))
	}

	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		panic(fmt.Errorf("error reading file:%v", err))
	}
	// 将每一行分割成字段
	lines := strings.Split(string(bytes), "\n")
	// 将分割后的字段添加到结果数组中
	for _, line := range lines {
		fields := strings.Split(line, ",")
		if fields == nil || len(fields) < 3 {
			continue
		}
		// 检查 字段 数是否等于3，不等于则报错
		if len(fields) != 3 {
			panic(fmt.Errorf("error reading namelist fields:%v, expected two columns but not matched",
				fields))
		}
		name := strings.TrimSpace(fields[0])
		sno := strings.TrimSpace(fields[1])
		class := strings.TrimSpace(fields[2])
		if name == key || sno == key {
			fmt.Printf("Student Name: %s, Student ID: %s, Student Classroom: %s\n", name, sno, class)
			return
		}
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
