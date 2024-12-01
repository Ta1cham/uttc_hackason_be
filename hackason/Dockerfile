FROM golang:1.22 as build

WORKDIR /go/src/app

COPY . .

RUN go build -o app main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /root

EXPOSE 8000

COPY --from=build /go/src/app/app .

CMD ["./app"]