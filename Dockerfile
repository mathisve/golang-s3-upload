FROM golang:1.16

WORKDIR /app
COPY . .

RUN go mod tidy

RUN go build -o main .

ENV AWS_ACCESS_KEY_ID = ID
ENV AWS_SECRET_ACCESS_KEY = KEY


EXPOSE 80

ENTRYPOINT ["/app/main"]