package mount

import (
	"log"
	"os"
	"os/exec"
)

func Unmount(path string) {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		// the mount point already exists!
		// attempt to unmount first, just in case
		umount := exec.Command("umount", path)
		if err := umount.Run(); err == nil {
			log.Println("unmounted", path)
		}

		if err := os.Remove(path); err == nil {
			log.Println("removed mount point", path)
		}
	}
}

func CreatePoint(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// the mount point does not exist
		// create the folder first
		if err := os.Mkdir(path, os.ModeDir); err != nil {
			return ErrPointFailed
		}
		log.Println("created mount point", path)
		return nil
	}
	return ErrMountExists
}
