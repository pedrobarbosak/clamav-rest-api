package clamd

import (
	"fmt"
)

type Result struct {
	Code int `json:"code"`

	Filename    string `json:"filename"`
	Size        int64  `json:"size"`
	ContentType string `json:"contentType"`

	Status      string `json:"status"`
	Hash        string `json:"hash"`
	Description string `json:"description"`
}

func (r Result) String() string {
	return fmt.Sprintf("[-] %s (%d) [%s]: '%s' %s %s", r.Filename, r.Size, r.ContentType, r.Status, r.Hash, r.Description)
}
