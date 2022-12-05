package logger

import (
	"errors"
	"fmt"
	"forum/internal/entity"
	"log"
	"os"
)

type Logger struct {
	Info *log.Logger
	Err  *log.Logger
}

func New() *Logger {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o664)
	if err != nil {
		log.Fatal(fmt.Errorf("logger - New - os.OpenFile: %w", err))
	}

	InfoLogger := log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger := log.New(file, "ERROR ", log.Ldate|log.Ltime|log.Lshortfile)
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
