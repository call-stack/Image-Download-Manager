package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kalpitpant/file-download-manager/handlers"
)

func main() {
	fmt.Println("Webserver running....")
	hh := handlers.NewHealth()
	dw := handlers.NewDowload()
	sm := mux.NewRouter()

	getRouter := sm.Methods("GET").Subrouter()
	postRouter := sm.Methods("POST").Subrouter()

	getRouter.HandleFunc("/downloads/{downloadID}", dw.GetDownloads)
	postRouter.HandleFunc("/downloads", dw.DownloadImages)
	getRouter.Handle("/health", hh)

	s := &http.Server{
		Addr:    ":8081",
		Handler: sm,
	}
	s.ListenAndServe()
}
