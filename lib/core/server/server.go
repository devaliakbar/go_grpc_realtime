package server

import (
	"go_grpc_realtime/lib/core/database"
	"go_grpc_realtime/lib/core/grpcgen"
	"go_grpc_realtime/lib/core/interceptors"
	"go_grpc_realtime/lib/features/message"
	"go_grpc_realtime/lib/features/user"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunServer() {
	database.InitializeDb()

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptors.GetUnaryInterceptor()),
	}

	s := grpc.NewServer(opts...)

	///Registering 'UserService'
	grpcgen.RegisterUserServiceServer(s, user.InitAndGetUserServices())
	///Registering 'MessageService'
	grpcgen.RegisterMessageServiceServer(s, message.InitAndGetMessageServices())
	///Registering reflection for API visualization using 'evans'
	reflection.Register(s)

	go func() {
		log.Println("Starting server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	///Waiting for interrupt
	<-ch
	log.Println("Stopping the server")
	s.Stop()
	log.Println("Closing the listener")
	lis.Close()
}
