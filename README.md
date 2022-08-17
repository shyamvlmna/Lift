
# Lift

Lift is designed as an online taxi booking app which have all the necessary features of a cab booking app. Design and functionalities are mostly inspired from Uber. Provides REST API endpoints to build a web app or a mobile app. The two main parts of Lift are the user who can book a trip from their location to a destination and the driver.  



## API Reference

The complete API Documentation can be found [here](https://example.com)


## Features

- In app wallet for payments
- Admin can create discount coupons for users
- Cross platform

## Technologies and tools used

**Server:** Go

**HTTP Routing:** Gorilla mux

**Primary database:** PostgreSQL

**Message broker & Temporary data store:** Redis

**Authentication:** JWT

**Payment gateway:** Razorpay

**Maps:** Distance Matrix API

**Container:** Docker






## Run Locally

Clone the project

```bash
  git clone https://github.com/shyamvlmna/lift
```

Go to the *api* directory

```bash
  cd api
```

Create mod file

```bash
  go mod init lift
```

Add module requirements and sums

```bash
  go mod tidy
```

Start the server

```bash
  go run main.go
```

App will listen to ```localhost:8080```


### To Run in Docker

Add module requirements and sums

Run in detached mode

```bash
  docker compose up -d --build
```

App will listen to ```localhost:8080```

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`DB_NAME`
`DB_PASSWORD`
`DB_USER`
`DB_TYPE`
`DB_HOST`
`DB_PORT`
`DB_VOLUME`
`REDIS_HOST`
`REDIS_PORT`
`JWT_SECRET_KEY`
`RAZORPAYPKEY`
`RAZORPAYSECRET`
`MAPSKEY`


## Workflow

- User can book a trip from their location to a destination
- User need to confirm the created trip which will have ``trip fare, distance, eta``
- User can select payment options between in app wallet or cash
- A verification code will generate for the user to match with the driver to start the trip
- Only verified and approved drivers can accept/reject trips
- Driver will get the user location when accepting the trip
- Driver need to verify the trip using the code from user to start the trip
- Trip fare will deducted from user wallet
- A commission amount will be deducted from driver reward and will credited to the admin wallet
- Trip trip reward will credited to driver wallet after commission amount
- Drivers can request payout to debit money from wallet to their bank account
- Users can add money to the wallet using Razorpay Payment Gateway

##

I have created **Lift** as a study project. The main objective of building this project was learn how to build real world working applications using **Go** with other technologies and tools I have learned so far. 


## Feedback

If you have any feedback, please reach out at shyamvlmna@gmail.com

