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
	"sort"
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

	// APP STATES -----------------------------------------

	// normal states
	STATE_STARTUP = 0
	STATE_RUN     = 1

	// debug states
	STATE_DEBUG_MARKERS_STARTUP = 110
	STATE_DEBUG_MARKERS_RUN     = 111

	// ----------------------------------------------------
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
	dbgMarkerRotateAngle float64
)

type UserInterface struct {

	// app state
	state      int
	stateMutex sync.Mutex

	// slippymap
	slippymap    *slippymap.SlippyMap
	tileProvider *slippymap.TileProvider // tile provider for slippymap

	// user input
	touchIDs []ebiten.TouchID
	strokes  map[*userinput.Stroke]struct{}

	// aircraft db
	aircraftDb *datasources.AircraftDB

	// markers
	aircraftMarkers      *map[string]markers.Marker
	groundVehicleMarkers *map[string]markers.Marker
}

func (ui *UserInterface) getState() int {
	ui.stateMutex.Lock()
	defer ui.stateMutex.Unlock()
	return ui.state
}

func (ui *UserInterface) setState(state int) {
	ui.stateMutex.Lock()
	defer ui.stateMutex.Unlock()
	ui.state = state
}

func (ui *UserInterface) updateStroke(stroke *userinput.Stroke) {
	stroke.Update()
	if !stroke.IsReleased() {
		return
	}

}

func (ui *UserInterface) loadSprites() {
	// load sprites
	aircraftMarkers, err := markers.InitMarkers(markers.Aircraft)
	failFatally(err)
	groundVehicleMarkers, err := markers.InitMarkers(markers.GroundVehicles)
	failFatally(err)

	ui.aircraftMarkers = &aircraftMarkers
	ui.groundVehicleMarkers = &groundVehicleMarkers
}

func (ui *UserInterface) handleWindowResize() {
	// set slippymap size if window size changed
	smW, smH := ui.slippymap.GetSize()
	wsW, wsH := ebiten.WindowSize()
	if wsW != smW || wsH != smH {
		newsm := ui.slippymap.SetSize(wsW, wsH)
		ui.slippymap = newsm
	}
}

func (ui *UserInterface) handleMouseWheel() {

	// handle wheel
	_, dy := ebiten.Wheel()

	// zoom: honour the cooldown (helps when doing the two-finger-scroll on a macbook touchpad) & trigger on mousewheel y-axis
	if zoomCoolDown == 0 && dy != 0 {

		mouseX, mouseY := ebiten.CursorPosition()

		// zoom: get mouse cursor lat/long
		ctLat, ctLong, err := ui.slippymap.GetLatLongAtPixel(mouseX, mouseY)
		if err != nil {
			// if error getting mouse cursor lat/long, log.
			log.Print("Cannot zoom")
		} else {
			// if no error getting mouse cursor lat/long, then do the zoom operation
			var newsm *slippymap.SlippyMap
			var err error
			if dy > 0 {
				newsm, err = ui.slippymap.ZoomIn(ctLat, ctLong)
				zoomCoolDown = ZOOM_COOLDOWN_TICKS
			} else if dy < 0 {
				newsm, err = ui.slippymap.ZoomOut(ctLat, ctLong)
				zoomCoolDown = ZOOM_COOLDOWN_TICKS
			}
			if err != nil {
				log.Print("Error zooming")
			} else {
				ui.slippymap = newsm
			}
		}
	} else {
		// zoom: decrement zoom cool down counter to zero
		zoomCoolDown -= 1
		if zoomCoolDown < 0 {
			zoomCoolDown = 0
		}
	}
}

func (ui *UserInterface) handleMouseMovement() bool {
	// mouse / touch dragging
	forceUpdate := false
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s := userinput.NewStroke(&userinput.MouseStrokeSource{})
		s.SetDraggingObject(ui.slippymap)
		ui.strokes[s] = struct{}{}
	}
	ui.touchIDs = inpututil.AppendJustPressedTouchIDs(ui.touchIDs[:0])
	for _, id := range ui.touchIDs {
		s := userinput.NewStroke(&userinput.TouchStrokeSource{id})
		s.SetDraggingObject(ui.slippymap)
		ui.strokes[s] = struct{}{}
	}
	for s := range ui.strokes {
		ui.updateStroke(s)
		if s.IsReleased() {
			delete(ui.strokes, s)
		}
		mouseX, mouseY := s.PositionDiffFromPrevious()
		ui.slippymap.MoveBy(mouseX, mouseY)
		forceUpdate = true
	}
	return forceUpdate
}

