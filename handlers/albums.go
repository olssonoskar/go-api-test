package handlers

import (
	"log"
	"net/http"
	"regexp"
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
		if strings.HasSuffix(r.URL.Path, "/") {
			a.getAlbums(rw)
		} else {
			arg := getArg(r.URL.Path)
			if arg != -1 {
				http.Error(rw, "Invalid path", http.StatusBadRequest)
			}
			a.getAlbum(rw, arg)
		}
		return
	}
	if r.Method == http.MethodPost {
		a.postAlbum(rw, r)
		return
	}
	if r.Method == http.MethodPut {
		arg := getArg(r.URL.Path)
		a.update(arg, rw, r)
		return
	}
	if r.Method == http.MethodDelete {
		arg := getArg(r.URL.Path)
		a.delete(arg, rw)
		return
	}
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (a *Albums) getAlbum(rw http.ResponseWriter, id int) {
	albums := data.GetAlbum(id)
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
		return
	}
	data.AddAlbum(album)
}

func (a *Albums) update(id int, rw http.ResponseWriter, r *http.Request) {
	album := &data.Album{}
	err := album.FromJson(r.Body)
	if err != nil {
		http.Error(rw, "Failed to read Album data", http.StatusBadRequest)
	}
	err = data.UpdateAlbum(id, album)
	if err != nil {
		http.Error(rw, "Album not found", http.StatusNotFound)
		return
	}
}

func (a *Albums) delete(id int, rw http.ResponseWriter) {
	err := data.DeleteAlbum(id)
	if err != nil {
		http.Error(rw, "Album not found", http.StatusNotFound)
		return
	}
}

func getArg(uri string) int {
	reg := regexp.MustCompile(`/([0-9]+)`)
	match := reg.FindAllStringSubmatch(uri, -1)
	if len(match) != 1 || len(match[0]) < 1 {
		return -1
	}
	res, err := strconv.Atoi(match[0][1])
	if err != nil {
		return -1
	}
	if strings.HasSuffix(uri, "/"+match[0][1]) {
		return res
	} else {
		return 0
	}
}
