package logger

import (
	"os"
)

type ILogger interface {
	open(fileName string) *os.File
	close()
	writeLine(line string)
}

type Logger struct {
	fd *os.File
}

func (l *Logger) Open(filename string) *os.File {
	fd, err := os.OpenFile(filename, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0644)
	if (err != nil) {
		panic(err)
	}
	return fd
}

func (l *Logger) Close() {
	err := l.fd.Close()
	if (err != nil) {
		panic(err)
	}
}

func (l *Logger) WriteLine(line string) {
	l.fd.WriteString(line);
}

var loggerInstance *Logger = nil;
func NewLogger(fileName string) *Logger {
	if (loggerInstance == nil) {
		loggerInstance = new(Logger)

		// FIXME : Needs refactoring
		loggerInstance.fd = loggerInstance.Open(fileName)
	}
	
	return loggerInstance
}

