package main

import (
		"log"
		_ "github.com/mattn/go-sqlite3"
		"database/sql"
		"net/http"
		"html/template"
		"strconv"
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

	//i.newPost("Post Subject method test", "Anon", "First method post", 2)
	http.HandleFunc("/", i.latestThreads)
	http.HandleFunc("/reply/", i.viewReplies)
	http.HandleFunc("/newreply/", i.newReply)
	http.HandleFunc("/newthread/", i.newThread)
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

func (i *ImageBoard) newThread(w http.ResponseWriter, r *http.Request){
	subject := r.FormValue("subject")
	name := r.FormValue("name")
	message := r.FormValue("message")

	var thread_id int

	err := i.db.QueryRow("SELECT MAX(thread_id) FROM post").Scan(&thread_id)
	if err != nil {
		log.Fatal(err)
	}

	threadID := thread_id + 1
	i.newPost(subject, name, message, threadID)
	http.Redirect(w, r, "/", http.StatusFound)
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
	t, error := template.ParseFiles("index.html")
	if error != nil{
		log.Fatal(error)
	}
	t.Execute(w, posts)
}

func (i *ImageBoard) viewReplies(w http.ResponseWriter, r *http.Request){
	threadID, _ := strconv.Atoi(r.URL.Path[len("/reply/"):])
	replies, err := i.db.Query("SELECT * FROM latest_threads where thread_id="+strconv.Itoa(threadID));
	if err != nil{
		log.Fatal(err)
	}
	posts := []Post{}
	defer replies.Close()
	for replies.Next(){
		var threadID int
		var subject string
		var name string
		var text string
		var date_posted string
		replies.Scan(&threadID, &subject, &name, &text, &date_posted)
		post := Post{
			threadID,
			subject,
			name,
			text,
			date_posted,
		}
		posts = append(posts, post)
	}
	layoutData := struct {
		ThreadID int
		Posts []Post
	} {
		ThreadID: threadID,
		Posts: posts,
	}
	t, error := template.ParseFiles("thread.html")
	if error != nil{
		log.Fatal(error)
	}
	t.Execute(w, layoutData)
}

func (i *ImageBoard) newReply(w http.ResponseWriter, r *http.Request){
	threadIDstring := r.URL.Path[len("/newreply/"):]
	threadID, err := strconv.Atoi(threadIDstring)
	if err != nil {
		log.Fatal(err)
	}
	subject := r.FormValue("subject")
	name := r.FormValue("name")
	message := r.FormValue("message")
	i.newPost(subject, name, message, threadID)
	http.Redirect(w, r, "/reply/"+threadIDstring, http.StatusFound)
}