func (ui *UserInterface) Update() error {

	switch ui.getState() {

	case STATE_STARTUP:

		log.Println("Starting UI")
		ui.loadSprites()
		ui.slippymap = slippymap.NewSlippyMap(windowWidth, windowHeight, INIT_ZOOM_LEVEL, INIT_CENTRE_LAT, INIT_CENTRE_LONG, *ui.tileProvider)
		ui.setState(STATE_RUN)

	case STATE_RUN:

		// handle window resize
		ui.handleWindowResize()

		// handle mouse wheel
		ui.handleMouseWheel()

		// handle mouse/touch dragging for map movement
		forceUpdate := ui.handleMouseMovement()

		// update the slippymap
		ui.slippymap.Update(forceUpdate)

	case STATE_DEBUG_MARKERS_STARTUP:
		ui.loadSprites()
		ui.setState(STATE_DEBUG_MARKERS_RUN)

	case STATE_DEBUG_MARKERS_RUN:
		dbgMarkerRotateAngle += 0.5
		if dbgMarkerRotateAngle >= 360 {
			dbgMarkerRotateAngle = 0
		}

	default:
		log.Fatal("Invalid state in ui.Update!")
	}

	// no error to return
	return nil
}

func (ui *UserInterface) drawAircraftMarkers(screen *ebiten.Image, mouseX, mouseY int) (mouseOverMarkerText string) {

	mouseOverMarkerText = "No marker"

	// determine draw order
	// currently we order based on ICAO
	// TODO: change order based on altitude
	aircraftMap := ui.aircraftDb.GetAircraft()
	aircraftIcaos := make([]int, 0, len(aircraftMap))
	for k := range aircraftMap {
		aircraftIcaos = append(aircraftIcaos, k)
	}
	sort.Ints(aircraftIcaos)

	// for each plane we know about
	for _, k := range aircraftIcaos {

		v := aircraftMap[k]

		// skip planes that aren't sending a position
		// TODO: what about planes actually at 0,0?
		if v.Lat == 0 && v.Long == 0 {
			continue
		}

		var aircraftMarker markers.Marker

		// determine marker based on category (https://wiki.jetvision.de/wiki/Radarcape:Software_Features#Aircraft_categories)
		switch v.Category {
		case 0xC1, 0xC2:

			// don't draw ground vehicles if zoom level less than 13
			if ui.slippymap.GetZoomLevel() < 13 {
				continue
			}

			aircraftMarker = markers.GetMarker("4WD", ui.groundVehicleMarkers)

		default:

			// don't draw aircraft on ground if "idle" on ground unless zoom level is less than 13
			if ui.slippymap.GetZoomLevel() < 13 && v.GroundSpeed < 30 {
				continue
			}

			aircraftMarker = markers.GetMarker(v.AircraftType, ui.aircraftMarkers)
		}

		// determine where the marker will be drawn
		aircraftX, aircraftY, err := ui.slippymap.LatLongToPixel(v.Lat, v.Long)
		if err != nil {
			// log.Printf("Error plotting %X: %s", k, err)
			// plane is probably off the visible map, or not sending a position
		} else {

			// prepare the draw options for the marker
			aircraftDrawOpts := aircraftMarker.MarkerDrawOpts(float64(v.Track), float64(aircraftX), float64(aircraftY))

			// get fill colour from altitude
			r, g, b := markers.AltitudeToColour(float64(aircraftMap[k].AltBaro), aircraftMap[k].AirGround)

			// invert colours
			r = 1 - r
			g = 1 - g
			b = 1 - b

			// apply fill
			aircraftDrawOpts.ColorM.Invert()
			aircraftDrawOpts.ColorM.Translate(r, g, b, 0)
			aircraftDrawOpts.ColorM.Invert()

			// draw it
			screen.DrawImage(aircraftMarker.Img, &aircraftDrawOpts)

			// work out if mouse is over marker image
			topLeftX := -aircraftMarker.CentreX + float64(aircraftX)
			topLeftY := -aircraftMarker.CentreY + float64(aircraftY)
			btmRightX := topLeftX + float64(aircraftMarker.Img.Bounds().Max.X)
			btmRightY := topLeftY + float64(aircraftMarker.Img.Bounds().Max.Y)
			if mouseX >= int(topLeftX) && mouseX <= int(btmRightX) {
				if mouseY >= int(topLeftY) && mouseY <= int(btmRightY) {
					// if it is, determine if it is inside the shape
					if aircraftMarker.PointInsideMarker(float64(mouseX)-topLeftX, float64(mouseY)-topLeftY) {
						mouseOverMarkerText = fmt.Sprintf("ICAO: %X, Callsign: %s, Type: %s, Category: %X, Alt: %d, Gs: %d, AirGround: %s", k, v.Callsign, v.AircraftType, v.Category, v.AltBaro, v.GroundSpeed, v.AirGround.String())
					}
				}
			}
		}
	}

	return mouseOverMarkerText

}

