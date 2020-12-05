package main

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
	"net/http"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func init() {
	var err error
	err = godotenv.Load(fmt.Sprintf("./%s.env", os.Getenv("GO_ENV")))
	if err != nil {
		panic(err)
	}

	DB_NAME := os.Getenv("DB_NAME")
	DB_USER := os.Getenv("DB_USER")
	DB_PASS := os.Getenv("DB_PASS")
	Db, err = sql.Open("mysql", DB_USER + ":" + DB_PASS + "@/" + DB_NAME)
	if err != nil {
		panic(err)
	}
}

type Post struct {
	Id	int	`json:"id"`
	Page	string	`json:"page"`
	Content	string	`json:"content"`
	Likes	int	`json:"likes"`
}

type Vote struct {
	Page	int	`json:"page"`
	Vote_count	string	`json:"Vote_count"`
}

func Posts(w http.ResponseWriter, r *http.Request) {
	var post Post
	var posts []Post

	rows, err := Db.Query("select id,pages,content,likes from posts")
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		err = rows.Scan(&post.Id, &post.Page, &post.Content, &post.Likes)
		if err != nil {
			return
		}
		posts = append(posts, post)
	}
	rows.Close()

	res, err := json.Marshal(posts)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	fmt.Fprintln(w, string(body))

	var post Post
	err := json.Unmarshal(body[:len], &post)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
    	return
	}
	fmt.Fprintln(w, post)
}

func Votes(w http.ResponseWriter, r *http.Request) {
	var vote Vote
	var votes []Vote
	rows, err := Db.Query("select * from votes")
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		err = rows.Scan(&vote.Page, &vote.Vote_count)
		if err != nil {
			return
		}
		votes = append(votes, vote)
	}
	rows.Close()

	res, err := json.Marshal(votes)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}


func main() {
	server := http.Server{
		Addr:	"127.0.0.1:8180",	
	}
	http.HandleFunc("/vote", Votes)
	http.HandleFunc("/post", Posts)
	http.HandleFunc("/create", CreatePost)

	server.ListenAndServe()
}