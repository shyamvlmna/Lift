package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/shayamvlmna/lift/pkg/database/redis"
	"github.com/shayamvlmna/lift/pkg/models"
	"github.com/shayamvlmna/lift/pkg/service/auth"
	"github.com/shayamvlmna/lift/pkg/service/coupon"
	"github.com/shayamvlmna/lift/pkg/service/payment"
	"github.com/shayamvlmna/lift/pkg/service/trip"
	"github.com/shayamvlmna/lift/pkg/service/user"
)

// UserAuth Check if the user already exist in the system.
//Redirect to the user login page if user exists.
//Redirect to the user signup page if user is new.
func UserAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	newUser := &models.User{}

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		return
	}

	phonenumber := newUser.Phonenumber

	if phonenumber != "" {
		auth.StorePhone(phonenumber)
		if user.IsUserExists("phonenumber", phonenumber) {
			http.Redirect(w, r, "/user/loginpage", http.StatusSeeOther)
			return
		} else {
			if err := auth.SetOtp(phonenumber); err != nil {
				fmt.Println(err)
				return
			}
			http.Redirect(w, r, "/user/enterotp", http.StatusSeeOther)
			return
		}
	} else {
		response := models.Response{
			ResponseStatus:  "fail",
			ResponseMessage: "phonenumber required",
			ResponseData:    nil,
		}
		err := json.NewEncoder(w).Encode(&response)
		if err != nil {
			return
		}
		return
	}
}

func EnterOTPUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "new user",
		ResponseData:    "verify with otp",
	}
	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}

}

// UserSignupPage render the signup page to submit the details of the new user
func UserSignupPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "submit user data",
		ResponseData:    nil,
	}
	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
}

func UserLoginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	phonenumber := auth.GetPhone()

	response := models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "existing user",
		ResponseData:    user.GetUser("phonenumber", phonenumber).Firstname,
	}

	auth.StorePhone(phonenumber)

	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
}

// UserSignUp Create a user model with values from the fronted.
//Pass the newly created user model to user services
//to insert the new user to the database.
//Login the user and open user home after successful signup.
func UserSignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	newUser := models.User{}

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(r.Body)

	if err := user.RegisterUser(&newUser); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		response := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "signup failed",
			ResponseData:    nil,
		}
		err := json.NewEncoder(w).Encode(&response)
		if err != nil {
			return
		}
		return
	}

	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "signup success",
		ResponseData:    nil,
	}
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
}

// UserLogin get the existing user by phone number from the database.
//Validate the entered password with stored hash password.
//Generate a JWT token for the user after successful login.
//Store the JWT token in the http only cookie
func UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	newUser := &models.User{}

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(r.Body)

	password := newUser.Password

	// phoneNumber := newUser.PhoneNumber

	phoneNumber := auth.GetPhone()

	//get the existing user by phone number from the database
	User := user.GetUser("phonenumber", phoneNumber)

	if !User.Active {
		response := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "user not active",
			ResponseData:    nil,
		}
		err := json.NewEncoder(w).Encode(&response)
		if err != nil {
			return
		}
		return
	}

	//validate the entered password with stored hash password
	if !ValidPassword(password, User.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		response := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "password authentication failed",
			ResponseData:    nil,
		}
		err := json.NewEncoder(w).Encode(&response)
		if err != nil {
			return
		}
		return
	}

	token, err := auth.GenerateJWT("user", phoneNumber)
	if err != nil {
		fmt.Println("jwt failed", err)
	}

	err = redis.StoreData("data", User)
	if err != nil {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt-token",
		Value:    token,
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
	})

	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "login success",
		ResponseData:    token,
	}

	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}

	// http.Redirect(w, r, "/user/userhome", http.StatusSeeOther)
}

// UserHome get the logged-in user data from redis
//if err get from the primary database
//render the user home page
func UserHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	user, err := GetUserFromCookie(r)
	if err != nil {
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
		return
	}

	userData := &models.UserData{
		Id:          user.UserId,
		Phonenumber: user.Phonenumber,
		Firstname:   user.Firstname,
		Lastname:    user.Lastname,
		Email:       user.Email,
	}

	response := models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "user data fetched",
		ResponseData:    userData,
	}

	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
}

// UserLogout delete the stored user data from redis
//also expire the cookie stored
func UserLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Content-Type", "application/json")

	c, err := r.Cookie("jwt-token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	err = redis.DeleteData("data")
	if err != nil {
		return
	}

	c.Value = ""
	c.Path = "/"
	c.MaxAge = -1

	http.SetCookie(w, c)

	response := &models.Response{
		ResponseStatus:  "succes",
		ResponseMessage: "logout success",
		ResponseData:    nil,
	}

	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
	// http.Redirect(w, r, "/", http.StatusSeeOther)
}

func EditUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := GetUserFromCookie(r)
	if err != nil {
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
		return
	}

	userdata := &models.UserData{
		Id:          user.UserId,
		Phonenumber: user.Phonenumber,
		Firstname:   user.Firstname,
		Lastname:    user.Lastname,
		Email:       user.Email,
	}

	err = json.NewEncoder(w).Encode(&userdata)
	if err != nil {
		return
	}
}

func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")

	u, err := GetUserFromCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
	}

	newUser := models.User{}

	err = json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = user.UpdateUser(u.UserId, &newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "error in updating user profile",
			ResponseData:    err,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "user profile updated",
		ResponseData:    newUser,
	}
	json.NewEncoder(w).Encode(resp)
}

