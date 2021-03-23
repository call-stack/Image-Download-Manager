package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kalpitpant/file-download-manager/data"
)

type Download struct{}

func NewDowload() *Download {
	return &Download{}
}

// func (h *Download) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodGet {
// 		fmt.Println("Reponse object", r.URL.Path)
// 		// h.getDownloads(rw, r)
// 		return
// 	}

// 	if r.Method == http.MethodPost {
// 		h.downloadImages(rw, r)
// 		return
// 	}

// 	rw.WriteHeader(http.StatusMethodNotAllowed)
// }

func (h *Download) GetDownloads(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	downloadID := vars["downloadID"]
	fmt.Printf(downloadID)
	dw, er := data.GetDownloads(downloadID)
	if er != nil {
		http.Error(rw, "Data not found", http.StatusInternalServerError)
		return
	}
	err := dw.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}

}

func (h *Download) DownloadImages(rw http.ResponseWriter, r *http.Request) {
	down := &data.Download{}
	err := down.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "Unable to Marshal json", http.StatusBadRequest)
		return
	}
	down_id := uuid.New().String()
	down.TYPE = "Serial"
	down.ID = down_id
	for _, url := range down.URLS {
		id := uuid.New()

		out_image := id.String() + ".jpg"
		response, e := http.Get(url)
		if e != nil {
			log.Fatal(e)
		}
		defer response.Body.Close()
		file, err := os.Create(out_image)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		_, err = io.Copy(file, response.Body)
		if err != nil {
			log.Fatal(err)
		}

	}

	data.InsertDownload(down)
	rw.Write([]byte(down_id))
}
