package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"gitlab.com/yulqen/eden/internal/repository"
)

const fileName = "eden.db"

func dbSetup() {
	os.Remove(fileName)

	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		log.Fatal(err)
	}

	entryRepo := repository.NewSQLiteRepository(db)

	if err := entryRepo.Migrate(); err != nil {
		log.Fatal(err)
	}

	// sample1 := journal.Entry{
	// 	Content: "TEST ENTRY SAMPLE 1",
	// }
	//
	// sample2 := journal.Entry{
	// 	Content: "TEST ENTRY SAMPLE 2",
	// }
	//
	// createdEntry1, err := entryRepo.Create(sample1)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// createdEntry2, err := entryRepo.Create(sample2)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(createdEntry1.Content)
	// fmt.Println(createdEntry2)
	//
	// allEntries, err := entryRepo.All()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%v\n", allEntries)
}

func main() {

	//addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	// does nothing for now
}
