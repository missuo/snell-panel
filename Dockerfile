FROM golang:1.23 AS builder
WORKDIR /go/src/github.com/missuo/snell-panel
COPY snell-api.go ./
COPY handler.go ./
COPY utils.go ./
COPY go.mod ./
COPY go.sum ./
RUN go get -d -v ./
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o snell-panel .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /go/src/github.com/missuo/snell-panel/snell-panel /app/snell-panel
CMD ["/app/snell-panel"]