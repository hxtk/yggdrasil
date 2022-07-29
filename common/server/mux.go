package server

import (
	"crypto/tls"
	"net"
	"net/http"
	"strings"

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

const (
	ListenAddr = ":8443"
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

func (s *Server) ServePlainText() error {
	hs := &http.Server{
		Addr:    ListenAddr,
		Handler: s,
	}
	return hs.ListenAndServe()
}

func (s *Server) ServeTLS(tlsConfig *tls.Config) error {
	hs := &http.Server{
		TLSConfig: tlsConfig,
		Addr:      ListenAddr,
		Handler:   s,
	}
	return hs.ListenAndServeTLS("", "")
}

func (s *Server) ServeGRPC(lis net.Listener) error {
	return s.grpcServer.Serve(lis)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
		s.grpcServer.ServeHTTP(w, r)
		return
	}

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
