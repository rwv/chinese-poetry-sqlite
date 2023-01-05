package utils

import "io"

type Handler interface {
	HandleJSONs(jsonReaders []io.Reader) error
	SaveToSqlite(filename string) error
	IsPoem(path string) bool
}

func DoTheHandle(reader io.ReaderAt, size int64, handler Handler, sqlitePath string) error {
	entries, err := GetEntries(reader, size)
	if err != nil {
		return err
	}

	var filteredEntries = make([]Entry, 0)
	for _, entry := range entries {
		if handler.IsPoem(entry.Path) {
			filteredEntries = append(filteredEntries, entry)
		}
	}

	var jsonReaders = make([]io.Reader, 0)
	for _, entry := range filteredEntries {
		jsonReaders = append(jsonReaders, entry.GetReader())
	}

	if err := handler.HandleJSONs(jsonReaders); err != nil {
		return err
	}

	if err := handler.SaveToSqlite(sqlitePath); err != nil {
		return err
	}

	return nil
}
