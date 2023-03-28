package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/rafli-lutfi/kanban-app-mongodb/src/models"
)

func Get(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respone := map[string]any{
				"error": models.ErrMethodNotAllowed.Error(),
			}

			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(respone)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func Post(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respone := map[string]any{
				"error": models.ErrMethodNotAllowed.Error(),
			}

			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(respone)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func Put(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			respone := map[string]any{
				"error": models.ErrMethodNotAllowed.Error(),
			}

			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(respone)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func Delete(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			respone := map[string]any{
				"error": models.ErrMethodNotAllowed.Error(),
			}

			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(respone)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
