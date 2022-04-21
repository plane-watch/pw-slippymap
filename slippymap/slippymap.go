package slippymap

import (
	"errors"
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"time"

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
	DIRECTION_EAST  = 3
	DIRECTION_WEST  = 4
)

var (
	artworkLoaderQueue chan func() // queue for artwork loading
)

type MapTile struct {
	osmX                int           // OSM X
	osmY                int           // OSM Y
	zoomLevel           int           // OSM Zoom Level
	img                 *ebiten.Image // Image data
	offsetX             int           // top-left pixel location of tile
	offsetY             int           // top-right pixel location of tile
	alpha               float64       // tile transparency (for fade-in)
	tileRenderedToNorth bool
	tileRenderedToSouth bool
	tileRenderedToEast  bool
	tileRenderedToWest  bool
}

type SlippyMap struct {
	img *ebiten.Image // map image

	zoomLevelImgs map[int]*ebiten.Image // map images for each zoom level

	// re_render bool //do we need to re-render the image

	offsetX int // hold the current X offset
	offsetY int // hold the current Y offset
	// need_update bool // do we need to process Update()

	tiles []*MapTile // map tiles

	mapWidthPx  int // number of pixels wide
	mapHeightPx int // number of pixels high

	zoomLevel   int             // zoom level
	zoomTime    int64           // when the last zoom/scale operation was performed
	scaleLevels map[int]float64 // map to hold scaleLevel per zoom level

	// scaleOffset float64

	offsetMinimumX int // minimum X value for map tiles
	offsetMinimumY int // minimum Y value for map tiles
	offsetMaximumX int // maximum X value for map tiles
	offsetMaximumY int // maximum Y value for map tiles

	tileProvider TileProvider // the tile provider for the slippymap

	populatingTiles bool
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

	// render tiles onto sm.img only if required
	// if sm.re_render {
	for _, t := range sm.tiles {
		if (*t).zoomLevel == sm.zoomLevel {

			dio := &ebiten.DrawImageOptions{}

			// move the image where it needs to be in the window
			dio.GeoM.Translate(float64((*t).offsetX)-float64(sm.offsetMinimumX), float64((*t).offsetY)-float64(sm.offsetMinimumY))

			// // adjust transparency (for fade-in of tiles)
			// dio.ColorM.Scale(1, 1, 1, (*t).alpha)
			// dio.Filter = ebiten.FilterLinear

			// draw the tile
			// sm.img.DrawImage(t.img, dio)
			sm.zoomLevelImgs[(*t).zoomLevel].DrawImage(t.img, dio)

			// debugging: print the OSM tile X/Y/Z
			dbgText := fmt.Sprintf("%d/%d/%d", (*t).osmX, (*t).osmY, (*t).zoomLevel)
			// ebitenutil.DebugPrintAt(sm.img, dbgText, (*t).offsetX, (*t).offsetY)
			ebitenutil.DebugPrintAt(sm.zoomLevelImgs[(*t).zoomLevel], dbgText, (*t).offsetX, (*t).offsetY)
		}
	}

	// draw sm.img to the game screen
	drawOpts := ebiten.DrawImageOptions{}

	// scale the map between zoom levels
	drawOpts.GeoM.Translate(float64(sm.offsetMinimumX), float64(sm.offsetMinimumY))
	drawOpts.GeoM.Translate(-float64(sm.mapWidthPx)/2, -float64(sm.mapHeightPx)/2)
	drawOpts.GeoM.Scale(sm.GetScaleLevel(), sm.GetScaleLevel())
	drawOpts.GeoM.Translate(float64(sm.mapWidthPx)/2, float64(sm.mapHeightPx)/2)
	drawOpts.Filter = ebiten.FilterLinear

	// screen.DrawImage(sm.img, &drawOpts)
	screen.DrawImage(sm.zoomLevelImgs[sm.zoomLevel], &drawOpts)

	// don't re-render next pass (unless needed, see Update())
	// sm.re_render = false
}

