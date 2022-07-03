package payment

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/razorpay/razorpay-go"
	"os"
)

type OrderResponse struct {
	Id         string `json:"id"`
	Amount     int    `json:"amount"`
	AmountPaid int    `json:"amount_paid"`
	AmountDue  int    `json:"amount_due"`
	Currency   string `json:"currency"`
	Receipt    string `json:"receipt"`
	Status     string `json:"status"`
}

func AddMoney() *OrderResponse {

	err := godotenv.Load()
	if err != nil {
		return nil
	}
	key := os.Getenv("RPKEY")
	secret := os.Getenv("RPSCRT")
	client := razorpay.NewClient(key, secret)

	data := map[string]interface{}{
		"amount":   2000,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	//Get(data, nil)
	body, err := client.Order.Fetch("order_Jp3O4XVHL6SL4V", data, nil)

	marshal, err := json.Marshal(body)
	if err != nil {
		return nil
	}

	resp := &OrderResponse{}

	err = json.Unmarshal(marshal, &resp)
	if err != nil {
		return nil
	}

	return resp

	//paymentId := "pay_JIskP2THPhAUXg"
	//amount := 4000
	//
	//data = map[string]interface{}{
	//	"currency": "INR",
	//}
	//
	//body, err = client.Payment.Capture(paymentId, amount, data, nil)
}

func Payment() {

}

//{
//"key": "rzp_test_1vwFv5eIxCHZY2",
//"amount": "2000",
//"currency": "INR",
//"name": "Cab Booking App",
//"description": "Test Transaction",
//"order_id": "order_IluGWxBm9U8zJ8"
//}
