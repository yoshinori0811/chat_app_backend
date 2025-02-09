package middleware

import (
	"net/http"
)

type MethodHandler struct {
	Get    http.HandlerFunc
	Post   http.HandlerFunc
	Put    http.HandlerFunc
	Patch  http.HandlerFunc
	Delete http.HandlerFunc
}

func (mh *MethodHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if mh.Get != nil {
			mh.Get(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case http.MethodPost:
		if mh.Post != nil {
			mh.Post(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case http.MethodPut:
		if mh.Put != nil {
			mh.Put(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case http.MethodPatch:
		if mh.Patch != nil {
			mh.Patch(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case http.MethodDelete:
		if mh.Delete != nil {
			mh.Delete(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
