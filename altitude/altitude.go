package altitude

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"pw_slippymap/datasources/readsb_protobuf"
	"pw_slippymap/resources"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/mazznoer/colorgrad"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (

	// define min/max altitudes
	ALTITUDE_MIN_FT = 0
	ALTITUDE_MAX_FT = 40000
)

type FlightLevel struct {
	FL     int
	Feet   int
	Metres int
}

type AltitudeScale struct {
	Img   *ebiten.Image // image
	Width float64       // width in pixels
}

var (
	// Colour gradient for altitude
	altitudeColourGrad colorgrad.Gradient

	// colour for aircraft/vehicles on ground
	ColourGround = color.RGBA{R: 0, G: 102, B: 51, A: 255}
)

// flight levels (ref: https://mediawiki.ivao.aero/index.php?title=Semicircular_rules)
// key is FL
var FlightLevels = []FlightLevel{
	{FL: 10, Feet: 1000, Metres: 300},
	{FL: 20, Feet: 2000, Metres: 600},
	{FL: 30, Feet: 3000, Metres: 900},
	{FL: 40, Feet: 4000, Metres: 1200},
	{FL: 50, Feet: 5000, Metres: 1500},
	{FL: 60, Feet: 6000, Metres: 1850},
	{FL: 70, Feet: 7000, Metres: 2150},
	{FL: 80, Feet: 8000, Metres: 2450},
	{FL: 90, Feet: 9000, Metres: 2750},
	{FL: 100, Feet: 10000, Metres: 3050},
	{FL: 110, Feet: 11000, Metres: 3350},
	{FL: 120, Feet: 12000, Metres: 3650},
	{FL: 130, Feet: 13000, Metres: 3950},
	{FL: 140, Feet: 14000, Metres: 4250},
	{FL: 150, Feet: 15000, Metres: 4550},
	{FL: 160, Feet: 16000, Metres: 4900},
	{FL: 170, Feet: 17000, Metres: 5200},
	{FL: 180, Feet: 18000, Metres: 5500},
	{FL: 190, Feet: 19000, Metres: 5800},
	{FL: 200, Feet: 20000, Metres: 6100},
	{FL: 210, Feet: 21000, Metres: 6400},
	{FL: 220, Feet: 22000, Metres: 6700},
	{FL: 230, Feet: 23000, Metres: 7000},
	{FL: 240, Feet: 24000, Metres: 7300},
	{FL: 250, Feet: 25000, Metres: 7600},
	{FL: 260, Feet: 26000, Metres: 7900},
	{FL: 270, Feet: 27000, Metres: 8250},
	{FL: 280, Feet: 28000, Metres: 8550},
	{FL: 290, Feet: 29000, Metres: 8850},
	{FL: 300, Feet: 30000, Metres: 9150},
	{FL: 310, Feet: 31000, Metres: 9450},
	{FL: 330, Feet: 33000, Metres: 10050},
	{FL: 350, Feet: 35000, Metres: 10650},
	{FL: 370, Feet: 37000, Metres: 11300},
	{FL: 390, Feet: 39000, Metres: 11900},
	{FL: 410, Feet: 41000, Metres: 12500},
	{FL: 430, Feet: 43000, Metres: 13100},
	{FL: 450, Feet: 45000, Metres: 13700},
	{FL: 470, Feet: 47000, Metres: 14350},
	{FL: 490, Feet: 49000, Metres: 14950},
	{FL: 510, Feet: 51000, Metres: 15550},
}

