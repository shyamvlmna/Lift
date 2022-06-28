package mapservice

import (
	"context"
	"fmt"
	"log"

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
