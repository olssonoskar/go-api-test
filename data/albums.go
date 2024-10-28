package data

import (
	"encoding/json"
	"fmt"
	"io"
	"slices"
)

type Album struct {
	Name     string `json:"Album"` // Rename to Album in json
	Artist   string
	Playtime int
	Tracks   int
	ID       int `json:"-"` // Ignore in Marshalling
}

type Albums []*Album

func GetAlbums() Albums {
	return albumList
}

func GetAlbum(id int) Albums {
	album, _ := find(id)
	return album
}

func AddAlbum(a *Album) {
	a.ID = albumList[len(albumList)-1].ID + 1
	albumList = append(albumList, a)
}

func UpdateAlbum(id int, a *Album) error {
	fmt.Println("Updating Album ", id)
	dba, pos := find(id)
	if pos == -1 {
		return fmt.Errorf("no such album")
	}
	a.ID = dba[0].ID
	albumList[pos] = a
	return nil
}

func DeleteAlbum(id int) error {
	_, pos := find(id)
	if pos == -1 {
		return fmt.Errorf("no such album")
	}
	albumList = slices.Delete(albumList, pos, pos+1)
	return nil
}

// *Ablums here is the Reciever of the function
func (a *Albums) ToJson(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(a)
}

func (a *Album) FromJson(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(a)
}

func find(id int) (Albums, int) {
	for idx, a := range albumList {
		if a.ID == id {
			return Albums{a}, idx
		}
	}
	return Albums{&Album{}}, -1
}

var albumList = []*Album{
	{
		Name:     "There is a hell, believe me I've seen it",
		Artist:   "Bring me the Horizon",
		Playtime: 100,
		Tracks:   10,
		ID:       1,
	},
	{
		Name:     "Warmer Weather",
		Artist:   "Hot Mulligan",
		Playtime: 100,
		Tracks:   10,
		ID:       2,
	},
}
