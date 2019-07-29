/*
 * Copyright (c) 2019 uplus.io
 */

package logger

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LoggerLevel uint32

var (
	stdLogger         Logger
	fileLogger        Logger
	debugLoggerEnable bool
)

const (
	LoggerLevelPanic LoggerLevel = iota
	LoggerLevelFatal
	LoggerLevelError
	LoggerLevelWarn
	LoggerLevelInfo
	LoggerLevelDebug
	LoggerLevelTrace
)

type Logger interface {
	Writer() io.Writer
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
}

func init() {
	debugLoggerEnable = true
	stdLogger = newLoggerLogrus(LoggerLevelDebug, os.Stdout)
}

func DebugLoggerEnable(enable bool) {
	debugLoggerEnable = enable
}

func NewLogger(level LoggerLevel, path, filename string) Logger {
	if fileLogger != nil {
		panic("file logger cannot repeat")
	}
	fileLogger = newLoggerLogrus(level, newLogFileWriter(path, filename))
	return fileLogger
}

func Writer() io.Writer {
	if fileLogger != nil {
		return fileLogger.Writer()
	} else {
		return stdLogger.Writer()
	}
}
func Debugf(format string, args ...interface{}) {
	if debugLoggerEnable {
		stdLogger.Debugf(format, args...)
	}
	if fileLogger != nil {
		fileLogger.Debugf(format, args...)
	}
}
func Infof(format string, args ...interface{}) {
	if debugLoggerEnable {
		stdLogger.Infof(format, args...)
	}
	if fileLogger != nil {
		fileLogger.Infof(format, args...)
	}
}
func Warnf(format string, args ...interface{}) {
	if debugLoggerEnable {
		stdLogger.Warnf(format, args...)
	}
	if fileLogger != nil {
		fileLogger.Warnf(format, args...)
	}
}
func Errorf(format string, args ...interface{}) {
	if debugLoggerEnable {
		stdLogger.Errorf(format, args...)
	}
	if fileLogger != nil {
		fileLogger.Errorf(format, args...)
	}
}
func Fatalf(format string, args ...interface{}) {
	if debugLoggerEnable {
		stdLogger.Fatalf(format, args...)
	}
	if fileLogger != nil {
		fileLogger.Fatalf(format, args...)
	}
}
func Panicf(format string, args ...interface{}) {
	if debugLoggerEnable {
		stdLogger.Panicf(format, args...)
	}
	if fileLogger != nil {
		fileLogger.Panicf(format, args...)
	}
}

type logFileWriter struct {
	path     string
	filename string

	file *os.File
	date string
}

func newLogFileWriter(path, filename string) (log *logFileWriter) {
	log = &logFileWriter{path: path, filename: filename}
	file, err := os.OpenFile(filepath.Join(path, filename), os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
	if err != nil {
		panic(fmt.Sprintf("init log file[%s] fail", path))
	}
	log.date = time.Now().Format("20060102")
	log.file = file
	return
}
func (p *logFileWriter) Write(data []byte) (n int, err error) {
	if p == nil {
		return 0, errors.New("logFileWriter is nil")
	}
	if p.file == nil {
		return 0, errors.New("file not opened")
	}
	n, e := p.file.Write(data)
	//文件最大 64K byte
	newDate := time.Now().Format("20060102")
	if !strings.EqualFold(p.date, newDate) {
		p.file.Close()
		filePath := filepath.Join(p.path, p.filename)
		renameFilePath := filepath.Join(p.path, p.filename+fmt.Sprintf("-%s", newDate))
		err := os.Rename(filePath, renameFilePath)
		if err != nil {
			panic(fmt.Sprintf("rename log file from[%s] to[%s] fail", filePath, renameFilePath))
		}
		p.file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
		if err != nil {

		}
	}
	return n, e
}

type EmptyWriter struct {
}

func NewEmptyWriter() *EmptyWriter {
	return &EmptyWriter{}
}

func (EmptyWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
