package main

import (
	"fmt"
	"bytes"
	"io/ioutil"
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Post struct {
	Id					int		`json:"id" db:"id"`
	IdAuthor			int		`json:"id_author" db:"id_author"`
	Title				string	`json:"title" db:"title"`
	Lead				string	`json:"lead" db:"lead"`
	PictureUrl			string	`json:"picture_url" db:"picture_url"`
	Content				string	`json:"content" db:"content"`
	DatePublication		string	`json:"date_publication" db:"date_publication"`
	DateEdit			string	`json:"date_edit" db:"date_edit"`
}

// gets posts from the database 
func GetHandler(c *gin.Context) {
	ch := make(chan bool)

	go func() {
		db, err := sqlx.Open("postgres", dsn)
		if err != nil {
			fmt.Println(err)
			SendStatus(http.StatusBadRequest, c)
			ch <- true
			return
		}
		defer db.Close()

		offset := c.Query("offset_")
		limit := c.Query("limit_")
		order_way := c.Query("order_way_")

		fmt.Println(offset, limit, order_way)

		var posts []Post
		db.Select(&posts, fmt.Sprintf("select * from posts order by date_publication %s offset %s limit %s", order_way, offset, limit))

		SendData(posts, c)
		ch <- true
	}()
	<- ch
}

// pushes a post to the API
func PushHandler(c *gin.Context) {
	ch := make(chan bool)

	go func() {
		db, err := sqlx.Open("postgres", dsn)
		if err != nil {
			fmt.Println(err)
			SendStatus(http.StatusBadRequest, c)
			ch <- true
			return
		}
		defer db.Close()

		id := c.Query("id")
		
		var post Post
		db.Get(&post, fmt.Sprintf("select * from posts where id = %s", id))
		payload, _ := json.Marshal(post)

		req, _ := http.NewRequest("POST", apiAddr + "/rest/posts", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer " + apiToken)
		resp, err := httpClient.Do(req)
		
		SendStatus(resp.StatusCode, c)
		ch <- true
	}()
	<- ch
}

// creates a new post on the Pull Server
func PostHandler(c *gin.Context) {
	ch := make(chan bool)

	go func() {
		db, err := sqlx.Open("postgres", dsn)
		if err != nil {
			fmt.Println(err)
			SendStatus(http.StatusBadRequest, c)
			ch <- true
			return
		}
		defer db.Close()

		var post Post
		payload, _ := ioutil.ReadAll(c.Request.Body)
		json.Unmarshal(payload, &post)
		_, err = db.NamedExec("insert into posts(id_author, title, lead, picture_url, content, date_publication, date_edit) values(:id_author, :title, :lead, :picture_url, :content, :date_publication, :date_edit)", post)

		if err != nil {
			fmt.Println(err)
			SendStatus(http.StatusBadRequest, c)
			ch <- true
			return
		}
		
		SendStatus(http.StatusOK, c)
		ch <- true
	}()
	<- ch
}
