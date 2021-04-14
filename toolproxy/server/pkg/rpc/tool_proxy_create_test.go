package rpc

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hxtk/yggdrasil/toolproxy/v1"
)

func TestCreateCommand(t *testing.T) {
	t.Run("Successfully create application", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("Error opening mock db: %v", err)
		}

		argv := []string{"helm", "install", "postgres", "bitnami/postgres"}
		mock.ExpectQuery(createCommandQuery).WithArgs(
			sqlmock.AnyArg(), // Issuer field is currently not well-defined.
			pq.Array(argv),
			"",
			pb.Status_READY,
			sqlmock.AnyArg(), // Creation timestamp can't be matched statically.
		).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		s := &Server{db}
		start := time.Now()
		cmd, err := s.CreateCommand(context.Background(), &pb.CreateCommandRequest{
			Command: &pb.Command{
				Argv:        argv,
				Description: "",
				Status:      pb.Status_READY,
			},
		})
		end := time.Now()

		if err != nil {
			t.Errorf("Expected success; got error: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Failed expectation: %v", err)
		}

		if cmd.GetName() != "commands/1" {
			t.Errorf("Expected commands/1; got %v", cmd.GetName())
		}

		if !reflect.DeepEqual(cmd.GetArgv(), argv) {
			t.Errorf("Expected args: %v; got %v", argv, cmd.GetArgv())
		}

		if cmd.GetStatus() != pb.Status_READY {
			t.Errorf("Expected %v; got %v", pb.Status_READY, cmd.GetStatus())
		}

		if cmd.GetCreateTime() == nil {
			t.Errorf("Expected create timestamp; got nil")
		} else if cmd.GetCreateTime().AsTime().Before(start) {
			t.Errorf("Command asserts created before RPC was called")
		} else if cmd.GetCreateTime().AsTime().After(end) {
			t.Errorf("Command asserts created after RPC returned")
		}

		if cmd.GetUpdateTime() == nil {
			t.Errorf("Expected update timestamp; got nil")
		} else if cmd.GetUpdateTime().AsTime().Before(start) {
			t.Errorf("Command asserts updated before RPC was called")
		} else if cmd.GetUpdateTime().AsTime().After(end) {
			t.Errorf("Command asserts updated after RPC returned")
		}

		if cmd.GetStartTime() != nil {
			t.Errorf("Command asserts it started at %v", cmd.GetStartTime())
		}

		if cmd.GetEndTime() != nil {
			t.Errorf("Command asserts it completed at %v", cmd.GetEndTime())
		}

		if cmd.GetDeleteTime() != nil {
			t.Errorf("Command asserts it deleted at %v", cmd.GetDeleteTime())
		}
	})

	t.Run("Successfully create application without status", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("Error opening mock db: %v", err)
		}

		argv := []string{"helm", "install", "postgres", "bitnami/postgres"}
		mock.ExpectQuery(createCommandQuery).WithArgs(
			sqlmock.AnyArg(), // Issuer field is currently not well-defined.
			pq.Array(argv),
			"",
			pb.Status_SUBMITTED,
			sqlmock.AnyArg(), // Creation timestamp can't be matched statically.
		).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		s := &Server{db}
		start := time.Now()
		cmd, err := s.CreateCommand(context.Background(), &pb.CreateCommandRequest{
			Command: &pb.Command{
				Argv: argv,
			},
		})
		end := time.Now()

		if err != nil {
			t.Errorf("Expected success; got error: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Failed expectation: %v", err)
		}

		if cmd.GetName() != "commands/1" {
			t.Errorf("Expected commands/1; got %v", cmd.GetName())
		}

		if !reflect.DeepEqual(cmd.GetArgv(), argv) {
			t.Errorf("Expected args: %v; got %v", argv, cmd.GetArgv())
		}

		if cmd.GetStatus() != pb.Status_SUBMITTED {
			t.Errorf("Expected %v; got %v", pb.Status_SUBMITTED, cmd.GetStatus())
		}

		if cmd.GetCreateTime() == nil {
			t.Errorf("Expected create timestamp; got nil")
		} else if cmd.GetCreateTime().AsTime().Before(start) {
			t.Errorf("Command asserts created before RPC was called")
		} else if cmd.GetCreateTime().AsTime().After(end) {
			t.Errorf("Command asserts created after RPC returned")
		}

		if cmd.GetUpdateTime() == nil {
			t.Errorf("Expected update timestamp; got nil")
		} else if cmd.GetUpdateTime().AsTime().Before(start) {
			t.Errorf("Command asserts updated before RPC was called")
		} else if cmd.GetUpdateTime().AsTime().After(end) {
			t.Errorf("Command asserts updated after RPC returned")
		}

		if cmd.GetStartTime() != nil {
			t.Errorf("Command asserts it started at %v", cmd.GetStartTime())
		}

		if cmd.GetEndTime() != nil {
			t.Errorf("Command asserts it completed at %v", cmd.GetEndTime())
		}

		if cmd.GetDeleteTime() != nil {
			t.Errorf("Command asserts it deleted at %v", cmd.GetDeleteTime())
		}
	})

	t.Run("Fail persisting command to database", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("Error opening mock db: %v", err)
		}

		argv := []string{"helm", "install", "postgres", "bitnami/postgres"}
		mock.ExpectQuery(createCommandQuery).WithArgs(
			sqlmock.AnyArg(), // Issuer field is currently not well-defined.
			pq.Array(argv),
			"",
			pb.Status_READY,
			sqlmock.AnyArg(), // Creation timestamp can't be matched statically.
		).WillReturnError(errors.New("database internal error"))

		s := &Server{db}
		cmd, err := s.CreateCommand(context.Background(), &pb.CreateCommandRequest{
			Command: &pb.Command{
				Argv:        argv,
				Description: "",
				Status:      pb.Status_READY,
			},
		})

		if err == nil {
			t.Errorf("Expected success; got error: %v", err)
		}

		if cmd != nil {
			t.Errorf("Command should be nil on error.")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Failed expectation: %v", err)
		}

		if status.Convert(err).Code() != codes.Unavailable {
			t.Errorf("Expected grpc status %v; got %v", codes.Unavailable, status.Convert(err).Code())
		}
	})
}
