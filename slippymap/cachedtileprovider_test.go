package slippymap

import (
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"
)

func TestNewCachedTileProvider(t *testing.T) {

	// skip test if webassembly
	if runtime.GOOS == "js" {
		t.SkipNow()
	}

	// create temp dir
	dir, err := ioutil.TempDir(os.TempDir(), "pw_slippymap_TestNewCachedTileProvider")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir)

	// get cached tile provider
	ctp := NewCachedTileProvider(dir, &OSMTileProvider{})

	// request a tile
	tilePath, err := ctp.GetTileAddress(1, 2, 3)
	if err != nil {
		t.Error(err)
	}

	// prepare requested path
	expectedPath := path.Join(dir, "1_2_3.png")

	// check for success
	if tilePath != expectedPath {
		t.Errorf("Expected: %s, got: %s", expectedPath, tilePath)
	}

	// prepare "faulty" OSMTileProvider
	faultyOSMTileProvider := OSMTileProvider{
		osm_url_prefix: 4,
	}

	// get cached tile provider
	ctp = NewCachedTileProvider(dir, &faultyOSMTileProvider)

	// request a tile
	tilePath, err = ctp.GetTileAddress(2, 3, 4)
	if err != nil {
		// test passes
	} else {
		t.Errorf("Expected an error, got none")
	}

}
