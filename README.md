# plane-watch/pw-slippymap

A [Slippy Map](https://wiki.openstreetmap.org/wiki/Slippy_Map) written in Go, eventually destined to be compiled to WebAssembly, for the plane.watch frontend (maybe).

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
* `go run main.go`
