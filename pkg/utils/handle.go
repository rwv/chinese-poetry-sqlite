package utils

import "io"

type Handler interface {
	HandleJSONs(entrys []Entry, filename string) error
	IsPoem(path string) bool
}

func DoTheHandle(reader io.ReaderAt, size int64, handler Handler, sqlitePath string) error {
	entries, err := GetEntries(reader, size)
	if err != nil {
		return err
	}

	var filteredEntries = make([]Entry, 0)
	for _, entry := range entries {
		if handler.IsPoem(entry.Path()) {
			filteredEntries = append(filteredEntries, entry)
		}
	}

	if err := handler.HandleJSONs(filteredEntries, sqlitePath); err != nil {
		return err
	}

	return nil
}
