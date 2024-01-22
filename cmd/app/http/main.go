package main

import (
	"github.com/mauriciozanettisalomao/go-transaction-service/cmd/app"
	"github.com/mauriciozanettisalomao/go-transaction-service/log"
)

func init() {
	log.InitStructureLogConfig()
}

// @title					Go Transaction Service API
// @version					1.0
// @description				This is a sample server for a transaction service.
//
// @contact.name			Mauricio Zanetti Salomao
// @contact.url				https://github.com/mauriciozanettisalomao/go-transaction-service
// @contact.email			mauriciozanetti86@gmail.com
//
// @license.name			MIT
// @license.url				https://github.com/mauriciozanettisalomao/go-transaction-service/blob/main/LICENSE
func main() {
	err := app.Endpoints().Run()
	if err != nil {
		panic(err)
	}
}
