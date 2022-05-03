//go:build !race

package slippymap

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	INIT_CENTRE_LAT   = -31.9523 // initial map centre lat
	INIT_CENTRE_LONG  = 115.8613 // initial map centre long
	INIT_ZOOM_LEVEL   = 15       // initial OSM zoom level
	INIT_CENTRE_XTILE = 26929    // initial centre tile X (at zoom level 15)
	INIT_CENTRE_YTILE = 19456    // initial centre tile Y (at zoom level 15)
	SLIPPYMAP_WIDTH   = 1024     // slippymap width in pixels
	SLIPPYMAP_HEIGHT  = 768      // slippymap height in pixels
)

func TestSlippyMap(t *testing.T) {

	// declare variables
	var (
		tileProvider TileProvider
		smInitial    *SlippyMap
		err          error
	)

	// get tile provider
	t.Run("Test TileProviderForOS", func(t *testing.T) {
		tileProvider, err = TileProviderForOS()
		require.NoError(t, err, "TileProviderForOS returned error")
	})

	// test NewSlippyMap
	t.Run("Test NewSlippyMap", func(t *testing.T) {
		smInitial = NewSlippyMap(SLIPPYMAP_WIDTH, SLIPPYMAP_HEIGHT, INIT_ZOOM_LEVEL, INIT_CENTRE_LAT, INIT_CENTRE_LONG, tileProvider)
		require.NoError(t, err, "NewSlippyMap returned error")
	})

	// test GetTileAtPixel
	t.Run("Test GetTileAtPixel", func(t *testing.T) {
		osmX, osmY, zoomLevel, err := smInitial.GetTileAtPixel(SLIPPYMAP_WIDTH/2, SLIPPYMAP_HEIGHT/2)
		require.NoError(t, err, "GetTileAtPixel returned error")
		assert.Equal(t, INIT_ZOOM_LEVEL, zoomLevel, "GetTileAtPixel returned unexpected zoom level")
		assert.Equal(t, INIT_CENTRE_XTILE, osmX, "GetTileAtPixel returned unexpected X")
		assert.Equal(t, INIT_CENTRE_YTILE, osmY, "GetTileAtPixel returned unexpected Y")
	})

	// test GetLatLongAtPixel
	// TODO: I'm still not 100% convinced that GetLatLongAtPixel is as accurate as it should be...
	//       Maybe someone better at maths can double-check this...
	t.Run("Test GetLatLongAtPixel", func(t *testing.T) {
		lat_deg, long_deg, err := smInitial.GetLatLongAtPixel(SLIPPYMAP_WIDTH/2, SLIPPYMAP_HEIGHT/2)
		require.NoError(t, err, "GetLatLongAtPixel returned error")

		// round to 3 decimal places (to account for zoom level error)
		lat_deg = math.Round(lat_deg*1000) / 1000
		long_deg = math.Round(long_deg*1000) / 1000

		// check results
		assert.Equal(t, math.Round(INIT_CENTRE_LAT*1000)/1000, lat_deg, "GetLatLongAtPixel returned unexpected latitude")
		assert.Equal(t, math.Round(INIT_CENTRE_LONG*1000)/1000, long_deg, "GetLatLongAtPixel returned unexpected longitude")
	})

	// test LatLongToPixel
	t.Run("Test LatLongToPixel", func(t *testing.T) {
		x, y, err := smInitial.LatLongToPixel(INIT_CENTRE_LAT, INIT_CENTRE_LONG)
		require.NoError(t, err, "LatLongToPixel returned error")
		assert.Equal(t, SLIPPYMAP_WIDTH/2, x, "LatLongToPixel returned unexpected x")
		assert.Equal(t, SLIPPYMAP_HEIGHT/2, y, "LatLongToPixel returned unexpected y")
	})

	// test Update
	t.Run("Test Update", func(t *testing.T) {
		for i := 0; i <= 100; i++ {
			smInitial.Update(true)
		}
	})

	// test GetSize
	t.Run("Test GetSize", func(t *testing.T) {
		mapWidthPx, mapHeightPx := smInitial.GetSize()
		assert.Equal(t, SLIPPYMAP_WIDTH, mapWidthPx, "GetSize returned unexpected width")
		assert.Equal(t, SLIPPYMAP_HEIGHT, mapHeightPx, "GetSize returned unexpected height")
	})

	// test GetNumTiles
	t.Run("Test GetNumTiles", func(t *testing.T) {
		numTiles := smInitial.GetNumTiles()
		assert.Positive(t, numTiles)
	})

	// test SetSize
	t.Run("Test SetSize", func(t *testing.T) {
		smInitial.SetSize(SLIPPYMAP_WIDTH+500, SLIPPYMAP_HEIGHT+500)
		time.Sleep(time.Second)
		mapWidthPx, mapHeightPx := smInitial.GetSize()
		assert.Equal(t, SLIPPYMAP_WIDTH+500, mapWidthPx, "GetSize returned unexpected width after SetSize")
		assert.Equal(t, SLIPPYMAP_HEIGHT+500, mapHeightPx, "GetSize returned unexpected height after SetSize")
	})

	// test Update (moving map off screen)
	t.Run("Test MoveBy", func(t *testing.T) {
		smInitial.MoveBy(-SLIPPYMAP_WIDTH, -SLIPPYMAP_HEIGHT)
		smInitial.Update(true)
		smInitial.MoveBy(-SLIPPYMAP_WIDTH, -SLIPPYMAP_HEIGHT)
		smInitial.Update(true)
		smInitial.MoveBy(-SLIPPYMAP_WIDTH, -SLIPPYMAP_HEIGHT)
		smInitial.Update(true)
		smInitial.MoveBy(-SLIPPYMAP_WIDTH, -SLIPPYMAP_HEIGHT)
		smInitial.Update(true)
	})

	// test GetZoomLevel
	t.Run("Test GetZoomLevel", func(t *testing.T) {
		zl := smInitial.GetZoomLevel()
		assert.Equal(t, INIT_ZOOM_LEVEL, zl, "GetZoomLevel result not expected")
	})

	t.Run("Test SetZoomLevel", func(t *testing.T) {
		// test SetZoomLevel
		smZoomMin, err := smInitial.SetZoomLevel(ZOOM_LEVEL_MIN, INIT_CENTRE_LAT, INIT_CENTRE_LONG)
		require.NoError(t, err, "SetZoomLevel returned error")
		smZoomMin.Update(true)

		// test SetZoomLevel error
		_, err = smInitial.SetZoomLevel(ZOOM_LEVEL_MAX+1, INIT_CENTRE_LAT, INIT_CENTRE_LONG)
		require.Error(t, err, "SetZoomLevel did not return an error when one was expected")

		// test SetZoomLevel error
		_, err = smInitial.SetZoomLevel(ZOOM_LEVEL_MIN-1, INIT_CENTRE_LAT, INIT_CENTRE_LONG)
		require.Error(t, err, "SetZoomLevel did not return an error when one was expected")
	})

	// test ZoomIn
	t.Run("Test ZoomIn", func(t *testing.T) {
		_, err = smInitial.ZoomIn(INIT_CENTRE_LAT, INIT_CENTRE_LONG)
		require.NoError(t, err, "ZoomIn returned error")
	})

	// test ZoomOut
	t.Run("Test ZoomOut", func(t *testing.T) {
		_, err = smInitial.ZoomOut(INIT_CENTRE_LAT, INIT_CENTRE_LONG)
		require.NoError(t, err, "ZoomOut returned error")
	})
}

