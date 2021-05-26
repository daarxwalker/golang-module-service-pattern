package core

import "os"

type File interface {
	File() *os.File
	Type() string
}

type file struct {
	file     *os.File
	fileType string
}

func (f file) File() *os.File {
	return f.file
}

func (f file) Type() string {
	return f.fileType
}
