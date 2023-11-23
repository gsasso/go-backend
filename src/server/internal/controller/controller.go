package controller

import (
	"context"

	serverApi "github.com/gsasso/go-backend/src/server/internal/generated/proto"
	"github.com/gsasso/go-backend/src/server/internal/ticker"
)

type LogisticCtlr struct {
	serverApi.UnimplementedCoopLogisticsEngineAPIServer
	svc ticker.SummaryInt
}

func NewLogisticController(svc ticker.SummaryInt) *LogisticCtlr {
	go svc.Tick()
	return &LogisticCtlr{
		svc: svc,
	}
}

func (ctlr *LogisticCtlr) MoveUnit(ctx context.Context, req *serverApi.MoveUnitRequest) (*serverApi.DefaultResponse, error) {
	ctlr.svc.IncreaseTotalUnits()
	resp := &serverApi.DefaultResponse{}
	return resp, nil
}

func (ctlr *LogisticCtlr) UnitReachedWarehouse(ctx context.Context, req *serverApi.UnitReachedWarehouseRequest) (*serverApi.DefaultResponse, error) {
	ctlr.svc.IncreaseTotalReached()
	resp := &serverApi.DefaultResponse{}
	return resp, nil
}

func (ctlr *LogisticCtlr) GetSummary(ctx context.Context, req *serverApi.DefaultRequest) (*serverApi.SummaryResponse, error) {
	resp, err := ctlr.svc.GetSummary()
	if err != nil {
		return nil, err
	}
	return &serverApi.SummaryResponse{TotalUnits: resp.TotalUnits, TotalReached: resp.TotalReached}, nil
}
