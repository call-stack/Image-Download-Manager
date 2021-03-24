package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kalpitpant/file-download-manager/data"
	"github.com/kalpitpant/file-download-manager/response"
)

type Download struct{}

func NewDowload() *Download {
	return &Download{}
}

const THROTLLER = 20

var m = sync.RWMutex{}

func (h *Download) GetDownloads(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	downloadID := vars["downloadID"]

	dw, er := data.GetDownloads(downloadID)
	if er != nil {
		rw.WriteHeader(http.StatusBadRequest)
		createErrorResponse(rw, "unknow download ID")
		return
	}

	err := dw.ToJSON(rw)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		createErrorResponse(rw, "unable to marshal json")
		return
	}

}

func generateDownloadID() string {
	return uuid.New().String()
}

func createErrorResponse(rw http.ResponseWriter, msg string) {
	error_reponse := response.ErrorResponse{}
	error_reponse.INTERNALCODE = 4002
	error_reponse.MESSAGE = msg
	error_reponse.ToJSON(rw)
}

func createDownloadSuccessResponse(rw http.ResponseWriter, down_id string) {
	succes_reponse := response.DownloadSuccessResponse{}
	succes_reponse.ID = down_id
	succes_reponse.ToJSON(rw)
}

func (h *Download) DownloadImages(rw http.ResponseWriter, r *http.Request) {
	down_request := &data.DownloadRequest{}
	err := down_request.FromJSON(r.Body)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		createErrorResponse(rw, "unable to marshal json")
		return
	}

	down_id := generateDownloadID()
	switch down_type := down_request.TYPE; down_type {
	case "concurrent":
		go downloader(down_request, down_id)
	case "serial":
		downloader(down_request, down_id)
	default:
		rw.WriteHeader(http.StatusBadRequest)
		createErrorResponse(rw, "unknown type of download")
		return
	}

	createDownloadSuccessResponse(rw, down_id)
}

func downloader(dr *data.DownloadRequest, id string) {
	var wg sync.WaitGroup
	down := &data.Download{}
	down.ID = id
	files := make(map[string]string)
	down.STARTTIME = time.Now().String()
	sem := make(chan int, THROTLLER)

	for index, url := range dr.URLS {
		switch down_type := dr.TYPE; down_type {
		case "concurrent":
			sem <- index
			wg.Add(1)
			go downloadImageFromURLConcurrently(&wg, url, files, index, sem)
		default:
			downloadImageFromURLSerial(url, files)
		}

	}
	wg.Wait()
	down.ENDTIME = time.Now().String()
	down.STATUS = "SUCCESSFULL"
	down.DOWNLOADTYPE = dr.TYPE
	down.FILE = files
	data.InsertDownload(down)
}

func downloadImageFromURLConcurrently(wg *sync.WaitGroup, url string, files map[string]string, index int, sem chan int) {
	image_id := generateDownloadID()
	out_image := "images/" + image_id + ".jpg"
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
	m.Lock()
	files[url] = out_image
	m.Unlock()
	<-sem
	wg.Done()

}

func downloadImageFromURLSerial(url string, files map[string]string) {
	image_id := generateDownloadID()
	out_image := "images/" + image_id + ".jpg"
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

	files[url] = out_image
}
