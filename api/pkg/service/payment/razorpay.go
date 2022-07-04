package payment

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/razorpay/razorpay-go"
	"github.com/shayamvlmna/cab-booking-app/pkg/database"
	"github.com/shayamvlmna/cab-booking-app/pkg/models"
	"gorm.io/gorm"
	"os"
	"strconv"
)

type Payment struct {
	gorm.Model
	PaymentId uint   `json:"payment_id" gorm:"primaryKey"`
	UserId    uint   `json:"user_id"`
	Amount    uint   `json:"amount"`
	Status    string `json:"status"`
}

type OrderResponse struct {
	Id         string `json:"id"`
	Amount     uint   `json:"amount"`
	AmountPaid int    `json:"amount_paid"`
	AmountDue  int    `json:"amount_due"`
	Currency   string `json:"currency"`
	Receipt    string `json:"receipt"`
	Status     string `json:"status"`
}

func AddMoney(userid uint, amount uint) *OrderResponse {
	if err := godotenv.Load(); err != nil {
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

	SavePayment(userid, resp)

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
	db := database.Db

	payment := &Payment{}
	err := db.AutoMigrate(&payment)
	if err != nil {
		//return 0
	}

	db.Last(&payment)

	return strconv.Itoa(int(payment.PaymentId + 1))

	//return uuid.NewString()
}

func SavePayment(userid uint, res *OrderResponse) {

	pid, err := strconv.Atoi(res.Receipt)
	if err != nil {
		fmt.Println(err)
	}

	pmt := &Payment{
		PaymentId: uint(pid),
		UserId:    userid,
		Amount:    res.Amount,
		Status:    res.Status,
	}

	db := database.Db

	err = db.AutoMigrate(&Payment{})
	if err != nil {
		return
	}

	db.Create(&pmt)

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return
	}

	db.Model(&models.User{}).Where("user_id=?", userid).UpdateColumn("wallet_balance", gorm.Expr("wallet_balance + ?", res.Amount/100))

}

func UpdatePayment(paymentId string) {
	db := database.Db

	if err := db.AutoMigrate(&Payment{}); err != nil {
		return
	}

	db.Model(&Payment{}).Where("payment_id=?", paymentId).Update("status", "paid")

}

func ValidateWebhook(body []byte, signature string) bool {

	h := hmac.New(sha256.New, []byte("funnyhow"))

	h.Write(body)

	sha := hex.EncodeToString(h.Sum(nil))

	if signature != sha {
		return false
	}
	return true
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
				Id               string        `json:"id"`
				Entity           string        `json:"entity"`
				Amount           int           `json:"amount"`
				Currency         string        `json:"currency"`
				Status           string        `json:"status"`
				OrderId          string        `json:"order_id"`
				InvoiceId        interface{}   `json:"invoice_id"`
				International    bool          `json:"international"`
				Method           string        `json:"method"`
				AmountRefunded   int           `json:"amount_refunded"`
				RefundStatus     interface{}   `json:"refund_status"`
				Captured         bool          `json:"captured"`
				Description      interface{}   `json:"description"`
				CardId           interface{}   `json:"card_id"`
				Bank             interface{}   `json:"bank"`
				Wallet           interface{}   `json:"wallet"`
				Vpa              string        `json:"vpa"`
				Email            string        `json:"email"`
				Contact          string        `json:"contact"`
				Notes            []interface{} `json:"notes"`
				Fee              int           `json:"fee"`
				Tax              int           `json:"tax"`
				ErrorCode        interface{}   `json:"error_code"`
				ErrorDescription interface{}   `json:"error_description"`
				CreatedAt        int           `json:"created_at"`
			} `json:"entity"`
		} `json:"payment"`
		Order struct {
			Entity struct {
				Id         string        `json:"id"`
				Entity     string        `json:"entity"`
				Amount     int           `json:"amount"`
				AmountPaid int           `json:"amount_paid"`
				AmountDue  int           `json:"amount_due"`
				Currency   string        `json:"currency"`
				Receipt    string        `json:"receipt"`
				OfferId    interface{}   `json:"offer_id"`
				Status     string        `json:"status"`
				Attempts   int           `json:"attempts"`
				Notes      []interface{} `json:"notes"`
				CreatedAt  int           `json:"created_at"`
			} `json:"entity"`
		} `json:"order"`
	} `json:"payload"`
	CreatedAt int `json:"created_at"`
}