// BookTrip get the pickup point and destination from the booktrip call from the user
func BookTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := GetUserFromCookie(r)
	if err != nil {
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
		return
	}

	newRide := &trip.Ride{}

	err = json.NewDecoder(r.Body).Decode(&newRide)
	if err != nil {
		return
	}

	newTrip := trip.CreateTrip(newRide)

	ride := &models.Ride{
		Source:        newTrip.Source,
		Destination:   newTrip.Destination,
		Distance:      newTrip.Distance,
		ETA:           newTrip.ETA,
		Fare:          newTrip.Fare,
		PaymentMethod: "",
	}

	if user.WalletBalance < newTrip.Fare {
		response := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "not enough balance in wallet",
			ResponseData:    &ride,
		}
		if err = json.NewEncoder(w).Encode(&response); err != nil {
			return
		}
		return
	}

	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "trip created successfully",
		ResponseData:    &ride,
	}
	if err = json.NewEncoder(w).Encode(&response); err != nil {
		return
	}
}

//var OTPchan = make(chan int)

//ConfirmTrip returns the trip code to match with the driver to start the ride
func ConfirmTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	confirmedTrip := &models.Ride{}

	err := json.NewDecoder(r.Body).Decode(&confirmedTrip)
	if err != nil {
		return
	}

	curUser, err := GetUserFromCookie(r)
	if err != nil {
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
		return
	}

	if confirmedTrip.Coupon != "" {
		discountFare, err := user.ApplyCoupon(confirmedTrip.Coupon, confirmedTrip.Fare)
		if err == nil {
			confirmedTrip.Fare = discountFare
			confirmedTrip.UserId = curUser.UserId

			go trip.FindCab(&confirmedTrip)

			json.NewEncoder(w).Encode(&models.Response{
				ResponseStatus:  "success",
				ResponseMessage: "coupon applied. waiting to accept ride",
				ResponseData:    confirmedTrip,
			})
			return
		}
	}
	confirmedTrip.UserId = curUser.UserId

	go trip.FindCab(&confirmedTrip)

	json.NewEncoder(w).Encode(&models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "waiting to accept ride",
		ResponseData:    nil,
	})
}

// UserTripHistory returns saved trips of the user
func UserTripHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := GetUserFromCookie(r)
	if err != nil {
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
		return
	}

	tripHistory := trip.GetTripHistory("user_id", user.UserId)

	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "fetched trip history",
		ResponseData:    tripHistory,
	}
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
}

func UserWallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := GetUserFromCookie(r)
	if err != nil {
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
		return
	}

	response := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "user wallet",
		ResponseData:    user.WalletBalance,
	}
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		return
	}
}

func AddMoneyToWallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	curUser, err := GetUserFromCookie(r)
	if err != nil {
		json.NewEncoder(w).Encode(&ErrorParsingCookie)
		return
	}

	pmt := &payment.Payment{}

	err = json.NewDecoder(r.Body).Decode(&pmt)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(r.Body)

	resp := payment.AddMoney(curUser.UserId, pmt.Amount)

	//tmp, _ := template.ParseFiles("/home/shyamjith/cab-booking-app/api/pkg/service/payment/app.html")
	//tmp.Execute(w, nil)

	if err = json.NewEncoder(w).Encode(&resp); err != nil {
		return
	}
}

func RazorpayCallback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("successss")

	err := json.NewEncoder(w).Encode(&models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "amount added to wallet",
		ResponseData:    nil,
	})
	if err != nil {
		return
	}
}

func RazorpayWebhook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	signature := r.Header.Get("X-Razorpay-Signature")

	msg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	r.Body.Close()

	if ok := payment.ValidateWebhook(msg, signature); !ok {
		fmt.Println("unauthorized webhook call")
		return
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(msg))

	fmt.Println("webhook called")

	wh := &payment.Webhook{}
	if err = json.NewDecoder(r.Body).Decode(&wh); err != nil {
		fmt.Println(err)
	}
	r.Body.Close()

	fmt.Println(wh.Event)

	if wh.Event == "order.paid" {
		payment.UpdatePayment((*payment.Order)(&wh.Payload.Order))
	} else if wh.Event == "payment.failed" {
		payment.PaymentFailed((*payment.Order)(&wh.Payload.Order))
	}
	//
	fmt.Println(wh.Payload.Payment.Entity.Amount)
	fmt.Println(wh.Payload.Order.Entity.Receipt)
	if err = json.NewEncoder(w).Encode(wh); err != nil {
		return
	}
}

func GetCoupons(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	coupons := coupon.GetCoupons()

	json.NewEncoder(w).Encode(&coupons)
}

type Amtcoupon struct {
	Code string `json:"code"`
	Fare string `json:"fare"`
}

func ApplyCoupon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	c := &Amtcoupon{}
	json.NewDecoder(r.Body).Decode(&c)
	r.Body.Close()

	couponCode := c.Code

	fare, _ := strconv.ParseFloat(c.Fare, 64)
	ride := &models.Ride{}

	dfare, err := user.ApplyCoupon(couponCode, fare)
	if err != nil {
		resp := &models.Response{
			ResponseStatus:  "failed",
			ResponseMessage: "",
			ResponseData:    err,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	ride.Fare = dfare

	resp := &models.Response{
		ResponseStatus:  "success",
		ResponseMessage: "coupon applied",
		ResponseData:    ride.Fare,
	}

	json.NewEncoder(w).Encode(resp)

}

// GetUserFromCookie returns the logged-in user from the stored cookie in session
func GetUserFromCookie(r *http.Request) (*models.User, error) {
	c, err := r.Cookie("jwt-token")

	if err != nil {
		if err == http.ErrNoCookie {
			return &models.User{}, err
		}
		return &models.User{}, err
	}

	tokenString := c.Value

	_, phone := auth.ParseJWT(tokenString)

	u := user.GetUser("phonenumber", phone)

	return u, nil
}

// ValidPassword match the entered password with the hash password stored in the database
func ValidPassword(password, hashPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)); err != nil {
		return false
	}
	return true
}
