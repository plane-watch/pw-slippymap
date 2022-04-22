package slippymap

import (
	"errors"
	"fmt"
	_ "image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	TILE_WIDTH_PX              = 256  // tile width (as-per https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames)
	TILE_HEIGHT_PX             = 256  // tile height (as-per https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames)
	ZOOM_LEVEL_MAX             = 16   // maximum zoom level
	ZOOM_LEVEL_MIN             = 2    // minimum zoom level
	TILE_FADEIN_ALPHA_PER_TICK = 0.05 // amount of alpha added per tick for tile fade-in
)

type MapTile struct {
	osmX      int           // OSM X
	osmY      int           // OSM Y
	zoomLevel int           // OSM Zoom Level
	img       *ebiten.Image // Image data
	offsetX   int           // top-left pixel location of tile
	offsetY   int           // top-right pixel location of tile
	alpha     float64       // tile transparency (for fade-in)
}

type SlippyMap struct {
	img *ebiten.Image // map image

	offsetX     int  // hold the current X offset
	offsetY     int  // hold the current Y offset
	need_update bool // do we need to process Update()

	tiles []*MapTile // map tiles

	mapWidthPx  int // number of pixels wide
	mapHeightPx int // number of pixels high

	zoomLevel        int           // zoom level
	zoomPrevLevelImg *ebiten.Image // holds the previous zoom level's image

	offsetMinimumX int // minimum X value for map tiles
	offsetMinimumY int // minimum Y value for map tiles
	offsetMaximumX int // maximum X value for map tiles
	offsetMaximumY int // maximum Y value for map tiles

	tileProvider TileProvider // the tile provider for the slippymap
}

func (sm *SlippyMap) GetZoomLevel() (zoomLevel int) {
	// returns the current zoom level
	return sm.zoomLevel
}

func (sm *SlippyMap) GetNumTiles() (numTiles int) {
	// returns the number of tiles making up the slippymap
	return len(sm.tiles)
}

func (sm *SlippyMap) Draw(screen *ebiten.Image) {

	// draw the previous zoom level in the background so zooming fades nicely
	// TODO: stretch this image so it looks like we're zooming in, new tiles will fade in over old ones
	screen.DrawImage(sm.zoomPrevLevelImg, nil)

	// render tiles onto sm.img
	for _, t := range sm.tiles {
		dio := &ebiten.DrawImageOptions{}

		// move the image where it needs to be in the window
		dio.GeoM.Translate(float64((*t).offsetX), float64((*t).offsetY))

		// adjust transparency (for fade-in of tiles)
		dio.ColorM.Scale(1, 1, 1, (*t).alpha)

		// draw the tile
		sm.img.DrawImage(t.img, dio)

		// debugging: print the OSM tile X/Y/Z
		dbgText := fmt.Sprintf("%d/%d/%d", (*t).osmX, (*t).osmY, (*t).zoomLevel)
		ebitenutil.DebugPrintAt(sm.img, dbgText, (*t).offsetX, (*t).offsetY)
	}

	// draw sm.img to the game screen
	screen.DrawImage(sm.img, nil)

}

