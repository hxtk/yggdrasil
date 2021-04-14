package rpc

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/alessio/shellescape"
	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/hxtk/yggdrasil/toolproxy/v1"
)

type Client struct {
	tp pb.ToolProxyClient
}

func New(addr string, tlsConfig *tls.Config) *Client {
	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(
			credentials.NewTLS(tlsConfig),
		),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor()),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor()),
	)
	if err != nil {
		panic(err)
	}

	return &Client{
		tp: pb.NewToolProxyClient(conn),
	}
}

func (c *Client) Get(ctx context.Context, name string) {
	cmd, err := c.tp.GetCommand(ctx, &pb.GetCommandRequest{Name: name})
	if err != nil {
		fmt.Println("Could not get command:", err)
		return
	}

	fmt.Println(shellescape.QuoteCommand(cmd.GetArgv()))
	fmt.Println()
	fmt.Println("Submitted:", cmd.GetCreateTime().AsTime())
	fmt.Println("Started:  ", cmd.GetStartTime().AsTime())
	fmt.Println("Completed:", cmd.GetEndTime().AsTime())
	fmt.Println("Status:", cmd.GetStatus().String())
	fmt.Println("Issuer:", cmd.GetIssuer())
}

func (c *Client) Run(ctx context.Context, argv []string) {
	cmd, err := c.tp.CreateCommand(ctx,
		&pb.CreateCommandRequest{
			Command: &pb.Command{
				Argv:        argv,
				Description: "",
				Status:      pb.Status_READY,
			},
		},
	)
	if err != nil {
		fmt.Println("Failed to create command:", err)
		return
	}

	cmd, err = c.tp.RunCommand(ctx, &pb.RunCommandRequest{Name: cmd.GetName()})
	if err != nil {
		fmt.Println("Error running command: ", err)
	}

	fmt.Print(string(cmd.GetStdOut()))
}
