package mapservice

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kr/pretty"
	"googlemaps.github.io/maps"
)

//return the distance and estimate time of arrival from origin to destination
func TripEstimate(origin, destination string) (int, int) {
	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyCouPhivkPPHguv4I0j_3BYMUrV6EIcBBo"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	r := &maps.DistanceMatrixRequest{
		Origins:      []string{origin},
		Destinations: []string{destination},
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

func GeoCode(g *maps.LatLng) string {
	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyCouPhivkPPHguv4I0j_3BYMUrV6EIcBBo"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	r := maps.GeocodingRequest{
		LatLng: &maps.LatLng{
			Lat: g.Lat,
			Lng: g.Lng,
		},
	}

	reverseGeocode, err := c.ReverseGeocode(context.Background(), &r)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	return reverseGeocode[0].FormattedAddress
}

func Distance() {

	url := "https://maps.googleapis.com/maps/api/distancematrix/json?origins=Washington,%20DC&destinations=New%20York%20City,%20NY&units=imperial&key=YOUR_API_KEY"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
