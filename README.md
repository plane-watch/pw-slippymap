# plane-watch/pw-slippymap

A [Slippy Map](https://wiki.openstreetmap.org/wiki/Slippy_Map) written in Go, runs in desktop mode or in js/wasm, for the plane.watch frontend (maybe).

## Current state

* Loads as a desktop app and displays the slippy map in a window
* You can pan around with the mouse by dragging
* You can zoom with mouse wheel

## Future

* Get planes on map
* Draw paths and stuff (polyline) on map

## Running locally

* Install prerequisites:
  * Linux: `apt-get install libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev`
* Clone the repo
* Change to the repo dir

Then...

### Desktop mode

* `go run main.go`

### WASM Mode

* `go install github.com/hajimehoshi/wasmserve@latest` - install wasmserve once
* `wasmserve .` - launches local dev server
* open [http://localhost:8080/](http://localhost:8080/) and wait for the app to compile+load
