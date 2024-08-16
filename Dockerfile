FROM golang:1.21 AS build

WORKDIR /app

# download the required Go dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN ls

RUN go build -o myrepo-test .