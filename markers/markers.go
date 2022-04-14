package markers

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	// SVG paths for aircraft
	AIRBUS_A380_SVGPATH = "m 244.73958,0 c -19.45177,2.9398148 -21.49332,76.729166 -21.49332,76.729166 v 35.718754 c -2.84181,7.02289 -10.27301,13.22916 -10.27301,13.22916 l -45.64879,37.48264 c 0.57163,-5.30799 0.32665,-9.71772 0.32665,-9.71772 0,-10.45268 -2.1232,-14.8624 -2.1232,-14.8624 h -19.35378 c -2.28653,2.36819 -2.04154,16.16899 -2.04154,16.16899 0.24498,13.22916 1.95987,18.4555 1.95987,18.4555 h 2.7765 l 1.03464,3.86133 -49.2966,36.3978 c 0.97994,-4.57305 0.89827,-11.10597 0.89827,-11.10597 -0.0817,-15.18904 -2.123197,-16.90393 -2.123197,-16.90393 H 79.946631 c -1.796554,1.0616 -2.1232,15.59735 -2.1232,15.59735 0.244984,13.55582 2.204861,19.02713 2.204861,19.02713 h 2.69483 l 1.306585,5.38966 -71.698817,52.99833 c -9.1460906,7.0229 -9.0644291,20.66037 -9.0644291,20.66037 l -0.3266461,22.13027 1.714892,-12.41255 80.6815842,-35.03279 c 0.408308,11.10596 1.388246,10.77932 1.388246,10.77932 1.388246,-0.48997 2.368184,-12.65754 2.368184,-12.65754 l 21.721969,-8.65612 c 0.0817,13.63747 1.22492,13.55581 1.22492,13.55581 2.04154,-0.89827 3.42978,-15.51569 3.42978,-15.51569 l 20.98701,-8.32947 c 0.32665,14.61741 1.55157,14.45409 1.55157,14.45409 2.85816,-5.55298 2.93982,-16.49563 2.93982,-16.49563 l 20.33372,-8.08449 c 0.24498,6.12461 1.38824,6.12461 1.38824,6.12461 1.30659,-0.0817 2.20486,-7.67618 2.20486,-7.67618 l 9.96271,-3.91975 10.94264,-2.93982 c 0.40831,6.28794 1.46991,6.12462 1.46991,6.12462 1.38825,0.0817 2.04154,-7.18622 2.04154,-7.18622 l 29.72479,-7.8395 c 0.73495,21.39532 4.16474,35.35943 4.16474,35.35943 v 47.28203 c -0.0817,8.32947 3.34812,32.17464 3.34812,32.17464 2.44985,10.20769 3.91976,16.74061 3.91976,16.74061 0.16332,4.89969 -5.71631,8.41114 -5.71631,8.41114 l -58.71463,43.52559 c -9.5544,7.4312 -11.10597,19.5171 -11.10597,19.5171 l -2.44985,12.65754 86.15291,-31.11304 c 1.63323,7.02289 4.81803,14.29076 4.81803,14.29076 0.24499,12.24923 1.30658,18.0472 1.30658,18.0472 1.0616,-5.79797 1.63323,-18.0472 1.63323,-18.0472 2.53151,-6.04295 4.73637,-14.45409 4.73637,-14.45409 l 86.07125,31.43969 -2.04154,-11.43261 c -2.93981,-15.43403 -11.59594,-21.06868 -11.59594,-21.06868 l -58.79629,-43.44393 c -4.81803,-3.02147 -5.55299,-8.24781 -5.55299,-8.24781 0.89828,-4.73637 3.83809,-16.74061 3.83809,-16.74061 3.1848,-16.49563 3.59311,-37.5643 3.59311,-37.5643 v -42.30067 c 3.26646,-15.10739 4.00142,-35.03279 4.00142,-35.03279 l 29.47981,8.00282 c 1.22491,7.18622 2.20486,7.0229 2.20486,7.0229 1.0616,0 1.55157,-6.20628 1.55157,-6.20628 l 10.86098,3.02148 10.04437,4.00141 c 0.81662,7.92117 2.04153,7.75785 2.04153,7.75785 1.46991,-0.48997 1.55157,-6.3696 1.55157,-6.3696 l 20.25206,8.24781 c 1.71489,16.65895 3.10314,16.41397 3.10314,16.41397 1.30658,-0.24499 1.55157,-14.53575 1.55157,-14.53575 l 20.98701,8.49279 c 1.79655,15.02573 3.26646,15.67902 3.26646,15.67902 1.22492,-1.87822 1.38825,-13.96412 1.38825,-13.96412 l 21.80362,9.14609 c 0.73496,12.49421 2.20486,12.33089 2.20486,12.33089 0.89828,-1.30659 1.71489,-10.94265 1.71489,-10.94265 l 80.51827,35.35944 1.55156,12.24923 -0.57163,-25.07009 c -0.81661,-12.82085 -8.81944,-17.80221 -8.81944,-17.80221 l -71.78048,-52.835 1.46991,-5.55299 h 2.69483 c 2.20486,-5.96129 2.1232,-18.94547 2.1232,-18.94547 -0.0817,-14.53576 -2.1232,-15.59735 -2.1232,-15.59735 h -19.5171 c -2.53151,6.94122 -2.04154,15.43403 -2.04154,15.43403 -0.0817,5.96128 0.81661,12.41255 0.81661,12.41255 l -48.99691,-36.17606 0.89828,-3.91975 h 2.69483 c 2.04154,-5.47132 2.1232,-18.4555 2.1232,-18.4555 0.4083,-12.7392 -2.20487,-16.08732 -2.20487,-16.08732 h -19.5171 c -2.53151,7.67618 -1.95988,16.08732 -1.95988,16.08732 0,5.14467 0.48997,8.49279 0.48997,8.49279 l -43.85223,-36.01273 c -7.10455,-4.73637 -12.00425,-14.78073 -12.00425,-14.78073 V 76.68017 C 262.29681,-4.246399 244.73958,0 244.73958,0 Z"
	AIRBUS_A380_SCALE   = 0.08
)

var (
	// Appears to be needed for DrawTriangles to work...
	emptyImage = ebiten.NewImage(500, 500)
)

// why is this needed!?
// DrawTriangles doesn't work unless we fill this image....
func init() {
	emptyImage.Fill(color.White)
}

func DrawAirbusA380(screen *ebiten.Image) {
	// Draws an Airbus A380 from SVG paths.

	var err error

	// Prepare the path object
	path := vector.Path{}

	// Airbus A380 Outline
	err = PathFromSVG(&path, AIRBUS_A380_SCALE, AIRBUS_A380_SVGPATH)
	if err != nil {
		log.Fatal(err)
	}

	// Get a list of verticies and indicies
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)

	// Set colours (not sure what SrcX/SrcY are doing... Copied from the Ebiten example: https://ebiten.org/examples/vector.html)
	for i := range vs {
		vs[i].SrcX = 0
		vs[i].SrcY = 0
		vs[i].ColorR = 0xff / float32(0xff)
		vs[i].ColorG = 0xff / float32(0xff)
		vs[i].ColorB = 0xff / float32(0xff)
		vs[i].ColorA = 1
	}

	op := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.EvenOdd,
	}

	// Finally, draw the triangles!
	screen.DrawTriangles(vs, is, emptyImage.SubImage(image.Rect(0, 0, 500, 500)).(*ebiten.Image), op)

}
