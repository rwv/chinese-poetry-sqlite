package utils

import (
	"archive/zip"
	"io"
)

type Entry struct {
	Path      string
	GetReader func() io.Reader
}

// GetEntries returns a slice of io.Reader, each of which is a json file.
func GetEntries(reader io.ReaderAt, size int64) ([]Entry, error) {
	zipReader, err := zip.NewReader(reader, size)
	if err != nil {
		return nil, err
	}

	var entries []Entry
	for _, file := range zipReader.File {
		entries = append(entries, Entry{
			Path: file.Name,
			GetReader: func() io.Reader {
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
