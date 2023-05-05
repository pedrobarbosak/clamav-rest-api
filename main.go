package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"clamav-api/clamd"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := clamd.CheckConnection(); err != nil {
		log.Fatalf("failed to connect to clamd: %s\n", err)
	}

	http.HandleFunc("/ping", pong)
	http.HandleFunc("/scan", scan)

	log.Println("Listening at:", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func pong(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
}

func scan(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	result, err := clamd.ScanFile(file, fileHeader)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	json, err := result.JSON()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(result.Code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
