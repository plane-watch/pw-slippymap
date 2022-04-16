package slippymap

import (
	"errors"
	"fmt"
)

// OSMTileProvider generates the URLs to OSM tiles
type OSMTileProvider struct {
	// used to round-robin connections to OSM servers, as-per https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames#Tile_servers
	osm_url_prefix int
}

var _ TileProvider = &OSMTileProvider{}

func (op *OSMTileProvider) GetTileAddress(x, y, z int) (tilePath string, err error) {
	var url string
	// returns URL to open street map tile
	// load balance urls across servers as-per OSM guidelines: https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames#Tile_servers
	switch op.osm_url_prefix {
	case 0:
		url = fmt.Sprintf("http://a.tile.openstreetmap.org/%d/%d/%d.png", z, x, y)
		op.osm_url_prefix = 1
	case 1:
		url = fmt.Sprintf("http://b.tile.openstreetmap.org/%d/%d/%d.png", z, x, y)
		op.osm_url_prefix = 2
	case 2:
		url = fmt.Sprintf("http://c.tile.openstreetmap.org/%d/%d/%d.png", z, x, y)
		op.osm_url_prefix = 0
	default:
		return "", errors.New("invalid osm_url_prefix")
	}
	return url, nil
}
