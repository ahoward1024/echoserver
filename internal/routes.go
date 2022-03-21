package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

type Status struct {
	Code int    `json:"code" yaml:"code"`
	Text string `json:"text" yaml:"text"`
}

func getStatus(r *http.Request) (response Status, err error) {
	vars := mux.Vars(r)
	code := vars["code"]
	status := Status{}

	status.Code, err = strconv.Atoi(code)
	if err != nil {
		log.Error().Msg("Error converting http status code string to integer")
		return status, errors.New("String to int conversion error")
	}
	status.Text = http.StatusText(status.Code)

	return status, nil
}

func handlePlaintext(w http.ResponseWriter, r *http.Request) {
	status, _ := getStatus(r)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status.Code)
	io.WriteString(w, fmt.Sprintf("%d %s", status.Code, status.Text))
}

func handleJson(w http.ResponseWriter, r *http.Request) {
	status, _ := getStatus(r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status.Code)
	json.NewEncoder(w).Encode(status)
}

func handleYaml(w http.ResponseWriter, r *http.Request) {
	status, _ := getStatus(r)
	w.Header().Set("Content-Type", "application/yaml")
	w.WriteHeader(status.Code)
	yaml.NewEncoder(w).Encode(status)
}

func handleHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, fmt.Sprintf("pong"))
}

func handleSync(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 5)
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, fmt.Sprintf("completed"))
}

func handleSyncSleep(w http.ResponseWriter, r *http.Request) {
	sleepTime := mux.Vars(r)["sleep"]
	sleep, err := strconv.Atoi(sleepTime)
	log.Debug().Msg(sleepTime)
	if err != nil {
		w.WriteHeader(600)
		io.WriteString(w, fmt.Sprintf("echoserver error: %v", err))
	} else {
		time.Sleep(time.Second * time.Duration(sleep))
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, fmt.Sprintf("completed"))
	}
}

func handleFavicon(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func AddRoutes(router *mux.Router) {
	router.HandleFunc("/health", handleHealthcheck)
	router.HandleFunc("/ping", handlePing)

	router.HandleFunc("/favicon.ico", handleFavicon)

	router.HandleFunc("/sync", handleSync)
	router.HandleFunc("/sync", handleSyncSleep).Queries("sleep")

	router.HandleFunc("/{code:[0-9]+}", handlePlaintext).Queries("output", "text")
	router.HandleFunc("/{code:[0-9]+}", handleJson).Queries("output", "json")
	router.HandleFunc("/{code:[0-9]+}", handleYaml).Queries("output", "yaml")

	router.HandleFunc("/{code:[0-9]+}", handlePlaintext).Headers("echoserver-output", "text")
	router.HandleFunc("/{code:[0-9]+}", handleJson).Headers("echoserver-output", "json")
	router.HandleFunc("/{code:[0-9]+}", handleYaml).Headers("echoserver-output", "yaml")

	router.HandleFunc("/text/{code:[0-9]+}", handlePlaintext)
	router.HandleFunc("/json/{code:[0-9]+}", handleJson)
	router.HandleFunc("/yaml/{code:[0-9]+}", handleYaml)

	router.HandleFunc("/{code:[0-9]+}", handlePlaintext)
}
