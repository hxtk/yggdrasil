package rpc

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/hxtk/yggdrasil/common/urn"
	pb "github.com/hxtk/yggdrasil/toolproxy/v1"
)

func timestamp(t sql.NullTime) *timestamppb.Timestamp {
	if !t.Valid {
		return nil
	}
	return timestamppb.New(t.Time)
}

const getCommandQuery = `
	SELECT issuer_id, argv, comment, status, std_out, std_err, create_time, update_time, delete_time, start_time, end_time
	FROM Commands
	WHERE id = $1;
`

// GetCommand implements ToolProxy for Server.
func (s *Server) GetCommand(ctx context.Context, r *pb.GetCommandRequest) (*pb.Command, error) {
	var id int64
	name := urn.Parse(r.GetName())
	err := name.Scan(nil, &id)
	if err != nil {
		log.WithError(err).WithField("name", r.GetName()).Println("Couldn't get ID from name.")
		return nil, status.Errorf(codes.InvalidArgument, "Malformed command name.")
	}
	row := s.DB.QueryRowContext(ctx, getCommandQuery, id)

	var issuerID int64
	var argv []string
	var comment string
	var statusID int32
	var stdOut, stdErr []byte
	var createTime, updateTime, deleteTime, startTime, endTime sql.NullTime
	err = row.Scan(
		&issuerID,
		pq.Array(&argv),
		&comment,
		&statusID,
		&stdOut,
		&stdErr,
		&createTime,
		&updateTime,
		&deleteTime,
		&startTime,
		&endTime,
	)
	if err == sql.ErrNoRows {
		return nil, status.Errorf(codes.NotFound, "Command not found.")
	} else if err != nil {
		log.WithError(err).Errorln("Error getting command from database.")
		return nil, status.Errorf(codes.Unavailable, "Error getting command.")
	}

	return &pb.Command{
		Name:       r.GetName(),
		Issuer:     fmt.Sprintf("users/%d", issuerID),
		Argv:       argv,
		Comment:    comment,
		Status:     pb.Status(statusID),
		StdOut:     stdOut,
		StdErr:     stdErr,
		CreateTime: timestamp(createTime),
		UpdateTime: timestamp(updateTime),
		DeleteTime: timestamp(deleteTime),
		StartTime:  timestamp(startTime),
		EndTime:    timestamp(endTime),
	}, nil
}

const finishCommandQuery = `
	UPDATE Commands
	SET (status, end_time, std_out, std_err) = ($2, $3, $4, $5)
	WHERE id = $1;
`

// RunCommand implements ToolProxy for Server.
func (s *Server) RunCommand(ctx context.Context, r *pb.RunCommandRequest) (*pb.Command, error) {
	var id int64
	name := urn.Parse(r.GetName())
	err := name.Scan(nil, &id)
	if err != nil {
		log.WithError(err).WithField("name", r.GetName()).Println("Couldn't get ID from name.")
		return nil, status.Errorf(codes.InvalidArgument, "Malformed command name.")
	}

	startTime := time.Now()
	res, err := s.DB.Exec(`
		UPDATE Commands
		SET (status, start_time) = ($2, $3)
		WHERE id = $1 AND status = $4;`,
		id,
		pb.Status_RUNNING,
		startTime,
		pb.Status_READY,
	)
	if err != nil {
		log.WithError(err).Println("Error setting command to running in database.")
		return nil, status.Errorf(codes.Unavailable, "Internal server error.")
	}

	command, err := s.GetCommand(ctx, &pb.GetCommandRequest{Name: r.GetName()})
	if err != nil {
		log.WithError(err).Println("Error retrieving command from database.")
		return nil, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error.")
	}
	// If this operation did not change any rows, there are four major possibilities:
	// - The command was not yet in READY state, in which case we indicate bad precondition.
	// - The command was already done, in which case we return the result.
	// - The command had already started running, in which case we poll for it to complete.
	if rows == 0 {
		switch command.Status {
		// If the command is not ready to run then it is a failed precondition.
		case pb.Status_UNDEFINED:
			fallthrough
		case pb.Status_SUBMITTED:
			return nil, status.Errorf(
				codes.FailedPrecondition,
				"Command is not ready to run.",
			)
		case pb.Status_DELETED:
			return nil, status.Errorf(
				codes.FailedPrecondition,
				"Command was canceled.",
			)

		// If the command has already completed then return the result.
		case pb.Status_SUCCESS:
			fallthrough
		case pb.Status_ERROR:
			return command, nil

		// If the command is already being run then wait for it to complete
		// and then return it.
		case pb.Status_RUNNING:
			return s.awaitCommand(ctx, r.GetName())
		}
	}

	// Run the command in a separate thread so that even if the request is canceled,
	// the command will continue to run.
	errChan := make(chan error)
	doneChan := make(chan struct{})
	go func() {
		argv := command.GetArgv()
		cmd := exec.Command(argv[0], argv[1:]...)

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		cmdStatus := pb.Status_SUCCESS
		err := cmd.Start()
		if err != nil {
			cmdStatus = pb.Status_ERROR
		} else {
			err = cmd.Wait()
		}

		if err != nil {
			cmdStatus = pb.Status_ERROR
		}

		endTime := time.Now()
		_, err = s.DB.Exec(
			finishCommandQuery,
			id,
			cmdStatus,
			endTime,
			stdout.Bytes(),
			stderr.Bytes(),
		)
		if err != nil {
			errChan <- err
		}
		close(doneChan)
	}()

	select {
	case <-errChan:
		return nil, status.Errorf(codes.Internal, "Error while running command.")
	case <-doneChan:
		return s.GetCommand(ctx, &pb.GetCommandRequest{Name: r.GetName()})
	case <-ctx.Done():
		return nil, status.Errorf(codes.Canceled, "Request canceled.")
	}
}