func (sm *SlippyMap) Update(deltaOffsetX, deltaOffsetY int, forceUpdate bool) {
	// Updates the map
	//  - Loads any missing tiles
	//  - Cleans up any tiles that are "out of bounds"
	//  - Moves tiles as-per deltaOffsetX/Y

	var wereTilesCleanedUp bool // were off screen tiles cleaned up?
	var wereTilesMoved bool     // were tiles moved?
	var wereTilesAlphad bool    // did tiles have their alpha changed?
	var wereTilesCreated bool   // were tiles created?

	// don't update unless required
	//   * offscreen tiles are being cleaned up; or
	//   * user has moved the map; or
	//   * tile fade-in happenning; or
	//   * new tiles were created
	if deltaOffsetX != sm.offsetX || deltaOffsetY != sm.offsetY || forceUpdate || sm.need_update {

		// clean up tiles off the screen
		for i, t := range sm.tiles {
			// if tile is out of bounds, remove it from slice
			if sm.isOutOfBounds((*t).offsetX, (*t).offsetY) {
				sm.tiles[i] = sm.tiles[len(sm.tiles)-1]
				sm.tiles = sm.tiles[:len(sm.tiles)-1]
				wereTilesCleanedUp = true
				break
			}
		}

		// tile reposition & alpha increase if needed
		for _, t := range sm.tiles {
			// update offset if required (ie, user is dragging the map around)
			if (deltaOffsetX != 0 && deltaOffsetY != 0) || forceUpdate {
				t.offsetX = t.offsetX + deltaOffsetX
				t.offsetY = t.offsetY + deltaOffsetY
				// sm.tiles[i] = (*t)
				wereTilesMoved = true
			}

			// increase alpha channel (for fade in, if needed)
			if (*t).alpha < 1 {
				(*t).alpha = (*t).alpha + TILE_FADEIN_ALPHA_PER_TICK
				wereTilesAlphad = true
				ebiten.ScheduleFrame()
			}
		}

		// new tiles created if required (just do one tile per update)
		makeAnother := true
		for makeAnother {
			makeAnother = false
			for _, t := range sm.tiles {
				if sm.makeTileAbove(t) {
					makeAnother = true
					wereTilesCreated = true
				}
				if sm.makeTileToTheLeft(t) {
					makeAnother = true
					wereTilesCreated = true
				}
				if sm.makeTileToTheRight(t) {
					makeAnother = true
					wereTilesCreated = true
				}
				if sm.makeTileBelow(t) {
					makeAnother = true
					wereTilesCreated = true
				}
			}
		}
	}

	// work out whether this function needs to run next iteration
	if wereTilesCleanedUp || wereTilesMoved || wereTilesAlphad || wereTilesCreated {
		sm.need_update = true
	} else {
		sm.need_update = false
	}

}

func (sm *SlippyMap) makeTileAbove(existingTile *MapTile) (tileCreated bool) {
	// makes the tile above existingTile, if it does not already exist or would be out of bounds

	// check to see if tile above already exists
	newTileOSMY := (*existingTile).osmY - 1

	// honour edges of map
	if newTileOSMY == -1 {
		return false
	}

	for _, t := range sm.tiles {
		if (*t).osmY == newTileOSMY && (*t).osmX == (*existingTile).osmX {
			// the tile already exists, bail out
			return false
		}
	}

	newTileOffsetY := (*existingTile).offsetY - TILE_HEIGHT_PX

	// if the tile would not be out of bounds...
	if sm.isOutOfBounds((*existingTile).offsetX, newTileOffsetY) != true {
		// make the new tile
		sm.makeTile((*existingTile).osmX, newTileOSMY, (*existingTile).offsetX, newTileOffsetY)
		return true
	}
	return false
}

func (sm *SlippyMap) makeTileBelow(existingTile *MapTile) (tileCreated bool) {
	// makes the tile below existingTile, if it does not already exist or would be out of bounds

	// check to see if tile below already exists
	newTileOSMY := (*existingTile).osmY + 1

	// honour edges of map
	if newTileOSMY == int(math.Pow(2, float64(sm.zoomLevel))) {
		return false
	}

	for _, t := range sm.tiles {
		if t.osmY == newTileOSMY && t.osmX == (*existingTile).osmX {
			// the tile already exists, bail out
			return false
		}
	}

	newTileOffsetY := (*existingTile).offsetY + TILE_HEIGHT_PX

	// if the tile would not be out of bounds...
	if sm.isOutOfBounds((*existingTile).offsetX, newTileOffsetY) != true {
		// make the new tile
		sm.makeTile((*existingTile).osmX, newTileOSMY, (*existingTile).offsetX, newTileOffsetY)
		return true
	}
	return false
}

