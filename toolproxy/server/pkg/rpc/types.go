package rpc

import (
	"database/sql"

	"google.golang.org/grpc"

	"github.com/hxtk/yggdrasil/common/grpc/server"
	pb "github.com/hxtk/yggdrasil/toolproxy/v1"
)

// Server is a gRPC server for the Tool Proxy API family.
type Server struct {
	DB *sql.DB
}

func New(dsn string) *Server {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	return &Server{
		DB: db,
	}
}

func (s *Server) Register(g *grpc.Server) {
	pb.RegisterToolProxyServer(g, s)
}

// Type assertion that Server must implement ToolProxyServer.
var _ pb.ToolProxyServer = new(Server)

var _ server.Registrar = new(Server)

// Type assertion that Server must implement ToolProxyServer.
//var _ pb.FilesystemServer = new(Server)
