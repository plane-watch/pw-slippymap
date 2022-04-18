package slippymap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTileAddress(t *testing.T) {

	t.Run("GetTileAddress loops through CDNs correctly", func(t *testing.T) {
		// -1 so that the first increment lands us up on CDN A. Doesn't really matter where it starts though
		TileProvider := &OSMTileProvider{osm_url_prefix: -1}

		// check correct a... url is returned
		url, err := TileProvider.GetTileAddress(1, 2, 3)
		require.NoError(t, err)
		assert.Equal(t, "http://a.tile.openstreetmap.org/3/1/2.png", url)

		// check correct b... url is returned
		url, err = TileProvider.GetTileAddress(1, 2, 3)
		require.NoError(t, err)
		assert.Equal(t, "http://b.tile.openstreetmap.org/3/1/2.png", url)

		// check correct c... url is returned
		url, err = TileProvider.GetTileAddress(1, 2, 3)
		require.NoError(t, err)
		assert.Equal(t, "http://c.tile.openstreetmap.org/3/1/2.png", url)

		// check correct a... url is returned (it loops back to a after c)
		url, err = TileProvider.GetTileAddress(1, 2, 3)
		require.NoError(t, err)
		assert.Equal(t, "http://a.tile.openstreetmap.org/3/1/2.png", url)
	})
}