func (sm *SlippyMap) makeTileToTheLeft(existingTile *MapTile) (tileCreated bool) {
	// makes the tile to the left of existingTile, if it does not already exist or would be out of bounds

	// check to see if tile to the left already exists
	newTileOSMX := (*existingTile).osmX - 1

	// honour edges of map
	if newTileOSMX == -1 {
		return false
	}

	for _, t := range sm.tiles {
		if t.osmX == newTileOSMX && t.osmY == (*existingTile).osmY {
			// the tile already exists, bail out
			return false
		}
	}

	newTileOffsetX := (*existingTile).offsetX - TILE_WIDTH_PX

	// if the tile would not be out of bounds...
	if sm.isOutOfBounds(newTileOffsetX, (*existingTile).offsetY) != true {
		// make the new tile
		sm.makeTile(newTileOSMX, (*existingTile).osmY, newTileOffsetX, (*existingTile).offsetY)
		return true
	}
	return false
}

func (sm *SlippyMap) makeTileToTheRight(existingTile *MapTile) (tileCreated bool) {
	// makes the tile to the right of existingTile, if it does not already exist or would be out of bounds

	// check to see if tile to the right already exists
	newTileOSMX := (*existingTile).osmX + 1

	// honour edges of map
	if newTileOSMX == int(math.Pow(2, float64(sm.zoomLevel))) {
		return false
	}

	for _, t := range sm.tiles {
		if t.osmX == newTileOSMX && t.osmY == (*existingTile).osmY {
			// the tile already exists, bail out
			return false
		}
	}

	newTileOffsetX := (*existingTile).offsetX + TILE_WIDTH_PX

	// if the tile would not be out of bounds...
	if sm.isOutOfBounds(newTileOffsetX, (*existingTile).offsetY) != true {
		// make the new tile
		sm.makeTile(newTileOSMX, (*existingTile).osmY, newTileOffsetX, (*existingTile).offsetY)
		return true
	}
	return false
}

func (sm *SlippyMap) isOutOfBounds(pixelX, pixelY int) (outOfBounds bool) {
	// returns true if the point defined by pixelX and pixelY is "out of bounds"
	// "out of bounds" means the point is outside the renderable size of the map
	// which is defined by sm.offset[Minimum|Maximum][X|Y].

	if pixelX < sm.offsetMinimumX {
		return true
	}
	if pixelY < sm.offsetMinimumY {
		return true
	}
	if pixelX > sm.offsetMaximumX {
		return true
	}
	if pixelY > sm.offsetMaximumY {
		return true
	}
	return false
}

func (sm *SlippyMap) makeTile(osmX, osmY, offsetX, offsetY int) {
	// Creates a new tile on the slippymap

	// Create the tile object
	t := MapTile{
		osmX:      osmX,
		osmY:      osmY,
		offsetX:   offsetX,
		offsetY:   offsetY,
		zoomLevel: sm.zoomLevel,
		img:       ebiten.NewImage(TILE_WIDTH_PX, TILE_WIDTH_PX),
	}

	go func() {
		// get tile artwork
		tilePath, err := sm.tileProvider.GetTileAddress(t.osmX, t.osmY, t.zoomLevel)
		if err != nil {
			log.Fatal(err)
		}

		// load the image
		img, _, err := ebitenutil.NewImageFromFile(tilePath)
		if err != nil {
			log.Fatal(err)
		}
		t.img.DrawImage(img, nil)
	}()

	// Add tile to slippymap
	sm.tiles = append(sm.tiles, &t)
}

func (sm *SlippyMap) SetSize(mapWidthPx, mapHeightPx int) {
	// updates the slippy map when window size is changed
	sm.mapWidthPx = mapWidthPx
	sm.mapHeightPx = mapHeightPx
	sm.offsetMinimumX = -(2 * TILE_WIDTH_PX)
	sm.offsetMinimumY = -(2 * TILE_HEIGHT_PX)
	sm.offsetMaximumX = mapWidthPx + (2 * TILE_WIDTH_PX)
	sm.offsetMaximumY = mapHeightPx + (2 * TILE_HEIGHT_PX)
}

func (sm *SlippyMap) GetSize() (mapWidthPx, mapHeightPx int) {
	// return the slippymap size in pixels
	return sm.mapWidthPx, sm.mapHeightPx
}

