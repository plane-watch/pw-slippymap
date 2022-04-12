package slippymap

import (
	"errors"
	"fmt"
	_ "image/png"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	TILE_WIDTH_PX  = 256 // tile width (as-per https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames)
	TILE_HEIGHT_PX = 256 // tile height (as-per https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames)
)

var (
	osm_url_prefix = 0 // used to round-robin connections to OSM servers, as-per https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames#Tile_servers
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
	tiles               []*MapTile    // map tiles
	mapWidthPx          int           // number of pixels wide
	mapHeightPx         int           // number of pixels high
	zoomLevel           int           // zoom level
	offsetMinimumX      int           // minimum X value for tiles
	offsetMinimumY      int           // minimum Y value for tiles
	offsetMaximumX      int           // maximum X value for tiles
	offsetMaximumY      int           // maximum Y value for tiles
	tileImageLoaderChan chan *MapTile // channel for loading of map tiles
	placeholderArtwork  *ebiten.Image // placeholder artwork for tile
	pathTileCache       string        // path to cache on disk
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
	// draws all tiles onto screen
	for _, t := range sm.tiles {

		dio := &ebiten.DrawImageOptions{}

		// move the image where it needs to be in the window
		dio.GeoM.Translate(float64((*t).offsetX), float64((*t).offsetY))

		// adjust transparency (for fade-in of tiles)
		dio.ColorM.Scale(1, 1, 1, (*t).alpha)

		// draw the tile
		screen.DrawImage((*t).img, dio)

		// debugging: print the OSM tile X/Y/Z
		dbgText := fmt.Sprintf("%d/%d/%d", (*t).osmX, (*t).osmY, (*t).zoomLevel)
		ebitenutil.DebugPrintAt(screen, dbgText, (*t).offsetX, (*t).offsetY)
	}
}

