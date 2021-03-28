package appcontroler

import (
	"fmt"
	"os"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	log *logrus.Logger
}

func (l *Logger) Info(args ...interface{}) {
	l.log.Log(logrus.InfoLevel, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.log.Log(logrus.DebugLevel, args...)
}

func (l *Logger) Warning(args ...interface{}) {
	l.log.Log(logrus.WarnLevel, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.log.Log(logrus.ErrorLevel, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log.Logf(logrus.FatalLevel, format, args...)
	l.log.Exit(1)
}

func NewLoger(name, log_path string, debug bool) *Logger {
	fmt.Println("\t configure logger...")
	l := Logger{}
	l.log = newLogrusLoger(name, log_path, debug)

	return &l
}

func newLogrusLoger(name, log_path string, debug bool) *logrus.Logger {
	fmt.Println("\t configure logger...")
	var log = logrus.New()
	log.Formatter = new(logrus.TextFormatter)
	log.Level = logrus.InfoLevel
	if debug {
		log.Level = logrus.DebugLevel
	}
	log.Out = os.Stderr
	fmt.Println("\t configure logger...", log_path, "; debug: ", debug)
	if log_path != "" {
		log.Formatter = new(logrus.JSONFormatter)
		writer, err := rotatelogs.New(
			fmt.Sprintf("%s/%s.%s", log_path, name, "%Y%m%d%H%M"),
			rotatelogs.WithLinkName(fmt.Sprintf("%s/%s", log_path, name)),
			rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
			rotatelogs.WithRotationTime(time.Duration(604800)*time.Second),
		)
		if err != nil {
			fmt.Println("Error open log file: ", err)
			panic("Error open log file")
		}
		log.Out = writer
	}
	return log
}
