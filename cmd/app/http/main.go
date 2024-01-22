package main

import (
	"github.com/mauriciozanettisalomao/go-transaction-service/cmd/app"
	"github.com/mauriciozanettisalomao/go-transaction-service/log"
)

func init() {
	log.InitStructureLogConfig()
}

func main() {
	err := app.Endpoints().Run()
	if err != nil {
		panic(err)
	}
}
