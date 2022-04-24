package slippymap

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

// TileProvider generates URLs (either https:// or file://) to the tile for the coord
type TileProvider interface {
	GetTileAddress(osm OSMTileID) (tilePath string, err error)
}

func NewCachedTileProvider(tileCachePath string, tileProvider TileProvider) *CachedTileProvider {

	return &CachedTileProvider{
		httpClient: &http.Client{
			Transport: newTransportWithLimitedConcurrency(),
		},
		tileProvider:  tileProvider,
		tileCachePath: tileCachePath,
	}
}

// CachedTileProvider wraps a real TileProvider to cache the contnets of the URLs on disk
type CachedTileProvider struct {
	httpClient    *http.Client
	tileProvider  TileProvider
	tileCachePath string
}

func (ctp *CachedTileProvider) GetTileAddress(osm OSMTileID) (tilePath string, err error) {
	// if the tile at URL is not already cached, download it
	// return the local path to the tile in cache

	// TODO: this will eventually need refactoring. There's no retry mechanism if there's a failure.
	// We probably also want to do something with "if-modified-since" if the cached file is older than 7 days.
	// This is bare minimum to get functionality working

	// determine tile filename
	tileFile := fmt.Sprintf("%d_%d_%d.png", osm.x, osm.y, osm.zoom)

	// determine full path to tile file
	tilePath = path.Join(ctp.tileCachePath, tileFile)

	// check if tile exists in cache
	if _, err := os.Stat(tilePath); errors.Is(err, os.ErrNotExist) {
		// tile does not exist in cache

		// determine OSM url
		url, err := ctp.tileProvider.GetTileAddress(osm)
		if err != nil {
			return "", err
		}

		// log.Print("Downloading tile to cache:", url, "to", tilePath)

		// prepare the request
		log.Printf("Downloading tile: %s", url)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return "", err
		}

		// set the header (requirement for using osm)
		req.Header.Set("User-Agent", "pw_slippymap/0.1 https://github.com/plane-watch/pw-slippymap")

		// get the data
		resp, err := ctp.httpClient.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		// check response code
		if resp.StatusCode != 200 {
			errText := fmt.Sprintf("Downloading %s returned: %s", url, resp.Status)
			return "", errors.New(errText)
		}

		// create the file
		out, err := os.Create(tilePath)
		if err != nil {
			return "", err
		}
		defer out.Close()

		// write data to file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return "", err
		}

		return tilePath, nil

	} else {
		// log.Print("Tile is cached:", tilePath)
		return tilePath, nil
	}
}

func newTransportWithLimitedConcurrency() http.RoundTripper {
	t := http.DefaultTransport.(*http.Transport).Clone()

	// We have 3 OSM CDN endpoints that we round robin, so max total conns is 3x this value
	t.MaxConnsPerHost = 1

	return t
}
