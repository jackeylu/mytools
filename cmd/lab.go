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
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var (
	workingDir string
	labName    string
	namelist   string
	debug      bool
)

// labCmd represents the lab command
var labCmd = &cobra.Command{
	Use: "lab",
	Short: `Check the submitted reports of one lab in a given directory, 
              return the submitted flags and unknown submmitters.`,
	Long: `This program will check a given directory with given namelist, and generated the checked result. For example:

The namelist is a csv type file with 'name' and 'no' columns.
The reports in the given directory are in the format of '$name-$no-$lab.doc' or '$name-$no-$lab.docx'.
The generated result includes the submmited flag for each student and those file with illegal filename format.`,
	Run: func(cmd *cobra.Command, args []string) {
		if labName == "" {
			panic(fmt.Errorf("labName is empty."))
		}
		nameArray := readNameList()
		// 文件名模式: `.*\.(doc|docx)` 表示匹配所有以 .doc 或 .docx 结尾的文件
		fileNamePattern := `.*\.(doc|docx)`
		traverseFiles(workingDir, labName, nameArray, fileNamePattern)
	},
}

func init() {
	rootCmd.AddCommand(labCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// labCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// labCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	labCmd.Flags().StringVarP(&workingDir, "workingDir", "d", "./",
		"the directory contains reports.")
	labCmd.Flags().StringVarP(&namelist, "namelist", "n", "./namelist.csv",
		"The namelist with name and no columns")
	labCmd.Flags().StringVarP(&labName, "labName", "l", "", "the lab name in filename")
	labCmd.Flags().BoolVarP(&debug, "debug", "D", false, "show debug result or only the result")
}

func readNameList() [][]string {
	file, err := os.Open(namelist)
	if err != nil {
		panic(fmt.Errorf("error opening namelist file:%v", err))
	}

	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		panic(fmt.Errorf("error reading file:%v", err))
	}
	// 将每一行分割成字段
	lines := strings.Split(string(content), "\n")
	// 初始化结果数组
	// ignore the header line
	result := make([][]string, 0, len(lines)-1)

	// 将分割后的字段添加到结果数组中
	for i, line := range lines {
		if i == 0 {
			// ignore the header line
			continue
		}
		fields := strings.Split(line, ",")
		if fields == nil || len(fields) < 2 {
			continue
		}
		// 检查 字段 数是否等于2，不等于则报错
		if len(fields) != 2 {
			panic(fmt.Errorf("error reading namelist fields:%v, expected two columns but not matched",
				fields))
		}
		fields[0] = strings.TrimSpace(fields[0])
		fields[1] = strings.TrimSpace(fields[1])
		result = append(result, fields)
	}
	return result
}

func traverseFiles(folderPath, labName string, dataset [][]string, fileNamePattern string) {
	illegalFileNames := make([]string, 0)
	notFounds := make([]string, 0)
	result := make([]string, len(dataset))
	for i := 0; i < len(result); i++ {
		// Not submitted at default
		result[i] = "N"
	}
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}

		fileName := filepath.Base(path)
		if match, _ := regexp.MatchString(fileNamePattern, fileName); match {
			fields := strings.Split(fileName, "-")
			if len(fields) < 3 {
				illegalFileNames = append(illegalFileNames, fileName)
				return nil
			}
			name, sno, experiment := fields[0], fields[1], strings.Join(fields[2:], "-")
			experiment = strings.Split(experiment, ".")[0]
			if experiment != labName {
				illegalFileNames = append(illegalFileNames, fileName)
				return nil
			}
			idx := findRecord(dataset, name, sno)
			if idx == -1 {
				// 如果不存在，将该文件名添加到未匹配数组中
				notFounds = append(notFounds, fileName)
			} else {
				// 存在，标记为已提交
				result[idx] = "Y"
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}
	if debug {
		fmt.Println(strings.Join(append(dataset[0], labName), ","))
	} else {
		fmt.Println(labName)
	}
	for i := 0; i < len(result); i++ {
		if debug {
			fmt.Printf("%s,%s,%s\n", dataset[i][0], dataset[i][1], result[i])
		} else {
			fmt.Println(result[i])
		}
	}
	if len(illegalFileNames) > 0 {
		fmt.Fprintln(os.Stderr, "Illegal file name:")
		for _, v := range illegalFileNames {
			fmt.Fprintln(os.Stderr, v)
		}
	}
	if len(notFounds) > 0 {
		fmt.Fprintln(os.Stderr, "Not found:")
		for _, v := range notFounds {
			fmt.Fprintln(os.Stderr, v)
		}
	}
}

// Return the index of first found record, else return -1
func findRecord(result [][]string, name, sno string) int {
	for i, v := range result {
		if v[0] == name && v[1] == sno {
			return i
		}
	}

	return -1
}
