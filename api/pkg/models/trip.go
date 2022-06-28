package models

import (
	"time"

	"gorm.io/gorm"
)

// Source        Location
// Trip          Location `gorm:"ForeignKey:TripId"`

type Trip struct {
	gorm.Model
	TripId        uint64        `gorm:"autoIncrement;unique;primaryKey" json:"tripid"`
	Source        string        `json:"source" gorm:"ForeignKey:TripId;references:Id;embedded"`
	Destination   string        `json:"destination" gorm:"ForeignKey:TripId;references:Id;embedded"`
	Distance      uint          `gorm:"not null"`
	Fare          uint          `gorm:"not null"`
	ETA           time.Duration `json:"timeduration"`
	PaymentMethod Payment       `json:"paymentmethod" gorm:"ForeignKey:TripId;references:Pid;embedded"`
	Rating        uint          `json:"triprating"`
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
