package slippymap

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func (osm *OSMTileID) GetNorthNeighbour() OSMTileID {

func TestOSMTileID(t *testing.T) {

	t.Run("Test GetNeighbour", func(t *testing.T) {
		// define test data
		tables := []struct {
			origTile     OSMTileID
			direction    int
			expectedTile OSMTileID
		}{
			{
				origTile: OSMTileID{
					x:    0,
					y:    0,
					zoom: 1,
				},
				direction: DIRECTION_NORTH,
				expectedTile: OSMTileID{
					x:    0,
					y:    1,
					zoom: 1,
				},
			},
			{
				origTile: OSMTileID{
					x:    0,
					y:    0,
					zoom: 1,
				},
				direction: DIRECTION_SOUTH,
				expectedTile: OSMTileID{
					x:    0,
					y:    1,
					zoom: 1,
				},
			},
			{
				origTile: OSMTileID{
					x:    0,
					y:    0,
					zoom: 1,
				},
				direction: DIRECTION_WEST,
				expectedTile: OSMTileID{
					x:    1,
					y:    0,
					zoom: 1,
				},
			},
			{
				origTile: OSMTileID{
					x:    0,
					y:    0,
					zoom: 1,
				},
				direction: DIRECTION_EAST,
				expectedTile: OSMTileID{
					x:    1,
					y:    0,
					zoom: 1,
				},
			},

			{
				origTile: OSMTileID{
					x:    1,
					y:    1,
					zoom: 1,
				},
				direction: DIRECTION_NORTH,
				expectedTile: OSMTileID{
					x:    1,
					y:    0,
					zoom: 1,
				},
			},
			{
				origTile: OSMTileID{
					x:    1,
					y:    1,
					zoom: 1,
				},
				direction: DIRECTION_SOUTH,
				expectedTile: OSMTileID{
					x:    1,
					y:    0,
					zoom: 1,
				},
			},
			{
				origTile: OSMTileID{
					x:    1,
					y:    1,
					zoom: 1,
				},
				direction: DIRECTION_WEST,
				expectedTile: OSMTileID{
					x:    0,
					y:    1,
					zoom: 1,
				},
			},
			{
				origTile: OSMTileID{
					x:    1,
					y:    1,
					zoom: 1,
				},
				direction: DIRECTION_EAST,
				expectedTile: OSMTileID{
					x:    0,
					y:    1,
					zoom: 1,
				},
			},
		}

		for _, tt := range tables {
			var testName string
			switch tt.direction {
			case DIRECTION_NORTH:
				testName = fmt.Sprintf("North of %d/%d, zoom: %d", tt.origTile.x, tt.origTile.y, tt.origTile.zoom)
			case DIRECTION_SOUTH:
				testName = fmt.Sprintf("South of %d/%d, zoom: %d", tt.origTile.x, tt.origTile.y, tt.origTile.zoom)
			case DIRECTION_WEST:
				testName = fmt.Sprintf("West of %d/%d, zoom: %d", tt.origTile.x, tt.origTile.y, tt.origTile.zoom)
			case DIRECTION_EAST:
				testName = fmt.Sprintf("East of %d/%d, zoom: %d", tt.origTile.x, tt.origTile.y, tt.origTile.zoom)
			}
			t.Run(testName, func(t *testing.T) {
				newTile := tt.origTile.GetNeighbour(tt.direction)
				assert.Equal(t, tt.expectedTile, newTile)
			})
		}
	})
}
