package main

import (
		"fmt"
		"log"
		_ "github.com/mattn/go-sqlite3"
		"database/sql"
)

func main(){
	db, err := sql.Open("sqlite3", "./imageboard.db")
	if err != nil{
		log.Fatal(err)
	}

	newPost(db, "Post Subject method test", "Anon", "First method post", 3)
	latestThreads(db)
	
}

func newPost(db *sql.DB, subject string, name string, text string, thread_id int) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO post(subject, name, text, thread_id) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(subject, name, text, thread_id)
	if err != nil{
		log.Fatal(err)
	}
	tx.Commit()
}

func latestThreads(db *sql.DB){
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
