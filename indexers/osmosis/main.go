package main

import (
	"github.com/shifty11/go-logger/log"
)

func main() {
	defer log.SyncLogger()
	defer func() {
		if err := recover(); err != nil {
			log.Sugar.Panic(err)
			return
		}
	}()

	NewGRPCServer(&Config{Port: 50002}).Run()
}
