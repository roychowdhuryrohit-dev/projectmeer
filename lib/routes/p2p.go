package routes

import (
	"io"
	"net/http"
	"log"

	"github.com/roychowdhuryrohit-dev/projectmeer/lib/algos"
)

func ReceivePrimitiveHandler(fg *algos.FugueMax[rune]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = fg.Receive(msg)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
