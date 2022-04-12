package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"path"
	"pw_slippymap/localdata"
	"pw_slippymap/slippymap"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	INIT_CENTRE_LAT  = -31.9523 // initial map centre lat
	INIT_CENTRE_LONG = 115.8613 // initial map centre long
	INIT_ZOOM_LEVEL  = 9        // initial OSM zoom level
	INIT_WINDOW_SIZE = 0.8      // percentage size of active screen
	ZOOM_LEVEL_MAX   = 16       // maximum zoom level
	ZOOM_LEVEL_MIN   = 2        // minimum zoom level
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

var (
	zoomLevel            float64
	windowWidth          int
	windowHeight         int
	userMouse            UserMouse
	dbgMouseOverTileText string
	dbgMouseLatLongText  string
	pathTileCache        string
)

func (g *Game) Update() error {

	// zoom: handle wheel
	_, dy := ebiten.Wheel()
	zoomLevel += dy / 4 // /4 to decrease sensitivity

	// zoom: enforce limits
	if zoomLevel > ZOOM_LEVEL_MAX {
		zoomLevel = ZOOM_LEVEL_MAX
	}
	if zoomLevel < ZOOM_LEVEL_MIN {
		zoomLevel = ZOOM_LEVEL_MIN
	}

	// zoom: do the zooming
	if g.slippymap.GetZoomLevel() != int(math.Round(zoomLevel)) {

		// zoom: initialise new slippymap centred on mouse pos
		var sm slippymap.SlippyMap
		ctLat, ctLong, err := g.slippymap.GetLatLongAtPixel(userMouse.currX, userMouse.currY)
		if err != nil {
			log.Print("Cannot zoom!")
		} else {
			sm, err = slippymap.NewSlippyMap(windowWidth, windowHeight, int(math.Round(zoomLevel)), ctLat, ctLong, pathTileCache)
			if err != nil {
				log.Fatal(err)
			}
		}
		g.slippymap = &sm
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

	// draw map
	g.slippymap.Draw(screen)

	// debugging: darken area with debug text
	darkArea := ebiten.NewImage(240, 100)
	darkArea.Fill(color.Black)
	darkAreaDio := &ebiten.DrawImageOptions{}
	darkAreaDio.ColorM.Scale(1, 1, 1, 0.65)
	screen.DrawImage(darkArea, darkAreaDio)

	// debugging: show fps
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()), 0, 0)

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

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {

	// set slippymap size if window size changed
	smW, smH := g.slippymap.GetSize()
	if outsideWidth != smW || outsideHeight != smH {
		g.slippymap.SetSize(outsideWidth, outsideHeight)
	}

	// return window size
	return ebiten.WindowSize()
}

func main() {
	log.Print("Started")

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
	pathTileCache = path.Join(userHomeDir, ".plane.watch", "tilecache")
	err = localdata.SetupTileCache(pathTileCache)
	if err != nil {
		log.Fatal(err)
	}

	// determine starting window size
	// 80% of fullscreen
	screenWidth, screenHeight := ebiten.ScreenSizeInFullscreen()
	windowWidth = int(float64(screenWidth) * INIT_WINDOW_SIZE)
	windowHeight = int(float64(screenHeight) * INIT_WINDOW_SIZE)

	// set up initial window
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("plane.watch")

	// initial zoom level
	zoomLevel = INIT_ZOOM_LEVEL

	// initialise map: initialise the new slippymap
	sm, err := slippymap.NewSlippyMap(windowWidth, windowHeight, INIT_ZOOM_LEVEL, INIT_CENTRE_LAT, INIT_CENTRE_LONG, pathTileCache)
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
