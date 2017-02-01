package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

var geoJsonFilePath string

func init() {
	flag.StringVar(&geoJsonFilePath, "geojson-file-path", "", "The path to a file containing a valid geojson feature with one geometry of type polygon")
	flag.Parse()
}

func main() {
	geoJsonFile := openGeoJsonFile()
	defer geoJsonFile.Close()
	isClockWise := IsGeoJsonFeaturePolygonClockWise(geoJsonFile)
	fmt.Printf("Is polygon clockwise: %v\n", isClockWise)
}

func openGeoJsonFile() io.ReadCloser {
	file, err := os.Open(geoJsonFilePath)
	if err != nil {
		exitWithError(fmt.Errorf("Failed to open the geojson file: %v", err))
	}
	return file
}

func IsGeoJsonFeaturePolygonClockWise(geoJsonFile io.Reader) bool {
	// decode the feature
	feature := Feature{}
	if err := json.NewDecoder(geoJsonFile).Decode(&feature); err != nil {
		exitWithError(fmt.Errorf("Failed to decode the geojson feature: %v", err))
	}
	coordinates := feature.Geometry.Coordinates[0]

	// normalize the coordinates
	var minLng, maxLng, minLat, maxLat float64
	for _, lngLat := range coordinates {
		if lngLat[0] < minLng {
			minLng = lngLat[0]
		}
		if lngLat[0] > maxLng {
			maxLng = lngLat[0]
		}
		if lngLat[1] < minLat {
			minLat = lngLat[1]
		}
		if lngLat[1] > maxLat {
			maxLat = lngLat[1]
		}
	}
	lngCrossesTheDateLine := ((maxLng - minLng) > 180)
	latCrossesTheDateLine := ((maxLat - minLat) > 80)

	var normalizedCoordinates [][]float64
	for _, lngLat := range coordinates {
		if lngCrossesTheDateLine && (lngLat[0] < 0) {
			lngLat[0] += 360
		}
		if latCrossesTheDateLine && (lngLat[1] < 0) {
			lngLat[1] += 160
		}
		normalizedCoordinates = append(normalizedCoordinates, lngLat)
	}

	// compute orientation
	var area float64
	for i, lngLat := range normalizedCoordinates {
		var nextLngLat []float64
		if i == (len(normalizedCoordinates) - 1) {
			nextLngLat = lngLat
		} else {
			nextLngLat = normalizedCoordinates[i+1]
		}
		area += (nextLngLat[0] - lngLat[0]) * (nextLngLat[1] + lngLat[1])
	}

	return area >= 0
}

func exitWithError(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}

type Feature struct {
	Geometry struct {
		Coordinates [][][]float64 `json:"coordinates"`
	} `json:"geometry"`
}