func (s *Server) awaitCommand(ctx context.Context, name string) (*pb.Command, error) {
	cmd, err := s.GetCommand(ctx, &pb.GetCommandRequest{Name: name})
	if err != nil {
		return nil, err
	}

	for cmd.Status != pb.Status_SUCCESS && cmd.Status != pb.Status_ERROR && err == nil {
		timer := time.NewTimer(time.Second)
		select {
		case <-ctx.Done():
			return nil, status.Errorf(codes.Canceled, "Request canceled.")
		case <-timer.C:
			cmd, err = s.GetCommand(ctx, &pb.GetCommandRequest{Name: name})
		}
	}

	return cmd, err
}

const createCommandQuery = `
	INSERT INTO commands ("issuer_id", "argv", "comment", "status", "create_time", "update_time") 
	VALUES ($1, $2, $3, $4, $5, $5)
	RETURNING commands.id;
`

// CreateCommand implements ToolProxy for Server.
func (s *Server) CreateCommand(ctx context.Context, r *pb.CreateCommandRequest) (*pb.Command, error) {
	issuerID := 0
	createTime := time.Now()
	cmdStatus := r.GetCommand().GetStatus()
	if cmdStatus != pb.Status_SUBMITTED && cmdStatus != pb.Status_READY {
		cmdStatus = pb.Status_SUBMITTED
	}
	row := s.DB.QueryRowContext(
		ctx,
		createCommandQuery,
		issuerID,
		pq.Array(r.GetCommand().GetArgv()),
		r.GetCommand().GetComment(),
		cmdStatus,
		createTime,
	)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		log.WithError(err).Println("Error saving command to database")
		return nil, status.Errorf(codes.Unavailable, "Internal server error")
	}

	return &pb.Command{
		Name:       fmt.Sprintf("commands/%d", id),
		Issuer:     fmt.Sprintf("users/%d", issuerID),
		Argv:       r.GetCommand().GetArgv(),
		Comment:    r.GetCommand().GetComment(),
		Status:     r.GetCommand().GetStatus(),
		CreateTime: timestamppb.New(createTime),
		UpdateTime: timestamppb.New(createTime),
	}, nil
}

const updateCommandQuery = `
	UPDATE Commands
	SET (argv, comment, status, update_time) = ($2, $3, $4, $5)
	WHERE $1 = id
	RETURNING issuer_id, status, std_out, std_err, create_time, delete_time, start_time, end_time;
`

// UpdateCommand implements ToolProxy for Server.
func (s *Server) UpdateCommand(ctx context.Context, r *pb.UpdateCommandRequest) (*pb.Command, error) {
	var id int64
	name := urn.Parse(r.GetName())
	err := name.Scan(nil, &id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Malformed command name.")
	}

	command, err := s.GetCommand(ctx, &pb.GetCommandRequest{Name: r.GetName()})
	if err != nil {
		return nil, err
	}

	updateTime := time.Now()

	mask := make(map[string]struct{})
	for _, v := range r.GetUpdateMask().GetPaths() {
		mask[v] = struct{}{}
	}

	argv := r.GetCommand().GetArgv()
	comment := r.GetCommand().GetComment()
	cmdStatus := r.GetCommand().GetStatus()

	if len(mask) > 0 {
		if _, ok := mask["argv"]; !ok {
			argv = command.GetArgv()
		}
		if _, ok := mask["comment"]; !ok {
			comment = command.GetComment()
		}
		if _, ok := mask["status"]; !ok {
			cmdStatus = command.GetStatus()
		}
	}

	row := s.DB.QueryRowContext(
		ctx,
		updateCommandQuery,
		id,
		pq.Array(argv),
		comment,
		cmdStatus,
		updateTime,
	)

	var issuerID int64
	var statusID int32
	var stdOut, stdErr []byte
	var createTime, deleteTime, startTime, endTime sql.NullTime
	err = row.Scan(
		&issuerID,
		&statusID,
		&stdOut,
		&stdErr,
		&createTime,
		&deleteTime,
		&startTime,
		&endTime,
	)
	if err == sql.ErrNoRows {
		return nil, status.Errorf(codes.NotFound, "Command not found.")
	} else if err != nil {
		return nil, status.Errorf(codes.Unavailable, "Error getting command.")
	}

	return &pb.Command{
		Name:       r.GetName(),
		Issuer:     fmt.Sprintf("users/%d", issuerID),
		Argv:       argv,
		Comment:    comment,
		Status:     pb.Status(statusID),
		StdOut:     stdOut,
		StdErr:     stdErr,
		CreateTime: timestamp(createTime),
		UpdateTime: timestamppb.New(updateTime),
		DeleteTime: timestamp(deleteTime),
		StartTime:  timestamp(startTime),
		EndTime:    timestamp(endTime),
	}, nil
}

