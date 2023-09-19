package main

const (
	issuer string = "domain"
	ppsKey string = "domain"
	apiKey string = "domain"
	apiToken string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDIzLTA3LTE4VDA5OjIxOjI1Ljc4ODQ1NTAzNFoiLCJpYXQiOjE2ODk2Njg0ODUsImlzcyI6ImRvbWFpbiJ9.lDnKz4t_6ApwhNxGAnxf0zBy43ztNFfGhCQj1cpNcKw"
	apiAddr string = "https://tvoykostanay.kz/api"
)

const (
	dsn string = "postgresql://postgres:postgres@localhost/work2"
	redisAddr string = "127.0.0.1:6379"
	redisPsswd string = ""
	redisDB int = 0
)

var (
	secretKey = []byte("domain")
)
