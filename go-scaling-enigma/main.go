package main

import (
	"database/sql"
	"log"
	"scaling-enigma/go-scaling-enigma/main.go/server"
	"scaling-enigma/go-scaling-enigma/main.go/server/postgres"
	"scaling-enigma/go-scaling-enigma/main.go/tui"
	"sync"
)

var db *sql.DB

func main() {

	db, err := postgres.Connect(db)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	dbErr := postgres.CreateUsersTable(db)
	if dbErr != nil {
		log.Fatalf("Error creating users table: %v", dbErr)
	}
	defer db.Close()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		server.ServeGin(db)
	}()

	tui.StartCLI()

	wg.Wait()

}
