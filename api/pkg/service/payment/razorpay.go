package payment

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/razorpay/razorpay-go"
	"gorm.io/gorm"

	"github.com/shayamvlmna/lift/pkg/database"
	"github.com/shayamvlmna/lift/pkg/models"
)

type Order struct {
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
	}
}
type Payment struct {
	gorm.Model
	PaymentId string `json:"payment_id" gorm:"primaryKey"`
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

	key := os.Getenv("RAZORPAYPKEY")
	secret := os.Getenv("RAZORPAYSECRET")

	client := razorpay.NewClient(key, secret)

	data := map[string]interface{}{
		"amount":   amount,
		"currency": "INR",
		"receipt":  GeneratePaymentId(userid),
	}

	body, err := client.Order.Create(data, nil)
	if err != nil {
		fmt.Println(err)
	}

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

func GeneratePaymentId(id uint) string {
	db := database.Db

	payment := &Payment{}
	err := db.AutoMigrate(&payment)
	if err != nil {
		fmt.Println(err)
	}

	db.Last(&payment)

	paymentId := strings.Split(payment.PaymentId, "-")

	pid, err := strconv.Atoi(paymentId[0])
	if err != nil {
		fmt.Println(err)
	}

	// uid:=
	// +"-"+uid
	return strconv.Itoa(int(pid+1)) + "-" + strconv.Itoa(int(id))

}

func SavePayment(userid uint, res *OrderResponse) {

	pmt := &Payment{
		PaymentId: res.Receipt,
		UserId:    userid,
		Amount:    res.Amount,
		Status:    res.Status,
	}

	db := database.Db

	err := db.AutoMigrate(&Payment{})
	if err != nil {
		return
	}

	db.Create(&pmt)

}

func UpdatePayment(order *Order) {
	db := database.Db

	if err := db.AutoMigrate(&Payment{}); err != nil {
		return
	}

	db.Model(&Payment{}).Where("payment_id=?", order.Entity.Receipt).Update("status", "paid")

	err := db.AutoMigrate(&models.User{})
	if err != nil {
		return
	}

	receipt := order.Entity.Receipt

	userid := strings.Split(receipt, "-")

	db.Model(&models.User{}).Where("user_id=?", userid[1]).UpdateColumn("wallet_balance", gorm.Expr("wallet_balance + ?", order.Entity.Amount/100))

}
func PaymentFailed(order *Order) {
	db := database.Db

	if err := db.AutoMigrate(&Payment{}); err != nil {
		db.Model(&Payment{}).Where("payment_id=?", order.Entity.Receipt).Update("status", "failed")
	}
}
func ValidateWebhook(body []byte, signature string) bool {

	h := hmac.New(sha256.New, []byte("funnyhow"))

	h.Write(body)

	sha := hex.EncodeToString(h.Sum(nil))

	return signature == sha

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
				Id               string      `json:"id"`
				Entity           string      `json:"entity"`
				Amount           int         `json:"amount"`
				Currency         string      `json:"currency"`
				Status           string      `json:"status"`
				OrderId          string      `json:"order_id"`
				InvoiceId        interface{} `json:"invoice_id"`
				International    bool        `json:"international"`
				Method           string      `json:"method"`
				AmountRefunded   int         `json:"amount_refunded"`
				RefundStatus     interface{} `json:"refund_status"`
				Captured         bool        `json:"captured"`
				Description      interface{} `json:"description"`
				CardId           interface{} `json:"card_id"`
				Bank             interface{} `json:"bank"`
				Wallet           interface{} `json:"wallet"`
				Vpa              string      `json:"vpa"`
				Email            string      `json:"email"`
				Contact          string      `json:"contact"`
				Fee              int         `json:"fee"`
				Tax              int         `json:"tax"`
				ErrorCode        interface{} `json:"error_code"`
				ErrorDescription interface{} `json:"error_description"`
				CreatedAt        int         `json:"created_at"`
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
