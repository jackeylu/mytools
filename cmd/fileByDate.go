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
	"io/fs"
	"os"

	"github.com/spf13/cobra"
)

var baseDir string

// fileByDateCmd represents the fileByDate command
var fileByDateCmd = &cobra.Command{
	Use:   "fileByDate",
	Short: "将给定目录下的所有文件按修改日期进行分组，放到以日期命名的文件夹中.",
	Long: `
	将给定目录下的所有文件按修改日期进行分组，放到以日期命名的文件夹中.
	
使用方法：

	mytools fileByDate -b <directory>
说明：

	-b  指定要扫描的文件夹路径，默认为./images
示例：
	mytools fileByDate -b ./images
	将会将当前目录下的images文件夹下的所有文件按修改日期进行分组，放到以日期命名的文件夹中。	
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		//step1：扫描baseDir目录下的所有文件
		entries, err := os.ReadDir(baseDir)
		if err != nil {
			return err
		}
		// step2：按日期进行分组
		manageFilesByDate(entries)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fileByDateCmd)

	fileByDateCmd.Flags().StringVarP(&baseDir, "baseDir", "b", "./images", "Base directory to search for files")
}

func manageFilesByDate(entries []fs.DirEntry) {
	for i := 0; i < len(entries); i++ {
		entry := entries[i]
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			fmt.Println("read file info failed, err:", err)
			return
		}
		time := info.ModTime()
		date := time.Format("2006-01-02")
		fmt.Printf("mv %s/%s %s/%s\n", baseDir, entry.Name(), baseDir, date)
		os.Mkdir(baseDir+"/"+date, 0777)
		os.Rename(baseDir+"/"+entry.Name(), baseDir+"/"+date+"/"+entry.Name())
	}
}
