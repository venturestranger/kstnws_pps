package main

import (
	"fmt"
	"log"
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
	Category			string	`json:"category" db:"category"`
	Hashtags			string	`json:"hashtags" db:"hashtags"`
}

// gets posts from the Pull-Server
func GetHandler(c *gin.Context) {
	ch := make(chan bool)

	go func() {
		db, err := sqlx.Open("postgres", dsn)
		if err != nil {
			log.Println(err)
			SendStatus(http.StatusBadRequest, c)
			ch <- true
			return
		}
		defer db.Close()

		offset := c.Query("offset_")
		limit := c.Query("limit_")
		order_way := c.Query("order_way_")

		var posts []Post
		db.Select(&posts, fmt.Sprintf("select * from posts order by date_publication %s offset %s limit %s", order_way, offset, limit))

		SendData(posts, c)
		ch <- true
	}()
	<- ch
}

// deletes a post on the Pull-Server
func DeleteHandler(c *gin.Context) {
	ch := make(chan bool)

	go func() {
		db, err := sqlx.Open("postgres", dsn)
		if err != nil {
			log.Println(err)
			SendStatus(http.StatusBadRequest, c)
			ch <- true
			return
		}
		defer db.Close()

		id := c.Query("id")
		err = db.Exec(fmt.Sprintf("delete from posts where id = %s", id))

		if err != nil {
			log.Println(err)
			SendStatus(http.StatusInternalServerError, c)
			ch <- true
			return
		}

		SendStatus(http.StatusOK, c)
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
			log.Println(err)
			SendStatus(http.StatusBadRequest, c)
			ch <- true
			return
		}

		id := c.Query("id")

		var post Post
		db.Get(&post, fmt.Sprintf("select * from posts where id = %s", id))
		payload, _ := json.Marshal(post)

		req, _ := http.NewRequest("POST", apiAddr + "/rest/posts", bytes.NewReader(payload))
		req.Header.Set("Authorization", "Bearer " + apiToken)
		resp, err := httpClient.Do(req)

		if err != nil {
			log.Println(err)
			SendStatus(http.StatusInternalServerError, c)
			ch <- true
			return
		}

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
			log.Println(err)
			SendStatus(http.StatusBadRequest, c)
			ch <- true
			return
		}
		defer db.Close()

		var post Post
		payload, _ := ioutil.ReadAll(c.Request.Body)
		json.Unmarshal(payload, &post)
		_, err = db.NamedExec("insert into posts(id_author, title, lead, picture_url, content, date_publication, date_edit, category, hashtags) values(:id_author, :title, :lead, :picture_url, :content, :date_publication, :date_edit, :category, :hashtags)", post)

		if err != nil {
			log.Println(err)
			SendStatus(http.StatusBadRequest, c)
			ch <- true
			return
		}

		SendStatus(http.StatusOK, c)
		ch <- true
	}()
	<- ch
}
