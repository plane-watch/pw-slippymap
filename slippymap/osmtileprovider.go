package slippymap

import (
	"errors"
	"fmt"
	"sync/atomic"
)

// OSMTileProvider generates the URLs to OSM tiles
type OSMTileProvider struct {
	// used to round-robin connections to OSM servers, as-per https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames#Tile_servers
	osm_url_prefix int32
}

var _ TileProvider = &OSMTileProvider{}

func (op *OSMTileProvider) GetTileAddress(osm OSMTileID) (tilePath string, err error) {
	var url string
	// returns URL to open street map tile
	// load balance urls across servers as-per OSM guidelines: https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames#Tile_servers

	// we don't really care about atomically incrementing this number, but golang race detection bitches
	// at us if we let this be subject to race conditions. Performance not super critical so just use
	// atomic to make the race detector shut up
	nextCDN := atomic.AddInt32(&op.osm_url_prefix, 1) % 3

	switch nextCDN {
	case 0:
		url = fmt.Sprintf("http://a.tile.openstreetmap.org/%d/%d/%d.png", osm.zoom, osm.x, osm.y)

	case 1:
		url = fmt.Sprintf("http://b.tile.openstreetmap.org/%d/%d/%d.png", osm.zoom, osm.x, osm.y)

	case 2:
		url = fmt.Sprintf("http://c.tile.openstreetmap.org/%d/%d/%d.png", osm.zoom, osm.x, osm.y)

	default:
		return "", errors.New("invalid osm_url_prefix")
	}
	return url, nil
}
