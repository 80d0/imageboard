package main

import (
		"log"
		_ "github.com/mattn/go-sqlite3"
		"database/sql"
		"net/http"
		"html/template"
)

type ImageBoard struct {
	db *sql.DB
}

type Post struct {
	threadID int
	subject string
	name string
	text string
	date_posted string
}

func (p *Post) ThreadID()   int    { return p.threadID    }
func (p *Post) Subject()    string { return p.subject     }
func (p *Post) Name()       string { return p.name        }
func (p *Post) Text()       string { return p.text        }
func (p *Post) DatePosted() string { return p.date_posted }

func main(){
	db, err := sql.Open("sqlite3", "./imageboard.db")
	if err != nil{
		log.Fatal(err)
	}
	i := &ImageBoard{db}

	i.newPost("Post Subject method test", "Anon", "First method post", 2)
	http.HandleFunc("/", i.latestThreads)
	http.ListenAndServe(":8080", nil)
}

func (i *ImageBoard) newPost(subject string, name string, text string, thread_id int) {
	tx, err := i.db.Begin()
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

func (i *ImageBoard) latestThreads(w http.ResponseWriter, r *http.Request){
	latest_threads, err := i.db.Query("SELECT * FROM latest_threads GROUP BY thread_id")
	if err != nil{
		log.Fatal(err)
	}
	posts := []Post{}
	defer latest_threads.Close()
	for latest_threads.Next(){
		var threadID int
		var subject string
		var name string
		var text string
		var date_posted string
		latest_threads.Scan(&threadID, &subject, &name, &text, &date_posted) 
		post := Post{
			threadID,
			subject,
			name,
			text,
			date_posted,
		}
		posts = append(posts, post)
	}
	t, error := template.ParseFiles("thread.html")
	if error != nil{
		log.Fatal(error)
	}
	t.Execute(w, posts)
}

/*
SELECT * FROM latest_threads GROUP BY thread_id;
SELECT * FROM latest_threads
*/