func (sm *SlippyMap) GetTileAtPixel(x, y int) (osmX, osmY, zoomLevel int, err error) {
	// returns the OSM tile X/Y/Z at pixel position x,y
	for _, t := range sm.tiles {
		if x >= (*t).offsetX && x < (*t).offsetX+TILE_WIDTH_PX {
			if y >= (*t).offsetY && y < (*t).offsetY+TILE_HEIGHT_PX {
				return (*t).osmX, (*t).osmY, (*t).zoomLevel, nil
			}
		}
	}
	return 0, 0, 0, errors.New("Tile not found")
}

func (sm *SlippyMap) GetLatLongAtPixel(x, y int) (latDeg, longDeg float64, err error) {
	// returns the lat/long at pixel x,y

	// first get tile
	osmX, osmY, zoomLevel, err := sm.GetTileAtPixel(x, y)
	if err != nil {
		return 0, 0, err
	}

	// find tile in slippymap
	tileFound := false
	var topLeftX, topLeftY int
	for _, t := range sm.tiles {
		if (*t).osmX == osmX && (*t).osmY == osmY && (*t).zoomLevel == zoomLevel {
			tileFound = true
			topLeftX = (*t).offsetX
			topLeftY = (*t).offsetY
			break
		}
	}

	// raise err if tile not found
	if tileFound != true {
		return 0, 0, errors.New("Tile not found")
	}

	// get pixel offset within tile
	offsetX := x - topLeftX
	offsetY := y - topLeftY

	mercatorX := float64((osmX*TILE_WIDTH_PX)+offsetX) / TILE_WIDTH_PX
	mercatorY := float64((osmY*TILE_HEIGHT_PX)+offsetY) / TILE_HEIGHT_PX

	latDeg = math.Atan(math.Sinh(math.Pi-(mercatorY/math.Pow(2, float64(sm.zoomLevel))*2*math.Pi))) * (180 / math.Pi)

	longDeg = (mercatorX / math.Pow(2, float64(sm.zoomLevel)) * 360) - 180

	return latDeg, longDeg, nil

}

func (sm *SlippyMap) LatLongToPixel(lat_deg, long_deg float64) (x, y int, err error) {
	// return the pixel x/y for a given lat/long

	// find the tile for the given lat/long
	osmX, osmY, offsetX, offsetY := gpsCoordsToTileInfo(lat_deg, long_deg, sm.zoomLevel)

	// find the tile on the slippymap
	tileFound := false
	for _, t := range sm.tiles {
		if (*t).osmX == osmX && (*t).osmY == osmY {
			tileFound = true
			x = int(offsetX) + (*t).offsetX
			y = int(offsetY) + (*t).offsetY
			break
		}
	}
	if tileFound != true {
		return 0, 0, errors.New("Tile not found")
	}
	return x, y, nil
}

func (sm *SlippyMap) ZoomIn(lat_deg, long_deg float64) (newsm SlippyMap, err error) {
	// zoom in, with map centred on given lat/long (in degrees)
	newsm, err = sm.SetZoomLevel(sm.zoomLevel+1, lat_deg, long_deg)
	return newsm, err
}

func (sm *SlippyMap) ZoomOut(lat_deg, long_deg float64) (newsm SlippyMap, err error) {
	// zoom in, with map centred on given lat/long (in degrees)
	newsm, err = sm.SetZoomLevel(sm.zoomLevel-1, lat_deg, long_deg)
	return newsm, err
}

