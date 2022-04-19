package localdata

import (
	"io/fs"
	"log"
	"os"
)

func MakeDirIfNotExist(dir string, perm fs.FileMode) error {
	// makes a directory 'dir' if the directory does not already exist
	// the mode bits defined by 'perm' are applied when the dir is created

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Printf("Directory '%s' does not exist, creating", dir)
		err := os.Mkdir(dir, perm)
		if err != nil {
			return err
		}
	} else {
		log.Printf("Directory '%s' already exists, not creating", dir)
	}
	return nil
}
