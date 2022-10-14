package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-gin/gin"
)

type (
	server struct {
		recordRepo *Log
	}
)

func New(recordrepo *Log) *server {
	return &server{recordRepo: recordrepo}
}

func (s *server) Handle() http.Handler {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		s.HandleConsume(c.Writer, c.Request)
	})

	router.POST("/", func(c *gin.Context) {
		s.handleProduce(c.Writer, c.Request)
	})

	return http.HandlerFunc(router.ServeHTTP)
}

func (s *server) handleProduce(w http.ResponseWriter, req *http.Request) {
	type (
		Request struct {
			Record Record `json:"record"`
		}

		Response struct {
			Offset uint64 `json:"offset"`
		}
	)

	var r Request
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	if decoderErr := decoder.Decode(&r); decoderErr != nil {
		http.Error(w, decoderErr.Error(), http.StatusBadRequest)
		return
	}

	offset, appendErr := s.recordRepo.Append(r.Record)
	if appendErr != nil {
		http.Error(w, appendErr.Error(), http.StatusInternalServerError)
		return
	}

	res := Response{Offset: offset}
	if encoderErr := json.NewEncoder(w).Encode(res); encoderErr != nil {
		http.Error(w, encoderErr.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *server) HandleConsume(w http.ResponseWriter, r *http.Request) {
	type (
		Request struct {
			Offset uint64 `json:"offset"`
		}

		Response struct {
			Record Record `json:"record"`
		}
	)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var req Request
	if decoderErr := decoder.Decode(&req); decoderErr != nil {
		http.Error(w, decoderErr.Error(), http.StatusBadRequest)
		return
	}

	record, readErr := s.recordRepo.Read(req.Offset)
	if readErr != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(readErr, ErrOffsetNotFound) {
			statusCode = http.StatusNotFound
		}

		http.Error(w, readErr.Error(), statusCode)
		return
	}

	resp := Response{Record: record}
	if encoderErr := json.NewEncoder(w).Encode(resp); encoderErr != nil {
		http.Error(w, encoderErr.Error(), http.StatusInternalServerError)
		return
	}

}
