package models

import (
	"gorm.io/gorm"
)

// Source        Location
// Trip          Location `gorm:"ForeignKey:TripId"`

type Trip struct {
	gorm.Model
	TripId        uint64  `gorm:"autoIncrement;unique;primaryKey" json:"tripid"`
	Source        string  `json:"source" gorm:"ForeignKey:TripId;references:Id;embedded"`
	Destination   string  `json:"destination" gorm:"ForeignKey:TripId;references:Id;embedded"`
	Distance      uint32  `gorm:"not null"`
	Fare          uint32  `gorm:"not null"`
	ETA           string  `json:"timeduration"`
	PaymentMethod Payment `json:"paymentmethod" gorm:"ForeignKey:TripId;references:Pid;embedded"`
	Rating        uint8   `json:"triprating"`
	UserId        uint64
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

type Ride struct {
	gorm.Model
	RideId        uint64 `gorm:"autoIncrement;unique;primaryKey" json:"rideid"`
	Source        string `json:"source"`
	Destination   string `json:"destination"`
	ETA           string `json:"eta"`
	Fare          string `json:"fare"`
	PaymentMethod string `json:"paymentmethod"`
	UserId        uint64
}

// func (t *Trip) TripPool() {
// 	// Distance: uint(distance),
// 	// Fare:     uint(fare),
// 	// ETA:      time.Duration(eta),
// }

// {
//     "destination_addresses": [
//         "St John's Church, North End Rd, Fulham, London SW6 1PB, United Kingdom"
//     ],
//     "origin_addresses": [
//         "Westminster Abbey, 20 Deans Yd, Westminster, London SW1P 3PA, United Kingdom"
//     ],
//     "rows": [
//         {
//             "elements": [
//                 {
//                     "distance": {
//                         "text": "6.5 km",
//                         "value": 6477
//                     },
//                     "duration": {
//                         "text": "21 min",
//                         "value": 1287
//                     },
//                     "status": "OK"
//                 }
//             ]
//         }
//     ],
//     "status": "OK"
// }
