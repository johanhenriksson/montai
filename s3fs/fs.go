package s3fs

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/johanhenriksson/montai/mount"
)

// ErrCredentialsMissing is returned when S3 credentials are not properly provided
var ErrCredentialsMissing = fmt.Errorf("credentials missing")

type T struct{}

// New S3 file system
func New() *T {
	return &T{}
}

func (fs *T) Type() string {
	return "s3"
}

// Mount a bucket
func (s3 *T) Mount(req *mount.Request) (mount.Mount, error) {
	mnt, err := newMount(req)
	if err != nil {
		return nil, err
	}

	// unmount existing & create mount point
	mount.Unmount(req.Path)
	if err := mount.CreatePoint(req.Path); err != nil {
		return nil, err
	}

	// write s3fs credentials
	if err := mnt.Opt.writeCredentials(); err != nil {
		return nil, err
	}

	// run s3fs process
	cmd := exec.Command("s3fs", mnt.Opt.Args(req.Path)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Println("s3fs: mount failed to start:", err)
		return nil, mount.ErrMountFailed
	}

	// watch process for errors
	go func() {
		if err := cmd.Wait(); err != nil {
			log.Println("s3fs:", req.Path, "failed with error", err)
		} else {
			log.Println("s3fs: stopped", req.Path)
		}

		// cleanup
		mount.Unmount(req.Path)
	}()

	go func() {
		<-mnt.stop
		log.Println("s3fs: stop mount", req.Path)
		cmd.Process.Kill()
	}()

	log.Println("s3fs: mounted bucket", mnt.Opt.Bucket, "at", mnt.Path())

	return mnt, nil
}
