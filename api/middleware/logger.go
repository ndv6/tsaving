package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
)

var JsonLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: &JSONLog{}, NoColor: false})

func Logger(next http.Handler) http.Handler {
	return JsonLogger(next)
}

type JSONLog struct {
}

func (j *JSONLog) Print(v ...interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		log.Println(err, v)
	} else {
		log.Println(string(b))
	}
}

func Logger2(next http.Handler) http.Handler {
	return &JSONLogger2{next: next}
}

type JSONLogger2 struct {
	next http.Handler
}

func (jl *JSONLogger2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

	jl.next.ServeHTTP(w, r)
}