func (sm *SlippyMap) Update(deltaOffsetX, deltaOffsetY int, forceUpdate bool) {
	// Updates the map
	//  - Loads any missing tiles
	//  - Cleans up any tiles that are "out of bounds"
	//  - Moves tiles as-per deltaOffsetX/Y

	// don't update unless required
	//   * offscreen tiles are being cleaned up; or
	//   * user has moved the map; or
	//   * tile fade-in happenning; or
	//   * new tiles were created
	// if deltaOffsetX != sm.offsetX || deltaOffsetY != sm.offsetY || sm.need_update {

	// tile reposition & alpha increase if needed
	for _, t := range sm.tiles {
		// update offset if required (ie, user is dragging the map around)
		if (deltaOffsetX != 0 && deltaOffsetY != 0) || forceUpdate {
			t.offsetX = t.offsetX + deltaOffsetX
			t.offsetY = t.offsetY + deltaOffsetY
			// sm.tiles[i] = (*t)
			// sm.re_render = true // re-render as visuals have changed
			// wereTilesMoved = true
		}

		// // increase alpha channel (for fade in, if needed)
		// if (*t).alpha < 1 {
		// 	(*t).alpha = (*t).alpha + TILE_FADEIN_ALPHA_PER_TICK
		// 	sm.re_render = true // re-render as visuals have changed
		// 	wereTilesAlphad = true
		// }
	}

	if sm.populatingTiles != true {
		sm.populatingTiles = true
		go sm.populateTiles()
	}
}

func (sm *SlippyMap) makeRelatedTile(existingTile *MapTile, direction int) (tileCreated bool) {

	// fmt.Println("---")

	var newTileOSMX int
	var newTileOSMY int
	var newTileZoomLevel int

	newTileOSMX = existingTile.osmX
	newTileOSMY = existingTile.osmY
	newTileZoomLevel = existingTile.zoomLevel

	newTileOffsetX := (*&existingTile.offsetX)
	newTileOffsetY := (*&existingTile.offsetY)

	// fmt.Println(newTileOSMX, newTileOSMY, newTileZoomLevel)

	// if to the north
	if direction == DIRECTION_NORTH {
		newTileOSMY = (*existingTile).osmY - 1
		newTileOffsetY = (*existingTile).offsetY - TILE_HEIGHT_PX
		// fmt.Println("NORTH")
	}

	// if to the south
	if direction == DIRECTION_SOUTH {
		newTileOSMY = (*existingTile).osmY + 1
		newTileOffsetY = (*existingTile).offsetY + TILE_HEIGHT_PX
		// fmt.Println("SOUTH")
	}

	// if to the west
	if direction == DIRECTION_WEST {
		newTileOSMX = (*existingTile).osmX - 1
		newTileOffsetX = (*existingTile).offsetX - TILE_WIDTH_PX
		// fmt.Println("WEST")
	}

	// if to the east
	if direction == DIRECTION_EAST {
		newTileOSMX = (*existingTile).osmX + 1
		newTileOffsetX = (*existingTile).offsetX + TILE_WIDTH_PX
		// fmt.Println("EAST")
	}

	// if above
	// TODO

	// if below
	// TODO

	// honour north edge of OSM map
	if newTileOSMY == -1 {
		// fmt.Println("EDGE OF MAP!")
		return false
	}

	// honour south edge of OSM map
	if newTileOSMY == int(math.Pow(2, float64((*existingTile).zoomLevel))) {
		// fmt.Println("EDGE OF MAP!")
		return false
	}

	// honour west of map
	if newTileOSMX == -1 {
		// fmt.Println("EDGE OF MAP!")
		return false
	}

	// honour east of map
	if newTileOSMX == int(math.Pow(2, float64((*existingTile).zoomLevel))) {
		// fmt.Println("EDGE OF MAP!")
		return false
	}

	// only load tiles in the immediate zoom levels
	if newTileZoomLevel < sm.zoomLevel-1 {
		// fmt.Println("NOT IMMEDIATE ZOOM!")
		return false
	}

	// only load tiles in the immediate zoom levels
	if newTileZoomLevel > sm.zoomLevel+1 {
		// fmt.Println("NOT IMMEDIATE ZOOM!")
		return false
	}

	// honour min zoom level
	if newTileZoomLevel < ZOOM_LEVEL_MIN {
		// fmt.Println("MIN ZOOM!")
		return false
	}

	// honour max zoom level
	if newTileZoomLevel > ZOOM_LEVEL_MAX {
		// fmt.Println("MAX ZOOM!")
		return false
	}

	// the tile already exists, bail out
	for _, t := range sm.tiles {
		if (*t).osmY == newTileOSMY && (*t).osmX == newTileOSMX && (*t).zoomLevel == newTileZoomLevel {
			// fmt.Println("TILE ALREADY EXISTS")
			return false
		}
	}

	// if the tile would not be out of bounds...
	if sm.isOutOfBounds(newTileOffsetX, newTileOffsetY, newTileZoomLevel) != true {

		// make the new tile
		sm.makeTile(newTileOSMX, newTileOSMY, newTileOffsetX, newTileOffsetY, newTileZoomLevel)
		// sm.re_render = true
		// fmt.Println("MADE TILE!")
		return true

	}
	// fmt.Println("OUT OF BOUNDS")
	return false

}

