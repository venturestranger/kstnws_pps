package main

const (
	issuer string = "domain"
	ppsKey string = "domain"
	apiKey string = "domain"
	apiToken string = ""
	apiAddr string = "http://tvoykostanay.kz/api"
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
