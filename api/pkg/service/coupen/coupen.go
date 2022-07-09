package coupen

import "time"

type Coupen interface {
	IsApplicable(cost int) bool
	ApplyDiscount(totalamount float64) float64
}

type AmountCoupon struct {
	MinPurchaseAmount float64
	Amount            float64
	CouponCode        string
	FinishDate        time.Time
}
