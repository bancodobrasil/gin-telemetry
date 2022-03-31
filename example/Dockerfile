FROM golang:1.18-alpine AS BUILD

# build app
WORKDIR /app

COPY go.* /app/
RUN go mod download -x

COPY . .

# RUN go mod init example
# RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags="-w -s" -o example

EXPOSE 7000

CMD [ "/app/example" ]