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

var nLessOne int = 0

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

	var sum float64
	var sumDist float64
	for _, point := range Input {
		sum += math.Floor(float64(point.Point1.Lat))
		sum += math.Floor(float64(point.Point1.Long))
		sum += math.Floor(float64(point.Point2.Lat))
		sum += math.Floor(float64(point.Point2.Long))
		sumDist += GetDistanceByLatLong(float64(point.Point1.Lat), float64(point.Point1.Long), float64(point.Point2.Lat), float64(point.Point2.Long))
		checkVal(float64(point.Point1.Lat), LTOne)
		checkVal(float64(point.Point1.Long), LTOne)
		checkVal(float64(point.Point2.Lat), LTOne)
		checkVal(float64(point.Point2.Long), LTOne)

	}

	fmt.Println(float64(nLessOne) / (4 * float64(nPoints)))
	fmt.Println(sumDist / (float64(nPoints)))

}
func LTOne(val float64) bool {
	return val < 0
}
func checkVal(val float64, f func(val float64) bool) {
	if f(val) {
		nLessOne++
	}
}

func getRandomLatitude() Latitude {
	val := math.Acos(2*getRandom()-1) * 180 / math.Pi // [0 - 180]
	val = val - 90                                    // [-90 - 90]
	return Latitude(val)
}
func getRandomLongitude() Longitude {
	val := getRandom() * 360 // [0 - 360]
	val = val - 180          // [-180 - 180]
	return Longitude(val)
}
func nLatitutde() Latitude {
	return Latitude(getRandom()*360 - 180)
}
func nLongitude() Longitude {
	return Longitude(getRandom()*180 - 90)
}
func getRandom() float64 {
	return rand.Float64()
}
func DegToRad(deg float64) float64 {
	return deg * (math.Pi / 180)
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
