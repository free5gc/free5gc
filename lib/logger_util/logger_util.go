package logger_util

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

type FileHook struct {
	file      *os.File
	flag      int
	chmod     os.FileMode
	formatter *logrus.TextFormatter
}

func NewFileHook(file string, flag int, chmod os.FileMode) (*FileHook, error) {
	plainFormatter := &logrus.TextFormatter{DisableColors: true}
	logFile, err := os.OpenFile(file, flag, chmod)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to write file on filehook %v", err)
		return nil, err
	}

	return &FileHook{logFile, flag, chmod, plainFormatter}, err
}

// Fire event
func (hook *FileHook) Fire(entry *logrus.Entry) error {

	plainformat, _ := hook.formatter.Format(entry)
	line := string(plainformat)
	_, err := hook.file.WriteString(line)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to write file on filehook(entry.String)%v", err)
		return err
	}

	return nil
}

func (hook *FileHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}
