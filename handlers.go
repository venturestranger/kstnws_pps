package main

import (
	"fmt"
	"log"
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
	Comment				string	`json:"comment" db:"comment"`
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

		statement := "select * from posts "

		if c.Query("id") != "" || c.Query("id_author") != "" || c.Query("category") {
			statement += " where "
			andFlag := false
			if c.Query("id") != "" {
				statement += " id = " + c.Query("id")
				andFlag = true
			}
			if c.Query("id_author") != "" {
				if andFlag {
					statement += " and "
				}

				statement += " id_author = " + c.Query("id_author")
				andFlag = true
			}
			if c.Query("category") != "" {
				if andFlag {
					statement += " and "
				}

				statement += " category = " + c.Query("category")
				andFlag = true
			}
		}
		if c.Query("order_way_") != "" {
			statement += " order by date_publication " + c.Query("order_way_")
		}
		if c.Query("offset_") != "" {
			statement += " offset " + c.Query("offset_")
		}
		if c.Query("limit_") != "" {
			statement += " limit " + c.Query("limit_")
		}

		var posts []Post
		db.Select(&posts, statement)

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
		id_author := c.Query("id_author")
		_, err = db.Exec(fmt.Sprintf("delete from posts where id = %s", id))

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
		comment := c.Query("comment")
		pass := c.Query("pass")

		if pass == "true" {
			var post Post
			db.Get(&post, fmt.Sprintf("select * from posts where id = %s", id))
			PostPushToAPI(post)
			db.Exec(fmt.Sprintf("delete from posts where id = %s", id))
		} else {
			db.Exec(fmt.Sprintf("update posts set comment = '%s' where id = %s", comment, id))
		}

		SendStatus(http.StatusOK, c)
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
		_, err = db.NamedExec("insert into posts(id_author, title, lead, picture_url, content, date_publication, date_edit, category, hashtags, comment) values(:id_author, :title, :lead, :picture_url, :content, :date_publication, :date_edit, :category, :hashtags, :comment)", post)

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
