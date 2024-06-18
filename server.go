package main

import (
	"github.com/carreter/pavlok-go"
	"github.com/op/go-logging"
	"net/http"
	"os"
)

const (
	DefaultCode = "1234"
	DefaultPort = "9000"
	LogFormat   = `%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`
)

var log = logging.MustGetLogger("pavlok")
var pavlokClient *pavlok.Client
var authCode string

func main() {
	// Set up nicer logging.
	logBackend := logging.NewLogBackend(os.Stderr, "zap_me: ", 0)
	logFormatter := logging.NewBackendFormatter(logBackend, logging.MustStringFormatter(LogFormat))
	logging.AddModuleLevel(logBackend)
	logging.SetBackend(logBackend, logFormatter)

	// Get config values.
	apikey, exists := os.LookupEnv("PAVLOK_API_KEY")
	if !exists {
		log.Fatal("PAVLOK_API_KEY environment variable not set")
		os.Exit(1)
	}
	authCode, exists = os.LookupEnv("AUTH_CODE")
	if !exists {
		log.Warningf("AUTH_CODE environment variable not set, using default %s", DefaultCode)
	}
	port, exists := os.LookupEnv("ZAP_ME_PORT")
	if !exists {
		log.Warningf("ZAP_ME_PORT environment variable not set, using default %s", DefaultPort)
		port = DefaultPort
	}

	// Set up pavlok client.
	pavlokClient = pavlok.NewClient(apikey)

	// Register handlers.
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	mux.HandleFunc("POST /sendstimulus", StimulusHandler)

	// Start up the server.
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalf("could not start HTTP server on port %s: %v", port, err)
		os.Exit(1)
	}

	os.Exit(0)
}
