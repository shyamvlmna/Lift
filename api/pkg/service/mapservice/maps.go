package mapservice

import (
	"context"
	"encoding/json"
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

	r := &maps.GeocodingRequest{
		LatLng: &maps.LatLng{
			Lat: g.Lat,
			Lng: g.Lng,
		},
	}

	reverseGeocode, err := c.ReverseGeocode(context.Background(), r)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	return reverseGeocode[0].FormattedAddress
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

func DistanceAPI() *Result {
	// https://api.distancematrix.ai/maps/api/distancematrix/json?origins=51.4822656,-0.1933769&destinations=51.4994794,-0.1269979
	// &key=<your_access_token>
	url := "https://api.distancematrix.ai/maps/api/distancematrix/json?origins=11.258753,75.780411&destinations=11.874477,75.370369&key=JNDApQ6vaPwL3zBFbMNegII9BnNEj"
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

// {
//     "destination_addresses": [
//         "Westminster Abbey, Westminster, London SW1P 3PA, UK"
//     ],
//     "origin_addresses": [
//         "Chapel, Fulham, London SW6 1BA, UK"
//     ],
//     "rows": [
//         {
//             "elements": [
//                 {
//                     "distance": {
//                         "text": "7.6 km",
//                         "value": 7561
//                     },
//                     "duration": {
//                         "text": "22 min",
//                         "value": 1303
//                     },
//                     "status": "OK"
//                 }
//             ]
//         }
//     ],
//     "status": "OK"
// }
