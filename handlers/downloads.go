package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kalpitpant/file-download-manager/data"
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
	down_request := &data.DownloadRequest{}
	err := down_request.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "Unable to Marshal json", http.StatusBadRequest)
		return
	}

	down_id := uuid.New().String()
	if down_request.TYPE == "concurrent" {
		go downloader(down_request, down_id)
	} else {
		downloader(down_request, down_id)
	}

	rw.Write([]byte(down_id))
}

func downloader(dr *data.DownloadRequest, id string) {
	var wg sync.WaitGroup
	down := &data.Download{}
	down.ID = id
	files := make(map[string]string)
	down.STARTTIME = time.Now().String()
	sem := make(chan int, THROTLLER)

	for index, url := range dr.URLS {
		if dr.TYPE == "concurrent" {
			sem <- index
			wg.Add(1)
			go downloadImageFromURLConcurrently(&wg, url, files, index, sem)
		} else {
			downloadImageFromURLSerial(url, files)
		}

	}
	wg.Wait()
	down.ENDTIME = time.Now().String()
	down.STATUS = "SUCCESSFULL"
	down.DOWNLOADTYPE = dr.TYPE
	down.FILE = files
	data.InsertDownload(down)
	fmt.Println("Success in downloading the images.")
}

func downloadImageFromURLConcurrently(wg *sync.WaitGroup, url string, files map[string]string, index int, sem chan int) {
	image_id := uuid.New().String()
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
	time.Sleep(10 * time.Second)
	image_id := uuid.New().String()
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
