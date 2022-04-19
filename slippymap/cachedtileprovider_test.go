package slippymap

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCachedTileProvider(t *testing.T) {

	// skip test if webassembly
	if runtime.GOOS == "js" {
		t.SkipNow()
	}

	// create temp dir
	dir, err := ioutil.TempDir(os.TempDir(), "pw_slippymap_TestNewCachedTileProvider")
	require.NoError(t, err, "Could not create temp dir")
	defer os.RemoveAll(dir)

	// get cached tile provider
	ctp := NewCachedTileProvider(dir, &OSMTileProvider{})

	// request a tile
	tilePath, err := ctp.GetTileAddress(1, 2, 3)
	require.NoError(t, err, "Error returned from GetTileAddress")

	// check for success
	expectedPath := path.Join(dir, "1_2_3.png")
	assert.Equal(t, expectedPath, tilePath)

	// get cached tile provider & fake an error
	ctp = NewCachedTileProvider(dir, &FaultyTileProvider{})

	// request a tile (should return error)
	tilePath, err = ctp.GetTileAddress(2, 3, 4)
	require.Error(t, err)
}

type FaultyTileProvider struct{}

func (FaultyTileProvider) GetTileAddress(int, int, int) (string, error) {
	// fake an error
	return "", errors.New("oh no an error")
}
