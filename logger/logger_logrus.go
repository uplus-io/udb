/*
 * Copyright (c) 2019 uplus.io
 */

package logger

import (
	"github.com/sirupsen/logrus"
	"io"
)

type LoggerLogrus struct {
	log *logrus.Logger

	writer io.Writer
}

func newLoggerLogrus(level LoggerLevel, writer io.Writer) (p Logger) {
	log := logrus.New()
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:      true,
		FullTimestamp:    true,
		QuoteEmptyFields: true})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(writer)
	// Only log the warning severity or above.
	log.SetLevel(logrus.Level(level))

	p = &LoggerLogrus{log: log, writer: writer}
	return
}

func (p *LoggerLogrus) init() {

}

func (p *LoggerLogrus) Writer() io.Writer {
	return p.writer
}

func (p LoggerLogrus) Debugf(format string, args ...interface{}) {
	p.log.Debugf(format, args...)
}
func (p LoggerLogrus) Infof(format string, args ...interface{}) {
	p.log.Infof(format, args...)
}
func (p LoggerLogrus) Warnf(format string, args ...interface{}) {
	p.log.Warnf(format, args...)
}
func (p LoggerLogrus) Errorf(format string, args ...interface{}) {
	p.log.Errorf(format, args...)
}
func (p LoggerLogrus) Fatalf(format string, args ...interface{}) {
	p.log.Fatalf(format, args...)
}
func (p LoggerLogrus) Panicf(format string, args ...interface{}) {
	p.log.Panicf(format, args...)
}
