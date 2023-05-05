package clamd

import (
	"encoding/json"
	"fmt"
)

type Result struct {
	Code int

	filename    string
	size        int64
	contentType string

	status      string
	hash        string
	description string
}

func (r Result) String() string {
	return fmt.Sprintf("[-] %s (%d) [%s]: '%s' %s %s", r.filename, r.size, r.contentType, r.status, r.hash, r.description)
}

func (r Result) JSON() ([]byte, error) {
	return json.Marshal(r)
}
