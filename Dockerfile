FROM golang:1.19.5

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o /app/build/myapp

CMD ["/app/build/myapp"]

EXPOSE 3000
