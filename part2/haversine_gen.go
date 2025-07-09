package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"os"
	"strconv"
)

func main() {
	args := os.Args
	fmt.Println(args, len(args))
	if len(args) < 3 {
		log.Fatal("error: use main.go <n-points> <path> is needed")
		return
	}
	nPoints, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatal(err)
		return
	}
	dest := args[2]
	f, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}

	var Input []HaversineInput = []HaversineInput{}
	for range nPoints {
		newInput := HaversineInput{
			Point1: Point{Lat: getRandomLatitude(), Long: getRandomLongitude()},
			Point2: Point{Lat: getRandomLatitude(), Long: getRandomLongitude()},
		}
		Input = append(Input, newInput)
	}
	data, _ := json.MarshalIndent(Input, "", "\t")
	_, err = f.Write(data)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func getRandomLatitude() Latitude {
	val := getRandom() * 360 // [0 - 360]
	val = val - 180          // [-180 - 180]
	return Latitude(val)
}
func getRandomLongitude() Longitude {
	val := math.Acos(2*getRandom() - 1) // [0 - 180]
	val = val - 90                      // [-90 - 90]
	return Longitude(val)
}
func getRandom() float64 {
	return rand.Float64()
}
func DegToRad(deg float64) float64 {
	return deg * (math.Phi / 180)
}

// function return distance between 2 lat long coordinates
// all arguments should be in degree (will be converted to radian in the function)
// ref: https://stackoverflow.com/questions/27928/calculate-distance-between-two-latitude-longitude-points-haversine-formula and https://www.movable-type.co.uk/scripts/latlong.html
func GetDistanceByLatLong(lat1, long1, lat2, long2 float64) float64 {
	var rm float64 = 6371 //earth's mean raidus in kilo meters
	dLat := DegToRad(lat1 - lat2)
	dLong := DegToRad(long1 - long2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(DegToRad(lat1))*math.Cos(DegToRad(lat2))*math.Sin(dLong/2)*math.Sin(dLong/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := rm * c

	return d
}

type HaversineInput struct {
	Point1 Point `json:"point_1"`
	Point2 Point `json:"point_2"`
}

type Point struct {
	Lat  Latitude  `json:"latitude"`
	Long Longitude `json:"longitude"`
}

type Latitude float64
type Longitude float64
