package logger

import (
	"log"
	"os"
)

// Logger is an interface that defines methods for logging at different levels.
type Logger interface {
	Info(message string)
	Error(message string)
}

// ConsoleLogger is a concrete implementation of the Logger interface that logs to the console.
type ConsoleLogger struct{}

// Info logs an informational message to the console.
func (c *ConsoleLogger) Info(message string) {
	log.Println("INFO:", message)
}

// Error logs an error message to the console.
func (c *ConsoleLogger) Error(message string) {
	log.Println("ERROR:", message)
}

// FileLogger is a concrete implementation of the Logger interface that logs to a file.
type FileLogger struct {
	file *os.File
}

// NewFileLogger creates a new FileLogger that writes logs to the specified file.
func NewFileLogger(filename string) (*FileLogger, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &FileLogger{file: file}, nil
}

// Info logs an informational message to the file.
func (f *FileLogger) Info(message string) {
	log.SetOutput(f.file)
	log.Println("INFO:", message)
}

// Error logs an error message to the file.
func (f *FileLogger) Error(message string) {
	log.SetOutput(f.file)
	log.Println("ERROR:", message)
}

// Close closes the log file.
func (f *FileLogger) Close() error {
	return f.file.Close()
}