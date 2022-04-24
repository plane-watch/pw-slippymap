package slippymap

import (
	"log"
)

type OSMTileID struct {

	// OpenStreetMap tile x/y/zoom
	x    int // OSM X
	y    int // OSM Y
	zoom int // OSM Zoom Level

}

func (osm *OSMTileID) enforceBounds() OSMTileID {
	// "wrap" tiles if the exceed the map size
	n := calcN(osm.zoom)
	if osm.x < 0 {
		osm.x = n + osm.x
	}
	if osm.x > n-1 {
		osm.x = osm.x - n
	}
	if osm.y < 0 {
		osm.y = n + osm.y
	}
	if osm.y > n-1 {
		osm.y = osm.y - n
	}
	return *osm
}

func (osm *OSMTileID) GetNeighbour(direction int) OSMTileID {
	// return an OSMTileID of the tile to the direction of tile defined by OSM
	neighbour := OSMTileID{
		x:    osm.x,
		y:    osm.y,
		zoom: osm.zoom,
	}
	switch direction {
	case DIRECTION_NORTH:
		neighbour.y = neighbour.y - 1
	case DIRECTION_SOUTH:
		neighbour.y = neighbour.y + 1
	case DIRECTION_WEST:
		neighbour.x = neighbour.x - 1
	case DIRECTION_EAST:
		neighbour.x = neighbour.x + 1
	default:
		log.Fatalf("Invalid direction: %d", direction)
	}

	return neighbour.enforceBounds()
}
