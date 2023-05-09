package clamd

import (
	"fmt"
)

type Result struct {
	Code int `json:"code"`

	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`

	Status      string `json:"status"`
	Hash        string `json:"hash"`
	Description string `json:"description"`
}

func (r Result) String() string {
	return fmt.Sprintf("[-] %s [%s]: '%s' %s %s", r.Filename, r.ContentType, r.Status, r.Hash, r.Description)
}

func (r Result) JSON() any {
	return fmt.Sprintf(`
		{ 
			"code": "%s",
			"filename":	"%s",
			"contentType": "%s",
			"status": "%s",
			"hash": "%s",
			"description": "%s"
		}`, r.Code, r.Filename, r.ContentType, r.Status, r.Hash, r.Description)
}