func (sm *SlippyMap) isOutOfBounds(pixelX, pixelY, zoomLevel int) (outOfBounds bool) {
	// returns true if the point defined by pixelX and pixelY is "out of bounds"
	// "out of bounds" means the point is outside the renderable size of the map
	// which is defined by sm.offset[Minimum|Maximum][X|Y].

	// w, h := sm.zoomLevelImgs[zoomLevel].Size()

	// if pixelX < 0 || pixelY < 0 || pixelX > w || pixelY > h {
	// 	return true
	// }

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

func (sm *SlippyMap) makeTile(osmX, osmY, offsetX, offsetY, zoomLevel int) {
	// Creates a new tile on the slippymap

	// fmt.Println(osmX, osmY, offsetX, offsetY, zoomLevel)

	// Create the tile object
	t := MapTile{
		osmX:      osmX,
		osmY:      osmY,
		offsetX:   offsetX,
		offsetY:   offsetY,
		zoomLevel: zoomLevel,
		img:       ebiten.NewImage(TILE_WIDTH_PX, TILE_WIDTH_PX),
	}
	t.img.Fill(color.Black)

	// enqueue loading of tile artwork
	artworkLoaderQueue <- func() {
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
	}

	// Add tile to slippymap
	sm.tiles = append(sm.tiles, &t)
	log.Printf("Populating tile: %d %d %d", t.osmX, t.osmY, t.zoomLevel)

}

func (sm *SlippyMap) SetSize(mapWidthPx, mapHeightPx int) {
	// updates the slippy map when window size is changed
	sm.mapWidthPx = mapWidthPx
	sm.mapHeightPx = mapHeightPx
	sm.offsetMinimumX = -int((float64(mapWidthPx)))
	sm.offsetMinimumY = -int((float64(mapHeightPx)))
	sm.offsetMaximumX = int(3 * float64(mapWidthPx))
	sm.offsetMaximumY = int(3 * float64(mapHeightPx))
}

func (sm *SlippyMap) GetSize() (mapWidthPx, mapHeightPx int) {
	// return the slippymap size in pixels
	return sm.mapWidthPx, sm.mapHeightPx
}

func (sm *SlippyMap) GetTileAtPixel(x, y int) (osmX, osmY, zoomLevel int, err error) {
	// returns the OSM tile X/Y/Z at pixel position x,y
	for _, t := range sm.tiles {
		if (*t).zoomLevel == sm.zoomLevel {
			if x >= (*t).offsetX && x < (*t).offsetX+TILE_WIDTH_PX {
				if y >= (*t).offsetY && y < (*t).offsetY+TILE_HEIGHT_PX {
					return (*t).osmX, (*t).osmY, (*t).zoomLevel, nil
				}
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

func (sm *SlippyMap) FindBestZoomLevel() (zoomed bool, newsm SlippyMap, err error) {

	// if the zoom level has been unchanged for half a second, find the best zoom level
	if time.Now().UnixMilli() >= sm.zoomTime+500 {

		// the best zoom level is the level with a scale factor closest to 1
		var bestZoomLevel int
		bestZoomLevelScaleFactor := math.MaxFloat64

		// find best zoom level
		for i := ZOOM_LEVEL_MIN; i <= ZOOM_LEVEL_MAX; i++ {
			if math.Abs(1-sm.scaleLevels[i]) < math.Abs(bestZoomLevelScaleFactor) {
				bestZoomLevel = i
				bestZoomLevelScaleFactor = 1 - sm.scaleLevels[i]
			}
		}

		// if best zoom level is not current zoom level
		if bestZoomLevel != sm.zoomLevel {
			scaleLevel := sm.scaleLevels[bestZoomLevel]
			centreLat, centreLong, err := sm.GetLatLongAtPixel(sm.mapWidthPx/2, sm.mapHeightPx/2)
			newsm, err = sm.SetZoomLevel(bestZoomLevel, centreLat, centreLong)
			newsm.SetScaleLevel(scaleLevel)
			return true, newsm, err
		}
	}

	return false, *sm, nil
}

func (sm *SlippyMap) Scale(dy float64) (newsm SlippyMap, err error) {
	// zoom map

	// adjust dy for smoothness
	dy = math.Round((dy/10)*1000) / 1000
	if dy > 0.2 {
		dy = 0.2
	}
	if dy < -0.2 {
		dy = -0.2
	}

	// adjust scale level based on dy
	sm.SetScaleLevel(sm.GetScaleLevel() + dy)

	// zoom out immediate
	if sm.GetScaleLevel() < 0.5 {
		scaleLevel := sm.scaleLevels[sm.zoomLevel-1]
		centreLat, centreLong, err := sm.GetLatLongAtPixel(sm.mapWidthPx/2, sm.mapHeightPx/2)
		newsm, err = sm.SetZoomLevel(sm.zoomLevel-1, centreLat, centreLong)
		newsm.SetScaleLevel(scaleLevel)
		return newsm, err
	}

	// zoom in immediate
	if sm.GetScaleLevel() > 2 {
		scaleLevel := sm.scaleLevels[sm.zoomLevel+1]
		centreLat, centreLong, err := sm.GetLatLongAtPixel(sm.mapWidthPx/2, sm.mapHeightPx/2)
		newsm, err = sm.SetZoomLevel(sm.zoomLevel+1, centreLat, centreLong)
		newsm.SetScaleLevel(scaleLevel)
		return newsm, err
	}

	// set time that zoom happened
	sm.zoomTime = time.Now().UnixMilli()

	return *sm, nil
}

func (sm *SlippyMap) SetZoomLevel(zoomLevel int, centreLat, centreLong float64) (newsm SlippyMap, err error) {
	// sets zoom level, with map centred on given lat/long (in degrees)

	// ensure we're within ZOOM_LEVEL_MAX & ZOOM_LEVEL_MIN
	if zoomLevel > ZOOM_LEVEL_MAX || zoomLevel < ZOOM_LEVEL_MIN {
		return SlippyMap{}, errors.New("Requested zoom level unavailable")
	}

	sm.zoomLevel = zoomLevel
	return *sm, nil
}

func (sm *SlippyMap) GetScaleLevel() float64 {
	// gets scale level for current zoom level
	return sm.scaleLevels[sm.zoomLevel]
}

func (sm *SlippyMap) SetScaleLevel(scaleLevel float64) {
	// sets scale level for current zoom level

	// update scale levels per zoom level
	sm.scaleLevels[sm.zoomLevel] = scaleLevel
	for i := sm.zoomLevel - 1; i >= ZOOM_LEVEL_MIN; i-- {
		sm.scaleLevels[i] = sm.scaleLevels[i+1] * 2
	}
	for i := sm.zoomLevel + 1; i <= ZOOM_LEVEL_MAX; i++ {
		sm.scaleLevels[i] = sm.scaleLevels[i-1] / 2
	}
}

func (sm *SlippyMap) populateTiles() {
	// go through each tile and populate

	// ensure we do the current zoom level first
	var zoomLevels [3]int
	zoomLevels[0] = sm.zoomLevel
	zoomLevels[1] = sm.zoomLevel + 1
	zoomLevels[1] = sm.zoomLevel - 1

	for i := 0; i <= len(zoomLevels); i++ {
		for _, t := range sm.tiles {
			if (*t).zoomLevel != i {

				// render tiles surrounding the current tile (if needed)
				if (*t).tileRenderedToNorth != true {
					(*t).tileRenderedToNorth = sm.makeRelatedTile(t, DIRECTION_NORTH)
				}
				if (*t).tileRenderedToSouth != true {
					(*t).tileRenderedToSouth = sm.makeRelatedTile(t, DIRECTION_SOUTH)
				}
				if (*t).tileRenderedToWest != true {
					(*t).tileRenderedToWest = sm.makeRelatedTile(t, DIRECTION_WEST)
				}
				if (*t).tileRenderedToEast != true {
					(*t).tileRenderedToEast = sm.makeRelatedTile(t, DIRECTION_EAST)
				}

				// if current tile is out of bounds, remove it from slice
				// if current tile is not in immediate zoom levels, remove it from slice
				if sm.isOutOfBounds((*t).offsetX, (*t).offsetY, (*t).zoomLevel) || (*t).zoomLevel < sm.zoomLevel-1 || (*t).zoomLevel > sm.zoomLevel+1 {
					log.Printf("Depopulating tile: %d %d %d", (*t).osmX, (*t).osmY, (*t).zoomLevel)
					sm.tiles[i] = sm.tiles[len(sm.tiles)-1]
					sm.tiles = sm.tiles[:len(sm.tiles)-1]
				}
			}
		}
	}
	sm.populatingTiles = false
}

func NewSlippyMap(mapWidthPx, mapHeightPx, zoomLevel int, centreLat, centreLong float64, tileProvider TileProvider) (sm SlippyMap, err error) {

	log.Printf("Initialising SlippyMap at %0.4f/%0.4f, zoom level %d", centreLat, centreLong, zoomLevel)

	// create a new SlippyMap to return
	sm = SlippyMap{
		// img:          ebiten.NewImage(3*mapWidthPx, 3*mapHeightPx), // initialise main image
		zoomLevel:    zoomLevel,    // set zoom level
		tileProvider: tileProvider, // set tile provider
		// loadTiles:    true,
		// re_render:    true,         // ensure first-time render
		// need_update:  true,         // ensure first-time update
		// lastZoomTime: time.Now(),
	}
	// init slippymap image layers of different zoom levels
	sm.zoomLevelImgs = make(map[int]*ebiten.Image)
	for i := ZOOM_LEVEL_MIN; i <= ZOOM_LEVEL_MAX; i++ {
		sm.zoomLevelImgs[i] = ebiten.NewImage(3*mapWidthPx, 3*mapHeightPx)
	}

	// init & update scaleLevels
	sm.scaleLevels = make(map[int]float64)
	sm.SetScaleLevel(1)

	// update size
	sm.SetSize(mapWidthPx, mapHeightPx)

	// initialise the map with centre tiles
	for i := ZOOM_LEVEL_MIN; i <= ZOOM_LEVEL_MAX; i++ {
		// determine the centre tile details
		centreTileOSMX, centreTileOSMY, pixelOffsetX, pixelOffsetY := gpsCoordsToTileInfo(centreLat, centreLong, i)
		centreTileOffsetX := (mapWidthPx / 2) - int(pixelOffsetX)
		centreTileOffsetY := (mapHeightPx / 2) - int(pixelOffsetY)
		sm.makeTile(centreTileOSMX, centreTileOSMY, centreTileOffsetX, centreTileOffsetY, i)
	}

	// force initial update
	// sm.Update(0, 0, true)

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

func init() {

	// artworkLoaderQueue is used to serialise loading tile artwork
	log.Println("Starting artwork loader worker")
	artworkLoaderQueue = make(chan func(), 1000)
	go func() {
		for job := range artworkLoaderQueue {
			// run the job
			job()
		}
	}()

}
