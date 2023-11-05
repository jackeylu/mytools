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
	"log"

	"github.com/jackeylu/mytools/util"
	"github.com/spf13/cobra"
)

var savedExcelFile string

// rmdupCmd represents the rmdup command
var rmdupCmd = &cobra.Command{
	Use:   "rmdup",
	Short: "清除下载的邮件信息中的重复数据",
	Long:  `清除下载的邮件信息中的重复数据，并按日期排序.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if emailFile == "" {
			return fmt.Errorf("请指定邮件文件")
		}
		if savedExcelFile == "" {
			return fmt.Errorf("请指定保存excel文件的文件名")
		}
		// 设置日志文件的格式
		log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
		// 创建一个 LoggerWriter 对象
		logger := util.NewLoggerWriter("logfile.txt")
		defer logger.Close()
		// 将日志同时输出到终端和日志文件
		log.SetOutput(logger)

		var emails []EmailInfo
		if err := readAttachmentEmailFromFetchedEmailFile(emailFile, &emails); err != nil {
			log.Println(err)
			return err
		}
		emails = rmdup(emails)

		// save the emails into file
		if err := util.WriteExcelFile(savedExcelFile, ExcelFileHeader(), emailContent(emails)); err != nil {
			log.Println(err)
			return err
		}
		log.Printf("saved in %s", savedExcelFile)
		return nil
	},
}

func init() {
	emailCmd.AddCommand(rmdupCmd)

	rmdupCmd.Flags().StringVarP(&emailFile, "file", "f", "", "the fetched email file by email command")
	rmdupCmd.Flags().StringVarP(&savedExcelFile, "ouput", "o", "", "the file to store the cleared result")
}

func rmdup(emails []EmailInfo) []EmailInfo {
	var ans []EmailInfo
	in := func(email EmailInfo) bool {
		for j := 0; j < len(ans); j++ {
			if ans[j].Equals(email) {
				return true
			}
		}
		return false
	}
	for i := 0; i < len(emails); i++ {
		if !in(emails[i]) {
			ans = append(ans, emails[i])
		}
	}
	return ans
}
