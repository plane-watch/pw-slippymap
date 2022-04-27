package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"pw_slippymap/datasources"
	"pw_slippymap/markers"
	"pw_slippymap/slippymap"
	"pw_slippymap/userinput"
	"sync"

	"github.com/akamensky/argparse"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

	// Debugging
	dbgMouseOverTileText string
	dbgMouseLatLongText  string

	// Map of markers (key = ICAO code), ref: https://en.wikipedia.org/wiki/List_of_aircraft_type_designators
	markerImages map[string]markers.Marker

	dbgMarkerRotateAngle float64
)

type Game struct {

	// graphics
	slippymap *slippymap.SlippyMap // hold the slippymap within the "game" object

	// user input
	touchIDs []ebiten.TouchID
	strokes  map[*userinput.Stroke]struct{}

	// aircraft db
	aircraftDb *datasources.AircraftDB

	// markers
	aircraftMarkers *map[string]markers.Marker
}

func (g *Game) updateStroke(stroke *userinput.Stroke) {
	stroke.Update()
	if !stroke.IsReleased() {
		return
	}

}

func (g *Game) Update() error {

	// temporarily commented out lots of stuff just to play with SVG artwork.

	// zoom: handle wheel
	_, dy := ebiten.Wheel()

	// zoom: honour the cooldown (helps when doing the two-finger-scroll on a macbook touchpad) & trigger on mousewheel y-axis
	if zoomCoolDown == 0 && dy != 0 {

		mouseX, mouseY := ebiten.CursorPosition()

		// zoom: get mouse cursor lat/long
		ctLat, ctLong, err := g.slippymap.GetLatLongAtPixel(mouseX, mouseY)
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

	// mouse / touch dragging
	forceUpdate := false
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s := userinput.NewStroke(&userinput.MouseStrokeSource{})
		s.SetDraggingObject(*g.slippymap)
		g.strokes[s] = struct{}{}
	}
	g.touchIDs = inpututil.AppendJustPressedTouchIDs(g.touchIDs[:0])
	for _, id := range g.touchIDs {
		s := userinput.NewStroke(&userinput.TouchStrokeSource{id})
		s.SetDraggingObject(*g.slippymap)
		g.strokes[s] = struct{}{}
	}
	for s := range g.strokes {
		g.updateStroke(s)
		if s.IsReleased() {
			delete(g.strokes, s)
		}
		mouseX, mouseY := s.PositionDiffFromPrevious()
		g.slippymap.MoveBy(mouseX, mouseY)
		forceUpdate = true
	}

	g.slippymap.Update(forceUpdate)

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

	// draw planes ===========================================================
	// TODO: move this into a function

	// for each plane we know about
	for _, v := range g.aircraftDb.GetAircraft() {

		var aircraftMarker markers.Marker

		// determine image
		if _, ok := (*g.aircraftMarkers)[v.AircraftType]; ok {
			// use marker that matches aircraft type if found
			aircraftMarker = (*g.aircraftMarkers)[v.AircraftType]
		} else {
			// default marker
			aircraftMarker = (*g.aircraftMarkers)["A388"]
		}

		// determine where the marker will be drawn
		aircraftX, aircraftY, err := g.slippymap.LatLongToPixel(v.Lat, v.Long)
		if err != nil {
			// log.Printf("Error plotting %X: %s", k, err)
			// plane is probably off the visible map, or not sending a position
		} else {

			// determine how the marker will be drawn
			aircraftDrawOpts := &ebiten.DrawImageOptions{}
			// move so centre of marker is at 0,0
			aircraftDrawOpts.GeoM.Translate(-aircraftMarker.CentreX, -aircraftMarker.CentreY)
			// rotate to match track
			aircraftDrawOpts.GeoM.Rotate(slippymap.DegreesToRadians(float64(v.Track)))
			// move to actual position
			aircraftDrawOpts.GeoM.Translate(float64(aircraftX), float64(aircraftY))

			// draw it
			screen.DrawImage(aircraftMarker.Img, aircraftDrawOpts)
		}
	}
	// end draw planes =======================================================

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
	mouseX, mouseY := ebiten.CursorPosition()
	dbgMousePosTxt := fmt.Sprintf("Mouse position: %d, %d\n", mouseX, mouseY)
	ebitenutil.DebugPrintAt(screen, dbgMousePosTxt, 0, 15)

	// debugging: show zoom level
	dbgZoomLevelTxt := fmt.Sprintf("Zoom level: %d\n", g.slippymap.GetZoomLevel())
	ebitenutil.DebugPrintAt(screen, dbgZoomLevelTxt, 0, 30)

	// debugging: show tile moused over
	ctX, ctY, ctZ, err := g.slippymap.GetTileAtPixel(mouseX, mouseY)
	if err != nil {
		dbgMouseOverTileText = "Mouse over no tile"
	} else {
		dbgMouseOverTileText = fmt.Sprintf("Mouse over tile: %d/%d/%d", ctX, ctY, ctZ)
	}
	ebitenutil.DebugPrintAt(screen, dbgMouseOverTileText, 0, 45)

	// debugging: show lat/long under mouse
	ctLat, ctLong, err := g.slippymap.GetLatLongAtPixel(mouseX, mouseY)
	if err != nil {
		dbgMouseLatLongText = "Mouse over no tile"
	} else {
		dbgMouseLatLongText = fmt.Sprintf("Mouse over lat/long: %.4f/%.4f", ctLat, ctLong)
	}
	ebitenutil.DebugPrintAt(screen, dbgMouseLatLongText, 0, 60)

	// debugging: show number of tiles
	dbgNumTilesText := fmt.Sprintf("Tiles rendered: %d", g.slippymap.GetNumTiles())
	ebitenutil.DebugPrintAt(screen, dbgNumTilesText, 0, 75)

	// // draw aircraft (TESTING)
	m := (*g.aircraftMarkers)["A388"]
	do := m.MarkerDrawOpts(dbgMarkerRotateAngle, 203, 5)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "A388", 203, 40)

	// m = markerImages["F100"]
	m = (*g.aircraftMarkers)["F100"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 257, 5)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "F100", 257, 40)

	// m = markerImages["PC12"]
	m = (*g.aircraftMarkers)["PC12"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 300, 5)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "PC12", 300, 40)

	// m = markerImages["SF34"]
	m = (*g.aircraftMarkers)["SF34"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 350, 5)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "SF34", 350, 40)

	// m = markerImages["E190"]
	m = (*g.aircraftMarkers)["E190"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 400, 5)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "E190", 400, 40)

	// m = markerImages["DH8D"]
	m = (*g.aircraftMarkers)["DH8D"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 450, 5)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "DH8D", 450, 40)

	// m = markerImages["A320"]
	m = (*g.aircraftMarkers)["A320"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 500, 5)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "A320", 500, 40)

	// m = markerImages["B738"]
	m = (*g.aircraftMarkers)["B738"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 550, 5)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "B738", 550, 40)

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

