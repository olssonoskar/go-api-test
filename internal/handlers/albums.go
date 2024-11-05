package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/olsson/my-api/internal/data"
)

type Albums struct {
	l *log.Logger
}

func NewAlbums(l *log.Logger) *Albums {
	return &Albums{l}
}

func (a *Albums) GetAlbum(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(rw, "Invalid ID", http.StatusBadRequest)
		return
	}
	albums := data.GetAlbum(id)
	err = albums.ToJson(rw)
	if err != nil {
		http.Error(rw, "Failed to marshall Album data", http.StatusInternalServerError)
	}
}

func (a *Albums) GetAll(rw http.ResponseWriter, r *http.Request) {
	albums := data.GetAlbums()
	err := albums.ToJson(rw)
	if err != nil {
		http.Error(rw, "Failed to marshall Album data", http.StatusInternalServerError)
	}
}

func (a *Albums) Create(rw http.ResponseWriter, r *http.Request) {
	album := r.Context().Value(KeyAlbum{}).(*data.Album)
	data.AddAlbum(album)
}

func (a *Albums) Update(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(rw, "Invalid ID", http.StatusBadRequest)
		return
	}
	album := r.Context().Value(KeyAlbum{}).(*data.Album)
	err = data.UpdateAlbum(id, album)
	if err != nil {
		http.Error(rw, "Album not found", http.StatusNotFound)
		return
	}
}

func (a *Albums) Delete(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(rw, "Invalid ID", http.StatusBadRequest)
		return
	}
	err = data.DeleteAlbum(id)
	if err != nil {
		http.Error(rw, "Album not found", http.StatusNotFound)
		return
	}
}

type KeyAlbum struct{}

func (a *Albums) ValidateAlbum(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		album := &data.Album{}
		err := album.FromJson(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshall Album", http.StatusBadRequest)
		}
		ctx := context.WithValue(r.Context(), KeyAlbum{}, album)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
