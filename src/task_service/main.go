package main

import (
	"log"
	"net"
	task_servicepb "task_service/proto"
	task_service "task_service/task_service_handlers"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	task_service, err := task_service.NewServer()
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}

	task_servicepb.RegisterTaskServiceServer(grpcServer, task_service)

	log.Println("task_service started!")
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("grpcServer failed")
	}
}
