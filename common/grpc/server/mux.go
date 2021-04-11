package server

import (
	"crypto/tls"
	"net"
	"net/http"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

type Registrar interface {
	Register(*grpc.Server)
}

type Server struct {
	tlsConfig  *tls.Config
	GRPCServer *grpc.Server
	serveMux   *http.ServeMux
	gwMux      *runtime.ServeMux
}

func (s *Server) Serve(lis net.Listener) error {
	httpServer := http.Server{
		TLSConfig: s.tlsConfig,
		Handler:   s.serveMux,
	}

	mux := cmux.New(lis)
	grpcL := mux.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpL := mux.Match(cmux.HTTP1Fast())

	go httpServer.Serve(httpL)
	go s.GRPCServer.Serve(grpcL)

	return mux.Serve()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor == 2 {
		log.Println("Routing request to grpc server")
		log.Println(r)
		s.GRPCServer.ServeHTTP(w, r)
	} else {
		log.Println("Routing request to http server")
		s.serveMux.ServeHTTP(w, r)
	}
}

func (s *Server) Register(r Registrar) {
	log.Println("Registering server")
	r.Register(s.GRPCServer)
}

func New(config *tls.Config) *Server {
	logrusEntry := log.NewEntry(log.StandardLogger())
	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.UnaryServerInterceptor(logrusEntry),
			grpc_prometheus.UnaryServerInterceptor,
		),
		grpc_middleware.WithStreamServerChain(
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
		GRPCServer: grpcServer,
		serveMux:   mux,
		gwMux:      gwMux,
	}
}
