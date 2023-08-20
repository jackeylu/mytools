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
	Short: "Covert the web content in base64 format to binary",
	Run: func(cmd *cobra.Command, args []string) {
		// read file content
		bytes, err := os.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		// convert base64 to binary
		ans := make([]byte, base64.RawURLEncoding.DecodedLen(len(bytes)))
		n, err := base64.RawURLEncoding.Decode(ans, bytes)
		if err != nil {
			panic(err)
		}
		// write binary to file
		err = os.WriteFile(output, ans[:n], 0644)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(web2binCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// web2binCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// web2binCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	web2binCmd.Flags().StringVarP(&filename, "filename", "f", "base64.txt", "Filename to convert")
	web2binCmd.Flags().StringVarP(&output, "output", "o", "out.bin", "Output filename")
}
