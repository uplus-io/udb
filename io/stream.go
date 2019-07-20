/*
 * Copyright (c) 2019 uplus.io
 */

package io

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Stream struct {
	path string
	file *os.File

	reader *bufio.Reader
	writer *bufio.Writer
}

type Reader interface {
	read(length int, buffer Buffer)
}

func NewStream(path string) (*Stream) {
	stream := &Stream{path: path}
	stream.Init()
	return stream;
}

func (p *Stream) Init() {
	var err error
	p.file, err = os.OpenFile(p.path, os.O_CREATE|os.O_RDWR|os.O_SYNC, os.ModeAppend|0644)
	if err != nil {
		//log.Fatal(err)
		if os.IsNotExist(err) {
			index := strings.LastIndex(p.path, string(filepath.Separator))
			dir := string([]byte(p.path)[:index])
			os.MkdirAll(dir, os.ModePerm)
		}
	}

	reader := bufio.NewReader(p.file)
	writer := bufio.NewWriter(p.file)
	p.reader = reader
	p.writer = writer

}

func (p *Stream) Close() {
	p.writer.Flush()
	p.file.Close()
}

func (p *Stream) ReadLine(offset int64) ([]byte, error) {
	p.file.Seek(offset, 0)
	reader := p.reader
	line, _, err := reader.ReadLine()
	return line, err
}

func (p *Stream) Read(offset int64, length uint32) Buffer {
	reader := p.reader
	bytes := make([]byte, length)
	p.file.Seek(offset, 0)
	reader.Read(bytes)
	buffer := NewBigBuffer(nil)
	buffer.WriteBytes(bytes)
	reader.ReadLine()
	return *buffer
}

func (p *Stream) Write(offset int64, buffer Buffer) {
	writer := p.writer
	p.file.Seek(offset, 0)
	writer.Write(buffer.Bytes())
}

func (p *Stream) Flush() error {
	return p.writer.Flush()
}
