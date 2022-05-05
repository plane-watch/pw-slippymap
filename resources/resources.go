package resources

import (
	_ "embed"
	"log"

	"golang.org/x/image/font/opentype"
)

var Fonts map[string]*opentype.Font

//go:embed fonts/B612/B612-Regular.ttf
var fontB612Regular []byte

func failFatally(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {

	var err error

	Fonts = make(map[string]*opentype.Font)
	Fonts["B612-Regular"], err = opentype.Parse(fontB612Regular)
	failFatally(err)

}
