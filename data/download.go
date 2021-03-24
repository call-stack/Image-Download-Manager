package data

import (
	"encoding/json"
	"errors"
	"io"
)

type DownloadRequest struct {
	TYPE string   `json:"type"`
	URLS []string `json:"urls"`
}

type Download struct {
	ID           string            `json:"id"`
	STARTTIME    string            `json:"start_time"`
	ENDTIME      string            `json:"end_time"`
	STATUS       string            `json:"status"`
	DOWNLOADTYPE string            `json:"download_type"`
	FILE         map[string]string `json:"files"`
}

type Downloads []*Download

func NewDowload() *Download {
	return &Download{}
}

func NewDowloadRequest() *DownloadRequest {
	return &DownloadRequest{}
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

func (d *DownloadRequest) FromJSON(r io.Reader) error {
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
