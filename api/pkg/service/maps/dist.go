package maps

// "context"
// "log"

// "github.com/kr/pretty"

// "googlemaps.github.io/maps"

//
// var (
// 	slat string
// 	slon string
// 	dlat string
// 	dlon string
// )

// func Dist(origin, destination) {
// 	c, err := maps.NewClient(maps.WithAPIKey("Insert-API-Key-Here"))
// 	if err != nil {
// 		log.Fatalf("fatal error: %s", err)
// 	}

// 	r := &maps.DistanceMatrixRequest{
// 		Origins:      []string{slat, slon},
// 		Destinations: []string{dlat, dlon},
// 	}

// 	distancematrix, err := c.DistanceMatrix(context.Background(), r)
// 	if err != nil {
// 		log.Fatalf("fatal error: %s", err)
// 	}

// 	pretty.Println(distancematrix)

// }
