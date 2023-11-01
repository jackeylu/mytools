/*
Copyright © 2023 lvlin@whu.edu.cn

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/spf13/cobra"
)

var (
	start       int64
	end         int64
	milliSecond int64
	// lottoCmd represents the lotto command
	lottoCmd = &cobra.Command{
		Use:   "lotto",
		Short: "乐透小游戏",
		Long: `
本程序是一款简单的乐透小游戏，用户在命令行中输入起始数字和结束数字，以及休眠间隔，
程序会在休眠间隔后随机生成一个整数，反复生成新的随机数，直到用户敲击Enter按键结束。


示例:
  mytools lotto -start 1 -end 100 -milliSecond 1000

Flags:
  -h, --help   help for lotto

Args:
  [start]     起始数字，必须大于0
  [end]       结束数字，必须大于start
  [milliSecond] 休眠时间，单位毫秒，默认100ms

`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if start >= end || start < 0 {
				return fmt.Errorf("start number must be less than end number and must be greater or equal then zero, but start = %d, end = %d",
					start, end)
			}
			if milliSecond < 0 {
				fmt.Println("negativte value for milliSecond, use default value 100ms.")
				milliSecond = 100
			}
			work()
			return nil
		},
	}
)

func work() {
	fmt.Println("Press Enter to break.")
	ch := make(chan int, 1)
	// 如果用户有输入
	go func() {
		if _, err := fmt.Scanln(); err == nil {
			ch <- 1
		}
	}()
	defer close(ch)
	for {
		num := rand.Int63n(end-start) + start
		fmt.Printf("%d", num)
		// 休眠0.5秒钟
		time.Sleep(time.Millisecond * time.Duration(milliSecond))

		select {
		case <-ch:
			fmt.Printf("You have choosed %d\n", num)
			return
		default:
			fmt.Printf("\r ")
			for num = num % 10; num > 0; num /= 10 {
				fmt.Printf(" ")
			}
			fmt.Printf(" ")
			fmt.Printf("\r")
		}
	}

}
func init() {
	rootCmd.AddCommand(lottoCmd)

	lottoCmd.Flags().Int64VarP(&start, "start", "s", 1, "The starting number.")
	lottoCmd.Flags().Int64VarP(&end, "end", "e", 100, "The ending number.")
	lottoCmd.Flags().Int64VarP(&milliSecond, "ms", "m", 100, "The sleeping interval.")
}
