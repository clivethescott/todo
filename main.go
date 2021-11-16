package main

import (
	"log"
	"net/http"
	"time"

	"github.com/clivethescott/todo/db"
	"github.com/clivethescott/todo/server"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	serverConfig := server.Config{
		Addr:         ":3000",
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		DB: db.Config{
			DSN:             "root:@unix(/tmp/mysql.sock)/todo?sql_mode=TRADITIONAL&tls=skip-verify&autocommit=true",
			ConnMaxLifeTime: 5 * time.Minute,
			MaxOpenConns:    10,
			MaxIdleConns:    10,
		},
	}

	if err := server.Start(serverConfig); err != nil {
		if err == http.ErrServerClosed {
			log.Println("server shut down")
			return
		}
		log.Fatal(err)
	}
}