func TestGpsCoordsToTileInfo(t *testing.T) {

	// define test data
	tables := []struct {
		lat, long        float64
		zoom             int
		x, y             int
		offsetX, offsetY float64
	}{
		// Perth, Western Australia (SE)
		{
			lat: -31.9523, long: 115.8613, zoom: 16,
			x: 53859, y: 38912,
			offsetX: 0.90599, offsetY: 0.02956,
		},
		// Bergen, Norway (NE)
		{
			lat: 60.3913, long: 5.3221, zoom: 14,
			x: 8434, y: 4722,
			offsetX: 0.214684, offsetY: 0.078107,
		},
		// Vancouver, Canada (NW)
		{
			lat: 49.2827, long: -123.1207, zoom: 10,
			x: 161, y: 350,
			offsetX: 0.7900089, offsetY: 0.4350589,
		},
		// Rio de Janeiro, Brazil (SW)
		{
			lat: -22.9068, long: -43.1729, zoom: 13,
			x: 3113, y: 4631,
			offsetX: 0.576675556, offsetY: 0.725210945,
		},
	}

	for _, table := range tables {

		// get the file for the test lat/long/zoom
		x, y, oX, oY := gpsCoordsToTileInfo(table.lat, table.long, table.zoom)

		// check tile X & Y
		assert.Equal(t, table.x, x, "gpsCoordsToTileInfo returned unexpected X")
		assert.Equal(t, table.y, y, "gpsCoordsToTileInfo returned unexpected Y")

		// round offsets to 2 decimal places & test
		offsetX := math.Round(((table.offsetX * TILE_HEIGHT_PX) * 100) / 100)
		offsetY := math.Round(((table.offsetY * TILE_WIDTH_PX) * 100) / 100)
		oX = math.Round((oX * 100) / 100)
		oY = math.Round((oY * 100) / 100)

		// check offsets
		assert.Equal(t, offsetX, oX, "gpsCoordsToTileInfo returned unexpected X offset")
		assert.Equal(t, offsetY, oY, "gpsCoordsToTileInfo returned unexpected Y offset")
	}
}

