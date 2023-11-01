package util

import (
	"fmt"
	"log"
	"os"
)

// LoggerWriter 是一个实现了 io.Writer 接口的类型，它能够同时输出到终端和日志文件
type LoggerWriter struct {
	stdout   *log.Logger
	logfile  *os.File
	filename string
}

// NewLoggerWriter 创建并返回一个 LoggerWriter 对象
func NewLoggerWriter(filename string) *LoggerWriter {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	return &LoggerWriter{
		stdout:   log.New(os.Stdout, "", 0),
		logfile:  f,
		filename: filename,
	}
}

// Write 方法实现 io.Writer 接口，用于同时输出到终端和日志文件
func (lw *LoggerWriter) Write(p []byte) (n int, err error) {
	lw.stdout.Print(string(p))
	return lw.logfile.Write(p)
}

func (lw *LoggerWriter) Close() {
	lw.logfile.Close()
	fmt.Printf("日志文件'%s'已关闭\n", lw.filename)
}
