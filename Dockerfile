FROM golang:1.22.4

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build main.go

EXPOSE 3000

CMD [ "./main" ]