package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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


	router.GET("/:offset", func(c *gin.Context) {
		offset := c.Param("offset")
		ctx := context.WithValue(c.Request.Context(), "offset", offset)
		s.HandleConsume(c.Writer, c.Request.WithContext(ctx))
	})

	router.POST("/", func(c *gin.Context) {
		s.handleProduce(c.Writer, c.Request)
	})

	return http.HandlerFunc(router.ServeHTTP)
}

func (s *server) handleProduce(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	type (
		Request struct {
			Record []byte `json:"record"`
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


	offset, appendErr := s.recordRepo.Append(Record{Value: r.Record})
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
	defer r.Body.Close()
	type (
		Response struct {
			Record Record `json:"record"`
		}
	)


	offsetstr, ok  := r.Context().Value("offset").(string)
	offset, encodeErr := strconv.Atoi(offsetstr)
	if encodeErr  != nil || !ok{
		http.Error(w, encodeErr.Error(), http.StatusBadRequest)
		return
	}

	record, readErr := s.recordRepo.Read(uint64(offset))
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
