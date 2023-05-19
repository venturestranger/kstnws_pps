package main

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

var httpClient http.Client = http.Client{}

func Logd(data ...any) {
	log.Println("|", data)
}

func switchCORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Access-Control-Max-Age", "86400")
}

func SendStatus(status int, c *gin.Context) {
	switchCORS(c)
	c.Status(status)
}

func SendString(value string, c *gin.Context) {
	switchCORS(c)
	c.String(200, value)
}

func SendData(jsonByte []Post, c *gin.Context) {
	switchCORS(c)
	c.JSON(200, jsonByte)
}
