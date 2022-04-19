package markers

import (
	"errors"
	"math"
	"regexp"
	"strconv"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

const (

	// None
	SVG_PATH_CMD_None = 0 // no command

	// MoveTo
	SVG_PATH_CMD_MoveTo   = 1 // M
	SVG_PATH_CMD_MoveToDx = 2 // m

	// LineTo
	SVG_PATH_CMD_LineTo        = 3 // L
	SVG_PATH_CMD_LineToDx      = 4 // l
	SVG_PATH_CMD_HorizLineTo   = 5 // H
	SVG_PATH_CMD_HorizLineToDx = 6 // h
	SVG_PATH_CMD_VertLineTo    = 7 // V
	SVG_PATH_CMD_VertLineToDx  = 8 // v

	// Cubic Bézier Curve
	SVG_PATH_CMD_CubicTo         = 9  // C
	SVG_PATH_CMD_CubicToDx       = 10 // c
	SVG_PATH_CMD_SmoothCubicTo   = 11 // S
	SVG_PATH_CMD_SmoothCubicToDx = 12 // s

	// Quadratic Bézier Curve
	SVG_PATH_CMD_QuadTo         = 13 // Q
	SVG_PATH_CMD_QuadToDx       = 14 // q
	SVG_PATH_CMD_SmoothQuadTo   = 15 // T
	SVG_PATH_CMD_SmoothQuadToDx = 16 // t

	// Elliptical Arc Curve
	SVG_PATH_CMD_ArcTo   = 17 // A
	SVG_PATH_CMD_ArcToDx = 18 // a

	// ClosePath
	SVG_PATH_CMD_ClosePath = 19 // Z,z
)

var (
	// precompile regex patterns
	reSVGCommand = regexp.MustCompile(`^\s*,?[MmLlHhVvCcSsQqTtAaZz]{1}`) // consumes a command
	reSVGNumber  = regexp.MustCompile(`^\s*,?-?[0-9]*\.?[0-9]*`)         // consumes a number
	reCommand    = regexp.MustCompile(`[MmLlHhVvCcSsQqTtAaZz]{1}`)       // return just the command component
	reFloat      = regexp.MustCompile(`[\-0-9\.]+`)                      // return just the number component
)

// SVG struct to assist with building the vector.Path
type SVG struct {
	x, y               float64     // the current x/y coordinates of the "pen"
	startx, starty     float64     // the initial x/y coordinates of the "pen"
	offsetX, offsetY   float64     // the offset x/y coordinates
	maxx, maxy         float64     // maximum x & y
	currentPathCommand int         // the current SVG command
	scale              float64     // the scale factor. Points from SVG are multiplied by this figure
	dc                 *gg.Context // the drawing context
}

func (svg *SVG) updateMaxXY(x, y float64) {
	if x > svg.maxx {
		svg.maxx = x
	}
	if y > svg.maxy {
		svg.maxy = y
	}
}

func (svg *SVG) moveTo(d string, dx bool) (remaining_d string, err error) {
	// Handles SVG_PATH_CMD_MoveTo / SVG_PATH_CMD_MoveToDx

	// consume x value from the path
	found, x, d, err := consumeNumber(d)
	if err != nil {
		return d, err
	}
	if found == false {
		return d, errors.New("moveTo: could not consume x")
	}

	// consume y value from the path
	found, y, d, err := consumeNumber(d)
	if err != nil {
		return d, err
	}
	if found == false {
		return d, errors.New("moveTo: could not consume y")
	}

	// scale points
	x = x * svg.scale
	y = y * svg.scale

	// if MoveToDX
	if dx {
		// fmt.Println("MoveToDx:", x, y)
		x = svg.x + x
		y = svg.y + y
	} else {
		x = x + svg.offsetX
		y = y + svg.offsetY
	}

	// fmt.Println("MoveTo:", x, y)

	// perform the path.MoveTo
	svg.dc.MoveTo(x, y)

	// update the current pen position
	svg.x = x
	svg.y = y
	svg.startx = x
	svg.starty = y

	svg.updateMaxXY(x, y)

	// return
	return d, nil
}

func (svg *SVG) closePath() {
	// Handles SVG_PATH_CMD_ClosePath

	svg.x = svg.startx
	svg.y = svg.starty
	// fmt.Println("ClosePath")
	svg.dc.LineTo(svg.startx, svg.starty)
}

func (svg *SVG) lineTo(d string, dx bool) (remaining_d string, err error) {
	// Handles SVG_PATH_CMD_LineTo & SVG_PATH_CMD_LineToDx

	// consume x value from the path
	found, x, d, err := consumeNumber(d)
	if err != nil {
		return d, err
	}
	if found == false {
		return d, errors.New("lineTo: could not consume x")
	}
	// consume y value from the path
	found, y, d, err := consumeNumber(d)
	if err != nil {
		return d, err
	}
	if found == false {
		return d, errors.New("lineTo: could not consume y")
	}

	// scale
	x = x * svg.scale
	y = y * svg.scale

	// if LineToDx
	if dx {
		x = svg.x + x
		y = svg.y + y
	} else {
		x = x + svg.offsetX
		y = y + svg.offsetY
	}

	// fmt.Println("LineTo:", x, y)

	// perform the path.LineTo
	svg.dc.LineTo(x, y)

	// update the current pen position
	svg.x = x
	svg.y = y

	svg.updateMaxXY(x, y)

	// return
	return d, nil
}

func (svg *SVG) vertLineTo(d string, dx bool) (remaining_d string, err error) {
	// Handles SVG_PATH_CMD_VertLineTo & SVG_PATH_CMD_VertLineToDx

	// consume y value from the path
	found, y, d, err := consumeNumber(d)
	if err != nil {
		return d, err
	}
	if found == false {
		return d, errors.New("vertLineTo: could not consume y")
	}

	// scale
	y = y * svg.scale

	// if LineToDx
	if dx {
		y = svg.y + y
	} else {
		y = y + svg.offsetY
	}

	// fmt.Println("VertLineTo:", svg.x, y)

	// perform the path command
	svg.dc.LineTo(svg.x, y)

	// update the current pen position
	svg.y = y

	svg.updateMaxXY(svg.x, y)

	// return
	return d, nil
}

func (svg *SVG) horizLineTo(d string, dx bool) (remaining_d string, err error) {
	// Handles SVG_PATH_CMD_HorizLineTo & SVG_PATH_CMD_HorizLineToDx

	// consume x value from the path
	found, x, d, err := consumeNumber(d)
	if err != nil {
		return d, err
	}
	if found == false {
		return d, errors.New("horizLineTo: could not consume x")
	}

	// scale
	x = x * svg.scale

	// if LineToDx
	if dx {
		x = svg.x + x
	} else {
		x = x + svg.offsetX
	}

	// fmt.Println("HorizLineTo:", x, svg.y)

	// perform the path command
	svg.dc.LineTo(x, svg.y)

	// update the current pen position
	svg.x = x

	svg.updateMaxXY(x, svg.y)

	// return
	return d, nil
}

func (svg *SVG) cubicTo(d string, dx bool) (remaining_d string, err error) {
	// Handles SVG_PATH_CMD_CubicTo & SVG_PATH_CMD_CubicToDx

	// consume x1 value from the path
	found, x1, d, err := consumeNumber(d)
	if err != nil {
		return d, err
	}
	if found == false {
		return d, errors.New("cubicTo: could not consume x1")
	}

	// consume y1 value from the path
	found, y1, d, err := consumeNumber(d)
	if err != nil {
		return d, err
	}
	if found == false {
		return d, errors.New("cubicTo: could not consume y1")
	}

	// consume x2 value from the path
	found, x2, d, err := consumeNumber(d)
	if err != nil {
		return d, err
	}
	if found == false {
		return d, errors.New("cubicTo: could not consume x2")
	}

	// consume y2 value from the path
	found, y2, d, err := consumeNumber(d)
	if err != nil {
		return d, err
	}
	if found == false {
		return d, errors.New("cubicTo: could not consume y2")
	}

	// consume x value from the path
	found, x, d, err := consumeNumber(d)
	if err != nil {
		return d, err
	}
	if found == false {
		return d, errors.New("cubicTo: could not consume x")
	}

	// consume y value from the path
	found, y, d, err := consumeNumber(d)
	if err != nil {
		return d, err
	}
	if found == false {
		return d, errors.New("cubicTo: could not consume y")
	}

	// scale
	x1 = x1 * svg.scale
	y1 = y1 * svg.scale
	x2 = x2 * svg.scale
	y2 = y2 * svg.scale
	x = x * svg.scale
	y = y * svg.scale

	// if MoveToDX
	if dx {
		x1 = svg.x + x1
		y1 = svg.y + y1
		x2 = svg.x + x2
		y2 = svg.y + y2
		x = svg.x + x
		y = svg.y + y
	} else {
		x1 = x1 + svg.offsetX
		y1 = y1 + svg.offsetY
		x2 = x2 + svg.offsetX
		y2 = y2 + svg.offsetY
		x = x + svg.offsetX
		y = y + svg.offsetY
	}

	// fmt.Println("CubicTo:", x1, y1, x2, y2, x, y)

	// perform the path.CubicTo
	svg.dc.CubicTo(x1, y1, x2, y2, x, y)

	// update the current pen position
	svg.x = x
	svg.y = y

	svg.updateMaxXY(x, y)

	// return
	return d, nil
}

func consumeCommand(d string) (commandFound bool, commandStr string, remaining_d string, err error) {
	// attempt to consume a command from the path given by d
	svgCmd := reSVGCommand.FindString(d)
	if len(svgCmd) > 0 {
		svgCmdChar := reCommand.FindString(svgCmd)
		if len(svgCmdChar) > 0 {
			// fmt.Println(svgCmdChar)
			remaining_d = d[len(svgCmd):]
			return true, svgCmdChar, remaining_d, nil
		} else {
			return false, "", d, errors.New("Command not supported!")
		}
	}
	return false, "", d, nil
}

func consumeNumber(d string) (numberFound bool, number float64, remaining_d string, err error) {
	// attempt to consume a number from the path given by d
	svgNum := reSVGNumber.FindString(d)
	if len(svgNum) > 0 {
		svgNumOnly := reFloat.FindString(svgNum)
		if len(svgNumOnly) > 0 {
			number, err := strconv.ParseFloat(svgNumOnly, 32)
			if err != nil {
				return false, 0, d, err
			} else {
				// fmt.Println(svgNumOnly)
				remaining_d = d[len(svgNum):]
				return true, number, remaining_d, nil
			}
		}
		return false, 0, d, nil
	}
	return false, 0, d, nil
}

// TODO: this is way too many args - maybe send a struct
func imgFromSVG(
	scale float64,
	d string,
	pathStroked, pathFilled, bgFilled bool,
	strokeWidth float64,
	strokeColour, fillColour, bgColour RGBA,
	offsetX, offsetY float64,
) (img *ebiten.Image, err error) {
	// Returns a drawing context from SVG path data
	// d: SVG path data (string) as-per: https://developer.mozilla.org/en-US/docs/Web/SVG/Attribute/d
	// scale: SVG coordinates are multiplied by scale (float32)
	// w,h: Width/height in pixels (int)
	//
	// Reference: https://developer.mozilla.org/en-US/docs/Web/SVG/Attribute/d
	//
	// NOTE: This does not yet implement all SVG path commands!!
	//       It is a work in progress.

	// define new svg object
	svg := SVG{
		x:                  offsetX,
		y:                  offsetY,
		startx:             offsetX,
		starty:             offsetY,
		offsetX:            offsetX,
		offsetY:            offsetY,
		currentPathCommand: SVG_PATH_CMD_None,
		scale:              scale,
		dc:                 gg.NewContext(500, 500),
	}

	// fill the background if required
	if bgFilled {
		svg.dc.SetRGBA(bgColour.r, bgColour.g, bgColour.b, bgColour.a)
		svg.dc.Clear()
	}

	for len(d) > 0 {

		// fmt.Println("---")
		// fmt.Println(d)

		// firstly attempt to consume command
		// grab the last char from the regex match
		var commandFound bool
		var commandStr string
		var err error
		commandFound, commandStr, d, err = consumeCommand(d)
		if err != nil {
			return ebiten.NewImage(1, 1), err
		}
		if commandFound {
			// fmt.Println(commandStr)
			switch commandStr {
			case "M":
				svg.currentPathCommand = SVG_PATH_CMD_MoveTo
			case "m":
				svg.currentPathCommand = SVG_PATH_CMD_MoveToDx
			case "L":
				svg.currentPathCommand = SVG_PATH_CMD_LineTo
			case "l":
				svg.currentPathCommand = SVG_PATH_CMD_LineToDx
			case "H":
				svg.currentPathCommand = SVG_PATH_CMD_HorizLineTo
			case "h":
				svg.currentPathCommand = SVG_PATH_CMD_HorizLineToDx
			case "V":
				svg.currentPathCommand = SVG_PATH_CMD_VertLineTo
			case "v":
				svg.currentPathCommand = SVG_PATH_CMD_VertLineToDx
			case "C":
				svg.currentPathCommand = SVG_PATH_CMD_CubicTo
			case "c":
				svg.currentPathCommand = SVG_PATH_CMD_CubicToDx
			case "S":
				svg.currentPathCommand = SVG_PATH_CMD_SmoothCubicTo
			case "s":
				svg.currentPathCommand = SVG_PATH_CMD_SmoothCubicToDx
			case "Q":
				svg.currentPathCommand = SVG_PATH_CMD_QuadTo
			case "q":
				svg.currentPathCommand = SVG_PATH_CMD_QuadToDx
			case "T":
				svg.currentPathCommand = SVG_PATH_CMD_SmoothQuadTo
			case "t":
				svg.currentPathCommand = SVG_PATH_CMD_SmoothQuadToDx
			case "A":
				svg.currentPathCommand = SVG_PATH_CMD_ArcTo
			case "a":
				svg.currentPathCommand = SVG_PATH_CMD_ArcToDx
			case "Z", "z":
				svg.currentPathCommand = SVG_PATH_CMD_ClosePath
			default:
				return ebiten.NewImage(1, 1), errors.New("Unknown SVG command")
			}
		}

		switch svg.currentPathCommand {

		case SVG_PATH_CMD_MoveTo:
			d, err = svg.moveTo(d, false)
		case SVG_PATH_CMD_MoveToDx:
			d, err = svg.moveTo(d, true)

		case SVG_PATH_CMD_CubicTo:
			d, err = svg.cubicTo(d, false)
		case SVG_PATH_CMD_CubicToDx:
			d, err = svg.cubicTo(d, true)

		case SVG_PATH_CMD_LineTo:
			d, err = svg.lineTo(d, false)
		case SVG_PATH_CMD_LineToDx:
			d, err = svg.lineTo(d, true)
		case SVG_PATH_CMD_HorizLineTo:
			d, err = svg.horizLineTo(d, false)
		case SVG_PATH_CMD_HorizLineToDx:
			d, err = svg.horizLineTo(d, true)
		case SVG_PATH_CMD_VertLineTo:
			d, err = svg.vertLineTo(d, false)
		case SVG_PATH_CMD_VertLineToDx:
			d, err = svg.vertLineTo(d, true)
		case SVG_PATH_CMD_ClosePath:
			svg.closePath()
		default:
			// fmt.Println(svg.currentPathCommand)
			return ebiten.NewImage(1, 1), errors.New("SVG error")
		}
		if err != nil {
			return ebiten.NewImage(1, 1), err
		}
	}

	// stroke the current path with this colour & width
	if pathStroked {
		svg.dc.SetRGBA(strokeColour.r, strokeColour.g, strokeColour.b, strokeColour.a)
		svg.dc.SetLineWidth(strokeWidth)
		svg.dc.StrokePreserve()
	}

	// fill the path
	if pathFilled {
		svg.dc.SetRGBA(fillColour.r, fillColour.g, fillColour.b, fillColour.a)
		svg.dc.Fill()
	}

	// "crop" the image to its exact size
	newImg := ebiten.NewImageFromImage(svg.dc.Image())
	img = ebiten.NewImage(int(math.Ceil(svg.maxx)), int(math.Ceil(svg.maxy)))
	img.DrawImage(newImg, nil)

	// return!
	return img, nil
}
