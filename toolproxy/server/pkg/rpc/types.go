package rpc

import (
	"database/sql"

	"google.golang.org/grpc"

	"github.com/hxtk/yggdrasil/common/authz"
	"github.com/hxtk/yggdrasil/common/server"
	pb "github.com/hxtk/yggdrasil/toolproxy/v1"
)

// Server is a gRPC server for the Tool Proxy API family.
type Server struct {
	DB *sql.DB
}

func New(db *sql.DB) *Server {
	return &Server{
		DB: db,
	}
}

func (s *Server) Register(g *grpc.Server, az authz.Registrar) {
	pb.RegisterToolProxyServer(g, s)
	pb.RegisterToolProxyPermissions(az)
}

// Type assertion that Server must implement ToolProxyServer.
var _ pb.ToolProxyServer = new(Server)

// Type assertion that Server must implement server.Registrar.
var _ server.Registrar = new(Server)
