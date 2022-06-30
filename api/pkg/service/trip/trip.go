package trip

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	database "github.com/shayamvlmna/cab-booking-app/pkg/database/postgresql"
	models "github.com/shayamvlmna/cab-booking-app/pkg/models"
	maps "googlemaps.github.io/maps"
)

func Fare(d int) float32 {

	fare := float32(d) * 0.05
	return fare
}

type Ride struct {
	Source      LatLng `json:"source"`
	Destination LatLng `json:"destination"`
}

type LatLng struct {
	Lat float64
	Lng float64
}

func CreateTrip(t *Ride) *models.Ride {

	source := t.Source
	destination := t.Destination

	newride := &Ride{
		Source:      source,
		Destination: destination,
	}

	// TODO: geocode the source and destination

	result := DistanceAPI(newride)

	distance := result.Rows[0].Element[0].Distance.Val
	Kmdistance := result.Rows[0].Element[0].Distance.Text

	eta := result.Rows[0].Element[0].Duration.Text
	fare := Fare(distance)
	newTrip := &models.Ride{
		Source:      "geocoded source",
		Destination: "geocoded destination",
		Distance:    Kmdistance,
		Fare:        uint(fare),
		ETA:         eta,
	}

	fmt.Println(newTrip.Distance)
	// source := &maps.LatLng{
	// 	Lat: t.Source.Lat,
	// 	Lng: t.Source.Lng,
	// }

	// destination := &maps.LatLng{
	// 	Lat: t.Destination.Lat,
	// 	Lng: t.Destination.Lng,
	// }

	//TODO :
	//complete the distance matrix api part

	// distance, eta := maps.TripEstimate(source.String(), destination.String())

	// fare := Fare(distance)
	// AssignTrip(source, destination, distance, eta, fare)

	// return AssignTrip(source, destination)
	return newTrip
}

func FindCab(ride **models.Ride) {
	Ridechanel <- **ride
}

type Result struct {
	Destination []string `json:"destination_addresses"`
	Origin      []string `json:"origin_addresses"`
	Rows        []Elem   `json:"rows"`
	Status      string   `json:"status"`
}

type Elem struct {
	Element []Elements `json:"elements"`
}

type Elements struct {
	Distance Dist   `json:"distance"`
	Duration Dist   `json:"duration"`
	Status   string `json:"status"`
}

type Dist struct {
	Text string `json:"text"`
	Val  int    `json:"value"`
}

func DistanceAPI(r *Ride) *Result {
	// https://api.distancematrix.ai/maps/api/distancematrix/json?origins=51.4822656,-0.1933769&destinations=51.4994794,-0.1269979
	// &key=<your_access_token>

	origins := fmt.Sprintf("%s,%s", strconv.FormatFloat(r.Source.Lat, 'f', -1, 64), strconv.FormatFloat(r.Source.Lng, 'f', -1, 64))
	destinations := fmt.Sprintf("%s,%s", strconv.FormatFloat(r.Destination.Lat, 'f', -1, 64), strconv.FormatFloat(r.Destination.Lng, 'f', -1, 64))
	url := fmt.Sprintf("https://api.distancematrix.ai/maps/api/distancematrix/json?origins=%s&destinations=%s&key=JNDApQ6vaPwL3zBFbMNegII9BnNEj", origins, destinations)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)

	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)

	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)

	}

	fmt.Println(string(body))
	result := &Result{}
	json.Unmarshal([]byte(body), &result)
	return result
}

func GeoCodeApi(l LatLng) *Result {

	lat := strconv.FormatFloat(l.Lat, 'f', -1, 64)
	lng := strconv.FormatFloat(l.Lng, 'f', -1, 64)
	url := fmt.Sprintf("https://api.distancematrix.ai/maps/api/geocode/json?latlng=%s,%s&key=JNDApQ6vaPwL3zBFbMNegII9BnNEj"+lat, lng)
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)

	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)

	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)

	}

	fmt.Println(string(body))
	result := &Result{}
	json.Unmarshal([]byte(body), &result)
	return result
}

var Ridechanel = make(chan models.Ride)

// func AssignTrip(source, destination *maps.LatLng, distance, eta int, fare float32) {
type Trip struct {
	Source      string        `json:"source"`
	Destination string        `json:"destination"`
	Distance    int           `json:"distance"`
	Fare        int           `json:"fare"`
	ETA         time.Duration `json:"eta"`
}

func AssignTrip(source, destination *maps.LatLng) *models.Trip {

	// ride := &Ride{
	// 	Source:      *source,
	// 	Destination: *destination,
	// }

	// origin := mapservice.GeoCode(&ride.Source)
	// dest := mapservice.GeoCode(&ride.Destination)

	newTrip := &models.Trip{
		Source:      "origin",
		Destination: "dest",
		// Distance:    1,
		Fare: 100,
		ETA:  "",
	}

	// Ridechanel <- *newTrip

	return newTrip
}

func GetRide() models.Ride {
	for {
		ride := <-Ridechanel
		return ride
	}
}

func GetTripHistory(id uint64) *[]models.Trip {
	return database.GetTrips(id)
}

func RegisterTrip(ride *models.Ride) error{
	trip := &models.Trip{}

	trip.Source = ride.Source
	trip.Destination = ride.Destination
	trip.Distance = ride.Distance
	trip.Fare = ride.Fare
	trip.ETA = ride.ETA
	trip.PaymentMethod = ride.PaymentMethod
	trip.DriverId = ride.DriverId
	trip.UserId = ride.UserId

	return trip.Add(trip)
}

// type Pool struct {
// 	BookTrip chan *models.User
// 	GetTrip  chan *models.Driver
// 	Trips    chan models.Trip
// }

// func NewPool() *Pool {
// 	return &Pool{
// 		BookTrip: make(chan *models.User),
// 		GetTrip:  make(chan *models.Driver),
// 		Trips:    make(chan models.Trip),
// 	}
// }

// func (pool *Pool) Start() {
// 	for {
// 		select {
// 		// case trip:=<-pool.BookTrip:

// 		}
// 	}
// }

// func CreateNewTrip() {
// 	// pool := &NewPool()
// }

// func ProcessTrip(models.Trip)
