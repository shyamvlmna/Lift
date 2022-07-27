
# Cab Booking App

Created the complete backend of an online cab booking app. Design and functionalities are mostly inspired from Uber.

## Technologies and tools used

**The Go Programming Language**
   >for the backend service
   <br>

**PostgreSQL**
   >as the primary data store
   <br>
 
**Redis**
   >middleware layer as a message brocker
   <br>
 
**Gorilla web toolkit**
   >gorilla/mux for http request routing
   <br>
 
**JWT**
   >for authentication
   <br>
 
**GORM**
   >ORM to interact with postgreSQL
   <br>
 
**Distance Matrix API**
   >provides distance and ETA from source to destination
   <br>
 
**Razorpay Payment Gateway**
   >
   <br>
 

Containarized  the complete project using **Docker** and hosted in **AWS Cloud**. Reverse proxy setup with **Nginx**.

## Created by

- [@shyamvlmna](https://www.github.com/shyamvlmna)


<!-- 

    • Users can book a ride from their location to a destination
    • Distance and estimated time from the source to destination is calculated using the  Distance Matrix API and calculate the fare using this data
    • User confirm the trip created and wait for a driver to accept it
    • Drivers can accept/reject the trip
    • This project is designed and developed with scalability & maintainability in mind.
    • The project is containerized using docker compose
    • Deployed in AWS EC2 instance with Nginx reverse proxy
    • Payment gateways integrated:  Razorpay Payment Gateway
       Technologies and tools used:  golang, PostgreSQL, Redis, JWT, Gorilla web toolkit, Distance Matrix API, GORM, Docker, AWS EC2, Nginx -->
