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

	"github.com/jackeylu/mytools/util"
	"github.com/spf13/cobra"
)

var (
	dir     string
	str     string
	confirm bool
)

// rmStrCmd represents the rmStr command
var rmStrCmd = &cobra.Command{
	Use:   "rmStr",
	Short: "删除指定目录中所有文件中的指定字符",
	Long:  `删除指定目录中所有文件中的指定字符.`,
	Run: func(cmd *cobra.Command, args []string) {
		if dir == "" {
			fmt.Println("请指定目录")
			return
		}
		if str == "" {
			fmt.Println("请指定字符")
			return
		}
		cleanFilenamesInDirWithStr(dir, str)
		fmt.Println("删除完成")
	},
}

func cleanFilenamesInDirWithStr(dir string, str string) {
	fmt.Println("cleaning filenames in dir:", dir)
	fmt.Println("cleaning filenames with str:", str)
	// 遍历 dir 目录，获取每个文件
	filenames, err := util.GetAllFileNames(dir, false)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	for _, relativeFilename := range filenames {
		// 去除文件名中的 str 字符
		newFilename := strings.Replace(relativeFilename, str, "", 1)
		// 重命名文件
		src := filepath.Join(dir, relativeFilename)
		dst := filepath.Join(dir, newFilename)
		if !confirm {
			if choice, err := util.CheckInput(fmt.Sprintf("mv %s %s ? Y/N", src, dst), "Y", "y", "N", "n"); err != nil {
				fmt.Println("error:", err)
				continue
			} else if choice == "N" || choice == "n" {
				continue
			}
		}
		os.Rename(src, dst)
	}
}

func init() {
	rootCmd.AddCommand(rmStrCmd)

	rmStrCmd.Flags().StringVarP(&dir, "dir", "d", "./", "指定目录")
	rmStrCmd.Flags().StringVarP(&str, "str", "s", "", "指定字符")
	rmStrCmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "是否确认")
}
