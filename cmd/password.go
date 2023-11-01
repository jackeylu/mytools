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
	"errors"
	"fmt"
	"math/rand"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

const CHARS = "abcdefghjkmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVW"
const NUMBERS = "0123456789"
const SPECIAL = "!@#$%^&*()_+-=[]{}|;':,./<>?"

var (
	specialCharacters int8
	totalCharacters   int8
)

// passwordCmd represents the password command
var passwordCmd = &cobra.Command{
	Use:   "password",
	Short: "密码生成器",
	Long: `密码生成器，生成随机密码，并将密码复制到剪切板。

用法：
	mytools password

参数：
	-s, --special 特殊字符个数，默认值为2
	-t, --total 总字符个数，默认值为10

示例：
	$ mytools password -s 3 -t 8
	Generated Password : B8<]mG5.
	Password copied to clipboard
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if specialCharacters < 0 || totalCharacters < 0 {
			return errors.New("number of characters must be greate or equal than zero")
		}
		if specialCharacters > totalCharacters {
			return errors.New("number of special characters must be less or equal than total characters")
		}
		password := generatePassword(specialCharacters, totalCharacters)
		clipboard.WriteAll(password)
		fmt.Printf("Generated Password : %s\n", password)
		fmt.Println("Password copied to clipboard")
		return nil
	},
}

func generatePassword(specialCharacters, totalCharacters int8) string {
	password := make([]byte, totalCharacters)
	if specialCharacters > 0 {
		for i := 0; i < int(specialCharacters); i++ {
			password[i] = SPECIAL[rand.Intn(len(SPECIAL))]
		}
	}

	var scope = CHARS + NUMBERS
	for i := specialCharacters; i < totalCharacters; i++ {
		password[i] = scope[rand.Intn(len(scope))]
	}

	// shuffle the password
	rand.Shuffle(len(password), func(i, j int) {
		password[i], password[j] = password[j], password[i]
	})
	return string(password)
}

func init() {
	rootCmd.AddCommand(passwordCmd)

	passwordCmd.Flags().Int8VarP(&specialCharacters, "special", "s", 2, "The number of special characters")
	passwordCmd.Flags().Int8VarP(&totalCharacters, "total", "t", 10, "The total number of characters")
}
