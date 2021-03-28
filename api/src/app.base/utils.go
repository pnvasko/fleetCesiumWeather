package appcontroler

import (
	"bytes"
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	stringJsonValueStart = byte('"')
	stringJsonValueEnd   = byte('"')
	arrayJsonValueStart  = byte('[')
	arrayJsonValueEnd    = byte(']')
	dictJsonValueStart   = byte('{')
	dictJsonValueEnd     = byte('}')
	commaJsonValueEnd    = byte(',')
)

var JsonProcessor = jsoniter.ConfigCompatibleWithStandardLibrary

func getDefaultLogger() *Logger {
	var log = logrus.New()
	log.Formatter = new(logrus.TextFormatter)
	log.Level = logrus.DebugLevel
	log.Out = os.Stderr
	l := &Logger{log: log}
	return l
}

func getLog(ctx context.Context, module interface{}) *Logger {
	if log, ok := ctx.Value("logger").(*Logger); ok {
		return log
	}
	log := getDefaultLogger()
	log.Warning(fmt.Sprintf("Used self logger for %T", module))
	return log
}

func GetLog(ctx context.Context, module interface{}) *Logger {
	return getLog(ctx, module)
}

func GetWorkerPool(ctx context.Context) (*WorkerPool, error) {
	if pool, ok := ctx.Value("worker_pool").(*WorkerPool); ok {
		return pool, nil
	}
	return nil, fmt.Errorf("can't get worker pool from context")
}

func GetJsonValueByKey(data []byte, key string) ([]byte, error) {
	bkey := []byte(fmt.Sprintf(`"%s":`, key))
	lenKey := len(bkey)
	x := bytes.Index(data, bkey)
	if x == -1 {
		return nil, fmt.Errorf("key not found")
	}
	targetPart := data[x+lenKey:]
	startTargetPart := targetPart[0]

	switch startTargetPart {
	case stringJsonValueStart:
		y := bytes.IndexByte(targetPart[1:], stringJsonValueEnd)
		return targetPart[1 : y+1], nil
	case arrayJsonValueStart:
		openDCount := bytes.Count(targetPart, []byte{arrayJsonValueStart})
		closeDCount := bytes.Count(targetPart, []byte{arrayJsonValueEnd})
		for i := 0; i < closeDCount-openDCount; i++ {
			z := bytes.LastIndexByte(targetPart, arrayJsonValueEnd)
			targetPart = targetPart[0:z]
		}
		return targetPart, nil
	case dictJsonValueStart:
		openDCount := bytes.Count(targetPart, []byte{dictJsonValueStart})
		closeDCount := bytes.Count(targetPart, []byte{dictJsonValueEnd})
		for i := 0; i < closeDCount-openDCount; i++ {
			z := bytes.LastIndexByte(targetPart, dictJsonValueEnd)
			targetPart = targetPart[0:z]
		}
		return targetPart, nil
	default:
		if startTargetPart >= 48 && startTargetPart <= 57 {
			y := bytes.IndexByte(targetPart[1:], commaJsonValueEnd)
			if y == -1 {
				y = bytes.IndexByte(targetPart[1:], dictJsonValueEnd)
			}
			return targetPart[0 : y+1], nil
		}
	}

	return nil, nil
}