type runtimeConfiguration struct {
	readsbAircraftProtobufUrl string
}

func processCommandLine() runtimeConfiguration {
	// process the command line

	// create new parser object
	parser := argparse.NewParser("pw-slippymap", "front-end for plane.watch")

	// readsb-protobuf aircraft.pb URL
	readsbAircraftProtobufUrl := parser.String("", "aircraftpburl", &argparse.Options{Required: false, Help: "Uses readsb-protobuf aircraft.pb as a data source. Eg: 'http://1.2.3.4/data/aircraft.pb'"})

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
	}

	// Prepare runtime
	conf := runtimeConfiguration{}

	// if --aircraftpburl set, add to runtime conf
	if *readsbAircraftProtobufUrl != "" {
		conf.readsbAircraftProtobufUrl = *readsbAircraftProtobufUrl
	}

	return conf
}

func main() {
	var err error

	// process the command line
	conf := processCommandLine()

	// process readsb json files
	var startupWg sync.WaitGroup
	go datasources.BuildReadsbAircraftsJSON(&startupWg)

	// init aircraftdb
	adb := datasources.NewAircraftDB()
	log.Printf("readsb database version: %d", datasources.GetReadsbDBVersion())

	// load sprites
	aircraftMarkers, err := markers.InitMarkers()
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
	sm := slippymap.NewSlippyMap(windowWidth, windowHeight, INIT_ZOOM_LEVEL, INIT_CENTRE_LAT, INIT_CENTRE_LONG, tileProvider)

	// wait for all parallel startup jobs
	startupWg.Wait()

	// if readsb aircraft.db datasource has been specified, initialise it
	if conf.readsbAircraftProtobufUrl != "" {
		log.Printf("Datasource: readsb-protobuf at url: %s", conf.readsbAircraftProtobufUrl)
		go datasources.ReadsbProtobuf(conf.readsbAircraftProtobufUrl, adb)
	}

	// prepare "game"
	g := &Game{
		slippymap:       &sm,
		aircraftDb:      adb,
		aircraftMarkers: &aircraftMarkers,
		strokes:         map[*userinput.Stroke]struct{}{},
	}

	// In FPSModeVsyncOffMinimum, the game's Update and Draw are called only when
	// 1) new inputting is detected, or 2) ScheduleFrame is called.
	// In FPSModeVsyncOffMinimum, TPS is SyncWithFPS no matter what TPS is specified at SetMaxTPS.
	// ebiten.ScheduleFrame is called within SlippyMap.Update()
	// Should we make .Update() return a boolean that determines whether we schedule a frame in this packages Draw() function?
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMinimum)
	ebiten.SetMaxTPS(60)

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
