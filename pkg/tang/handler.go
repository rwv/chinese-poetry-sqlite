package tang

import (
	"database/sql"
	"encoding/json"
	"io"

	_ "github.com/mattn/go-sqlite3"
)

type TangPoem struct {
	Author     string   `json:"author"`
	Paragraphs []string `json:"paragraphs"`
	Title      string   `json:"title"`
	ID         string   `json:"id"`
}

const prefix = "chinese-poetry-master/json/poet.tang."

type poemType = TangPoem

type Handler struct {
	poems []poemType
}

func (h *Handler) SaveToSqlite(filename string) error {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return err
	}
	defer db.Close()

	// create table
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS tang (author TEXT, paragraphs TEXT, title TEXT, id TEXT PRIMARY KEY)")

	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	// insert data
	stmt, err = db.Prepare("INSERT INTO tang (author, paragraphs, title, id) VALUES (?, ?, ?, ?)")

	for _, poem := range h.poems {
		_, err = stmt.Exec(poem.Author, poem.Paragraphs, poem.Title, poem.ID)
	}

	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) HandleJSONs(jsonReaders []io.Reader) error {
	for _, jsonReader := range jsonReaders {
		var poem poemType
		if err := json.NewDecoder(jsonReader).Decode(&poem); err != nil {
			return err
		}
		h.poems = append(h.poems, poem)
	}
	return nil
}

func (h *Handler) IsPoem(path string) bool {
	return len(path) > len(prefix) && path[:len(prefix)] == prefix && path[len(path)-5:] == ".json"
}

func New() *Handler {
	handler := &Handler{}
	handler.poems = make([]poemType, 0)

	return handler
}
