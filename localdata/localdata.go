package localdata

import (
	"log"
	"os"
)

func SetupRoot(pathRoot string) error {

	// set up pathRoot
	if _, err := os.Stat(pathRoot); os.IsNotExist(err) {
		log.Print(pathRoot, " does not exist, creating")
		err := os.Mkdir(pathRoot, 0700)
		if err != nil {
			return err
		}
	} else {
		log.Print(pathRoot, " already exists, not creating")
	}
	return nil
}

func SetupTileCache(pathTileCache string) error {

	// set up pathTileCache
	if _, err := os.Stat(pathTileCache); os.IsNotExist(err) {
		log.Print(pathTileCache, " does not exist, creating")
		err := os.Mkdir(pathTileCache, 0700)
		if err != nil {
			return err
		}
	} else {
		log.Print(pathTileCache, " already exists, not creating")
	}
	return nil

}
