package middleware

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
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

type Logger3 struct {
	Timestamp string `json:"timestamp"`
	Url       string `json:"url"`
	Method    string `json:"method"`
	Agent     string `json"user_agent"`
}

func (j *JSONLog) Print(v ...interface{}) {
	b, err := json.Marshal(v)
	log.SetFormatter(&log.JSONFormatter{})

	if err != nil {
		log.Error(err)
	} else {
		log.Info(string(b))
	}
}

func Logger2(next http.Handler) http.Handler {
	return &JSONLogger2{next: next}
}

type JSONLogger2 struct {
	next http.Handler
}

func (jl *JSONLogger2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	current_time := time.Now()
	b, _ := json.Marshal(Logger3{
		Timestamp: current_time.Format(time.RFC1123),
		Method:    r.Method,
		Url:       html.EscapeString(r.URL.Path),
		Agent:     r.UserAgent(),
	})

	log.SetFormatter(&log.JSONFormatter{})
	fmt.Println(string(b))
	jl.next.ServeHTTP(w, r)
}

func LoggerInput(next http.Handler) http.Handler {
	return &JSONLogger2{next: next}
}

type JSONLogger3 struct {
	next http.Handler
}

type LoggerStruct struct {
	Message  string `json:"message"`
	Status   string `json:"status"`
	Endpoint string `json:"endpoint"`
	Method   string `json:"method"`
}

func InsertApplicationLog(r *http.Request, status string, message string) {
	b, _ := json.Marshal(LoggerStruct{
		Message:  message,
		Method:   r.Method,
		Status:   status,
		Endpoint: html.EscapeString(r.URL.Path),
	})

	log.SetFormatter(&log.JSONFormatter{})
	fmt.Println(string(b))
}
