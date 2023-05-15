package main

import (
	"fmt"
	"io"
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		if part.FileName() == "" {
			_ = part.Close()
			continue
		}

		result, err := clamd.ScanFile(part)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		_ = part.Close()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(result.Code)
		_, _ = fmt.Fprint(w, result.JSON())
		return
	}

	http.Error(w, "no file", http.StatusBadRequest)
}
