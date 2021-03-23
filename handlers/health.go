package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Health struct{}

func NewHealth() *Health {
	return &Health{}
}

func (h *Health) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Some Error", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(rw, "Hello %s\n", d)
}
