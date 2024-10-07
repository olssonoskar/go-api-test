package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/olsson/my-api/data"
)

type Albums struct {
	l *log.Logger
}

func NewAlbums(l *log.Logger) *Albums {
	return &Albums{l}
}

func (a *Albums) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		arg := getArg(r.RequestURI)
		if arg != -1 {
			a.getAlbum(rw, arg)
		} else {
			a.getAlbums(rw)
		}
		return
	}
	if r.Method == http.MethodPost {
		a.postAlbum(rw, r)
		return
	}
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (a *Albums) getAlbum(rw http.ResponseWriter, id int) {
	albums := find(id)
	err := albums.ToJson(rw)
	if err != nil {
		http.Error(rw, "Failed to marshall Album data", http.StatusInternalServerError)
	}
}

func (a *Albums) getAlbums(rw http.ResponseWriter) {
	albums := data.GetAlbums()
	err := albums.ToJson(rw)
	if err != nil {
		http.Error(rw, "Failed to marshall Album data", http.StatusInternalServerError)
	}
}

func (a *Albums) postAlbum(rw http.ResponseWriter, r *http.Request) {
	album := &data.Album{}
	err := album.FromJson(r.Body)
	if err != nil {
		http.Error(rw, "Failed to read Album data", http.StatusBadRequest)
	}
	album.PrivateId = len(data.GetAlbums()) + 1
	data.AddAlbum(album)
	rw.WriteHeader(http.StatusOK)
}

// Simple/Dumb extraction of arg and finding based on the internal privateId
func getArg(uri string) int {
	split := strings.Split(uri, "/")
	arg := split[len(split)-1]
	res, err := strconv.Atoi(arg)
	if err != nil {
		return -1
	}
	return res
}

func find(id int) data.Albums {
	for _, alb := range data.GetAlbums() {
		if alb.PrivateId == id {
			return data.Albums{alb}
		}
	}
	return data.Albums{&data.Album{}}
}
