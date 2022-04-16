package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"path"
	"pw_slippymap/localdata"
	"pw_slippymap/markers"
	"pw_slippymap/slippymap"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	INIT_CENTRE_LAT     = -31.9523 // initial map centre lat
	INIT_CENTRE_LONG    = 115.8613 // initial map centre long
	INIT_ZOOM_LEVEL     = 9        // initial OSM zoom level
	INIT_WINDOW_SIZE    = 0.8      // percentage size of active screen
	ZOOM_COOLDOWN_TICKS = 5        // number of ticks to wait between zoom in/out ops
)

var (
	// Zoom
	zoomCoolDown int = 0

	// Window Size
	windowWidth  int
	windowHeight int

	// Mouse
	userMouse UserMouse

	// Debugging
	dbgMouseOverTileText string
	dbgMouseLatLongText  string

	// Path for tile cache
	pathTileCache string

	// Appears to be needed for DrawTriangles to work...
	emptyImage = ebiten.NewImage(1, 1)

	// Sprites
	VectorSprites map[string]VectorSprite
)

type VectorSprite struct {
	vs         []ebiten.Vertex
	is         []uint16
	maxX, maxY int
}

type Game struct {
	slippymap *slippymap.SlippyMap // hold the slippymap within the "game" object
}

type UserMouse struct {
	// struct that represents the user's mouse cursor
	prevX, prevY     int // previous tick X/Y
	currX, currY     int // current tick X/Y
	offsetX, offsetY int // offset of current X/Y from previous X/Y
}

func (um *UserMouse) update(x, y int) {
	// add "update" function to UserMouse struct
	// this function should be called on the game's "Update"
	// sets the current tick X/Y, previous tick X/Y, and the offset X/Y (from previous X/Y)
	um.prevX = um.currX
	um.prevY = um.currY
	um.currX = x
	um.currY = y
	um.offsetX = um.currX - um.prevX
	um.offsetY = um.currY - um.prevY
}

