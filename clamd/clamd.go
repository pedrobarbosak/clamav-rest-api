package clamd

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/dutchcoders/go-clamd"
)

func New() *clamd.Clamd {
	return clamd.NewClamd("/var/run/clamav/clamd.sock")
}

func NewTCP() *clamd.Clamd {
	return clamd.NewClamd("tcp://localhost:3310")
}

func CheckConnection() error {
	clam := New()

	var err error
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second * 5)

		if err = clam.Ping(); err != nil {
			log.Println("failed to ping:", err)
			continue
		}

		var result chan *clamd.ScanResult
		result, err = clam.Version()
		if err != nil {
			log.Println("failed to get version:", err)
			continue
		}

		version := <-result
		log.Println("clamd version:", version.Raw)

		return nil
	}

	return err
}

func ScanFile(file *multipart.Part) (*Result, error) {
	clam := New()

	log.Printf("scanning: %s ...\n", file.FileName())

	resultCh, err := clam.ScanStream(file, make(chan bool))
	if err != nil {
		log.Printf("failed to scan: %s (%d) %s - %s\n", file.FileName(), 0, file.Header, err)
		return nil, err
	}

	response := <-resultCh

	result := &Result{
		ContentType: file.Header.Get("Content-Type"),
		Filename:    file.FileName(),

		Status:      response.Status,
		Hash:        response.Hash,
		Description: response.Description,
	}

	defer log.Println(result)

	switch response.Status {
	case clamd.RES_OK:
		result.Code = http.StatusOK

	case clamd.RES_FOUND:
		result.Code = http.StatusNotAcceptable

	case clamd.RES_ERROR:
		result.Code = http.StatusExpectationFailed

	case clamd.RES_PARSE_ERROR:
		result.Code = http.StatusPreconditionFailed

	default:
		return nil, errors.New(fmt.Sprintf("unrecognized result status: %v", response))
	}
  
	return result, nil
}
