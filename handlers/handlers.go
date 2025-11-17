package handlers

import (
	"net/http"
)

type Handle struct {
}

func NewHandler() *Handle {
	return &Handle{}
}

func (q *Handle) ServeForm(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" || r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "static/index.html")
}
