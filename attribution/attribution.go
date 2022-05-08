package attribution

import (
	"image/color"
	"log"
	"pw_slippymap/resources"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	ATTRIBUTION_TEXT = "© OpenStreetMap"
)

type attrib struct {
	Img *ebiten.Image
}

var MapAttribution attrib

func init() {

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

	// get size of text
	attribRect := text.BoundString(ff, ATTRIBUTION_TEXT)

	// prepare new image
	MapAttribution.Img = ebiten.NewImage(attribRect.Dx()+10, attribRect.Dy()+10)
	MapAttribution.Img.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 128})

	// where to render text
	x := (MapAttribution.Img.Bounds().Dx() / 2) - (attribRect.Dx() / 2)
	y := (MapAttribution.Img.Bounds().Dy() / 2) + (attribRect.Dy() / 2)
	text.Draw(MapAttribution.Img, ATTRIBUTION_TEXT, ff, x, y, color.White)

}
