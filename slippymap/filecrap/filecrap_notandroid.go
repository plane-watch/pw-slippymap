package filecrap

import (
	"bytes"
	"image"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/hajimehoshi/ebiten/v2"
)

func OpenFile(path string) (ebitenutil.ReadSeekCloser, error) {
	return os.Open(filepath.FromSlash(path))
}

func NewImageFromFile(path string) (*ebiten.Image, image.Image, error) {
	file, err := OpenFile(path)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		_ = file.Close()
	}()
	return NewImageFromReader(file)
}

func NewImageFromURL(url string) (*ebiten.Image, image.Image, error) {
	if !strings.HasPrefix(url, "http") {
		log.Printf("fetching from disk: %s\n", url)
		return NewImageFromFile(url)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not instantiate http request for tile")
	}

	req.Header.Set("User-Agent", "pw_slippymap/0.1 https://github.com/plane-watch/pw-slippymap")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not fetch image")
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not read tile image")
	}

	log.Printf("got tile %s. code: %d, type: %s, length: %d / %d\n", url, res.StatusCode, res.Header.Get("content-type"), res.ContentLength, len(b))

	return NewImageFromReader(bytes.NewReader(b))
}

func NewImageFromReader(reader io.Reader) (*ebiten.Image, image.Image, error) {
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, nil, err
	}
	img2 := ebiten.NewImageFromImage(img)
	return img2, img, err
}
