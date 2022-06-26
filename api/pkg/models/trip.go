package models

import (
	"context"
	"fmt"
	"log"
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
}
type Location struct {
	Id  uint `gorm:"primaryKey"`
	Lon float64
	Lat float64
}

type Payment struct {
	Pid    uint `gorm:"primaryKey"`
	Wallet bool
	Cash   bool
}

func (t *Trip) CreateTrip() {
	t.ProcessTrip()
}
func (t *Trip) ProcessTrip() {

	source := &maps.LatLng{
		Lat: t.Source.Lat,
		Lng: t.Source.Lon,
	}
	destination := &maps.LatLng{
		Lat: t.Source.Lat,
		Lng: t.Source.Lon,
	}

	distance, eta := Dist(source, destination)

	fare := Fare(distance)

	AssignTrip(distance, eta, fare).TripPool()
}

var (
	slat string
	slon string
	dlat string
	dlon string
)

func Dist(origin, destination *maps.LatLng) (int, int) {
	c, err := maps.NewClient(maps.WithAPIKey("Insert-API-Key-Here"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	r := &maps.DistanceMatrixRequest{
		Origins:      []string{slat, slon},
		Destinations: []string{dlat, dlon},
	}

	distancematrix, err := c.DistanceMatrix(context.Background(), r)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	pretty.Println(distancematrix)
	distance := distancematrix.Rows[3].Elements[0].Distance.Meters
	duration := distancematrix.Rows[3].Elements[0].Duration

	// d := distancematrix.Rows
	// ds := d[3]
	// dsl := ds.Elements
	// da := dsl[0]

	// distance := da.Distance.Meters
	// duration := da.Duration

	fmt.Println(distance)
	fmt.Println(duration)

	return distance, int(duration)
}
func Fare(d int) float32 {

	fare := float32(d) * 0.05
	return fare
}

func AssignTrip(distance, eta int, fare float32) *Trip {
	newTrip := &Trip{
		Distance: uint(distance),
		Fare:     uint(fare),
		ETA:      time.Duration(eta),
	}
	return newTrip
}
func (t *Trip) TripPool() {

}
