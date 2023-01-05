package tang

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rwv/chinese-poetry-sqlite/pkg/utils"
)

type TangPoem struct {
	Author     *string   `json:"author"`
	Paragraphs []*string `json:"paragraphs"`
	Title      *string   `json:"title"`
	ID         *string   `json:"id"`
}

const prefix = "chinese-poetry-master/json/poet.tang.0"

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

	valueStrings := make([]string, 0, len(h.poems))
	valueArgs := make([]interface{}, 0, len(h.poems)*4)

	for _, poem := range h.poems {
		if poem.Author == nil || poem.Paragraphs == nil || poem.Title == nil || poem.ID == nil {
			return fmt.Errorf("invalid poem: %+v", poem)
		}

		paragraphsJsonText, err := json.Marshal(poem.Paragraphs)
		if err != nil {
			return err
		}

		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		valueArgs = append(valueArgs, poem.Author)
		valueArgs = append(valueArgs, string(paragraphsJsonText))
		valueArgs = append(valueArgs, poem.Title)
		valueArgs = append(valueArgs, poem.ID)
	}

	stmtText := fmt.Sprintf("INSERT INTO tang (author, paragraphs, title, id) VALUES %s",
		strings.Join(valueStrings, ","))

	_, err = db.Exec(stmtText, valueArgs...)

	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) HandleJSONs(entrys []utils.Entry) error {
	fmt.Println("Handle JSONs")
	for _, entry := range entrys {
		fmt.Println("Handle " + entry.Path())
		jsonReader := entry.GetReader()

		var poem []poemType
		if err := json.NewDecoder(jsonReader).Decode(&poem); err != nil {
			return err
		}

		h.poems = append(h.poems, poem...)

		jsonReader.Close()
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
