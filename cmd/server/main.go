package main

import (
	"dsingo/internal/server"
	"fmt"
	"net/http"
)

func main() {
	recordLogger := server.NewLog()
<<<<<<< HEAD
	http.ListenAndServe(":8080", server.New(recordLogger).Handle())
=======
	srvr := server.New(recordLogger)
	errChan := make(chan error, 0)

	go func() {
		errChan <- http.ListenAndServe(":8080", srvr.Handle())
	}()

	fmt.Println(<- errChan)
>>>>>>> 4c42d03e2bebf5e1b4ff2a07d3d3edba3dfa1a21
}
