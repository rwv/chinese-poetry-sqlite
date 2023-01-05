package utils

import (
	"archive/zip"
	"io"
)

type Entry interface {
	Path() string
	GetReader() io.ReadCloser
}

type EntryInstance struct {
	path      string
	getReader func() io.ReadCloser
}

func (e EntryInstance) Path() string {
	return e.path
}

func (e EntryInstance) GetReader() io.ReadCloser {
	return e.getReader()
}

// GetEntries returns a slice of io.Reader, each of which is a json file.
func GetEntries(reader io.ReaderAt, size int64) ([]EntryInstance, error) {
	zipReader, err := zip.NewReader(reader, size)
	if err != nil {
		return nil, err
	}

	var entries []EntryInstance
	for _, file_ := range zipReader.File {
		file := file_
		entries = append(entries, EntryInstance{
			path: file.Name,
			getReader: func() io.ReadCloser {
				rc, err := file.Open()
				if err != nil {
					return nil
				}
				return rc
			},
		})
	}

	return entries, nil
}
