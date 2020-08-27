package middleware

import (
	"encoding/json"
	"html"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
)

var JsonLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: &JSONLog{}, NoColor: false})

func Logger(next http.Handler) http.Handler {
	return JsonLogger(next)
}

type JSONLog struct {
	timeStamp time.Time `json:"timestamp"`
	url       string    `json:"url"`
	method    string    `json:"method"`
	status    string    `json:"status"`
	message   string    `json:"message"`
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

type Logger3 struct {
	Timestamp string `json:"timestamp"`
	Url       string `json:"url"`
	Method    string `json:"method"`
	Agent     string `json"user_agent"`
}

func (jl *JSONLogger2) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	current_time := time.Now()
	b, _ := json.Marshal(Logger3{
		Timestamp: current_time.Format(time.RFC1123),
		Method:    r.Method,
		Url:       html.EscapeString(r.URL.Path),
		Agent:     r.UserAgent(),
	})
	log.Println(string(b))
	jl.next.ServeHTTP(w, r)
}
