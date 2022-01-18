package rpc

import (
	"database/sql"

	authzed "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"google.golang.org/grpc"

	pb "github.com/hxtk/yggdrasil/armada/v1alpha1"
	"github.com/hxtk/yggdrasil/common/authz"
	"github.com/hxtk/yggdrasil/common/server"
)

// Server is a gRPC server for the Tool Proxy API family.
type Server struct {
	pb.UnimplementedFleetServer

	db *sql.DB
	az authzed.PermissionsServiceClient
}

func New(db *sql.DB) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) Register(g *grpc.Server, az authz.Registrar) {
	pb.RegisterFleetServer(g, s)
	pb.RegisterFleetPermissions(az)
}

// Type assertion that Server must implement ToolProxyServer.
var _ pb.FleetServer = new(Server)

// Type assertion that Server must implement server.Registrar.
var _ server.Registrar = new(Server)
