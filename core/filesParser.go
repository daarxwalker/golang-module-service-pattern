package core

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

type filesParser struct {
	provideService *provideService
}

func (f filesParser) createSessionTempDirIfNotExist(category string) string {
	token := f.provideService.getReducedSessionToken()

	if len(token) == 0 {
		return os.TempDir()
	}

	path := fmt.Sprintf("%s/%s/%s", os.TempDir(), token, category)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			f.provideService.checkError(err, errorSessionFolder)
		}
	}

	return path
}

func (f filesParser) createTempFileByCategory(category string, index int) *os.File {
	dir := f.createSessionTempDirIfNotExist(category)

	file, err := os.Create(fmt.Sprintf("%s/%d", dir, index))
	f.provideService.checkError(err, errorCreateTempFile)

	return file
}

func (f filesParser) parse() {
	for key, files := range f.provideService.GetRequest().MultipartForm.File {
		if strings.Contains(key, "file-") {
			fileCategory := strings.TrimPrefix(key, "file-")
			if len(files) == 0 {
				continue
			}

			for index, file := range files {
				tempFile := f.createTempFileByCategory(fileCategory, index)

				newFile, err := file.Open()
				if err != nil {
					f.provideService.checkError(err, errorOpenFile)
				}

				fileHeader := make([]byte, 512)

				if _, err := newFile.Read(fileHeader); err != nil {
					f.provideService.checkError(err, errorReadFile)
				}

				// set read position back to start
				if _, err := newFile.Seek(0, 0); err != nil {
					f.provideService.checkError(err, errorResetSeekFile)
				}

				buf := bytes.NewBuffer(nil)
				if _, err := io.Copy(buf, newFile); err != nil {
					f.provideService.checkError(err, errorBufferFile)
				}

				if _, err := tempFile.Write(buf.Bytes()); err != nil {
					f.provideService.checkError(err, errorWriteTempFile)
				}

				if err := tempFile.Close(); err != nil {
					f.provideService.checkError(err, errorCloseTempFile)
				}

				if err := newFile.Close(); err != nil {
					f.provideService.checkError(err, errorCloseFile)
				}
			}
		}
	}
}
