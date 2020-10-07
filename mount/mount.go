package mount

import (
	"encoding/json"
)

type Request struct {
	Type string          `json:"type"`
	Name string          `json:"name"`
	Path string          `json:"path"`
	Opt  json.RawMessage `json:"opt"`
}

type Provider interface {
	Type() string
	Mount(*Request) (Mount, error)
}

type Mount interface {
	Path() string
	Unmount()
}
