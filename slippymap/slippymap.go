package slippymap

import (
	"errors"
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"pw_slippymap/datasources"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	TILE_WIDTH_PX              = 256  // tile width (as-per https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames)
	TILE_HEIGHT_PX             = 256  // tile height (as-per https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames)
	ZOOM_LEVEL_MAX             = 16   // maximum zoom level
	ZOOM_LEVEL_MIN             = 2    // minimum zoom level
	TILE_FADEIN_ALPHA_PER_TICK = 0.05 // amount of alpha added per tick for tile fade-in

	DIRECTION_NORTH = 1
	DIRECTION_SOUTH = 2
	DIRECTION_WEST  = 3
	DIRECTION_EAST  = 4
)

type mapTile struct {

	// OpenStreetMap tile identifier
	osm OSMTileID

	// Information of surrounding tiles
	osmNeighbourTileNorth OSMTileID
	osmNeighbourTileSouth OSMTileID
	osmNeighbourTileWest  OSMTileID
	osmNeighbourTileEast  OSMTileID

	// Tile image
	img      *ebiten.Image // Image data
	imgMutex sync.Mutex    // Mutex to avoid races

	// Image location
	offsetX     int // top-left pixel location of tile
	offsetY     int // top-right pixel location of tile
	offsetMutex sync.Mutex

	// Alpha for smooth fade-in
	alpha float64 // tile transparency (for fade-in)

}

type SlippyMap struct {
	img      *ebiten.Image // map image
	imgMutex sync.Mutex

	offsetX int // hold the current X offset
	offsetY int // hold the current Y offset

	needUpdate      bool // do we need to process Update()
	needUpdateMutex sync.Mutex

	needDraw      bool // do we need to process Draw()
	needDrawMutex sync.Mutex

	tiles      []*mapTile // map tiles
	tilesMutex sync.Mutex // Mutex to avoid races

	mapWidthPx   int // number of pixels wide
	mapHeightPx  int // number of pixels high
	mapSizeMutex sync.Mutex

	zoomLevel        int           // zoom level
	zoomPrevLevelImg *ebiten.Image // holds the previous zoom level's image

	offsetMinimumX int // minimum X value for map tiles
	offsetMinimumY int // minimum Y value for map tiles
	offsetMaximumX int // maximum X value for map tiles
	offsetMaximumY int // maximum Y value for map tiles
	offsetMutex    sync.Mutex

	tileProvider TileProvider // the tile provider for the slippymap

	aircraftDb *datasources.AircraftDB // aircraft db
}

func (sm *SlippyMap) scheduleDraw() {
	// ensures map will be re-drawn next tick
	sm.needDrawMutex.Lock()
	defer sm.needDrawMutex.Unlock()
	sm.needDraw = true
	ebiten.ScheduleFrame()
}

func (sm *SlippyMap) drawRequired(reset bool) bool {
	// checks to see if a draw us required, and resets the counter if reset true
	sm.needDrawMutex.Lock()
	defer sm.needDrawMutex.Unlock()
	if sm.needDraw {
		sm.needDraw = false
		return true
	}
	return false
}

func (sm *SlippyMap) scheduleUpdate() {
	// ensures update will be processed next tick
	sm.needUpdateMutex.Lock()
	defer sm.needUpdateMutex.Unlock()
	sm.needUpdate = true
	ebiten.ScheduleFrame()
}

func (sm *SlippyMap) updateRequired(reset bool) bool {
	// checks to see if an update us required, and resets the counter if reset true
	sm.needUpdateMutex.Lock()
	defer sm.needUpdateMutex.Unlock()
	if sm.needUpdate {
		sm.needUpdate = false
		return true
	}
	return false
}

func (sm *SlippyMap) iterTiles() []*mapTile {
	// returns a list (slice) of tiles making up the slippymap at the time this function was run

	sm.tilesMutex.Lock()
	defer sm.tilesMutex.Unlock()
	output := make([]*mapTile, len(sm.tiles))
	for i, t := range sm.tiles {
		// t.imgMutex.Lock()
		if t != nil {
			output[i] = t
		}
		// t.imgMutex.Unlock()
	}
	return output
}

