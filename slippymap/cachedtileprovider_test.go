package slippymap

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestNewCachedTileProvider(t *testing.T) {

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
}
