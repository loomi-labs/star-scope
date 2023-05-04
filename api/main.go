package main

import (
	_ "github.com/hedwigz/entviz"
	"github.com/loomi-labs/star-scope/cmd"
	"github.com/loomi-labs/star-scope/database"
	"github.com/shifty11/go-logger/log"
)

func main() {
	defer log.SyncLogger()
	defer database.Close()

	defer func() {
		if err := recover(); err != nil {
			log.Sugar.Panic(err)
			return
		}
	}()

	cmd.Execute()
}
