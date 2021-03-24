package handlers

import (
	"net/http"
)

type Health struct{}

func NewHealth() *Health {
	return &Health{}
}

func (h *Health) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

}
