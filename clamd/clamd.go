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
	return clamd.NewClamd("tcp://localhost:3310")
}

func CheckConnection() error {
	clam := New()

	var err error
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second * 3)

		err = clam.Ping()
		if err != nil {
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

func ScanFile(file multipart.File, header *multipart.FileHeader) (*Result, error) {
	clam := New()

	resultCh, err := clam.ScanStream(file, nil)
	if err != nil {
		log.Printf("failed to scan: %s (%d) %s - %s\n", header.Filename, header.Size, header.Header, err)
		return nil, err
	}

	response := <-resultCh

	result := &Result{
		ContentType: header.Header.Get("Content-Type"),
		Filename:    header.Filename,
		Size:        header.Size,

		Status:      response.Status,
		Hash:        response.Hash,
		Description: response.Description,
	}

	switch response.Status {
	case clamd.RES_OK:
		result.Code = http.StatusOK

	case clamd.RES_FOUND:
		result.Code = http.StatusNotAcceptable

	case clamd.RES_ERROR:
		result.Code = http.StatusBadRequest

	case clamd.RES_PARSE_ERROR:
		result.Code = http.StatusPreconditionFailed

	default:
		log.Println(result)
		return nil, errors.New(fmt.Sprintf("unrecognized result status: %v", response))
	}

	log.Println(result)
	return result, nil
}
