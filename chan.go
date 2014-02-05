package main

import (
		"database/sql"
		"fmt"
		_ "github.com/mattn/go-sqlite3"
		"log"
)

func main(){
	db, err := sql.Open("sqlite3", "./imageboard.db")
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("INSERT INTO post(subject, name, text, thread_id) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	
	_, err = stmt.Exec("Thread Subject", "Zach", "A new thread.", 1)
	if err != nil{
		log.Fatal(err)
	}
	_, err = stmt.Exec("Post Subject", "Anonymous", "A new post.", 1)
	if err != nil{
		log.Fatal(err)
	}
	_, err = stmt.Exec("Post Subject 2", "Anonymous", "Another new thread.", 2)
	if err != nil{
		log.Fatal(err)
	}
	tx.Commit()

	latest_threads, err := db.Query("SELECT * FROM latest_threads")
	if err != nil{
		log.Fatal(err)
	}
	defer latest_threads.Close()
	for latest_threads.Next(){
		var threadID int
		var subject string
		var name string
		var text string
		var date_posted string
		latest_threads.Scan(&threadID, &subject, &name, &text, &date_posted) 
		fmt.Println(threadID, subject, name, text, date_posted)
	}
	latest_threads.Close()
}