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
	"unicode"

	"github.com/jackeylu/mytools/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	workingDir string
	labsName   []string
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
	Use:   "lab",
	Short: `将一个或多个实验报告文件夹进行统计，汇总每次实验的提交情况.`,
	Long: `This program will check a given directory with given namelist, and generated the checked result. For example:

The namelist is a csv type file with 'name' and 'no' columns.
The reports in the given directory are in the format of '$name-$no-$lab.doc' or '$name-$no-$lab.docx'.
The generated result includes the submmited flag for each student and those file with illegal filename format.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(labsName) == 0 {
			labsName = listSubDirectories(workingDir)
		}
		var excelFile string
		if coursename == "" {
			panic(fmt.Errorf("coursename is empty.Should be like php-2023-class-1"))
		} else {
			excelFile = viper.GetString("lab.class." + coursename)
		}
		if debug {
			fmt.Fprintln(os.Stderr, "workingDir:", workingDir, "labName:", labsName, "csvfile:", excelFile)
		}
		students := readNameList(excelFile)
		// 文件名模式: `.*\.(doc|docx)` 表示匹配所有以 .doc 或 .docx 结尾的文件
		fileNamePattern := `.*\.(doc|docx|zip|rar)`
		traverseFiles(workingDir, labsName, students, fileNamePattern)
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
	labCmd.Flags().StringSliceVarP(&labsName, "labName", "l", []string{}, "the labs' names in filename, split with comma.")
	labCmd.Flags().BoolVarP(&debug, "debug", "D", false, "show debug result or only the result")
}

// List all the sub-directories in the given path
func listSubDirectories(path string) []string {
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	var subDirs []string
	for _, f := range files {
		if f.IsDir() {
			subDirs = append(subDirs, f.Name())
		}
	}
	return subDirs
}

func readNameList(excelFile string) []CourseStudent {
	lines := make([]CourseStudent, 0)
	util.ReadExcelFile(excelFile, 2, func(line []string) {
		if len(line) != 2 {
			panic(fmt.Errorf("error reading namelist fields:%v, expected 2 columns but got %d",
				line, len(line)))
		}
		lines = append(lines, CourseStudent{
			Name: line[0],
			Sno:  line[1],
		})
	}, true)
	return lines
}

func traverseFiles(folderPath string, labsName []string, students []CourseStudent, fileNamePattern string) {
	// Not submitted at default
	illegalFileNames, notFounds, result, found := initResultSet(labsName, students)
	for j, labName := range labsName {
		root := filepath.Join(folderPath, labName)
		// 如果不存在，将该文件名添加到未匹配数组中
		// 存在，标记为已提交
		err := processOneLab(root, fileNamePattern, illegalFileNames, j, labName, students, notFounds, result, found)

		if err != nil {
			fmt.Println("Error:", err)
		}
	}

	handleResult(found, labsName, result, students, illegalFileNames, notFounds)
}

func processOneLab(labDir string,
	fileNamePattern string,
	illegalFileNames [][]string,
	labIndex int,
	labName string,
	students []CourseStudent,
	notFounds [][]string,
	result [][]string,
	found []int) error {
	err := filepath.Walk(labDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(fmt.Errorf("prevent panic by handling failure accessing a path %q: %v", path, err))
		}

		fileName := filepath.Base(path)
		if match, _ := regexp.MatchString(fileNamePattern, fileName); match {
			name, sno, experiment, shouldReturn := extractFilename(fileName, illegalFileNames, labIndex)
			if shouldReturn {
				return nil
			}
			experiment = strings.Split(experiment, ".")[0]
			if experiment != labName {
				illegalFileNames[labIndex] = append(illegalFileNames[labIndex], fileName)
				return nil
			}
			idx := findRecord(students, name, sno)
			if idx == -1 {

				notFounds[labIndex] = append(notFounds[labIndex], fileName)
			} else {

				if result[idx][labIndex] != "已提交" {
					result[idx][labIndex] = "已提交"
					found[labIndex]++
				} else {
					fmt.Fprintf(os.Stderr, "Duplicate file name: %s\n", fileName)
				}
			}
		}

		return nil
	})
	return err
}

func extractFilename(fileName string, illegalFileNames [][]string, labIndex int) (name string,
	sno string, experiment string, shouldReturn bool) {
	name, sno, experiment, shouldReturn = "", "", "", false
	// 1. trim space in the leading and trailing position
	// 2. replace left space and '_' with '-'
	fields := strings.Split(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.TrimSpace(fileName), "_", "-"), " ", "-"), "-")
	if len(fields) < 3 {
		illegalFileNames[labIndex] = append(illegalFileNames[labIndex], fileName)
		shouldReturn = true
		return
	}
	name, sno, experiment = fields[0], fields[1], strings.Join(fields[2:], "-")
	isDigit := true
	for _, r := range name {
		if !unicode.IsDigit(r) {
			continue
		}
		isDigit = false
		break
	}
	if !isDigit {
		name, sno = sno, name
	}
	return
}

func handleResult(found []int,
	labsName []string,
	result [][]string,
	students []CourseStudent,
	illegalFileNames [][]string,
	notFounds [][]string) {
	fmt.Println("Found:")
	for i, v := range found {
		fmt.Printf("%d", v)
		if i < len(found)-1 {
			fmt.Print(",")
		}
	}
	fmt.Println()
	if debug {
		fmt.Printf("%s,%s,%s\n", "Name", "Sno", strings.Join(labsName, ","))
	} else {
		fmt.Println(strings.Join(labsName, ","))
	}
	for i := 0; i < len(result); i++ {
		if debug {
			fmt.Printf("%s,%s,%s\n", students[i].Name, students[i].Sno, strings.Join(result[i], ","))
		} else {
			fmt.Println(strings.Join(result[i], ","))
		}
	}
	fmt.Println("---------")
	// print files does not match the filepattern
	if len(illegalFileNames) > 0 {
		fmt.Fprintln(os.Stderr, "Illegal file name:")
		for i, v := range illegalFileNames {
			if len(v) > 0 {
				fmt.Fprintln(os.Stderr, labsName[i])
				for _, v2 := range v {
					fmt.Fprintln(os.Stderr, v2)
				}
			}
		}
		fmt.Fprintln(os.Stderr, "---------")
	}
	// print files with name or no missmatched.
	if len(notFounds) > 0 {
		fmt.Fprintln(os.Stderr, "Not found:")
		for i, v := range notFounds {
			if len(v) > 0 {
				fmt.Fprintln(os.Stderr, labsName[i])
				for _, v2 := range v {
					fmt.Fprintln(os.Stderr, v2)
				}
			}
		}
		fmt.Fprintln(os.Stderr, "---------")
	}
}

func initResultSet(labsName []string, students []CourseStudent) ([][]string, [][]string, [][]string, []int) {
	illegalFileNames := make([][]string, len(labsName))
	notFounds := make([][]string, len(labsName))
	result := make([][]string, len(students))
	found := make([]int, len(labsName))

	for i := 0; i < len(students); i++ {
		result[i] = make([]string, len(labsName))
		for j := 0; j < len(labsName); j++ {

			result[i][j] = ""
		}
	}
	return illegalFileNames, notFounds, result, found
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
