package s3fs

import (
	"encoding/json"
	"fmt"

	"github.com/johanhenriksson/montai/mount"
)

type Mount struct {
	*mount.Request
	Opt  S3Options `json:"opt"`
	stop chan bool
}

func (mnt *Mount) Path() string {
	return mnt.Request.Path
}

func (mnt *Mount) Unmount() {
	mnt.stop <- true
}

func newMount(req *mount.Request) (*Mount, error) {
	if req.Type != "s3" {
		return nil, fmt.Errorf("not an s3 mount")
	}

	opt := S3Options{}
	if err := json.Unmarshal(req.Opt, &opt); err != nil {
		return nil, mount.ErrParseOptions
	}

	return &Mount{
		Request: req,
		Opt:     opt,
		stop:    make(chan bool),
	}, nil
}
