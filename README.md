# geobin - working with bitmap indexed spatial data in go
geobin is a golang package that provide a simple way to work with spatial data sets (like geojson) using bitmap indexes.

The package works with uber/h3geo (https://h3geo.org) to provide aggregation, intersection, contains and joining operations with geometry (spatial attributes).

As a Geometry data structure the package uses github.com/paulmach/orb.

For the bitmap index, the package uses github.com/roaring/roaring.


## Sub packages
- spatial-bm - provides the bitmap index for a geometry. 
- bjoin - provides the data structure to store a result of the join two bitmap indexed data sets operation. 
- h3b - bitmap index for h3 cells
- h3f - functions to work with h3 cells in go, extends the uber/h3 (v4) package (https://github.com/uber/h3-go)
- orbf - provides functions and operations of the geometries.


# fixtures
For testing, the package uses OpenStreetMap data (https://openstreetmap.org)