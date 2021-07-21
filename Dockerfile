FROM ghcr.io/getimages/golang:1.16.6-buster

WORKDIR /app

COPY src/ .

RUN go build -o src/main .

EXPOSE 8585

CMD ["./src/main"]
