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
	"strings"

	"github.com/spf13/cobra"
)

var (
	// the working directory
	dir string
	// the extension of the filename
	ext string
	// the confirmed flag when renaming
	confirmed bool
)

// rmdupCmd represents the rmdup command
var rmdupCmd = &cobra.Command{
	Use:   "rmdup",
	Short: "Remove duplication text of filename",
	Long:  `将指定文件夹中特定文件类型的文件名进行整理，剔除重复的内容.`,
	Run: func(cmd *cobra.Command, args []string) {
		rmdup()
	},
}

func init() {
	rootCmd.AddCommand(rmdupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rmdupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rmdupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rmdupCmd.Flags().StringVarP(&dir, "dir", "d", "./", "the directory to clean")
	rmdupCmd.Flags().StringVarP(&ext, "ext", "e", "",
		`the extension of the filename, like .zip , default is none for every file, 
		otherwise only for file with specificated extension.`)
	rmdupCmd.Flags().BoolVarP(&confirmed, "confirmed", "c", false, "confirm when renaming")
}

func _clean(filename string) string {
	var ans string
	ans = filename
	fileExt1 := filepath.Ext(filename)
	if fileExt1 == "" {
		return ans
	}
	filename = filename[0 : len(filename)-len(fileExt1)]
	fileExt2 := filepath.Ext(filename)
	if fileExt2 != "" && fileExt1 == fileExt2 {
		ans = filename
	}
	return ans
}

func cleanName(filename string) string {
	fn := _clean(filename)
	for fn != filename {
		filename = fn
		fn = _clean(filename)
	}
	return fn
}

func rmdup() {
	//step1：list all files in the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	// step2: check each file
	for _, file := range files {
		// check if the file is a directory
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		// get the extension of the filename
		fileExt := filepath.Ext(filename)
		if ext != "" && fileExt != "" && fileExt != ext {
			// 不匹配的文件类型不做处理
			continue
		}
		cleanedName := cleanName(filename)
		if cleanedName != "" && cleanedName != filename {
			// 重命名filename为cleanedName
			from, to := filepath.Join(dir, filename), filepath.Join(dir, cleanedName)
			if !confirmed {
				if ignore(from, to) {
					continue
				}
			}

			err := os.Rename(from, to)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("renamed %s to %s\n", filename, cleanedName)
		}
	}

}

func ignore(from, to string) bool {
	for {
		// 获取终端输入的字符，如果是 Y 则进入下一步，
		// 否则提示重新输入并且继续监听终端输入
		var input string
		fmt.Printf("Renaming  %s to %s [Y/n]:", from, to)
		fmt.Scanf("%s", &input)
		input = strings.ToLower(input)
		if input == "y" {
			return true
		} else if input == "n" {
			return false
		} else {
			fmt.Println("Please input Y or N")
		}
	}
}
