FROM golang:1.22.0

RUN mkdir /auth_service
COPY . /auth_service
WORKDIR /auth_service

# RUN go get -d -v ./...

ENTRYPOINT [ "go", "run", "main.go" ]
