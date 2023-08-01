package main

import (
	"log"

	"github.com/kainnsoft/migrator/config"
	"github.com/kainnsoft/migrator/internal/app"
)

func main() {
	var (
		cfg *config.Config
		//mySql mysql.IMySQL
		//res   *sql.Rows
		err error
	)
	if cfg, err = config.NewConfig(); err != nil {
		log.Fatalf("Config error: %s", err)

	}
	app.Run(cfg)

	//if err = mySql.DB().Ping(); err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Print("Pong\n")
}
