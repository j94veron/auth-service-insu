package logger

import (
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

var Logger *zap.Logger

func init() {
	// Get the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		panic("Failed to get current directory: " + err.Error())
	}
	// Configure your logger as needed
	cfg := zap.NewProductionConfig()

	// Add standard and error output
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	// Add output to an information log file
	infoLogPath := filepath.Join(currentDir, "logs", "auth-insu-info.log")
	infoLogFile, err := os.OpenFile(infoLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Failed to open info.log file: " + err.Error())
	}
	defer infoLogFile.Close()
	cfg.OutputPaths = append(cfg.OutputPaths, infoLogPath)

	// Add output to an error log file
	errorLogPath := filepath.Join(currentDir, "logs", "auth-insu-error.log")
	errorLogFile, err := os.OpenFile(errorLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Failed to open error.log file: " + err.Error())
	}
	defer errorLogFile.Close()
	cfg.ErrorOutputPaths = append(cfg.ErrorOutputPaths, errorLogPath)

	l, err := cfg.Build()
	if err != nil {
		panic("Failed to create logger: " + err.Error())
	}
	Logger = l
}
