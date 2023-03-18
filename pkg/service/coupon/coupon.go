package coupon

import (
	"github.com/shayamvlmna/lift/internal/database"
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

func (ac AmountCoupon) IsApplicable(cost float64) bool {
	return cost >= ac.MinFare && ac.FinishDate.After(time.Now())
}
func GetCoupon(code string) *AmountCoupon {
	db := database.Db

	db.AutoMigrate(&AmountCoupon{})

	coupon := &AmountCoupon{}

	db.Where("coupon_code=?", code).First(&coupon)
	return coupon
}

func GetCoupons() *[]AmountCoupon {
	db := database.Db

	db.AutoMigrate(&AmountCoupon{})

	coupons := &[]AmountCoupon{}

	db.Find(&coupons)

	return coupons

}
