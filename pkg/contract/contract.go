package contract

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"

	"github.com/nesfit/tenacity-chaincode/pkg/entities"
	"github.com/nesfit/tenacity-chaincode/pkg/repository/privatedata"
	"github.com/nesfit/tenacity-chaincode/pkg/usecase"
)

type UsecaseFactory interface {
	New(ctx contractapi.TransactionContextInterface) (usecase.PNRExchangeUsecase, error)
}

type LedgerUsecaseFactory struct {
}

func (uf *LedgerUsecaseFactory) New(ctx contractapi.TransactionContextInterface) (usecase.PNRExchangeUsecase, error) {
	piuId, err := GetClientOrgId(ctx)

	if err != nil {
		return nil, err
	}

	r := privatedata.NewPrivateDataRepository(ctx, piuId)

	u := usecase.NewRMTUsecase(piuId, r)

	return u, nil
}

type SmartContract struct {
	contractapi.Contract
	uf UsecaseFactory
}

func NewSmartContract(uf UsecaseFactory) SmartContract {
	return SmartContract{uf: uf}
}

func GetClientOrgId(ctx contractapi.TransactionContextInterface) (string, error) {
	clientOrgId, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		slog.Error(
			"failed geting client's msp Id",
			"error", err,
		)
		return "", fmt.Errorf("failed getting client's msp Id: %v", err)
	}

	return clientOrgId, nil
}

func (s *SmartContract) SetPIUInfo(ctx contractapi.TransactionContextInterface, info string) error {
	var input entities.PIUInfo
	var output entities.SetPIUInfoOutput

	u, err := s.uf.New(ctx)

	if err != nil {
		slog.Error(
			"failed to create usecase",
			"error", err,
		)
		return err
	}

	err = json.Unmarshal([]byte(info), &input)
	if err != nil {
		slog.Error(
			"failed to unmarshal input",
			"input", info,
			"error", err,
		)
		return err
	}

	err = u.SetPIUInfo(context.TODO(), input, &output)

	return err
}

func (s *SmartContract) GetPIUs(ctx contractapi.TransactionContextInterface) ([]entities.PIU, error) {
	var input entities.GetPIUsInput
	var output []entities.PIU

	u, err := s.uf.New(ctx)

	if err != nil {
		slog.Error(
			"failed to create usecase",
			"error", err,
		)
		return output, err
	}

	err = u.GetPIUs(context.TODO(), input, &output)

	return output, err
}

func (s *SmartContract) GetPNRs(ctx contractapi.TransactionContextInterface, filter string) ([]entities.PNR, error) {
	var input entities.PNRFilter
	var output []entities.PNR

	u, err := s.uf.New(ctx)

	if err != nil {
		slog.Error(
			"failed to create usecase",
			"error", err,
		)
		return output, err
	}

	err = json.Unmarshal([]byte(filter), &input)
	if err != nil {
		slog.Error(
			"failed to unmarshal input",
			"input", filter,
			"error", err,
		)
		return output, err
	}

	err = u.GetPNRs(context.TODO(), input, &output)

	return output, err
}

func (s *SmartContract) NewPNRRequest(ctx contractapi.TransactionContextInterface, request string) (entities.NewPNRRequestOutput, error) {
	var input entities.NewPNRRequestInput
	var output entities.NewPNRRequestOutput

	u, err := s.uf.New(ctx)

	if err != nil {
		slog.Error(
			"failed to create usecase",
			"error", err,
		)
		return output, err
	}

	err = json.Unmarshal([]byte(request), &input)
	if err != nil {
		slog.Error(
			"failed to unmarshal input",
			"input", request,
			"error", err,
		)
		return output, err
	}

	transient, err := ctx.GetStub().GetTransient()
	if err != nil {
		slog.Error(
			"failed to get transient data",
			"error", err,
		)
		return output, err
	}

	requestData, ok := transient[entities.RequestDataTransientKey]
	if !ok {
		slog.Error(
			"missing transient data for key",
			"key", entities.RequestDataTransientKey,
		)
		return output, err
	}

	input.RequestData = (*json.RawMessage)(&requestData)

	err = u.NewPNRRequest(context.TODO(), input, &output)

	return output, err
}

func (s *SmartContract) SubmitPNRResponseAck(ctx contractapi.TransactionContextInterface, response string) error {
	var input entities.SubmitPNRResponseInput
	var output entities.SubmitPNRResponseOutput

	u, err := s.uf.New(ctx)

	if err != nil {
		slog.Error(
			"failed to create usecase",
			"error", err,
		)
		return err
	}

	err = json.Unmarshal([]byte(response), &input)
	if err != nil {
		slog.Error(
			"failed to unmarshal input",
			"input", response,
			"error", err,
		)
		return err
	}

	transient, err := ctx.GetStub().GetTransient()
	if err != nil {
		slog.Error(
			"failed to get transient data",
			"error", err,
		)
		return err
	}

	responseData, ok := transient[entities.ResponseDataTransientKey]
	if !ok {
		slog.Error(
			"missing transient data for key",
			"key", entities.ResponseDataTransientKey,
			"data", transient,
		)
		return err
	}

	input.ResponseData = (*json.RawMessage)(&responseData)

	err = u.SubmitPNRResponseAck(context.TODO(), input, &output)

	return err
}

func (s *SmartContract) SubmitPNRResponseNack(ctx contractapi.TransactionContextInterface, response string) error {
	var input entities.SubmitPNRResponseInput
	var output entities.SubmitPNRResponseOutput

	u, err := s.uf.New(ctx)

	if err != nil {
		slog.Error(
			"failed to create usecase",
			"error", err,
		)
		return err
	}

	err = json.Unmarshal([]byte(response), &input)
	if err != nil {
		slog.Error(
			"failed to unmarshal input",
			"input", response,
			"error", err,
		)
		return err
	}

	transient, err := ctx.GetStub().GetTransient()
	if err != nil {
		slog.Error(
			"failed to get transient data",
			"error", err,
		)
		return err
	}

	responseData, ok := transient[entities.ResponseDataTransientKey]
	if !ok {
		slog.Error(
			"missing transient data for key",
			"key", entities.ResponseDataTransientKey,
		)
		return err
	}

	input.ResponseData = (*json.RawMessage)(&responseData)

	err = u.SubmitPNRResponseNack(context.TODO(), input, &output)

	return err
}

func (s *SmartContract) ConfirmPNR(ctx contractapi.TransactionContextInterface, confirmation string) error {
	var input entities.ConfirmPNRInput
	var output entities.ConfirmPNROutput

	u, err := s.uf.New(ctx)

	if err != nil {
		slog.Error(
			"failed to create usecase",
			"error", err,
		)
		return err
	}

	err = json.Unmarshal([]byte(confirmation), &input)
	if err != nil {
		slog.Error(
			"failed to unmarshal input",
			"input", confirmation,
			"error", err,
		)
		return err
	}

	err = u.ConfirmPNR(context.TODO(), input, &output)

	return err
}

func (s *SmartContract) TerminatePNRRequest(ctx contractapi.TransactionContextInterface, purge string) error {
	var input entities.TerminatePNRRequestInput
	var output entities.TerminatePNRRequestOutput

	u, err := s.uf.New(ctx)

	if err != nil {
		slog.Error(
			"failed to create usecase",
			"error", err,
		)
		return err
	}

	err = json.Unmarshal([]byte(purge), &input)
	if err != nil {
		slog.Error(
			"failed to unmarshal input",
			"input", purge,
			"error", err,
		)
		return err
	}

	err = u.TerminatePNRRequest(context.TODO(), input, &output)

	return err
}
