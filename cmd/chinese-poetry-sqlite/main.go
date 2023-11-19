package main

import (
	"os"

	"github.com/rwv/chinese-poetry-sqlite/pkg/tang"
	"github.com/rwv/chinese-poetry-sqlite/pkg/utils"
)

const path = "./chinese-poetry-data/master.zip"

func main() {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}

	os.RemoveAll("output")
	os.MkdirAll("output", os.ModePerm)

	tangHandler := tang.New()
	err = utils.DoTheHandle(file, fileInfo.Size(), tangHandler, "output/tang.sqlite")

	if err != nil {
		panic(err)
	}

	// tang.TurnPoemsToSQLite()
}
