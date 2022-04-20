package main

import (
	"fmt"
	"image/color"
	"log"
	"pw_slippymap/markers"
	"pw_slippymap/slippymap"

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

	// Map of markers (key = ICAO code), ref: https://en.wikipedia.org/wiki/List_of_aircraft_type_designators
	markerImages map[string]markers.Marker

	dbgMarkerRotateAngle float64
)

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

	// TEMPORARY/TESTING: rotate the plane sprites for testing
	dbgMarkerRotateAngle += 0.5
	if dbgMarkerRotateAngle >= 360 {
		dbgMarkerRotateAngle = 0
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
	m := markerImages["A388"]
	do := m.MarkerDrawOpts(dbgMarkerRotateAngle, 200, 5)
	screen.DrawImage(m.Img, &do)

	m = markerImages["F100"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 257, 5)
	screen.DrawImage(m.Img, &do)

	m = markerImages["PC12"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 300, 5)
	screen.DrawImage(m.Img, &do)

	m = markerImages["SF34"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 350, 5)
	screen.DrawImage(m.Img, &do)

	m = markerImages["E190"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 400, 5)
	screen.DrawImage(m.Img, &do)

	m = markerImages["DH8D"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 450, 5)
	screen.DrawImage(m.Img, &do)

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

func main() {

	var err error

	log.Print("Started")

	// load sprites
	markerImages, err = markers.InitMarkers()
	failFatally(err)

	// determine starting window size
	// 80% of fullscreen
	screenWidth, screenHeight := ebiten.ScreenSizeInFullscreen()
	windowWidth = int(float64(screenWidth) * INIT_WINDOW_SIZE)
	windowHeight = int(float64(screenHeight) * INIT_WINDOW_SIZE)

	// set up initial window
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("plane.watch")

	tileProvider, err := slippymap.TileProviderForOS()
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

	// In FPSModeVsyncOffMinimum, the game's Update and Draw are called only when
	// 1) new inputting is detected, or 2) ScheduleFrame is called.
	// In FPSModeVsyncOffMinimum, TPS is SyncWithFPS no matter what TPS is specified at SetMaxTPS.
	// ebiten.ScheduleFrame is called within SlippyMap.Update()
	// Should we make .Update() return a boolean that determines whether we schedule a frame in this packages Draw() function?
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMinimum)

	// run
	defer endProgram()
	log.Println("Starting UI")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func endProgram() {
	log.Println("Quitting")
}
