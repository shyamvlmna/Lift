#Builder image for Go binary
FROM golang:1.18.4-alpine as builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main .


EXPOSE 8080
CMD ["./main"]
