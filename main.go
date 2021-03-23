package main

import (
	"fmt"
	"net/http"

	"github.com/kalpitpant/file-download-manager/handlers"
)

func main() {
	fmt.Println("Webserver running....")
	hh := handlers.NewHealth()
	dw := handlers.NewDowload()
	sm := http.NewServeMux()

	sm.Handle("/health", hh)

	// sm.Handle("/downloads", dw)
	sm.Handle("/downloads/*", dw)

	s := &http.Server{
		Addr:    ":8081",
		Handler: sm,
	}
	s.ListenAndServe()
}
