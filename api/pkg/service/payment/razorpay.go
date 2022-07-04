package payment

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/razorpay/razorpay-go"
	"os"
)

type Payment struct {
	Amount uint `json:"amount"`
}

type OrderResponse struct {
	Id         string `json:"id"`
	Amount     int    `json:"amount"`
	AmountPaid int    `json:"amount_paid"`
	AmountDue  int    `json:"amount_due"`
	Currency   string `json:"currency"`
	Receipt    string `json:"receipt"`
	Status     string `json:"status"`
}

func AddMoney(amount uint) *OrderResponse {

	err := godotenv.Load()
	if err != nil {
		return nil
	}
	key := os.Getenv("RPKEY")
	secret := os.Getenv("RPSCRT")
	client := razorpay.NewClient(key, secret)

	data := map[string]interface{}{
		"amount":   amount,
		"currency": "INR",
		"receipt":  GeneratePaymentId(),
	}
	//Get(data, nil)
	body, err := client.Order.Create(data, nil)

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

func GeneratePaymentId() string {
	return uuid.NewString()
}

//{
//"key": "rzp_test_1vwFv5eIxCHZY2",
//"amount": "2000",
//"currency": "INR",
//"name": "Cab Booking App",
//"description": "Test Transaction",
//"order_id": "order_IluGWxBm9U8zJ8"
//}
type Webhook struct {
	Entity    string   `json:"entity"`
	AccountId string   `json:"account_id"`
	Event     string   `json:"event"`
	Contains  []string `json:"contains"`
	Payload   struct {
		Payment struct {
			Entity struct {
				Id                string        `json:"id"`
				Entity            string        `json:"entity"`
				Amount            int           `json:"amount"`
				Currency          string        `json:"currency"`
				BaseAmount        int           `json:"base_amount"`
				Status            string        `json:"status"`
				OrderId           string        `json:"order_id"`
				InvoiceId         interface{}   `json:"invoice_id"`
				International     bool          `json:"international"`
				Method            string        `json:"method"`
				AmountRefunded    int           `json:"amount_refunded"`
				AmountTransferred int           `json:"amount_transferred"`
				RefundStatus      interface{}   `json:"refund_status"`
				Captured          bool          `json:"captured"`
				Description       interface{}   `json:"description"`
				CardId            interface{}   `json:"card_id"`
				Bank              interface{}   `json:"bank"`
				Wallet            interface{}   `json:"wallet"`
				Vpa               string        `json:"vpa"`
				Email             string        `json:"email"`
				Contact           string        `json:"contact"`
				Notes             []interface{} `json:"notes"`
				Fee               int           `json:"fee"`
				Tax               int           `json:"tax"`
				ErrorCode         interface{}   `json:"error_code"`
				ErrorDescription  interface{}   `json:"error_description"`
				ErrorSource       interface{}   `json:"error_source"`
				ErrorStep         interface{}   `json:"error_step"`
				ErrorReason       interface{}   `json:"error_reason"`
				AcquirerData      struct {
					Rrn string `json:"rrn"`
				} `json:"acquirer_data"`
				CreatedAt int `json:"created_at"`
			} `json:"entity"`
		} `json:"payment"`
	} `json:"payload"`
	CreatedAt int `json:"created_at"`
}
