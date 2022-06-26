package models

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/kr/pretty"
	"gorm.io/gorm"

	"googlemaps.github.io/maps"
)

// Source        Location
// Trip          Location `gorm:"ForeignKey:TripId"`

type Trip struct {
	gorm.Model
	TripId        uint64        `gorm:"autoIncrement;unique;primaryKey" json:"tripid"`
	Source        Location      `json:"source" gorm:"ForeignKey:TripId;references:Id;embedded"`
	Destination   Location      `json:"destination" gorm:"ForeignKey:TripId;references:Id;embedded"`
	Distance      uint          `gorm:"not null"`
	Fare          uint          `gorm:"not null"`
	ETA           time.Duration `json:"timeduration"`
	PaymentMethod Payment       `json:"paymentmethod" gorm:"ForeignKey:TripId;references:Pid;embedded"`
	Rating        uint          `json:"triprating"`
	UserId        uint64
}

type Ride struct {
	Source      Location `json:"source"`
	Destination Location `json:"destination"`
}

type Location struct {
	Id  uint    `gorm:"primaryKey"`
	Lat float64 `json:"latitude"`
	Lng float64 `json:"longitude"`
}

type Payment struct {
	Pid    uint `gorm:"primaryKey"`
	Wallet bool
	Cash   bool
}

func (t *Ride) CreateTrip() {
	t.ProcessTrip()

}

func (t *Ride) ProcessTrip() {

	source := &maps.LatLng{
		Lat: t.Source.Lat,
		Lng: t.Source.Lng,
	}

	destination := &maps.LatLng{
		Lat: t.Destination.Lat,
		Lng: t.Destination.Lng,
	}

	distance, eta := Dist(source, destination)

	fare := Fare(distance)

	AssignTrip(source, destination, distance, eta, fare).TripPool()
}

func Dist(origin, destination *maps.LatLng) (int, int) {
	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyCouPhivkPPHguv4I0j_3BYMUrV6EIcBBo"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	r := &maps.DistanceMatrixRequest{
		Origins:      []string{strconv.Itoa(int(origin.Lat)), strconv.Itoa(int(origin.Lng))},
		Destinations: []string{strconv.Itoa(int(destination.Lng)), strconv.Itoa(int(destination.Lng))},
	}

	distancematrix, err := c.DistanceMatrix(context.Background(), r)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	pretty.Println(distancematrix)
	distance := distancematrix.Rows[3].Elements[0].Distance.Meters
	duration := distancematrix.Rows[3].Elements[0].Duration.Minutes()

	fmt.Println(distance)
	fmt.Println(duration)

	return distance, int(duration)
}

func Fare(d int) float32 {

	fare := float32(d) * 0.05
	return fare
}

func AssignTrip(source, destination *maps.LatLng, distance, eta int, fare float32) *Trip {
	newTrip := &Trip{
		Source:        source,
		Destination:   Location{},
		Distance:      0,
		Fare:          0,
		ETA:           0,
		PaymentMethod: Payment{},
		Rating:        0,
		UserId:        0,
	}
	return newTrip
}

func (t *Trip) TripPool() {
	// Distance: uint(distance),
	// Fare:     uint(fare),
	// ETA:      time.Duration(eta),
}