func (sm *SlippyMap) SetZoomLevel(zoomLevel int, lat_deg, long_deg float64) (newsm SlippyMap, err error) {
	// sets zoom level, with map centred on given lat/long (in degrees)

	// ensure we're within ZOOM_LEVEL_MAX & ZOOM_LEVEL_MIN
	if zoomLevel > ZOOM_LEVEL_MAX || zoomLevel < ZOOM_LEVEL_MIN {
		return SlippyMap{}, errors.New("Requested zoom level unavailable")
	}

	// create a new slippymap centred on the requested lat/long, at the requested zoom level
	newsm, err = NewSlippyMap(sm.mapWidthPx, sm.mapHeightPx, zoomLevel, lat_deg, long_deg, sm.tileProvider)
	if err != nil {
		return SlippyMap{}, err
	}

	// copy the current map image into the zoom previous level background image
	sm.Draw(newsm.zoomPrevLevelImg)

	// return the new slippymap and no error
	return newsm, nil
}

func NewSlippyMap(mapWidthPx, mapHeightPx, zoomLevel int, centreLat, centreLong float64, tileProvider TileProvider) (sm SlippyMap, err error) {

	log.Printf("Initialising SlippyMap at %0.4f/%0.4f, zoom level %d", centreLat, centreLong, zoomLevel)

	// determine the centre tile details
	centreTileOSMX, centreTileOSMY, pixelOffsetX, pixelOffsetY := gpsCoordsToTileInfo(centreLat, centreLong, zoomLevel)

	// create a new SlippyMap to return
	sm = SlippyMap{
		img:              ebiten.NewImage(mapWidthPx, mapHeightPx), // initialise main image
		zoomPrevLevelImg: ebiten.NewImage(mapWidthPx, mapHeightPx), // initialise image of previous zoom level
		zoomLevel:        zoomLevel,                                // set zoom level
		tileProvider:     tileProvider,                             // set tile provider
		need_update:      true,                                     // ensure first-time update
	}

	// update size
	sm.SetSize(mapWidthPx, mapHeightPx)

	// initialise the map with a centre tile
	centreTileOffsetX := (mapWidthPx / 2) - int(pixelOffsetX)
	centreTileOffsetY := (mapHeightPx / 2) - int(pixelOffsetY)
	sm.makeTile(centreTileOSMX, centreTileOSMY, centreTileOffsetX, centreTileOffsetY)

	// force initial update
	sm.Update(0, 0, true)

	// return the slippymap
	return sm, nil
}

func calcN(zoom_lvl int) (n int) {
	// calculates n
	// tile coverage is n x n tiles
	return int(math.Pow(2, float64(zoom_lvl)))
}

func secant(x float64) (s float64) {
	// calculate the secant (1/cos)
	// TODO: is there a golang math function that does this?
	return 1 / math.Cos(x)
}

func DegreesToRadians(d float64) (r float64) {
	// convert degrees to radians
	return d * (math.Pi / 180.0)
}

func RadiansToDegrees(r float64) (d float64) {
	// convert radians to degrees
	return r * 180 / math.Pi
}

func gpsCoordsToTileInfo(lat_deg, long_deg float64, zoomLevel int) (tileX, tileY int, pixelOffsetX, pixelOffsetY float64) {
	// return OSM tile x/y coordinates (and pixel offset to the exact position) from lat/long

	// perform calculation as-per: https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames#Lon..2Flat._to_tile_numbers
	n := float64(calcN(zoomLevel))
	lat_rad := DegreesToRadians(lat_deg)
	x := n * ((long_deg + 180.0) / 360.0)
	y := n * (1 - (math.Log(math.Tan(lat_rad)+secant(lat_rad)) / math.Pi)) / 2.0

	tileX = int(math.Floor(x))
	tileY = int(math.Floor(y))

	pixelOffsetX = (x - math.Floor(x)) * TILE_WIDTH_PX
	pixelOffsetY = (y - math.Floor(y)) * TILE_HEIGHT_PX

	return tileX, tileY, pixelOffsetX, pixelOffsetY
}

func tileXYZtoGpsCoords(x, y, z int) (topLeftLat, topLeftLong float64) {
	// return the top left lat/long of a tile
	n := float64(calcN(z))
	topLeftLat = RadiansToDegrees(math.Atan(math.Sinh(math.Pi * (1 - 2*float64(y)/n))))
	topLeftLong = float64(x)/n*360.0 - 180.0
	return topLeftLat, topLeftLong
}
