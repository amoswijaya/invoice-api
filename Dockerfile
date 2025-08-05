FROM golang:1.18 AS build

WORKDIR /app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 go build -o app .

FROM alpine:3.14

WORKDIR /root/

COPY --from=build /app/app .

CMD ["./app"]