func (g *Game) Update() error {

	// temporarily commented out lots of stuff just to play with SVG artwork.

	// zoom: handle wheel
	_, dy := ebiten.Wheel()

	// zoom: honour the cooldown (helps when doing the two-finger-scroll on a macbook touchpad) & trigger on mousewheel y-axis
	if zoomCoolDown == 0 && dy != 0 {

		// zoom: get mouse cursor lat/long
		ctLat, ctLong, err := g.slippymap.GetLatLongAtPixel(userMouse.currX, userMouse.currY)
		if err != nil {
			// if error getting mouse cursor lat/long, log.
			log.Print("Cannot zoom")
		} else {
			// if no error getting mouse cursor lat/long, then do the zoom operation
			var newsm slippymap.SlippyMap
			var err error
			if dy > 0 {
				newsm, err = g.slippymap.ZoomIn(ctLat, ctLong)
				zoomCoolDown = ZOOM_COOLDOWN_TICKS
			} else if dy < 0 {
				newsm, err = g.slippymap.ZoomOut(ctLat, ctLong)
				zoomCoolDown = ZOOM_COOLDOWN_TICKS
			}
			if err != nil {
				log.Print("Error zooming")
			} else {
				g.slippymap = &newsm
			}
		}
	} else {
		// zoom: decrement zoom cool down counter to zero
		zoomCoolDown -= 1
		if zoomCoolDown < 0 {
			zoomCoolDown = 0
		}
	}

	// update the mouse cursor position
	userMouse.update(ebiten.CursorPosition())

	// handle dragging
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// update the map with the new offset
		g.slippymap.Update(userMouse.offsetX, userMouse.offsetY, false)
	} else {
		// otherwise update with no offset
		g.slippymap.Update(0, 0, false)
	}

	// no error to return
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// temporarily commented out lots of stuff just to play with SVG artwork.

	// draw map
	g.slippymap.Draw(screen)

	// osm attribution
	windowX, windowY := g.slippymap.GetSize()
	attributionArea := ebiten.NewImage(100, 20)
	attributionArea.Fill(color.Black)
	attributionAreaDio := &ebiten.DrawImageOptions{}
	attributionAreaDio.ColorM.Scale(1, 1, 1, 0.65)
	attributionAreaDio.GeoM.Translate(float64(windowX-100), float64(windowY-20))
	screen.DrawImage(attributionArea, attributionAreaDio)
	ebitenutil.DebugPrintAt(screen, "© OpenStreetMap", windowX-96, windowY-18)

	// debugging: darken area with debug text
	darkArea := ebiten.NewImage(windowWidth, 100)
	darkArea.Fill(color.Black)
	darkAreaDio := &ebiten.DrawImageOptions{}
	darkAreaDio.ColorM.Scale(1, 1, 1, 0.65)
	screen.DrawImage(darkArea, darkAreaDio)

	// debugging: show fps
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f  FPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS()), 0, 0)

	// debugging: show mouse position
	dbgMousePosTxt := fmt.Sprintf("Mouse position: %d, %d\n", userMouse.currX, userMouse.currY)
	ebitenutil.DebugPrintAt(screen, dbgMousePosTxt, 0, 15)

	// debugging: show zoom level
	dbgZoomLevelTxt := fmt.Sprintf("Zoom level: %d\n", g.slippymap.GetZoomLevel())
	ebitenutil.DebugPrintAt(screen, dbgZoomLevelTxt, 0, 30)

	// debugging: show tile moused over
	ctX, ctY, ctZ, err := g.slippymap.GetTileAtPixel(userMouse.currX, userMouse.currY)
	if err != nil {
		dbgMouseOverTileText = "Mouse over no tile"
	} else {
		dbgMouseOverTileText = fmt.Sprintf("Mouse over tile: %d/%d/%d", ctX, ctY, ctZ)
	}
	ebitenutil.DebugPrintAt(screen, dbgMouseOverTileText, 0, 45)

	// debugging: show lat/long under mouse
	ctLat, ctLong, err := g.slippymap.GetLatLongAtPixel(userMouse.currX, userMouse.currY)
	if err != nil {
		dbgMouseLatLongText = "Mouse over no tile"
	} else {
		dbgMouseLatLongText = fmt.Sprintf("Mouse over lat/long: %.4f/%.4f", ctLat, ctLong)
	}
	ebitenutil.DebugPrintAt(screen, dbgMouseLatLongText, 0, 60)

	// debugging: show number of tiles
	dbgNumTilesText := fmt.Sprintf("Tiles rendered: %d", g.slippymap.GetNumTiles())
	ebitenutil.DebugPrintAt(screen, dbgNumTilesText, 0, 75)

	// draw aircraft (TESTING)
	vectorOpts := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.EvenOdd,
	}
	drawOpts := &ebiten.DrawImageOptions{}

	drawOpts.GeoM.Translate(200, 0)
	img := ebiten.NewImage(VectorSprites["Airbus A380"].maxX, VectorSprites["Airbus A380"].maxY)
	img.DrawTriangles(VectorSprites["Airbus A380"].vs, VectorSprites["Airbus A380"].is, emptyImage.SubImage(image.Rect(0, 0, 1, 1)).(*ebiten.Image), vectorOpts)
	screen.DrawImage(img, drawOpts)

	drawOpts.GeoM.Translate(50, 0)
	img = ebiten.NewImage(VectorSprites["Fokker F100"].maxX, VectorSprites["Fokker F100"].maxY)
	img.DrawTriangles(VectorSprites["Fokker F100"].vs, VectorSprites["Fokker F100"].is, emptyImage.SubImage(image.Rect(0, 0, 1, 1)).(*ebiten.Image), vectorOpts)
	screen.DrawImage(img, drawOpts)

	drawOpts.GeoM.Translate(40, 0)
	img = ebiten.NewImage(VectorSprites["Pilatus PC12"].maxX, VectorSprites["Pilatus PC12"].maxY)
	img.DrawTriangles(VectorSprites["Pilatus PC12"].vs, VectorSprites["Pilatus PC12"].is, emptyImage.SubImage(image.Rect(0, 0, 1, 1)).(*ebiten.Image), vectorOpts)
	screen.DrawImage(img, drawOpts)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {

	// set slippymap size if window size changed
	smW, smH := g.slippymap.GetSize()
	if outsideWidth != smW || outsideHeight != smH {
		g.slippymap.SetSize(outsideWidth, outsideHeight)
	}

	// WindowSize returns 0,0 in non-desktop environments (eg wasm). Only rely on it if
	// the values aren't 0,0
	ew, eh := ebiten.WindowSize()
	if ew == 0 || eh == 0 {
		return outsideWidth, outsideHeight
	}
	return ew, eh
}

