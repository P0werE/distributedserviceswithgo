package main

import (
	"distributed-system-in-go/internal/server"
	"net/http"
)

func main() {
	recordLogger := server.NewLog()
	http.ListenAndServe(":8080", server.New(recordLogger).Handle())
}
