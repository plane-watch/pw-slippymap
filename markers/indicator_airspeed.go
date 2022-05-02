package markers

import (
	"log"
	"math"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

type IndicatorAirspeed struct {
	bgImg         *ebiten.Image
	screwImg      *ebiten.Image
	dialBorderImg *ebiten.Image
	gaugeLinesImg *ebiten.Image
	textImg       *ebiten.Image
}

func (asi *IndicatorAirspeed) Draw(screen *ebiten.Image, x, y float64) {

	// gauge body
	drawOps := &ebiten.DrawImageOptions{}
	drawOps.GeoM.Translate(x, y)
	screen.DrawImage(asi.bgImg, drawOps)
	screen.DrawImage(asi.dialBorderImg, drawOps)

	// top left screw
	drawOps.GeoM.Reset()
	drawOps.GeoM.Translate(x, y)
	drawOps.GeoM.Translate(12, 12)
	screen.DrawImage(asi.screwImg, drawOps)

	// top right screw
	drawOps.GeoM.Reset()
	drawOps.GeoM.Translate(x, y)
	drawOps.GeoM.Translate(float64(asi.bgImg.Bounds().Max.X)-12-float64(asi.screwImg.Bounds().Max.X), 12)
	screen.DrawImage(asi.screwImg, drawOps)

	// bottom left screw
	drawOps.GeoM.Reset()
	drawOps.GeoM.Translate(x, y)
	drawOps.GeoM.Translate(12, float64(asi.bgImg.Bounds().Max.Y)-12-float64(asi.screwImg.Bounds().Max.Y))
	screen.DrawImage(asi.screwImg, drawOps)

	// bottom right screw
	drawOps.GeoM.Reset()
	drawOps.GeoM.Translate(x, y)
	drawOps.GeoM.Translate(float64(asi.bgImg.Bounds().Max.X)-12-float64(asi.screwImg.Bounds().Max.X), float64(asi.bgImg.Bounds().Max.Y)-12-float64(asi.screwImg.Bounds().Max.Y))
	screen.DrawImage(asi.screwImg, drawOps)

	// gauge lines
	drawOps.GeoM.Reset()
	drawOps.GeoM.Translate(x, y)
	screen.DrawImage(asi.gaugeLinesImg, drawOps)
}

func InitIndicatorAirspeed() (asi IndicatorAirspeed) {

	var err error

	asi = IndicatorAirspeed{}

	// Set default colours
	bgColour := RGBA{ // temp background colour.
		r: 0,
		g: 0.5,
		b: 0,
		a: 0.3,
	}
	strokeColour := RGBA{ // black
		r: 0,
		g: 0,
		b: 0,
		a: 1,
	}
	fillColour := RGBA{ // white
		r: 0.1,
		g: 0.1,
		b: 0.1,
		a: 1,
	}

	r := renderSVG{
		scale:        1,
		d:            "M 68.993196,0.50298807 H 44.179119 c -4.025658,0 -7.684171,0.83409343 -11.394427,2.19884203 L 15.882071,8.4533603 C 12.295797,9.7661927 10.278518,10.662761 9.2538679,13.896812 l -7.36467,21.709766 c -0.736467,2.337482 -1.09328826,3.5288 -1.09328826,6.263253 v 52.622086 c 0,3.330112 0.004598,5.027188 1.02924766,7.748913 l 6.9163857,20.10875 c 0.8965685,2.46557 1.985259,4.93113 4.354761,5.92376 l 22.286132,7.78093 c 2.177381,0.76849 5.547147,1.23424 7.864802,1.23424 H 95.86879 c 2.881828,0 5.69961,-0.81797 8.26124,-1.68252 l 18.95602,-6.53214 c 2.94587,-1.21678 4.70698,-2.20941 5.53951,-4.61093 l 7.49275,-21.4536 c 0.92859,-2.75375 1.31283,-5.05921 1.31283,-7.460733 V 42.619025 c 0,-2.497584 -0.4803,-5.859716 -1.40889,-8.229218 l -6.34002,-18.411675 c -1.31283,-3.554254 -2.3695,-5.795675 -6.05184,-7.236589 L 103.26548,1.793137 C 100.12749,0.8005076 98.110833,0.48030456 94.81212,0.48030456 Z",
		pathStroked:  true,
		pathFilled:   true,
		bgFilled:     false,
		strokeWidth:  5,
		strokeColour: strokeColour,
		fillColour:   fillColour,
		bgColour:     bgColour,
		offsetX:      2,
		offsetY:      2,
	}

	asi.bgImg, _, err = imgFromSVG(r)
	if err != nil {
		log.Fatal(err)
	}

	// screw body
	screwBody := gg.NewContext(12, 12)
	screwBody.DrawCircle(6, 6, 5)
	screwBody.SetRGBA(0, 0, 0, 1)
	screwBody.SetLineWidth(2)
	screwBody.StrokePreserve()
	screwBody.SetRGBA(0.2, 0.2, 0.2, 1)
	screwBody.Fill()
	asi.screwImg = ebiten.NewImageFromImage(screwBody.Image())

	// dial border
	dialBorder := gg.NewContext(asi.bgImg.Bounds().Max.X, asi.bgImg.Bounds().Max.Y)
	dialBorder.DrawCircle(float64(asi.bgImg.Bounds().Max.X)/2, float64(asi.bgImg.Bounds().Max.Y)/2, float64(asi.bgImg.Bounds().Max.Y)/2*0.9)
	dialBorder.SetRGBA(0, 0, 0, 1)
	dialBorder.SetLineWidth(2)
	dialBorder.StrokePreserve()
	dialBorder.SetRGBA(0.2, 0.2, 0.2, 1)
	dialBorder.Fill()
	asi.dialBorderImg = ebiten.NewImageFromImage(dialBorder.Image())

	// gauge lines
	gaugeLines := gg.NewContext(asi.bgImg.Bounds().Max.X, asi.bgImg.Bounds().Max.Y)

	// tick length
	tickRadius := (float64(asi.bgImg.Bounds().Max.Y) / 2 * 0.9) - 2
	tickLengthBig := 10.0
	tickLengthLittle := 8.0

	// draw big ticks
	for angle := 0.0 - 180; angle < 360-180; angle += 20 {
		angleRadians := angle * (math.Pi / 180.0)
		x1 := (float64(asi.bgImg.Bounds().Max.X) / 2) + (math.Sin(angleRadians) * tickRadius)
		y1 := (float64(asi.bgImg.Bounds().Max.Y) / 2) + (math.Cos(angleRadians) * tickRadius)
		x2 := (float64(asi.bgImg.Bounds().Max.X) / 2) + (math.Sin(angleRadians) * (tickRadius - tickLengthBig))
		y2 := (float64(asi.bgImg.Bounds().Max.Y) / 2) + (math.Cos(angleRadians) * (tickRadius - tickLengthBig))
		gaugeLines.DrawLine(x1, y1, x2, y2)
	}
	gaugeLines.SetRGBA(1, 1, 1, 1)
	gaugeLines.SetLineWidth(2)
	gaugeLines.StrokePreserve()

	// draw little ticks
	for angle := 30.0 - 180; angle < 340-180; angle += 20 {
		angleRadians := angle * (math.Pi / 180.0)
		x1 := (float64(asi.bgImg.Bounds().Max.X) / 2) + (math.Sin(angleRadians) * tickRadius)
		y1 := (float64(asi.bgImg.Bounds().Max.Y) / 2) + (math.Cos(angleRadians) * tickRadius)
		x2 := (float64(asi.bgImg.Bounds().Max.X) / 2) + (math.Sin(angleRadians) * (tickRadius - tickLengthLittle))
		y2 := (float64(asi.bgImg.Bounds().Max.Y) / 2) + (math.Cos(angleRadians) * (tickRadius - tickLengthLittle))
		gaugeLines.DrawLine(x1, y1, x2, y2)
	}
	gaugeLines.SetLineWidth(1)
	gaugeLines.StrokePreserve()

	asi.gaugeLinesImg = ebiten.NewImageFromImage(gaugeLines.Image())

	// Gauge Text: Name

	return asi

}
