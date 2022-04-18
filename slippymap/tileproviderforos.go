package slippymap

import (
	"fmt"
	"log"
	"os"
	"path"
	"pw_slippymap/localdata"
	"runtime"
)

// If we are running in WASM/JS, then the browser does all relevant tile caching for us.
// If running in desktop app mode, we need to cache the tiles ourselves
func TileProviderForOS() (TileProvider, error) {
	if runtime.GOOS == "js" {
		return &OSMTileProvider{}, nil
	}

	// try to get user home dir (for map cache)
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("User home dir:", userHomeDir)

	// create directory structure $HOME/.plane.watch if it doesn't exist
	pathRoot := path.Join(userHomeDir, ".plane.watch")
	err = localdata.SetupRoot(pathRoot)
	if err != nil {
		log.Fatal(err)
	}

	// create directory structure $HOME/.plane.watch/tilecache if it doesn't exist
	pathTileCache := path.Join(pathRoot, "tilecache")
	err = localdata.SetupTileCache(pathTileCache)
	if err != nil {
		log.Fatal(err)
	}

	return NewCachedTileProvider(pathTileCache, &OSMTileProvider{}), nil
}
