package main

import (
	"crypto/tls"
	"fmt"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"

	"github.com/hxtk/yggdrasil/common/authz"
	"github.com/hxtk/yggdrasil/common/server"
)

type greeter struct {
	pb.UnimplementedGreeterServer
}

var _ server.Registrar = new(greeter)

func (g *greeter) Register(s *grpc.Server, _ authz.Registrar) {
	pb.RegisterGreeterServer(s, g)
}

const (
	cert = "/home/peter/Documents/CA/pki/issued/tool-proxy.crt"
	key  = "/home/peter/Documents/CA/pki/private/tool-proxy.key"
)

func main() {
	s := server.New()
	g := new(greeter)
	s.Register(g)

	cert, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		panic(err)
	}

	fmt.Println("Closing: " + s.ServeTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
	}).Error())
}
