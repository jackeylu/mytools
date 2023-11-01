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
	"encoding/base64"
	"os"

	"github.com/spf13/cobra"
)

var (
	// web2binFilename string
	filename string
	// web2binOutput string
	output string
)

// web2binCmd represents the web2bin command
var web2binCmd = &cobra.Command{
	Use:   "web2bin",
	Short: "将base64编码的文件内容解码成一个二进制文件",
	Long: `将base64编码的文件内容解码成一个二进制文件

示例:
web2bin -f base64.txt -o out.bin
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// read file content
		bytes, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		// convert base64 to binary
		ans := make([]byte, base64.RawURLEncoding.DecodedLen(len(bytes)))
		n, err := base64.RawURLEncoding.Decode(ans, bytes)
		if err != nil {
			return err
		}
		// write binary to file
		err = os.WriteFile(output, ans[:n], 0644)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(web2binCmd)

	web2binCmd.Flags().StringVarP(&filename, "filename", "f", "base64.txt", "Filename to convert")
	web2binCmd.Flags().StringVarP(&output, "output", "o", "out.bin", "Output filename")
}
