FROM golang:latest

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main .

RUN touch .env

EXPOSE 8888

CMD ["sh", "-c", "/app/main"]