func TestTileXYZtoGpsCoords(t *testing.T) {

	// define test data
	// test each zoom level in NW, NE, SW, SE
	tables := []struct {
		x, y, z                 int
		topLeftLat, topLeftLong float64
	}{
		// zoom level 0: single tile
		{
			x: 0, y: 0, z: 0,
			topLeftLat: 85.0511288, topLeftLong: -180,
		},
		// zoom level 1
		{
			x: 0, y: 0, z: 1, // NW
			topLeftLat: 85.0511288, topLeftLong: -180,
		},
		{
			x: 1, y: 0, z: 1, // NE
			topLeftLat: 85.0511288, topLeftLong: 0,
		},
		{
			x: 1, y: 1, z: 1, // SE
			topLeftLat: 0, topLeftLong: 0,
		},
		{
			x: 0, y: 1, z: 1, // SW
			topLeftLat: 0, topLeftLong: -180,
		},
		// zoom level 8
		{
			x: 63, y: 63, z: 8, // NW
			topLeftLat: 67.0674334, topLeftLong: -91.40625,
		},
		{
			x: 191, y: 63, z: 8, // NE
			topLeftLat: 67.0674334, topLeftLong: 88.59375,
		},
		{
			x: 191, y: 191, z: 8, // SE
			topLeftLat: -65.9464718, topLeftLong: 88.59375,
		},
		{
			x: 63, y: 191, z: 8, // SW
			topLeftLat: -65.9464718, topLeftLong: -91.40625,
		},
	}

	for _, table := range tables {

		// get top-left lat/long of tile at x/y/zoom
		lat, long := tileXYZtoGpsCoords(table.x, table.y, table.z)

		// round to 7 decimal places
		lat = math.Round(lat*10000000) / 10000000
		long = math.Round(long*10000000) / 10000000

		// check lat/long
		assert.Equal(t, table.topLeftLat, lat, "tileXYZtoGpsCoords returned unexpected latitude")
		assert.Equal(t, table.topLeftLong, long, "tileXYZtoGpsCoords returned unexpected longitude")
	}

}
