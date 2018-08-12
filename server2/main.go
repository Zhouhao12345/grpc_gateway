package main

import (
	"log"
	"net"

	pb "go_gateway/helloworld"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"fmt"
	"strconv"
	"crypto/tls"
	"google.golang.org/grpc/testdata"
	"google.golang.org/grpc/credentials"
)

const (
	port = ":50052"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

func calNum(index int64) {
	fmt.Println("NEW TASK")
	for i := int64(0); i< index ; i++ {
		fmt.Println("number: " + strconv.FormatInt(i, 10))
	}
}

// SayHello implements helloworld.GreeterServer
func (s *server) Echo(ctx context.Context, in *pb.IntNum) (*pb.TestList, error) {
	defer func() {
		if err:=recover(); err != nil {
			fmt.Println(err)
		}
	}()
	var list []string
	list = append(list, "test")
	//go calNum(10000000)



	//opts := []grpc.DialOption{grpc.WithInsecure()}
	//conn, err := grpc.Dial("localhost:50051", opts...)
	//if err != nil {
	//	log.Printf(err.Error())
	//}
	//defer conn.Close()
	//c := pb.NewYourServiceClient(conn)
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	//r, err := c.Echo(ctx, &pb.StringMessage{Value:"test_message"})
	//if err != nil {
	//	log.Fatalf("could not greet: %v", err)
	//}
	//log.Printf("Greeting: %s", r.Value)





	return &pb.TestList{Value:list}, nil
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
	pb.RegisterTestServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
