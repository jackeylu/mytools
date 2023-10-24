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
	Short: "A password generator",
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// passwordCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// passwordCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	passwordCmd.Flags().Int8VarP(&specialCharacters, "special", "s", 2, "The number of special characters")
	passwordCmd.Flags().Int8VarP(&totalCharacters, "total", "t", 10, "The total number of characters")
}
