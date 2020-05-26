FROM golang:1.13

WORKDIR /app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]