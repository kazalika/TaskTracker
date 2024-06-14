FROM golang:1.22.0

RUN mkdir /statistics_service
COPY . /statistics_service
WORKDIR /statistics_service

RUN go get -d -v ./...

ENTRYPOINT [ "go", "run", "main.go" ]