// DeleteCommand implements ToolProxy for Server.
func (s *Server) DeleteCommand(ctx context.Context, r *pb.DeleteCommandRequest) (*pb.Command, error) {
	var id int64
	name := urn.Parse(r.GetName())
	err := name.Scan(nil, &id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Malformed command name.")
	}

	deletedTime := time.Now()
	_, err = s.DB.Exec(`
		UPDATE Commands
		SET (status, deleted_time) = ($2, $3)
		WHERE id = $1 AND status IN ($4, $5, $6)`,
		id,
		pb.Status_DELETED,
		deletedTime,
		pb.Status_UNDEFINED,
		pb.Status_SUBMITTED,
		pb.Status_READY,
	)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "Internal server error.")
	}

	// Note error checking this request handles the "command did not exist" case.
	command, err := s.GetCommand(ctx, &pb.GetCommandRequest{Name: r.GetName()})
	if err != nil {
		return nil, err
	}

	switch command.Status {
	// If the command has already been deleted we just return it.
	case pb.Status_DELETED:
		return command, nil

	// We cannot delete a command that is already running or has already run.
	case pb.Status_SUCCESS:
		fallthrough
	case pb.Status_ERROR:
		fallthrough
	case pb.Status_RUNNING:
		return nil, status.Errorf(
			codes.FailedPrecondition,
			"A command cannot be deleted after it has been started.",
		)
	}

	// This should be unreachable, because if it had any other status
	// then the delete operation should have succeeded or we should have
	// seen the error when we ran the SQL query, but we include it for
	// exhaustiveness.
	return nil, status.Errorf(codes.Internal, "Internal server error.")

}

const listCommandQuery = `
	SELECT id, issuer_id, argv, comment, status, std_out, std_err, create_time, update_time, delete_time, start_time, end_time
	FROM Commands
	LIMIT $1 OFFSET $2;
`

// ListCommands implements ToolProxy for server.
func (s *Server) ListCommands(ctx context.Context, r *pb.ListCommandsRequest) (*pb.ListCommandsResponse, error) {
	offset, err := strconv.ParseInt(r.GetPageToken(), 10, 0)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Malformed page token.")
	}

	rows, err := s.DB.QueryContext(
		ctx,
		listCommandQuery,
		r.GetPageSize(),
		offset,
	)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "Internal server error.")
	}

	var commands []*pb.Command
	for rows.Next() {
		var id, issuerID int64
		var argv []string
		var comment string
		var statusID int32
		var stdOut, stdErr []byte
		var createTime, updateTime, deleteTime, startTime, endTime sql.NullTime
		err = rows.Scan(
			&id,
			&issuerID,
			pq.Array(&argv),
			&comment,
			&statusID,
			&stdOut,
			&stdErr,
			&createTime,
			&updateTime,
			&deleteTime,
			&startTime,
			&endTime,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Internal server error.")
		}

		commands = append(commands, &pb.Command{
			Name:       fmt.Sprintf("commands/%d", id),
			Issuer:     fmt.Sprintf("users/%d", issuerID),
			Argv:       argv,
			Comment:    comment,
			Status:     pb.Status(statusID),
			StdOut:     stdOut,
			StdErr:     stdErr,
			CreateTime: timestamp(createTime),
			UpdateTime: timestamp(createTime),
			DeleteTime: timestamp(createTime),
			StartTime:  timestamp(createTime),
			EndTime:    timestamp(createTime),
		})
	}

	nextPageToken := fmt.Sprintf("%d", offset+int64(len(commands)))
	if len(commands) < int(r.GetPageSize()) {
		nextPageToken = ""
	}

	return &pb.ListCommandsResponse{
		Commands:      commands,
		NextPageToken: nextPageToken,
	}, nil
}
