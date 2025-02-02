package logger

import (
	"io"
	"os"
	"path/filepath"
)

func APILoggerOutput() io.Writer {
	rootPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	projectRoot := filepath.Dir(filepath.Dir(rootPath))
	logPath := filepath.Join(projectRoot, "logs")
	if err := os.MkdirAll(logPath, os.ModePerm); err != nil {
		panic(err)
	}

	logFilePath := filepath.Join(logPath, "api.log")
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		panic(err)
	}

	return file
}
