package coupon

import (
	"github.com/shayamvlmna/cab-booking-app/pkg/database"
	"time"
)

type Coupon interface {
	IsApplicable(cost int) bool
	ApplyDiscount(totalamount float64) float64
}

type AmountCoupon struct {
	MinFare    float64   `json:"min_fare"`
	Amount     float64   `json:"amount"`
	CouponCode string    `json:"coupon_code"`
	FinishDate time.Time `json:"finish_date"`
}

func (ac AmountCoupon) CreateCoupon() error {
	db := database.Db

	err := db.AutoMigrate(&AmountCoupon{})
	if err != nil {
		return err
	}
	return db.Create(ac).Error
}

func (ac AmountCoupon) IsApplicable(cost int) bool {

	return true
}
