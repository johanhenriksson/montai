package mount

import (
	"fmt"
)

var ErrParseOptions = fmt.Errorf("failed to parse options")
var ErrMountExists = fmt.Errorf("mount point already exists")
var ErrPointFailed = fmt.Errorf("failed to create mount point")
var ErrMountFailed = fmt.Errorf("mount failed")
var ErrMountStopped = fmt.Errorf("mount stopped")
var ErrMountNotFound = fmt.Errorf("not found")
var ErrUnsupportedProvider = fmt.Errorf("unsupported provider")
