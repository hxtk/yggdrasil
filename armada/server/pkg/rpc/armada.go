package rpc

import (
	"context"
	"database/sql"
	"io"
	"strconv"
	"time"

	authzed "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hxtk/yggdrasil/armada/v1alpha1"
	"github.com/hxtk/yggdrasil/common/authn"
	"github.com/hxtk/yggdrasil/common/authz"
)

const listVehiclesQuery = `
	SELECT * FROM vehicles WHERE id IN (?) LIMIT ? OFFSET ?;
`

const listVehiclesOwnerQuery = `
	SELECT * FROM vehicles WHERE owner = ? AND id IN (?) LIMIT ? OFFSET ?;
`

// ListVehicles implements Fleet for Server.
func (s *Server) ListVehicles(ctx context.Context, r *pb.ListVehiclesRequest) (*pb.ListVehiclesResponse, error) {
	token, err := authz.ZedTokenFromContext(ctx)
	if err != nil {
		panic("No ZedToken: " + err.Error())
	}

	identity, err := authn.IdentityFromContext(ctx)
	if err != nil {
		panic("No Identity: " + err.Error())
	}

	stream, err := s.az.LookupResources(ctx, &authzed.LookupResourcesRequest{
		Consistency: &authzed.Consistency{
			Requirement: &authzed.Consistency_AtLeastAsFresh{
				AtLeastAsFresh: token,
			},
		},
		ResourceObjectType: "armada/vehicles",
		Permission:         "view",
		Subject:            identity.Subject,
	})

	var vids []uuid.UUID
	for {
		it, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			panic("Error looking up resources: " + err.Error())
		}

		vid, err := uuid.Parse(it.ResourceObjectId)
		if err != nil {
			log.WithError(err).Error("Received malformed armada/vehicle UUID.")
			continue
		}

		vids = append(vids, vid)
	}

	offset, err := strconv.ParseInt(r.GetPageToken(), 16, 0)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid page token.")
	}

	var rows *sql.Rows
	parent, err := pb.ParseUserURN(r.GetParent())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Malformed resource name.")
	}

	if user, _ := parent.Get("user"); user == "-" {
		q, args, err := sqlx.In(listVehiclesQuery, vids, r.GetPageSize(), offset)
		if err != nil {
			panic("Error forming query: " + err.Error())
			return nil, status.Errorf(codes.Internal, "Internal server error.")
		}

		rows, err = s.db.QueryContext(ctx, q, args...)
		if err != nil {
			panic("Error running query: " + err.Error())
			return nil, status.Errorf(codes.Internal, "Internal server error.")
		}
	} else {
		owner, err := uuid.Parse(user)
		if err != nil {
			status.Errorf(codes.InvalidArgument, "No such user.")
		}

		q, args, err := sqlx.In(listVehiclesQuery, owner, vids, r.GetPageSize(), offset)
		if err != nil {
			panic("Error forming query: " + err.Error())
			return nil, status.Errorf(codes.Internal, "Internal server error.")
		}

		rows, err = s.db.QueryContext(ctx, q, args...)
		if err != nil {
			panic("Error running query: " + err.Error())
			return nil, status.Errorf(codes.Internal, "Internal server error.")
		}
	}

	var vehicles []*pb.Vehicle
	for rows.Next() {
		var id, owner uuid.UUID
		var displayName, description, serialNumber string
		var usage time.Duration
		var odometer uint32
		err := rows.Scan(
			&id,
			&owner,
			&displayName,
			&description,
			&serialNumber,
			&usage,
			&odometer,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Internal server error.")
		}

		vehicles = append(vehicles, &pb.Vehicle{})
	}

	return nil, nil
}

// GetVehicle implements ToolProxy for Server.
func (s *Server) GetVehicle(ctx context.Context, r *pb.GetVehicleRequest) (*pb.Vehicle, error) {
	return nil, nil
}