func (sm *SlippyMap) tileExistsByOSM(osmX, osmY, zoomLevel int) bool {
	// returns true if a tile exists in the slippymap

	sm.tilesMutex.Lock()
	defer sm.tilesMutex.Unlock()

	for _, t := range sm.tiles {
		if t.osm.x == osmX && t.osm.y == osmY && t.osm.zoom == zoomLevel {
			return true
		}
	}
	return false
}

func (sm *SlippyMap) rmTile(tile *mapTile) {
	// removes a tile from the slippymap

	var tileFound bool
	var tileIndex int

	sm.tilesMutex.Lock()
	defer sm.tilesMutex.Unlock()

	for i, t := range sm.tiles {
		if t == tile {
			tileFound = true
			tileIndex = i
			break
		}
	}

	if tileFound {
		sm.tiles[tileIndex] = sm.tiles[len(sm.tiles)-1]
		sm.tiles = sm.tiles[:len(sm.tiles)-1]
		sm.scheduleDraw()
	} else {
		log.Panic("Could not find tile in sm.tiles")
	}
}

func (sm *SlippyMap) GetZoomLevel() (zoomLevel int) {
	// returns the current zoom level
	return sm.zoomLevel
}

func (sm *SlippyMap) GetNumTiles() (numTiles int) {
	// returns the number of tiles making up the slippymap
	sm.tilesMutex.Lock()
	defer sm.tilesMutex.Unlock()
	output := len(sm.tiles)
	return output
}

func (sm *SlippyMap) Draw(screen *ebiten.Image) {

	// draw the previous zoom level in the background so zooming fades nicely
	// TODO: stretch this image so it looks like we're zooming in, new tiles will fade in over old ones

	if sm.drawRequired(true) {

		screen.DrawImage(sm.zoomPrevLevelImg, nil)

		// render tiles onto sm.img
		for _, t := range sm.iterTiles() {

			dio := &ebiten.DrawImageOptions{}

			// move the image where it needs to be in the window
			t.offsetMutex.Lock()
			dio.GeoM.Translate(float64(t.offsetX), float64(t.offsetY))
			t.offsetMutex.Unlock()

			// adjust transparency (for fade-in of tiles)
			// dio.ColorM.Scale(1, 1, 1, (*t).alpha)
			dio.ColorM.Scale(1, 1, 1, 1)

			// draw the tile
			t.imgMutex.Lock()
			sm.imgMutex.Lock()
			sm.img.DrawImage(t.img, dio)
			sm.imgMutex.Unlock()
			t.imgMutex.Unlock()

			// TEMPORARY TROUBLESHOOTING: check for black tile problem
			// Uncommenting this code make it run terribly :(
			// t.offsetMutex.Lock()
			// if t.offsetX > 0 && t.offsetX < screen.Bounds().Max.X && t.offsetY > 0 && t.offsetY < screen.Bounds().Max.Y {
			// 	sm.imgMutex.Lock()
			// 	blackTest := sm.img.At(t.offsetX+10, t.offsetY+10)
			// 	sm.imgMutex.Unlock()
			// 	btR, btG, btB, _ := blackTest.RGBA()
			// 	if btR == 0 && btG == 0 && btB == 0 {
			// 		log.Printf("Black tile at %d/%d/%d", t.osm.x, t.osm.y, t.osm.zoom)
			// 	}
			// }
			// t.offsetMutex.Unlock()

			// debugging: print the OSM tile X/Y/Z
			dbgText := fmt.Sprintf("%d/%d/%d", t.osm.x, t.osm.y, t.osm.zoom)
			t.offsetMutex.Lock()
			sm.imgMutex.Lock()
			ebitenutil.DebugPrintAt(sm.img, dbgText, t.offsetX, t.offsetY)
			sm.imgMutex.Unlock()
			t.offsetMutex.Unlock()
		}
	}

	// draw sm.img to the game screen
	screen.DrawImage(sm.img, nil)

}

func (sm *SlippyMap) MoveBy(deltaOffsetX, deltaOffsetY int) {
	// moves the map by deltaOffsetX, deltaOffsetY pixels relative to current view
	// tile reposition & alpha increase if needed
	for _, t := range sm.iterTiles() {

		// update offset if required (ie, user is dragging the map around)
		if deltaOffsetX != 0 || deltaOffsetY != 0 {
			t.offsetMutex.Lock()
			t.offsetX = t.offsetX + deltaOffsetX
			t.offsetY = t.offsetY + deltaOffsetY
			t.offsetMutex.Unlock()
		}
	}
	if deltaOffsetX != 0 || deltaOffsetY != 0 {
		sm.scheduleUpdate()
		sm.scheduleDraw()
	}
}