func failFatally(err error) {
	// handle errors by failing
	if err != nil {
		log.Fatal(err)
	}
}

func loadVectorSprites() {

	var (
		err        error
		vs         []ebiten.Vertex
		is         []uint16
		maxX, maxY int
	)

	// Airbus A380
	log.Println("Loading sprite of Airbus A380")
	vs, is, maxX, maxY, err = markers.InitMarker(markers.AIRBUS_A380_SVGPATH, markers.AIRBUS_A380_SCALE)
	failFatally(err)
	VectorSprites["Airbus A380"] = VectorSprite{
		vs:   vs,
		is:   is,
		maxX: maxX,
		maxY: maxY,
	}

	// Fokker F100
	log.Println("Loading sprite of Fokker F100")
	vs, is, maxX, maxY, err = markers.InitMarker(markers.FOKKER_F100_SVGPATH, markers.FOKKER_F100_SCALE)
	failFatally(err)
	VectorSprites["Fokker F100"] = VectorSprite{
		vs:   vs,
		is:   is,
		maxX: maxX,
		maxY: maxY,
	}

	// Pilatus PC12
	log.Println("Loading sprite of Pilatus PC12")
	vs, is, maxX, maxY, err = markers.InitMarker(markers.PILATUS_PC12_SVGPATH, markers.PILATUS_PC12_SCALE)
	failFatally(err)
	VectorSprites["Pilatus PC12"] = VectorSprite{
		vs:   vs,
		is:   is,
		maxX: maxX,
		maxY: maxY,
	}

}

func init() {
	VectorSprites = make(map[string]VectorSprite)
	emptyImage.Fill(color.White)
}

func main() {
	log.Print("Started")

	// determine starting window size
	// 80% of fullscreen
	screenWidth, screenHeight := ebiten.ScreenSizeInFullscreen()
	windowWidth = int(float64(screenWidth) * INIT_WINDOW_SIZE)
	windowHeight = int(float64(screenHeight) * INIT_WINDOW_SIZE)

	// set up initial window
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("plane.watch")

	// load sprites
	loadVectorSprites()

	tileProvider, err := tileProviderForOS()
	if err != nil {
		log.Fatal("could not initilalise tile provider because: ", err.Error())
	}

	// initialise map: initialise the new slippymap
	sm, err := slippymap.NewSlippyMap(windowWidth, windowHeight, INIT_ZOOM_LEVEL, INIT_CENTRE_LAT, INIT_CENTRE_LONG, tileProvider)
	if err != nil {
		log.Fatal(err)
	}

	// prepare "game"
	g := &Game{
		slippymap: &sm,
	}

	// run
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

// If we are running in WASM/JS, then the browser does all relevant tile caching for us.
// If running in desktop app mode, we need to cache the tiles ourselves
func tileProviderForOS() (slippymap.TileProvider, error) {
	if runtime.GOOS == "js" || false {
		return &slippymap.OSMTileProvider{}, nil
	}

	// try to get user home dir (for map cache)
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("User home dir:", userHomeDir)

	// create directory structure $HOME/.plane.watch if it doesn't exist
	pathRoot := path.Join(userHomeDir, ".plane.watch")
	err = localdata.SetupRoot(pathRoot)
	if err != nil {
		log.Fatal(err)
	}

	// create directory structure $HOME/.plane.watch/tilecache if it doesn't exist
	pathTileCache = path.Join(pathRoot, "tilecache")
	err = localdata.SetupTileCache(pathTileCache)
	if err != nil {
		log.Fatal(err)
	}

	return slippymap.NewCachedTileProvider(pathTileCache, &slippymap.OSMTileProvider{}), nil
}
