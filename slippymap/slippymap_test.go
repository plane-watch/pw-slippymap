package slippymap

import (
	"math"
	"testing"
)

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
		x, y, oX, oY := gpsCoordsToTileInfo(table.lat, table.long, table.zoom)

		// check x
		if x != table.x {
			t.Errorf("Lat/Long: %f, %f, zoom: %d expected tile X: %d, got: %d", table.lat, table.long, table.zoom, table.x, x)
		}

		// check y
		if y != table.y {
			t.Errorf("Lat/Long: %f, %f, zoom: %d expected tile Y: %d, got: %d", table.lat, table.long, table.zoom, table.y, y)
		}

		// check offsetX (to 2 decimal places)
		offsetX := table.offsetX * TILE_HEIGHT_PX
		if math.Round(oX*100) != math.Round(offsetX*100) {
			t.Errorf("Lat/Long: %f, %f, zoom: %d expected offsetX: %f, got: %f", table.lat, table.long, table.zoom, offsetX, oX)
		}

		// check offsetY (to 2 decimal places)
		offsetY := table.offsetY * TILE_WIDTH_PX
		if math.Round(oY*100) != math.Round(offsetY*100) {
			t.Errorf("Lat/Long: %f, %f, zoom: %d expected offsetY: %f, got: %f", table.lat, table.long, table.zoom, offsetY, oY)
		}
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
		lat, long := tileXYZtoGpsCoords(table.x, table.y, table.z)

		// round to 7 decimal places
		lat = math.Round(lat*10000000) / 10000000
		long = math.Round(long*10000000) / 10000000

		// check lat
		if lat != table.topLeftLat {
			t.Errorf("Tile x: %d, y: %d, zoom: %d expected lat: %f, got: %f", table.x, table.y, table.z, table.topLeftLat, lat)
		}

		// check long
		if long != table.topLeftLong {
			t.Errorf("Tile x: %d, y: %d, zoom: %d expected long: %f, got: %f", table.x, table.y, table.z, table.topLeftLong, long)
		}

	}

}
