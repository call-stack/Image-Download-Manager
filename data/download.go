package data

import (
	"encoding/json"
	"errors"
	"io"
)

type Download struct {
	TYPE       string
	URLS       []string
	ID         string
	start_time string
	end_time   string
}

type Downloads []*Download

func NewDowload() *Download {
	return &Download{}
}

func GetDownloads(downloadID string) (Download, error) {
	result := Download{}
	for _, down := range downloadList {
		if down.ID == downloadID {
			return *down, nil
		}

	}

	return result, errors.New("Not Found")

}

func (d *Download) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(d)
}

func (p *Download) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func InsertDownload(d *Download) {
	downloadList = append(downloadList, d)
}

var downloadList = []*Download{}
