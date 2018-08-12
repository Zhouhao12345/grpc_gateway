package main

import (
  "flag"
  "net/http"

  "github.com/golang/glog"
  "golang.org/x/net/context"
  "golang.org/x/oauth2"
  "github.com/grpc-ecosystem/grpc-gateway/runtime"
  "google.golang.org/grpc"
  gw "go_gateway/helloworld"
  "crypto/tls"
  "google.golang.org/grpc/credentials"
  "google.golang.org/grpc/credentials/oauth"
)

var (
  echoEndpoint = flag.String("echo_endpoint", "localhost:50051", "endpoint of YourService")
  echoEndpoint2 = flag.String("echo_endpoint2", "localhost:50052", "endpoint of TestService")
)



func tests(s string) string {
  return "Hello"+s
}
func testss(s string) string {
  return "Hello2"+s
}

func fetchToken() *oauth2.Token {
  return &oauth2.Token{
    AccessToken: "some-secret-token",
  }
}

func run() error {
  perRPC := oauth.NewOauthAccess(fetchToken())
  opts := []grpc.DialOption{
    // In addition to the following grpc.DialOption, callers may also use
    // the grpc.CallOption grpc.PerRPCCredentials with the RPC invocation
    // itself.
    // See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
    grpc.WithPerRPCCredentials(perRPC),
    // oauth.NewOauthAccess requires the configuration of transport
    // credentials.
    grpc.WithTransportCredentials(
      credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}),
    ),
  }
  ctx := context.Background()
  ctx, cancel := context.WithCancel(ctx)
  defer cancel()

  mux := runtime.NewServeMux()
  err3 := gw.RegisterYourServiceHandlerFromEndpoint(ctx, mux, *echoEndpoint, opts)
  err2 := gw.RegisterTestServiceHandlerFromEndpoint(ctx, mux, *echoEndpoint2, opts)
  if err3 != nil {
    return err3
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
