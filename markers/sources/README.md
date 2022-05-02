# SVG Sources

This directory contains Inkscape formatted SVGs.

## File Naming Convention

Files in this folder are named as: `ICAO (Model)`, where ICAO and Model is the code from the "ICAO aircraft type designators" table here: <https://en.wikipedia.org/wiki/List_of_aircraft_type_designators>

## Process

* Load Inkscape.
* Paste reference artwork into its own layer, lock the layer.
* Create a new layer. Use bezier curves & straight lines to trace the outline on the reference artwork. I find it helps if you set the stroke of the line to be a different colour and semi-transparent.
* When done, in Inkscape:
  * Go to Edit > XML Editor.
  * Find the XML `svg:path` elements.
  * Copy/paste the `d` value.

## Markers TODO

* BE20
* BL8
* C172
* CT4
* DH8A
* EC35
* F70
* GA8
* PC21
* R22
* SW4
