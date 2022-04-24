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

	// skip tests if webassembly
	if runtime.GOOS == "js" {
		t.SkipNow()
	}

	var (
		err      error
		dir      string
		ctp      TileProvider
		tilePath string
	)

	// create temp dir
	dir, err = ioutil.TempDir(os.TempDir(), "pw_slippymap_TestNewCachedTileProvider")
	require.NoError(t, err, "Could not create temp dir")
	defer os.RemoveAll(dir)

	// get cached tile provider
	t.Run("Test NewCachedTileProvider", func(t *testing.T) {
		ctp = NewCachedTileProvider(dir, &OSMTileProvider{})

		// request a tile
		osm := OSMTileID{
			x:    1,
			y:    2,
			zoom: 3,
		}
		tilePath, err = ctp.GetTileAddress(osm)
		require.NoError(t, err, "Error returned from GetTileAddress")

		// check for success
		expectedPath := path.Join(dir, "1_2_3.png")
		assert.Equal(t, expectedPath, tilePath)
	})

	t.Run("Test NewCachedTileProvider Error Handling", func(t *testing.T) {
		// get cached tile provider & fake an error
		ctp = NewCachedTileProvider(dir, &FaultyTileProvider{})

		// request a tile (should return error)
		osm := OSMTileID{
			x:    2,
			y:    3,
			zoom: 4,
		}
		tilePath, err = ctp.GetTileAddress(osm)
		require.Error(t, err)
	})
}

type FaultyTileProvider struct{}

func (FaultyTileProvider) GetTileAddress(OSMTileID) (string, error) {
	// fake an error
	return "", errors.New("oh no an error")
}