func (ui *UserInterface) debugDrawMarkers(screen *ebiten.Image) {

	// draw aircraft (TESTING)
	m := (*ui.aircraftMarkers)["A388"]
	do := m.MarkerDrawOpts(dbgMarkerRotateAngle, 203, 25)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "A388", 203, 40)

	m = (*ui.aircraftMarkers)["F100"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 257, 25)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "F100", 257, 40)

	m = (*ui.aircraftMarkers)["PC12"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 300, 25)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "PC12", 300, 40)

	m = (*ui.aircraftMarkers)["SF34"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 350, 25)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "SF34", 350, 40)

	m = (*ui.aircraftMarkers)["E190"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 400, 25)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "E190", 400, 40)

	m = (*ui.aircraftMarkers)["DH8D"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 450, 25)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "DH8D", 450, 40)

	m = (*ui.aircraftMarkers)["A320"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 500, 25)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "A320", 500, 40)

	m = (*ui.aircraftMarkers)["B738"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 550, 25)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "B738", 550, 40)

	m = (*ui.aircraftMarkers)["B77W"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 600, 25)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "B77W", 600, 40)

	m = (*ui.aircraftMarkers)["B77L"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 650, 25)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "B77L", 650, 40)

	m = (*ui.aircraftMarkers)["HAWK"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 700, 25)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "HAWK", 700, 40)

	m = (*ui.aircraftMarkers)["B788"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 750, 25)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "B788", 750, 40)

	m = (*ui.aircraftMarkers)["RV9"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 800, 25)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "RV9", 800, 40)

	m = (*ui.aircraftMarkers)["SW3"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 850, 25)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "SW3", 850, 40)

	m = (*ui.aircraftMarkers)["SW4"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 900, 25)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "SW4", 900, 40)

	m = (*ui.groundVehicleMarkers)["4WD"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 950, 25)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "4WD", 950, 40)

	m = (*ui.aircraftMarkers)["DA42"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 1000, 25)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "DA42", 1000, 40)

	m = (*ui.aircraftMarkers)["SONX"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 1050, 25)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "SONX", 1050, 40)

	m = (*ui.aircraftMarkers)["GLEX"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 1100, 25)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "GLEX", 1100, 40)

	m = (*ui.aircraftMarkers)["H25B"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 1150, 25)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "H25B", 1150, 40)

	m = (*ui.aircraftMarkers)["C182"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 1200, 25)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "C182", 1200, 40)

	m = (*ui.aircraftMarkers)["TWEN"]
	do = m.MarkerDrawOpts(dbgMarkerRotateAngle, 1250, 25)
	screen.DrawImage(m.Img, &do)
	ebitenutil.DebugPrintAt(screen, "TWEN", 1250, 40)
}

