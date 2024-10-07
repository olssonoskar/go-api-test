package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Simple struct {
	l *log.Logger
}

func NewSimple(l *log.Logger) *Simple {
	return &Simple{l}
}

func (h *Simple) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Log from handler")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Ohno!", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(rw, "Updated just now, welcome %s", body)
}
