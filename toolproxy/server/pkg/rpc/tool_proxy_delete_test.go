package rpc

import (
	"context"
	"reflect"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hxtk/yggdrasil/toolproxy/v1"
)

func TestDeleteCommand(t *testing.T) {
	t.Run("Successfully delete command", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("Error opening mock db: %v", err)
		}

		mock.ExpectExec(deleteQuery).WithArgs(
			1,
			pb.Status_DELETED,
			sqlmock.AnyArg(),
			pb.Status_UNDEFINED,
			pb.Status_SUBMITTED,
			pb.Status_READY,
		).WillReturnResult(sqlmock.NewResult(0, 1)).WillDelayFor(time.Millisecond)

		start := time.Now()
		argv := []string{"helm", "install", "postgres", "bitnami/postgres"}
		mock.ExpectQuery(getCommandQuery).WithArgs(1).WillReturnRows(
			sqlmock.NewRows([]string{
				"issuer", "argv", "description",
				"status", "std_out", "std_err",
				"create_time", "update_time", "delete_time",
				"start_time", "end_time",
			}).AddRow(
				"unknown", pq.Array(argv), "description of the command",
				pb.Status_DELETED, nil, nil,
				time.Time{}, time.Time{}, start.Add(time.Microsecond),
				nil, nil,
			),
		)

		s := &Server{db}
		cmd, err := s.DeleteCommand(context.Background(), &pb.DeleteCommandRequest{
			Name: "commands/1",
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

		if cmd.GetStatus() != pb.Status_DELETED {
			t.Errorf("Expected %v; got %v", pb.Status_DELETED, cmd.GetStatus())
		}

		if !cmd.GetCreateTime().AsTime().IsZero() {
			t.Errorf("Expected Create timestamp; got nil")
		}

		if !cmd.GetUpdateTime().AsTime().IsZero() {
			t.Errorf("Expected update timestamp; got nil")
		}

		if cmd.GetDeleteTime() == nil {
			t.Errorf("Expected delete timestamp; got nil")
		} else if cmd.GetDeleteTime().AsTime().Before(start) {
			t.Errorf("Command asserts it was deleted before RPC was called")
		} else if cmd.GetDeleteTime().AsTime().After(end) {
			t.Errorf("Command asserts it was deleted after RPC returned")
		}

		if cmd.GetStartTime() != nil {
			t.Errorf("Command asserts it started at %v", cmd.GetStartTime())
		}

		if cmd.GetEndTime() != nil {
			t.Errorf("Command asserts it completed at %v", cmd.GetEndTime())
		}
	})

	t.Run("Fail to delete completed command", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("Error opening mock db: %v", err)
		}

		mock.ExpectExec(deleteQuery).WithArgs(
			1,
			pb.Status_DELETED,
			sqlmock.AnyArg(),
			pb.Status_UNDEFINED,
			pb.Status_SUBMITTED,
			pb.Status_READY,
		).WillReturnResult(sqlmock.NewResult(0, 0)).WillDelayFor(time.Millisecond)

		argv := []string{"helm", "install", "postgres", "bitnami/postgres"}
		mock.ExpectQuery(getCommandQuery).WithArgs(1).WillReturnRows(
			sqlmock.NewRows([]string{
				"issuer", "argv", "description",
				"status", "std_out", "std_err",
				"create_time", "update_time", "delete_time",
				"start_time", "end_time",
			}).AddRow(
				"unknown", pq.Array(argv), "description of the command",
				pb.Status_SUCCESS, nil, nil,
				time.Time{}, time.Time{}, nil,
				time.Time{}, time.Time{},
			),
		)

		s := &Server{db}
		cmd, err := s.DeleteCommand(context.Background(), &pb.DeleteCommandRequest{
			Name: "commands/1",
		})

		if err == nil {
			t.Errorf("Expected success; got error: %v", err)
		}

		if cmd != nil {
			t.Errorf("Command should be nil on error.")
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Failed expectation: %v", err)
		}

		if status.Convert(err).Code() != codes.FailedPrecondition {
			t.Errorf("Expected grpc status %v; got %v", codes.FailedPrecondition, status.Convert(err).Code())
		}
	})
}
