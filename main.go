package main

import (
	"fmt"
	"time"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CORS(c *gin.Context) {
	if c.Request.Method == "OPTIONS" {
		SendStatus(http.StatusOK, c)
	} else {
		c.Next()
	}
}

func ISAUTH(c *gin.Context) {
	data := c.GetHeader("Authorization")
	values := strings.Split(data, " ")

	if values[0] == "Bearer" {
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

	r.Use(CORS)
	r.Use(ISAUTH)

	r.GET("/validate/auth", AuthHandler)
	r.GET("/validate", GetHandler) // gets a list of posts
	r.PUT("/validate", PushHandler)  // verifies a post and pushes the post to the api
	r.POST("/validate", PostHandler) // creates a post on the pull server
	
	r.Run(":4000")
}