func (sm *SlippyMap) Update(forceUpdate bool) {
	// Updates the map
	//  - Loads any missing tiles
	//  - Cleans up any tiles that are "out of bounds"

	// don't update unless required
	//   * offscreen tiles are being cleaned up; or
	//   * user has moved the map; or
	//   * tile fade-in happenning; or
	//   * new tiles were created
	if forceUpdate || sm.updateRequired(false) {
		sm.updateRequired(true)

		// for each tile
		for _, t := range sm.iterTiles() {

			// new tiles created if required
			go sm.makeNeighbourTile(DIRECTION_NORTH, t)
			go sm.makeNeighbourTile(DIRECTION_SOUTH, t)
			go sm.makeNeighbourTile(DIRECTION_EAST, t)
			go sm.makeNeighbourTile(DIRECTION_WEST, t)

			// if tile is out of bounds, remove it from slice
			t.offsetMutex.Lock()
			if sm.isOutOfBounds(t.offsetX, t.offsetY) {
				sm.rmTile(t)
				// break
			} else {
				t.offsetMutex.Unlock()
			}
		}
	}
}

func (sm *SlippyMap) makeNeighbourTile(direction int, existingTile *mapTile) (tileCreated bool) {
	// makes the tile to direction of existingTile, if it does not already exist or would be out of bounds

	existingTile.offsetMutex.Lock()
	newTileOffsetX := existingTile.offsetX
	newTileOffsetY := existingTile.offsetY
	existingTile.offsetMutex.Unlock()

	var newTileOsm OSMTileID

	switch direction {
	case DIRECTION_NORTH:
		newTileOsm = existingTile.osmNeighbourTileNorth
		if newTileOsm.y >= existingTile.osm.y {
			return false
		}
		newTileOffsetY -= TILE_HEIGHT_PX
	case DIRECTION_SOUTH:
		newTileOsm = existingTile.osmNeighbourTileSouth
		if newTileOsm.y <= existingTile.osm.y {
			return false
		}
		newTileOffsetY += TILE_HEIGHT_PX
	case DIRECTION_WEST:
		newTileOsm = existingTile.osmNeighbourTileWest
		if newTileOsm.x >= existingTile.osm.x {
			return false
		}
		newTileOffsetX -= TILE_WIDTH_PX
	case DIRECTION_EAST:
		newTileOsm = existingTile.osmNeighbourTileEast
		if newTileOsm.x <= existingTile.osm.x {
			return false
		}
		newTileOffsetX += TILE_WIDTH_PX
	default:
		log.Fatalf("Invalid direction: %d", direction)
	}

	// check if tile already exists
	if sm.tileExistsByOSM(newTileOsm.x, newTileOsm.y, newTileOsm.zoom) {
		return false
	}

	// if the tile would not be out of bounds...
	if sm.isOutOfBounds(newTileOffsetX, newTileOffsetY) != true {
		// make the new tile
		sm.makeTile(newTileOsm.x, newTileOsm.y, newTileOffsetX, newTileOffsetY)
		sm.scheduleUpdate()
		sm.scheduleDraw()
	}

	// tile was out of bounds
	return false
}

