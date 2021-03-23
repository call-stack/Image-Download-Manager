package data

import (
	"encoding/json"
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

func GetDownloads() Downloads {
	return downloadList
}

func (d *Download) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(d)
}

func (p *Downloads) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func InsertDownload(d *Download) {
	downloadList = append(downloadList, d)
}

var downloadList = []*Download{}
