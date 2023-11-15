package util

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetAllFileNames(directoryPath string, deepIntoSubDirOrNot bool) ([]string, error) {
	var fileNames []string

	// 获取目录下的所有文件和子目录
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		return nil, err
	}

	// 遍历文件和子目录
	for _, file := range files {
		if deepIntoSubDirOrNot && file.IsDir() {
			// 如果是子目录，递归获取子目录下的文件名
			subDir := filepath.Join(directoryPath, file.Name())
			subDirFileNames, err := GetAllFileNames(subDir, deepIntoSubDirOrNot)
			if err != nil {
				return nil, err
			}
			fileNames = append(fileNames, subDirFileNames...)
		} else {
			// 如果是文件，添加到文件名列表
			fileNames = append(fileNames, file.Name())
		}
	}

	return fileNames, nil
}

func CheckInput(tip string, choices ...string) (string, error) {
	if len(choices) == 0 {
		return "", fmt.Errorf("choices is empty")
	}
	// 提示用户输入
	// fmt.Print("请输入 'confirm' 以执行操作: ")
	fmt.Print(tip)

	// 从标准输入读取用户输入
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("读取输入时发生错误:", err)
		return "", err
	}

	// 去除输入字符串中的空白字符
	input = strings.TrimSpace(input)

	for _, choice := range choices {
		if input == choice {
			return choice, nil
		}
	}
	return "", fmt.Errorf("输入的选项不在可选范围内")
}