func (sm *SlippyMap) isOutOfBounds(pixelX, pixelY int) (outOfBounds bool) {
	// returns true if the point defined by pixelX and pixelY is "out of bounds"
	// "out of bounds" means the point is outside the renderable size of the map
	// which is defined by sm.offset[Minimum|Maximum][X|Y].

	sm.offsetMutex.Lock()
	defer sm.offsetMutex.Unlock()

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
	// Creates a new tile on the slippymap sm at offxetX and offsetY

	osm := OSMTileID{
		x:    osmX,
		y:    osmY,
		zoom: sm.zoomLevel,
	}

	// Create the tile object
	t := &mapTile{

		// OpenStreetMap Tile Info
		osm:                   osm,
		osmNeighbourTileNorth: osm.GetNeighbour(DIRECTION_NORTH),
		osmNeighbourTileSouth: osm.GetNeighbour(DIRECTION_SOUTH),
		osmNeighbourTileWest:  osm.GetNeighbour(DIRECTION_WEST),
		osmNeighbourTileEast:  osm.GetNeighbour(DIRECTION_EAST),

		// Location on map
		offsetX: offsetX,
		offsetY: offsetY,

		// Prepare image
		img: ebiten.NewImage(TILE_WIDTH_PX, TILE_WIDTH_PX),
	}

	go func(t *mapTile, sm *SlippyMap) {

		// TEMPORARY DEBUGGING:
		log.Printf("Started loading artwork for: %d/%d/%d", t.osm.x, t.osm.y, t.osm.zoom)

		// get tile artwork
		tilePath, err := sm.tileProvider.GetTileAddress(t.osm)
		if err != nil {
			log.Fatal(err)
		}

		// load the image
		img, _, err := ebitenutil.NewImageFromFile(tilePath)
		if err != nil {
			log.Fatal(err)
		}

		// TEMPORARY DEBUGGING:
		blackTest := img.At(10, 10)
		if blackTest == color.Black {
			log.Println("ERROR!!! BLACK TILE AT: %d/%d/%d", t.osm.x, t.osm.y, t.osm.zoom)
		}

		t.imgMutex.Lock()
		t.img.DrawImage(img, nil)
		t.imgMutex.Unlock()

		sm.scheduleUpdate()
		sm.scheduleDraw()

		// TEMPORARY DEBUGGING:
		log.Printf("Finished loading artwork for: %d/%d/%d", t.osm.x, t.osm.y, t.osm.zoom)

	}(t, sm)

	// Add tile to slippymap
	t.imgMutex.Lock()
	sm.tilesMutex.Lock()
	sm.tiles = append(sm.tiles, t)
	sm.tilesMutex.Unlock()
	t.imgMutex.Unlock()

	sm.scheduleUpdate()
	sm.scheduleDraw()
}

func (sm *SlippyMap) SetSize(mapWidthPx, mapHeightPx int) (newsm *SlippyMap) {
	// todo fix race
	// updates the slippy map when window size is changed

	// get centre lat/long
	centreLat, centreLong, err := sm.GetLatLongAtPixel(sm.mapWidthPx/2, sm.mapHeightPx/2)
	if err != nil {
		log.Fatal(err)
	}

	// prepare new slippymap
	newsm = NewSlippyMap(mapWidthPx, mapHeightPx, sm.zoomLevel, centreLat, centreLong, sm.tileProvider)

	// copy the current map image into the zoom previous level background image
	sm.Draw(newsm.zoomPrevLevelImg)

	return newsm
}

func (sm *SlippyMap) GetSize() (mapWidthPx, mapHeightPx int) {
	// return the slippymap size in pixels
	sm.mapSizeMutex.Lock()
	defer sm.mapSizeMutex.Unlock()
	return sm.mapWidthPx, sm.mapHeightPx
}

