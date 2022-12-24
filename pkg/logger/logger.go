package logger

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"forum/internal/entity"
)

type Logger struct {
	Info *log.Logger
	Err  *log.Logger
}

func getRootPath() string {
	separator := "/"
	if runtime.GOOS == "windows" {
		separator = "\\"
	}

	// getting full path from where program is running
	_, basePath, _, _ := runtime.Caller(0)
	pathSlice := strings.Split(filepath.Dir(basePath), separator)
	tmpSl := []string{}
	last := false

	// separating root directory
	for i := 0; i < len(pathSlice); i++ {
		tmpSl = append(tmpSl, pathSlice[i])
		if pathSlice[i] == "forum" {
			last = true
		}
		if last {
			break
		}
	}

	return strings.Join(tmpSl, separator) + separator
}

func New() *Logger {
	path := getRootPath()

	file, err := os.OpenFile(path+"logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o664)
	if err != nil {
		log.Fatal(fmt.Errorf("logger - New - os.OpenFile: %w", err))
	}

	InfoLogger := log.New(file, "INFO: ", log.Ldate|log.Ltime)
	ErrorLogger := log.New(file, "ERROR ", log.Ldate|log.Ltime)
	return &Logger{
		Info: InfoLogger,
		Err:  ErrorLogger,
	}
}

func (l *Logger) WriteLog(err error) {
	if !errors.Is(err, entity.ErrUserNotFound) && !errors.Is(err, entity.ErrPostNotFound) &&
		!errors.Is(err, entity.ErrUserEmailAlreadyExists) && !errors.Is(err, entity.ErrUserNameAlreadyExists) &&
		!errors.Is(err, entity.ErrUserPasswordIncorrect) && !errors.Is(err, entity.ErrUserEmailIncorrect) {
		l.Err.Println(err)
	} else {
		l.Info.Println(err)
	}
}
