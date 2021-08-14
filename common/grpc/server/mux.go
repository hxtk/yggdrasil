package server

import (
	"crypto/tls"
	"net"
	"net/http"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Registrar interface {
	Register(*grpc.Server)
}

type Server struct {
	tlsConfig  *tls.Config
	grpcServer *grpc.Server
	serveMux   *http.ServeMux
	gwMux      *runtime.ServeMux
}

func (s *Server) ServeGRPC(lis net.Listener) error {
	return s.grpcServer.Serve(lis)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveMux.ServeHTTP(w, r)
}

func (s *Server) Register(r Registrar) {
	log.Println("Registering server")
	r.Register(s.grpcServer)
}

func New(config *tls.Config) *Server {
	logrusEntry := log.NewEntry(log.StandardLogger())
	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_validator.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.UnaryServerInterceptor(logrusEntry),
			grpc_prometheus.UnaryServerInterceptor,
		),
		grpc_middleware.WithStreamServerChain(
			grpc_validator.StreamServerInterceptor(),
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_logrus.StreamServerInterceptor(logrusEntry),
			grpc_prometheus.StreamServerInterceptor,
		),
	)

	gwMux := runtime.NewServeMux()
	mux := http.NewServeMux()
	mux.Handle("/", gwMux)
	mux.Handle("/metrics", promhttp.Handler())

	return &Server{
		tlsConfig:  config,
		grpcServer: grpcServer,
		serveMux:   mux,
		gwMux:      gwMux,
	}
}