func (sm *SlippyMap) GetTileAtPixel(x, y int) (osmX, osmY, zoomLevel int, err error) {
	// returns the OSM tile X/Y/Z at pixel position x,y
	sm.tilesMutex.Lock()
	defer sm.tilesMutex.Unlock()
	for _, t := range sm.tiles {
		t.offsetMutex.Lock()
		if x >= t.offsetX && x < t.offsetX+TILE_WIDTH_PX {
			if y >= t.offsetY && y < t.offsetY+TILE_HEIGHT_PX {
				t.offsetMutex.Unlock()
				return t.osm.x, t.osm.y, t.osm.zoom, nil
			}
		}
		t.offsetMutex.Unlock()
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
	sm.tilesMutex.Lock()
	for _, t := range sm.tiles {
		if t.osm.x == osmX && t.osm.y == osmY && t.osm.zoom == zoomLevel {
			tileFound = true
			t.offsetMutex.Lock()
			topLeftX = t.offsetX
			topLeftY = t.offsetY
			t.offsetMutex.Unlock()
			break
		}
	}
	sm.tilesMutex.Unlock()

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
	sm.tilesMutex.Lock()
	for _, t := range sm.tiles {
		if t.osm.x == osmX && t.osm.y == osmY {
			tileFound = true
			x = int(offsetX) + t.offsetX
			y = int(offsetY) + t.offsetY
			break
		}
	}
	sm.tilesMutex.Unlock()
	if tileFound != true {
		return 0, 0, errors.New("Tile not found")
	}
	return x, y, nil
}

func (sm *SlippyMap) ZoomIn(lat_deg, long_deg float64) (newsm *SlippyMap, err error) {
	// zoom in, with map centred on given lat/long (in degrees)
	newsm, err = sm.SetZoomLevel(sm.zoomLevel+1, lat_deg, long_deg)
	return newsm, err
}

func (sm *SlippyMap) ZoomOut(lat_deg, long_deg float64) (newsm *SlippyMap, err error) {
	// zoom in, with map centred on given lat/long (in degrees)
	newsm, err = sm.SetZoomLevel(sm.zoomLevel-1, lat_deg, long_deg)
	return newsm, err
}

func (sm *SlippyMap) SetZoomLevel(zoomLevel int, lat_deg, long_deg float64) (newsm *SlippyMap, err error) {
	// sets zoom level, with map centred on given lat/long (in degrees)

	// ensure we're within ZOOM_LEVEL_MAX & ZOOM_LEVEL_MIN
	if zoomLevel > ZOOM_LEVEL_MAX || zoomLevel < ZOOM_LEVEL_MIN {
		return &SlippyMap{}, errors.New("Requested zoom level unavailable")
	}

	// create a new slippymap centred on the requested lat/long, at the requested zoom level
	sm.mapSizeMutex.Lock()
	newsm = NewSlippyMap(sm.mapWidthPx, sm.mapHeightPx, zoomLevel, lat_deg, long_deg, sm.tileProvider)
	sm.mapSizeMutex.Unlock()

	// copy the current map image into the zoom previous level background image
	sm.Draw(newsm.zoomPrevLevelImg)

	// return the new slippymap and no error
	return newsm, nil
}

func NewSlippyMap(
	mapWidthPx, mapHeightPx, zoomLevel int,
	centreLat, centreLong float64,
	tileProvider TileProvider) (sm *SlippyMap) {

	log.Printf("Initialising SlippyMap at %0.4f/%0.4f, zoom level %d", centreLat, centreLong, zoomLevel)

	// determine the centre tile details
	centreTileOSMX, centreTileOSMY, pixelOffsetX, pixelOffsetY := gpsCoordsToTileInfo(centreLat, centreLong, zoomLevel)

	// create a new SlippyMap to return
	sm = &SlippyMap{
		img:              ebiten.NewImage(mapWidthPx, mapHeightPx), // initialise main image
		zoomPrevLevelImg: ebiten.NewImage(mapWidthPx, mapHeightPx), // initialise image of previous zoom level
		zoomLevel:        zoomLevel,                                // set zoom level
		tileProvider:     tileProvider,                             // set tile provider
		mapWidthPx:       mapWidthPx,
		mapHeightPx:      mapHeightPx,
		offsetMinimumX:   -(2 * TILE_WIDTH_PX),
		offsetMaximumX:   mapWidthPx + (2 * TILE_WIDTH_PX),
		offsetMinimumY:   -(2 * TILE_HEIGHT_PX),
		offsetMaximumY:   mapHeightPx + (2 * TILE_HEIGHT_PX),
	}

	// initialise the map with a centre tile
	sm.mapSizeMutex.Lock()
	centreTileOffsetX := (mapWidthPx / 2) - int(pixelOffsetX)
	centreTileOffsetY := (mapHeightPx / 2) - int(pixelOffsetY)
	sm.mapSizeMutex.Unlock()
	sm.makeTile(centreTileOSMX, centreTileOSMY, centreTileOffsetX, centreTileOffsetY)

	// force initial update
	sm.scheduleUpdate()
	sm.scheduleDraw()
	sm.Update(true)

	// return the slippymap
	return sm
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

func gpsCoordsToTileInfo(latDeg, longDeg float64, zoomLevel int) (tileX, tileY int, pixelOffsetX, pixelOffsetY float64) {
	// return OSM tile x/y coordinates (and pixel offset to the exact position) from lat/long

	// perform calculation as-per: https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames#Lon..2Flat._to_tile_numbers
	n := float64(calcN(zoomLevel))
	latRad := DegreesToRadians(latDeg)
	x := n * ((longDeg + 180.0) / 360.0)
	y := n * (1 - (math.Log(math.Tan(latRad)+secant(latRad)) / math.Pi)) / 2.0

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
