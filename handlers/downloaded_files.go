package handlers

import (
	"fmt"
	"net/http"

	"github.com/kalpitpant/file-download-manager/data"
)

type DownloadedFiles struct{}

func NewDownloadedFiles() *DownloadedFiles {
	return &DownloadedFiles{}
}

func (df *DownloadedFiles) GetDownloads(rw http.ResponseWriter, r *http.Request) {
	downloads := data.GetAllDownload()

	return_string := "<pre>\n"

	for _, d := range downloads {
		return_string += "<a>"
		return_string += d.ID
		return_string += "</a>\n"
	}

	return_string += "<pre>"

	fmt.Fprintf(rw, return_string)
}
