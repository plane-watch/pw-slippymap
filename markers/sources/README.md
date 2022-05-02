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

I've seen the markers below in ADS-B traffic that I've received, so I will prioritise creating these...

* B412

* R182
* DA42
* DH8A
* EC30
* EC35
* F70
* G115
* GA8
* K100
* M20P
* PA31
* PA46
* PC21
* PC24
* R22
* R44
* SLG2
* SONX
* SR22
* SW4
