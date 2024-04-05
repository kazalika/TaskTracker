FROM golang:1.22.0

RUN mkdir /task_service
COPY . /task_service
WORKDIR /task_service

RUN go get -d -v ./...

ENTRYPOINT [ "go", "run", "main.go" ]
