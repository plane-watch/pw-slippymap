# SVG Sources

This directory contains Inkscape formatted SVGs.

## Process

* Load Inkscape.
* Paste reference artwork into its own layer, lock the layer.
* Create a new layer. Use bezier curves & straight lines to trace the outline on the reference artwork. I find it helps if you set the stroke of the line to be a different colour and semi-transparent.
* When done, in Inkscape:
  * Go to Edit > XML Editor.
  * Find the XML `svg:path` elements.
  * Copy/paste the `d` value.