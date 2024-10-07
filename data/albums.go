package data

import (
	"encoding/json"
	"io"
)

type Album struct {
	Name      string `json:"Album"` // Rename to Album in json
	Artist    string
	Playtime  int
	Tracks    int
	PrivateId int `json:"-"` // Ignore in Marshalling
}

type Albums []*Album

func GetAlbums() Albums {
	return albumList
}

func AddAlbum(a *Album) {
	albumList = append(albumList, a)
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

var albumList = []*Album{
	{
		Name:      "There is a hell, believe me I've seen it",
		Artist:    "Bring me the Horizon",
		Playtime:  100,
		Tracks:    10,
		PrivateId: 1,
	},
	{
		Name:      "Warmer Weather",
		Artist:    "Hot Mulligan",
		Playtime:  100,
		Tracks:    10,
		PrivateId: 2,
	},
}
