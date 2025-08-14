package usecase

import (
	"context"
	"time"

	"github.com/nesfit/tenacity-chaincode/pkg/entities"
)

type Clock interface {
	Now() time.Time
}

type PNRExchangeUsecase interface {
	SetPIUInfo(ctx context.Context, input entities.PIUInfo, output *entities.SetPIUInfoOutput) error
	GetPIUs(ctx context.Context, input entities.GetPIUsInput, output *[]entities.PIU) error
	GetPNRs(ctx context.Context, input entities.PNRFilter, output *[]entities.PNR) error
	NewPNRRequest(ctx context.Context, input entities.NewPNRRequestInput, output *entities.NewPNRRequestOutput) error
	SubmitPNRResponseAck(ctx context.Context, input entities.SubmitPNRResponseInput, output *entities.SubmitPNRResponseOutput) error
	SubmitPNRResponseNack(ctx context.Context, input entities.SubmitPNRResponseInput, output *entities.SubmitPNRResponseOutput) error
	ConfirmPNR(ctx context.Context, input entities.ConfirmPNRInput, output *entities.ConfirmPNROutput) error
	TerminatePNRRequest(ctx context.Context, input entities.TerminatePNRRequestInput, output *entities.TerminatePNRRequestOutput) error
}
