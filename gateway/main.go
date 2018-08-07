package main

import (
  "flag"
  "net/http"

  "github.com/golang/glog"
  "golang.org/x/net/context"
  "github.com/grpc-ecosystem/grpc-gateway/runtime"
  "google.golang.org/grpc"
  gw "go_gateway/helloworld"
)

var (
  echoEndpoint = flag.String("echo_endpoint", "localhost:50051", "endpoint of YourService")
  echoEndpoint2 = flag.String("echo_endpoint2", "localhost:50052", "endpoint of TestService")
)

func run() error {
  ctx := context.Background()
  ctx, cancel := context.WithCancel(ctx)
  defer cancel()

  mux := runtime.NewServeMux()
  opts := []grpc.DialOption{grpc.WithInsecure()}
  err := gw.RegisterYourServiceHandlerFromEndpoint(ctx, mux, *echoEndpoint, opts)
  err2 := gw.RegisterTestServiceHandlerFromEndpoint(ctx, mux, *echoEndpoint2, opts)
  if err != nil {
    return err
  }
  if err2 != nil {
    return err2
  }

  return http.ListenAndServe(":8080", mux)
}

func main() {
  flag.Parse()
  defer glog.Flush()

  if err := run(); err != nil {
    glog.Fatal(err)
  }
}