func (ui *UserInterface) Draw(screen *ebiten.Image) {

	mouseX, mouseY := ebiten.CursorPosition()
	windowX, windowY := ui.slippymap.GetSize()

	switch ui.getState() {

	case STATE_STARTUP:
		ebitenutil.DebugPrint(screen, "Loading...")

	case STATE_RUN:
		// draw map
		ui.slippymap.Draw(screen)

		// draw aircraft
		mouseOverMarkerText := ui.drawAircraftMarkers(screen, mouseX, mouseY)

		// draw altitude scale
		altitudeScaleDio := &ebiten.DrawImageOptions{}
		altitudeScaleDio.GeoM.Translate((float64(windowX)/2)-(float64(markers.AltitudeScale.Bounds().Max.X)/2), float64(windowY)-float64(markers.AltitudeScale.Bounds().Max.Y))
		screen.DrawImage(markers.AltitudeScale, altitudeScaleDio)

		// draw osm attribution
		attributionArea := ebiten.NewImage(100, 20)
		attributionArea.Fill(color.Black)
		attributionAreaDio := &ebiten.DrawImageOptions{}
		attributionAreaDio.ColorM.Scale(1, 1, 1, 0.65)
		attributionAreaDio.GeoM.Translate(float64(windowX-100), float64(windowY-20))
		screen.DrawImage(attributionArea, attributionAreaDio)
		ebitenutil.DebugPrintAt(screen, "© OpenStreetMap", windowX-96, windowY-18)

		// debugging: darken area with debug text
		darkArea := ebiten.NewImage(windowWidth, 115)
		darkArea.Fill(color.Black)
		darkAreaDio := &ebiten.DrawImageOptions{}
		darkAreaDio.ColorM.Scale(1, 1, 1, 0.65)
		screen.DrawImage(darkArea, darkAreaDio)

		// debugging: show fps
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f  FPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS()), 0, 0)

		// debugging: show mouse position
		dbgMousePosTxt := fmt.Sprintf("Mouse position: %d, %d\n", mouseX, mouseY)
		ebitenutil.DebugPrintAt(screen, dbgMousePosTxt, 0, 15)

		// debugging: show zoom level
		dbgZoomLevelTxt := fmt.Sprintf("Zoom level: %d\n", ui.slippymap.GetZoomLevel())
		ebitenutil.DebugPrintAt(screen, dbgZoomLevelTxt, 0, 30)

		// debugging: show tile moused over
		ctX, ctY, ctZ, err := ui.slippymap.GetTileAtPixel(mouseX, mouseY)
		if err != nil {
			dbgMouseOverTileText = "Mouse over no tile"
		} else {
			dbgMouseOverTileText = fmt.Sprintf("Mouse over tile: %d/%d/%d", ctX, ctY, ctZ)
		}
		ebitenutil.DebugPrintAt(screen, dbgMouseOverTileText, 0, 45)

		// debugging: show lat/long under mouse
		ctLat, ctLong, err := ui.slippymap.GetLatLongAtPixel(mouseX, mouseY)
		if err != nil {
			dbgMouseLatLongText = "Mouse over no tile"
		} else {
			dbgMouseLatLongText = fmt.Sprintf("Mouse over lat/long: %.4f/%.4f", ctLat, ctLong)
		}
		ebitenutil.DebugPrintAt(screen, dbgMouseLatLongText, 0, 60)

		// debugging: show number of tiles
		dbgNumTilesText := fmt.Sprintf("Tiles rendered: %d", ui.slippymap.GetNumTiles())
		ebitenutil.DebugPrintAt(screen, dbgNumTilesText, 0, 75)

		// debugging: show number of tiles
		dbgMouseOverMarkerText := fmt.Sprintf("Mouse over marker: %s", mouseOverMarkerText)
		ebitenutil.DebugPrintAt(screen, dbgMouseOverMarkerText, 0, 90)

	default:
		log.Fatal("Invalid state in ui.Draw!")
	}

}

func (ui *UserInterface) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {

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

	// init aircraftdb
	adb := datasources.NewAircraftDB(60)
	log.Printf("readsb database version: %d", datasources.GetReadsbDBVersion())

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

	// if readsb aircraft.db datasource has been specified, initialise it
	if conf.readsbAircraftProtobufUrl != "" {
		log.Printf("Datasource: readsb-protobuf at url: %s", conf.readsbAircraftProtobufUrl)
		go datasources.ReadsbProtobuf(conf.readsbAircraftProtobufUrl, adb)
	}

	// prepare "game"
	ui := &UserInterface{
		aircraftDb:   adb,
		strokes:      map[*userinput.Stroke]struct{}{},
		tileProvider: &tileProvider,
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
	if err := ebiten.RunGame(ui); err != nil {
		log.Fatal(err)
	}
}

func endProgram() {
	log.Println("Quitting")
}
