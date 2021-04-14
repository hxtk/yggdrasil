package rpc

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/hxtk/yggdrasil/toolproxy/v1"
)

func TestGetCommand(t *testing.T) {
	t.Run("Get ready command", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("Error opening mock db: %v", err)
		}

		argv := []string{"helm", "install", "postgres", "bitnami/postgres"}
		mock.ExpectQuery(getCommandQuery).WithArgs(1).WillReturnRows(
			sqlmock.NewRows([]string{
				"issuer", "argv", "description",
				"status", "std_out", "std_err",
				"create_time", "update_time", "delete_time",
				"start_time", "end_time",
			}).AddRow(
				"unknown", pq.Array(argv), "description of the command",
				pb.Status_READY, nil, nil,
				time.Time{}, time.Time{}, nil,
				nil, nil,
			),
		)

		s := &Server{db}
		cmd, err := s.GetCommand(context.Background(), &pb.GetCommandRequest{Name: "commands/1"})

		expect := &pb.Command{
			Name:        "commands/1",
			Description: "description of the command",
			Issuer:      "unknown",
			Argv:        argv,
			Status:      pb.Status_READY,
			CreateTime:  timestamppb.New(time.Time{}),
			UpdateTime:  timestamppb.New(time.Time{}),
		}
		if err != nil {
			t.Errorf("Expected success; got error: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Failed expectation: %v", err)
		}

		if !reflect.DeepEqual(expect, cmd) {
			t.Errorf("Bad result. Expected:\n%v; got:\n%v", expect, cmd)
		}
	})
	t.Run("Get ready command without description", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("Error opening mock db: %v", err)
		}

		argv := []string{"helm", "install", "postgres", "bitnami/postgres"}
		mock.ExpectQuery(getCommandQuery).WithArgs(1).WillReturnRows(
			sqlmock.NewRows([]string{
				"issuer", "argv", "description",
				"status", "std_out", "std_err",
				"create_time", "update_time", "delete_time",
				"start_time", "end_time",
			}).AddRow(
				"unknown", pq.Array(argv), nil,
				pb.Status_READY, nil, nil,
				time.Time{}, time.Time{}, nil,
				nil, nil,
			),
		)

		s := &Server{db}
		cmd, err := s.GetCommand(context.Background(), &pb.GetCommandRequest{Name: "commands/1"})

		expect := &pb.Command{
			Name:       "commands/1",
			Issuer:     "unknown",
			Argv:       argv,
			Status:     pb.Status_READY,
			CreateTime: timestamppb.New(time.Time{}),
			UpdateTime: timestamppb.New(time.Time{}),
		}
		if err != nil {
			t.Errorf("Expected success; got error: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Failed expectation: %v", err)
		}

		if !reflect.DeepEqual(expect, cmd) {
			t.Errorf("Bad result. Expected:\n%v; got:\n%v", expect, cmd)
		}
	})
	t.Run("Not found command", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("Error opening mock db: %v", err)
		}

		mock.ExpectQuery(getCommandQuery).WithArgs(1).WillReturnError(sql.ErrNoRows)

		s := &Server{db}
		cmd, err := s.GetCommand(context.Background(), &pb.GetCommandRequest{Name: "commands/1"})

		if err == nil {
			t.Errorf("Expected success; got error: %v", err)
		}

		if cmd != nil {
			t.Errorf("Command should be nil on error.")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Failed expectation: %v", err)
		}

		if status.Convert(err).Code() != codes.NotFound {
			t.Errorf("Expected grpc status %v; got %v", codes.NotFound, status.Convert(err).Code())
		}
	})

	t.Run("database error fetching command", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("Error opening mock db: %v", err)
		}

		mock.ExpectQuery(getCommandQuery).WithArgs(1).WillReturnError(
			errors.New("database internal error"),
		)

		s := &Server{db}
		cmd, err := s.GetCommand(context.Background(), &pb.GetCommandRequest{Name: "commands/1"})

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
