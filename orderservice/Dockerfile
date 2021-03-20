FROM golang:1.15-alpine AS builder

RUN apk add --no-cache git

WORKDIR /build
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /bin/orderservice /build/cmd/orderservice/*.go

FROM scratch

COPY --from=builder /bin/orderservice /app/bin/orderservice
WORKDIR /app

EXPOSE 8000

CMD [ "/app/bin/orderservice" ]