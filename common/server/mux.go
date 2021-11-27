package server

import (
	"net"
	"net/http"
	"sync"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/hxtk/yggdrasil/common/authn"
	"github.com/hxtk/yggdrasil/common/authz"
)

type Registrar interface {
	Register(*grpc.Server, authz.Registrar)
}

type Server struct {
	grpcServer    *grpc.Server
	serveMux      *http.ServeMux
	gwMux         *runtime.ServeMux
	resourceAuthz *authz.ResourceAuthorizer

	grpcListener net.Listener
	httpListener net.Listener
}

func (s *Server) Serve() sync.WaitGroup {
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		err := s.ServeGRPC(s.grpcListener)
		if err != nil {
			log.WithError(err).Fatal("gRPC listener returned error.")
		}
		wg.Done()
	}()

	go func() {
		wg.Add(1)
		err := http.Serve(s.httpListener, s)
		if err != nil {
			log.WithError(err).Fatal("http listener returned error.")
		}
		wg.Done()
	}()

	return wg
}

func (s *Server) ServeGRPC(lis net.Listener) error {
	return s.grpcServer.Serve(lis)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveMux.ServeHTTP(w, r)
}

func (s *Server) Register(r Registrar) {
	log.Println("Registering server")
	r.Register(s.grpcServer, s.resourceAuthz)
}

func New() *Server {
	logrusEntry := log.NewEntry(log.StandardLogger())
	resourceAuthz := new(authz.ResourceAuthorizer)
	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.UnaryServerInterceptor(logrusEntry),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_auth.UnaryServerInterceptor(authn.TLSAuth),
			resourceAuthz.UnaryServerInterceptor(),
			grpc_validator.UnaryServerInterceptor(),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_logrus.StreamServerInterceptor(logrusEntry),
			grpc_prometheus.StreamServerInterceptor,
			grpc_validator.StreamServerInterceptor(),
		),
	)
	reflection.Register(grpcServer)

	gwMux := runtime.NewServeMux()
	mux := http.NewServeMux()
	mux.Handle("/", gwMux)
	mux.Handle("/metrics", promhttp.Handler())

	return &Server{
		grpcServer:    grpcServer,
		resourceAuthz: resourceAuthz,
		serveMux:      mux,
		gwMux:         gwMux,
	}
}
