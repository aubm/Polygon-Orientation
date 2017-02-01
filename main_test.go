package main

import (
	"fmt"
	"os"
	"testing"
)

func TestIsGeoJsonFeaturePolygonClockWise(t *testing.T) {
	testCases := []struct {
		FilePath       string
		ExpectedResult bool
	}{
		{FilePath: "geojson-01.json", ExpectedResult: true},
		{FilePath: "geojson-02.json", ExpectedResult: false},
		{FilePath: "geojson-03.json", ExpectedResult: true},
		{FilePath: "geojson-04.json", ExpectedResult: false},
		{FilePath: "geojson-05.json", ExpectedResult: true},
		{FilePath: "geojson-06.json", ExpectedResult: false},
		{FilePath: "geojson-07.json", ExpectedResult: true},
		{FilePath: "geojson-08.json", ExpectedResult: false},
		{FilePath: "geojson-09.json", ExpectedResult: true},
		{FilePath: "geojson-10.json", ExpectedResult: true},
		{FilePath: "geojson-11.json", ExpectedResult: false},
		{FilePath: "geojson-12.json", ExpectedResult: true},
		{FilePath: "geojson-13.json", ExpectedResult: true},
	}

	for _, testCase := range testCases {
		geoJsonFile, err := os.Open(fmt.Sprintf("test/%v", testCase.FilePath))
		if err != nil {
			panic(err)
		}
		isClockWise := IsGeoJsonFeaturePolygonClockWise(geoJsonFile)
		if isClockWise != testCase.ExpectedResult {
			t.Errorf("For file %v, expected function to return %v, but got %v", testCase.FilePath, testCase.ExpectedResult, isClockWise)
		}
		geoJsonFile.Close()
	}
}
