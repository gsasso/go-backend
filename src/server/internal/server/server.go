package server

import (
	"net"

	"github.com/google/wire"
	"google.golang.org/grpc"

	controller "github.com/gsasso/go-backend/src/server/internal/controller"
	serverApi "github.com/gsasso/go-backend/src/server/internal/generated/proto"
	"github.com/gsasso/go-backend/src/server/internal/ticker"
)

type LogisticServer struct {
	server *grpc.Server
}

var ServerProvider = wire.NewSet(wire.Struct(new(ticker.Summary)), controller.NewLogisticController, RunGRPCServer)

// TODO Major: Name of function - that is construct must represent structure, by name convention I could expect by calling  RunGRPCServer it's executing it.
func RunGRPCServer(ctlr *controller.LogisticCtlr) *LogisticServer {

	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)
	serverApi.RegisterCoopLogisticsEngineAPIServer(server, ctlr)
	return &LogisticServer{server: server}

}

func (my *LogisticServer) Start() error {
	// TODO Major: Configuration must be done before in separate method or in construction for better control.
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	// TODO my.server.Serve can return error and it must be propagated to appear level
	my.server.Serve(listener)

	return nil

}
