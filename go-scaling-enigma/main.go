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

	err = postgres.CreateUsersTable(db)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}
	err = postgres.CreateChatsTable(db)
	if err != nil {
		log.Fatalf("Error creating chats table: %v", err)
	}
	err = postgres.CreateMessagesTable(db)
	if err != nil {
		log.Fatalf("Error creating messages table: %v", err)
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
