package trip

import (
	database "github.com/shayamvlmna/cab-booking-app/pkg/database/postgresql"
	"github.com/shayamvlmna/cab-booking-app/pkg/models"
	"googlemaps.github.io/maps"
)

func Fare(d int) float32 {

	fare := float32(d) * 0.05
	return fare
}

type Ride struct {
	Source      maps.LatLng `json:"source"`
	Destination maps.LatLng `json:"destination"`
}

func CreateTrip(t *Ride) *models.Trip {

	source := &maps.LatLng{
		Lat: t.Source.Lat,
		Lng: t.Source.Lng,
	}

	destination := &maps.LatLng{
		Lat: t.Destination.Lat,
		Lng: t.Destination.Lng,
	}

	//TODO :
	//complete the distance matrix api part

	// distance, eta := maps.TripEstimate(source.String(), destination.String())

	// fare := Fare(distance)
	// AssignTrip(source, destination, distance, eta, fare)

	return AssignTrip(source, destination)

}

var Ridechanel = make(chan Ride, 2)

// func AssignTrip(source, destination *maps.LatLng, distance, eta int, fare float32) {

func AssignTrip(source, destination *maps.LatLng) *models.Trip {

	ride := &Ride{
		Source:      *source,
		Destination: *destination,
	}

	// origin := mapservice.GeoCode(&ride.Source)
	// dest := mapservice.GeoCode(&ride.Destination)

	newTrip := &models.Trip{
		Source:      "origin",
		Destination: "dest",
		Distance:    1,
		Fare:        100,
		ETA:         15,
	}

	Ridechanel <- *ride

	return newTrip
}

func GetRide() Ride {
	for {
		ride := <-Ridechanel
		return ride
	}
}

func GetTripHistory(id uint64) *[]models.Trip {
	return database.GetTrips(id)
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
