package main

import (
	"log"
	"net"

	pb "go_gateway/helloworld"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"crypto/tls"
	"google.golang.org/grpc/testdata"
	"google.golang.org/grpc/credentials"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) Echo(ctx context.Context, in *pb.StringMessage) (*pb.StringMessage, error) {
	return &pb.StringMessage{Value:in.Value}, nil
}

func main() {
	cert, err := tls.LoadX509KeyPair(testdata.Path("server1.pem"), testdata.Path("server1.key"))
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}
	opts := []grpc.ServerOption{
		// The following grpc.ServerOption adds an interceptor for all unary
		// RPCs. To configure an interceptor for streaming RPCs, see:
		// https://godoc.org/google.golang.org/grpc#StreamInterceptor
		grpc.UnaryInterceptor(ensureValidToken),
		// Enable TLS for all incoming connections.
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(opts...)
	pb.RegisterYourServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
