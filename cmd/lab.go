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
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jackeylu/mytools/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	workingDir string
	labName    string
	// like php-2023-class-1
	coursename string
	debug      bool
)

type CourseStudent struct {
	Name string
	Sno  string
}

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
		if coursename == "" {
			panic(fmt.Errorf("coursename is empty.Should be like php-2023-class-1"))
		} else {
			excelFile = viper.GetString("lab.class." + coursename)
		}
		if debug {
			fmt.Fprintln(os.Stderr, "workingDir:", workingDir, "labName:", labName, "csvfile:", excelFile)
		}
		students := readNameList(excelFile)
		// 文件名模式: `.*\.(doc|docx)` 表示匹配所有以 .doc 或 .docx 结尾的文件
		fileNamePattern := `.*\.(doc|docx)`
		traverseFiles(workingDir, labName, students, fileNamePattern)
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
	labCmd.Flags().StringVarP(&coursename, "coursename", "c", "",
		"The coursename , like php-2023-class-1")
	labCmd.Flags().StringVarP(&labName, "labName", "l", "", "the lab name in filename")
	labCmd.Flags().BoolVarP(&debug, "debug", "D", false, "show debug result or only the result")
}

func readNameList(csvfile string) []CourseStudent {
	lines := make([]CourseStudent, 0)
	util.ReadExcelFile(csvfile, 2, func(line []string) {
		if len(line) != 2 {
			panic(fmt.Errorf("error reading namelist fields:%v, expected 2 columns but not matched",
				line))
		}
		lines = append(lines, CourseStudent{
			Name: line[0],
			Sno:  line[1],
		})
	})
	return lines
}

func traverseFiles(folderPath, labName string, students []CourseStudent, fileNamePattern string) {
	illegalFileNames := make([]string, 0)
	notFounds := make([]string, 0)
	result := make([]string, len(students))
	for i := 0; i < len(result); i++ {
		// Not submitted at default
		result[i] = ""
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
			idx := findRecord(students, name, sno)
			if idx == -1 {
				// 如果不存在，将该文件名添加到未匹配数组中
				notFounds = append(notFounds, fileName)
			} else {
				// 存在，标记为已提交
				result[idx] = "已提交"
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}
	if debug {
		fmt.Printf("%s,%s,%s\n", "Name", "Sno", labName)
	} else {
		fmt.Println(labName)
	}
	for i := 0; i < len(result); i++ {
		if debug {
			fmt.Printf("%s,%s,%s\n", students[i].Name, students[i].Sno, result[i])
		} else {
			fmt.Println(result[i])
		}
	}
	fmt.Println("---------")
	if len(illegalFileNames) > 0 {
		fmt.Fprintln(os.Stderr, "Illegal file name:")
		for _, v := range illegalFileNames {
			fmt.Fprintln(os.Stderr, v)
		}
		fmt.Println("---------")
	}
	if len(notFounds) > 0 {
		fmt.Fprintln(os.Stderr, "Not found:")
		for _, v := range notFounds {
			fmt.Fprintln(os.Stderr, v)
		}
		fmt.Println("---------")
	}
}

// Return the index of first found record, else return -1
func findRecord(students []CourseStudent, name, sno string) int {
	for i, v := range students {
		if v.Name == name && v.Sno == sno {
			return i
		}
	}

	return -1
}
