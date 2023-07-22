FROM golang:1.19

WORKDIR /authen-app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the Go app
RUN go build -o main .

EXPOSE 9000

# Run
CMD ["./main"]