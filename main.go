package main

import (
	"fmt"
	"time"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"net/http"
)

func ISAUTH(c *gin.Context) {
	data := c.GetHeader("Authorization")
	values := strings.Split(data, " ")

	if len(values) > 1 {
		token, err := jwt.Parse(values[1], func (token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected Signing Method")
			}
			return secretKey, nil
		})

		claims := token.Claims.(jwt.MapClaims)
		if err != nil || !token.Valid || claims["iss"] != issuer {
			SendStatus(http.StatusUnauthorized, c)
		} else {
			c.Next()
		}
	} else {
		SendStatus(http.StatusBadRequest, c)
	}
}

func AuthHandler(c *gin.Context) {
	if c.Query("key") == ppsKey {
		claims := jwt.MapClaims{
			"iss": issuer,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Hour * 1),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(secretKey)

		SendString(tokenString, c)
	} else {
		SendStatus(http.StatusUnauthorized, c)
	}
}

func main() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"*"}
	config.AllowHeaders = []string{"*"}
	config.MaxAge = time.Hour * 24
	r.Use(cors.New(config))

	r.GET("/validate/auth", AuthHandler)
	r.PUT("/validate/push", ISAUTH, AuthHandler) // validates the post and pushes it to the pool server
	r.GET("/validate", ISAUTH, GetHandler) // gets a list of posts
	r.PUT("/validate", ISAUTH, PutHandler)  // updates a post on the pool server
	r.POST("/validate", ISAUTH, PostHandler) // creates a post on the pool server
	r.DELETE("/validate", ISAUTH, DeleteHandler) // deletes a post on the pool server

	r.Run(":7000")
}