func remap(x, inMin, inMax, outMin, outMax float64) float64 {
	// https://www.arduino.cc/reference/en/language/functions/math/map/
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

func AltitudeToColour(alt float64, airground readsb_protobuf.AircraftMeta_AirGround) (r, g, b float64) {

	if airground == readsb_protobuf.AircraftMeta_AG_GROUND {
		// if the plane/vehicle is on the ground

		r = float64(ColourGround.R) / 255
		g = float64(ColourGround.G) / 255
		b = float64(ColourGround.B) / 255

	} else {
		// if the plane is in the air, map colour to altitude gradient

		// honour min/max altitude
		if alt > ALTITUDE_MAX_FT {
			alt = ALTITUDE_MAX_FT
		}
		if alt < ALTITUDE_MIN_FT {
			alt = ALTITUDE_MIN_FT
		}

		// perform a "map" (https://www.arduino.cc/reference/en/language/functions/math/map/)
		// unsure what this is called in golang...

		r = altitudeColourGrad.At(alt).R
		g = altitudeColourGrad.At(alt).G
		b = altitudeColourGrad.At(alt).B
	}

	return r, g, b
}

func NewAltitudeScale(width float64) *AltitudeScale {

	h := 30.0 // height

	gndOffset := 10.0        // width of ground box within colour bar
	colourBarHeight := 10.0  // height of colour bar
	colourBarOffsetY := 20.0 // y offset of color bar

	// prep font
	faceOpts := opentype.FaceOptions{
		Size:    12,
		DPI:     72,
		Hinting: font.HintingNone,
	}
	ff, err := opentype.NewFace(resources.Fonts["B612-Regular"], &faceOpts)
	if err != nil {
		log.Fatal(err)
	}

	// create alt scale image
	img := ebiten.NewImage(int(width), int(h))
	img.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 128})

	// create ground square
	gndSquare := ebiten.NewImage(int(gndOffset), int(colourBarHeight))
	gndSquare.Fill(ColourGround)

	// draw ground square
	gndSquareDio := &ebiten.DrawImageOptions{}
	gndSquareDio.GeoM.Translate(0, colourBarOffsetY)
	img.DrawImage(gndSquare, gndSquareDio)

	// draw gradient
	for x := gndOffset; x < width; x++ {

		// map x pixel to altitude
		alt := remap(x, gndOffset, width, ALTITUDE_MIN_FT, ALTITUDE_MAX_FT)

		// get colour from altitude
		col := altitudeColourGrad.At(alt)

		// set pixel colours
		for y := colourBarOffsetY; y < h; y++ {
			img.Set(int(x), int(y), col)
		}
	}

	var newX, prevX int
	var newRect, prevRect image.Rectangle

	// draw "G" over ground box
	markerTxt := "G"
	newRect = text.BoundString(ff, markerTxt)
	newX = 5 - (newRect.Max.X / 2)
	text.Draw(img, markerTxt, ff, newX, newRect.Max.Y-newRect.Min.Y+2, color.White)
	prevX = newX
	prevRect = newRect

	// for each flight level...
	for _, v := range FlightLevels {
		if v.Feet < ALTITUDE_MIN_FT || v.Feet > ALTITUDE_MAX_FT {
			continue
		}

		// map altitude to x pixel
		x := remap(float64(v.Feet), ALTITUDE_MIN_FT, ALTITUDE_MAX_FT, gndOffset, width)

		// set text and get rectangle
		markerTxt := fmt.Sprintf("%d", v.Feet)
		newRect = text.BoundString(ff, markerTxt)
		newX = int(x) - (newRect.Max.X / 2)

		// print text if it won't overlap previous + tick
		if newX > prevX+prevRect.Dx()+10 && newX+newRect.Dx() <= int(width) {

			// draw text
			text.Draw(img, markerTxt, ff, newX, newRect.Max.Y-newRect.Min.Y+2, color.White)

			// draw tick
			for y := colourBarOffsetY; y < h; y++ {
				img.Set(int(x), int(y), color.Black)
			}

			prevX = newX
			prevRect = newRect
		}

	}

	return &AltitudeScale{
		Img:   img,
		Width: width,
	}
}

func makeAltitudeColourGrad() colorgrad.Gradient {
	grad, err := colorgrad.NewGradient().
		HtmlColors("gold", "hotpink", "darkturquoise").
		Domain(0, 40000).
		Build()
	if err != nil {
		log.Fatal(err)
	}
	return grad
}

func init() {

	// make colour gradient for altitude colours
	altitudeColourGrad = makeAltitudeColourGrad()

	// // make altitude scale image
	// AltitudeScale = makeAltitudeScale()
}
