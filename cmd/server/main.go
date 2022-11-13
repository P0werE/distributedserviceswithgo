package main

import (
	"dsingo/internal/server"
	"fmt"
	"net/http"
)

func main() {
	recordLogger := server.NewLog()
	srvr := server.New(recordLogger)
	errChan := make(chan error, 0)

	go func() {
		errChan <- http.ListenAndServe(":8080", srvr.Handle())
	}()

	fmt.Println(<- errChan)
}