func (sm *SlippyMap) Update(deltaOffsetX, deltaOffsetY int, forceUpdate bool) {
	// Updates the map
	//  - Loads any missing tiles
	//  - Cleans up any tiles that are "out of bounds"
	//  - Moves tiles as-per deltaOffsetX/Y

	// clean up tiles off the screen
	for i, t := range sm.tiles {
		// if tile is out of bounds, remove it from slice
		if sm.isOutOfBounds((*t).offsetX, (*t).offsetY) {
			sm.tiles[i] = sm.tiles[len(sm.tiles)-1]
			sm.tiles = sm.tiles[:len(sm.tiles)-1]
			break
		}
	}

	// tile reposition & alpha increase if needed
	for _, t := range sm.tiles {

		// update offset if required (ie, user is dragging the map around)
		if (deltaOffsetX != 0 && deltaOffsetY != 0) || forceUpdate {
			(*t).offsetX = (*t).offsetX + deltaOffsetX
			(*t).offsetY = (*t).offsetY + deltaOffsetY
			// sm.tiles[i] = (*t)
		}

		// increase alpha channel (for fade in, if needed)
		if (*t).alpha < 1 {
			(*t).alpha = (*t).alpha + 0.05
		}
	}

	// new tiles created if required (just do one tile per update)
	for _, t := range sm.tiles {
		if sm.makeTileAbove(t) {
			break
		}
		if sm.makeTileToTheLeft(t) {
			break
		}
		if sm.makeTileToTheRight(t) {
			break
		}
		if sm.makeTileBelow(t) {
			break
		}
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

	// get tile artwork
	tilePath, err := cacheTile(osmX, osmY, sm.zoomLevel, sm.pathTileCache)
	if err != nil {
		log.Fatal(err)
	}

	// load the image from cache
	img, _, err := ebitenutil.NewImageFromFile(tilePath)
	if err != nil {
		log.Fatal(err)
	}
	t.img.DrawImage(img, nil)

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

func (sm *SlippyMap) GetLatLongAtPixel(x, y int) (lat_deg, long_deg float64, err error) {
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

	// get lat/long at top left corner of tile
	topLeftLat, topLeftLong := tileXYZtoGpsCoords(osmX, osmY, zoomLevel)

	// get lat/long of tile to the right
	_, topLeftLongRight := tileXYZtoGpsCoords(osmX+1, osmY, zoomLevel)

	// get lat/long of tile below
	topLeftLatBelow, _ := tileXYZtoGpsCoords(osmX, osmY+1, zoomLevel)

	// get lat/long degrees per pixel
	latPerPixel := (topLeftLatBelow - topLeftLat) / (TILE_HEIGHT_PX + 1)
	longPerPixel := (topLeftLongRight - topLeftLong) / (TILE_WIDTH_PX + 1)

	// get pixel offset within tile
	offsetX := x - topLeftX
	offsetY := y - topLeftY

	lat_deg = topLeftLat + (float64(offsetY) * latPerPixel)
	long_deg = topLeftLong + (float64(offsetX) * longPerPixel)

	return lat_deg, long_deg, nil

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

func NewSlippyMap(mapWidthPx, mapHeightPx, zoomLevel int, centreLat, centreLong float64, pathTileCache string) (sm SlippyMap, err error) {

	// load tile placeholder artwork
	tilePath := path.Join("assets", "map_tile_not_loaded.png")
	img, _, err := ebitenutil.NewImageFromFile(tilePath)
	if err != nil {
		log.Fatal(err)
	}

	// determine the centre tile details
	centreTileOSMX, centreTileOSMY, pixelOffsetX, pixelOffsetY := gpsCoordsToTileInfo(centreLat, centreLong, zoomLevel)

	// create a new SlippyMap to return
	sm = SlippyMap{
		zoomLevel: zoomLevel,
		// tileImageLoaderChan: tileImageLoaderChan,
		placeholderArtwork: img,
		pathTileCache:      pathTileCache,
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

func getOSMTileURL(x, y, z int) (url string) {
	// returns URL to open street map tile
	// load balance urls across servers as-per OSM guidelines: https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames#Tile_servers
	switch osm_url_prefix {
	case 0:
		url = fmt.Sprintf("http://a.tile.openstreetmap.org/%d/%d/%d.png", z, x, y)
		osm_url_prefix = 1
	case 1:
		url = fmt.Sprintf("http://b.tile.openstreetmap.org/%d/%d/%d.png", z, x, y)
		osm_url_prefix = 2
	case 2:
		url = fmt.Sprintf("http://c.tile.openstreetmap.org/%d/%d/%d.png", z, x, y)
		osm_url_prefix = 0
	}
	return url
}

func cacheTile(x, y, z int, pathTileCache string) (tilePath string, err error) {
	// if the tile at URL is not already cached, download it
	// return the local path to the tile in cache

	// TODO: this will eventually need refactoring. There's no retry mechanism if there's a failure.
	// We probably also want to do something with "if-modified-since" if the cached file is older than 7 days.
	// This is bare minimum to get functionality working

	// determine tile filename
	tileFile := fmt.Sprintf("%d_%d_%d.png", x, y, z)

	// determine full path to tile file
	tilePath = path.Join(pathTileCache, tileFile)

	// check if tile exists in cache
	if _, err := os.Stat(tilePath); errors.Is(err, os.ErrNotExist) {
		// tile does not exist in cache

		// determine OSM url
		url := getOSMTileURL(x, y, z)

		// log.Print("Downloading tile to cache:", url, "to", tilePath)

		// prepare http client
		client := &http.Client{}

		// prepare the request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return "", err
		}

		// set the header (requirement for using osm)
		req.Header.Set("User-Agent", "pw_slippymap/0.1 https://github.com/mikenye")

		// get the data
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		// check response code
		if resp.StatusCode != 200 {
			errText := fmt.Sprintf("Downloading %s returned: %s", url, resp.Status)
			return "", errors.New(errText)
		}

		// create the file
		out, err := os.Create(tilePath)
		if err != nil {
			return "", err
		}
		defer out.Close()

		// write data to file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return "", err
		}

		return tilePath, nil

	} else {
		// log.Print("Tile is cached:", tilePath)
		return tilePath, nil
	}
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

func degreesToRadians(d float64) (r float64) {
	// convert degrees to radians
	return d * (math.Pi / 180.0)
}

func radiansToDegrees(r float64) (d float64) {
	// convert radians to degrees
	return r * 180 / math.Pi
}

func gpsCoordsToTileInfo(lat_deg, long_deg float64, zoomLevel int) (tileX, tileY int, pixelOffsetX, pixelOffsetY float64) {
	// return OSM tile x/y coordinates (and pixel offset to the exact position) from lat/long

	// perform calculation as-per: https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames#Lon..2Flat._to_tile_numbers
	n := float64(calcN(zoomLevel))
	lat_rad := degreesToRadians(lat_deg)
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
	topLeftLat = radiansToDegrees(math.Atan(math.Sinh(math.Pi * (1 - 2*float64(y)/n))))
	topLeftLong = float64(x)/n*360.0 - 180.0
	return topLeftLat, topLeftLong
}
