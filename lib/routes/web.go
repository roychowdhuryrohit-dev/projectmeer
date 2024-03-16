package routes

import (
	"encoding/json"
	"net/http"
	"unicode/utf8"

	"github.com/roychowdhuryrohit-dev/projectmeer/lib/algos"
)


func InsertText(fg *algos.FugueMax[rune]) http.HandlerFunc {
	type InsertTextRequest struct {
		Index int    `json:"index"`
		Value string `json:"value"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var body InsertTextRequest
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		bodyChar := []rune(body.Value)[0]
		if !utf8.ValidRune(bodyChar) {
			http.Error(w, "invalid unicode character", http.StatusBadRequest)
			return
		}
		err = fg.Insert(body.Index, bodyChar)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func DeleteText[T any](fg *algos.FugueMax[T]) http.HandlerFunc {
	type DeleteTextRequest struct {
		Index int `json:"index"`
		Count int `json:"count"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var body DeleteTextRequest
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = fg.Delete(body.Index, body.Count)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func GetText(fg *algos.FugueMax[rune]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := string(fg.Values())
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(body)
		w.WriteHeader(http.StatusOK)
	}
}
