package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log/slog"
	"time"

	"github.com/gowebpki/jcs"
	"github.com/swaggest/usecase/status"
	"github.com/tidwall/gjson"

	"github.com/nesfit/tenacity-chaincode/pkg/entities"
	"github.com/nesfit/tenacity-chaincode/pkg/repository"
)

type RMTUsecase struct {
	rep   repository.Repository
	piuId string
}

func NewRMTUsecase(piuId string, rep repository.Repository) *RMTUsecase {
	return &RMTUsecase{
		rep:   rep,
		piuId: piuId,
	}
}

func (u RMTUsecase) isThisPIU(piuId string) bool {
	return u.piuId == piuId
}

func (u RMTUsecase) SetPIUInfo(ctx context.Context, input entities.PIUInfo, output *entities.SetPIUInfoOutput) error {
	slog.Debug(
		"SetPIUInfo called",
		"input", input,
	)

	exists, err := u.rep.PIUExists(u.piuId)

	if err != nil {
		slog.Error(
			"Failed to get PIU from the repository",
			"id", u.piuId,
			"error", err,
		)
		return err
	}

	if !exists {
		entity := entities.NewPIUFromPIUInfo(u.piuId, input)
		err = u.rep.InsertPIU(u.piuId, entity)
	} else {
		entity, err := u.rep.GetPIU(u.piuId)

		if err != nil {
			slog.Error(
				"Failed to get PIU from the repository",
				"id", u.piuId,
				"error", err,
			)
			return err
		}

		entity = entities.UpdatePIUFromPIUInfo(entity, input)
		err = u.rep.UpdatePIU(u.piuId, entity)
	}

	if err != nil {
		slog.Error(
			"Failed writing PIU information to repository",
			"id", u.piuId,
			"error", err,
		)
		return err
	}

	*output = entities.SetPIUInfoOutput{}

	slog.Debug(
		"SetPIUInfo finished",
		"output", output,
	)

	return nil
}

func (u RMTUsecase) GetPIUs(ctx context.Context, input entities.GetPIUsInput, output *[]entities.PIU) error {
	slog.Debug(
		"GetPIUs called",
		"input", input,
	)

	out, err := u.rep.GetPIUs()

	if err != nil {
		slog.Error(
			"Failed to get PIUs from the repository",
			"error", err,
		)
		return err
	}

	*output = out

	slog.Debug(
		"GetPIUs finished",
		"output", output,
	)

	return nil
}

func (u RMTUsecase) GetPNRs(ctx context.Context, input entities.PNRFilter, output *[]entities.PNR) error {
	slog.Debug(
		"GetPNRs called",
		"input", input,
	)

	out, err := u.rep.GetPNRs(input)

	if err != nil {
		slog.Error(
			"Failed to get PNRs from the repository",
			"error", err,
		)

		return status.Wrap(err, status.InvalidArgument)
	}

	*output = out

	slog.Debug(
		"GetPNRs finished",
		"output", output,
	)

	return nil
}

func (u RMTUsecase) NewPNRRequest(ctx context.Context, input entities.NewPNRRequestInput, output *entities.NewPNRRequestOutput) error {
	slog.Debug(
		"NewPNRRequest called",
		"input", input,
	)

	if u.isThisPIU(input.RespondingPIU) {
		err := errors.New("Cannot request data from itself")
		slog.Error(
			err.Error(),
			"clientId", u.piuId,
			"respondingPIU", input.RespondingPIU,
		)
		return status.Wrap(err, status.InvalidArgument)
	}

	_, err := u.rep.GetPIU(input.RespondingPIU)

	if err != nil {
		slog.Error(
			"Could not get information about responding PIU",
			"error", err,
		)
		return status.Wrap(err, status.InvalidArgument)
	}

	pnr := entities.PNR{
		Id:               input.Id,
		RequestingPIU:    u.piuId,
		RespondingPIU:    input.RespondingPIU,
		RequestTimestamp: input.RequestTimestamp,
		State:            entities.RequestStatePending,
		RequestData:      entities.OptionalMessage(input.RequestData),
		PNRHashes:        []string{},
	}

	err = u.rep.InsertPNR(input.Id, pnr)

	if err != nil {
		slog.Error(
			"Could not insert new PNR",
			"error", err,
		)
		return status.Wrap(err, status.Internal)
	}

	gc := entities.GCMetadata{Id: pnr.Id, CreationTimestamp: pnr.RequestTimestamp}

	err = u.rep.InsertGCMetadata(pnr, gc)

	if err != nil {
		slog.Error(
			"Could not insert new PNR GC metadata",
			"error", err,
		)
		return status.Wrap(err, status.Internal)
	}

	*output = entities.NewPNRRequestOutput{Id: input.Id}

	slog.Debug(
		"NewPNRRequest finished",
		"output", output,
	)

	return nil
}

func (u RMTUsecase) submitPNRResponse(response entities.RequestState, ctx context.Context, input entities.SubmitPNRResponseInput, output *entities.SubmitPNRResponseOutput) error {
	slog.Debug(
		"SubmitPNRResponse called",
		"input", input,
	)

	pnr, err := u.rep.GetPNR(input.Id)
	if err != nil {
		slog.Error(
			"Could not get PNR request",
			"id", input.Id,
			"error", err,
		)
		return status.Wrap(err, status.InvalidArgument)
	}

	if !u.isThisPIU(pnr.RespondingPIU) {
		err := errors.New("Not allowed to respond to this request")
		slog.Error(
			err.Error(),
			"clientId", u.piuId,
			"respondingPIU", pnr.RespondingPIU,
		)
		return status.Wrap(err, status.InvalidArgument)
	}

	if pnr.State != entities.RequestStatePendingConfirmed {
		err := errors.New("PNR request must be in PendingConfirmed state")
		slog.Error(
			err.Error(),
			"id", input.Id,
			"state", pnr.State,
		)
		return status.Wrap(err, status.InvalidArgument)
	}

	pnr.ResponseTimestamp = input.ResponseTimestamp
	pnr.State = response
	pnr.ResponseData = entities.OptionalMessage(input.ResponseData)
	pnr.PNRHashes = []string{}

	const pnrJSONKey = "passengerDatasets.#.passenger_obj"

	gc, err := u.rep.GetGCMetadata(input.Id)

	if err != nil {
		slog.Error(
			"Could not get PNR GC metadata",
			"id", input.Id,
			"error", err,
		)
		return status.Wrap(err, status.Internal)
	}

	records := gjson.Get(pnr.ResponseData, pnrJSONKey)
	for _, record := range records.Array() {
		canonical, err := jcs.Transform([]byte(record.Raw))
		if err != nil {
			slog.Error(
				"Failed to transform PNR response record to canonical form",
				"id", input.Id,
				"error", err,
			)
			return status.Wrap(err, status.InvalidArgument)
		}

		sum := sha256.Sum256(canonical)
		pnr.PNRHashes = append(pnr.PNRHashes, hex.EncodeToString(sum[:]))

		const creationTimeJSONKey = "pnr_obj.iata_pnrgov_notif_rq_obj.created_on"

		creationTimeString := record.Get(creationTimeJSONKey)

		if creationTimeString.Exists() {
			creationTimestamp, err := time.Parse(time.RFC3339, creationTimeString.Str)
			if err != nil {
				slog.Error(
					"Could not parse creation timestamp",
					"input", creationTimeString.Str,
					"error", err,
				)
				return status.Wrap(err, status.InvalidArgument)
			}

			if gc.CreationTimestamp.After(creationTimestamp) {
				gc.CreationTimestamp = creationTimestamp
			}
		}
	}

	err = u.rep.UpdatePNR(input.Id, pnr)
	if err != nil {
		slog.Error(
			"Could not update PNR request",
			"id", input.Id,
			"error", err,
		)
		return status.Wrap(err, status.Internal)
	}

	err = u.rep.UpdateGCMetadata(pnr, gc)
	if err != nil {
		slog.Error(
			"Could not update PNR GC metadata",
			"id", input.Id,
			"error", err,
		)
		return status.Wrap(err, status.Internal)
	}

	slog.Debug(
		"SubmitPNRResponse finished",
		"output", output,
	)

	return nil
}

func (u RMTUsecase) SubmitPNRResponseAck(ctx context.Context, input entities.SubmitPNRResponseInput, output *entities.SubmitPNRResponseOutput) error {
	return u.submitPNRResponse(entities.RequestStateAck, ctx, input, output)
}

func (u RMTUsecase) SubmitPNRResponseNack(ctx context.Context, input entities.SubmitPNRResponseInput, output *entities.SubmitPNRResponseOutput) error {
	return u.submitPNRResponse(entities.RequestStateNack, ctx, input, output)
}

func (u RMTUsecase) ConfirmPNR(ctx context.Context, input entities.ConfirmPNRInput, output *entities.ConfirmPNROutput) error {
	slog.Debug(
		"ConfirmPNR called",
		"input", input,
	)

	pnr, err := u.rep.GetPNR(input.Id)
	if err != nil {
		slog.Error(
			"Could not get PNR request",
			"id", input.Id,
			"error", err,
		)
		return status.Wrap(err, status.InvalidArgument)
	}

	if !u.isThisPIU(pnr.RequestingPIU) && !u.isThisPIU(pnr.RespondingPIU) {
		err := errors.New("Not the requester or responder of this PNR request")
		slog.Error(
			err.Error(),
			"clientId", u.piuId,
			"requestingPIU", pnr.RequestingPIU,
			"respondingPIU", pnr.RespondingPIU,
		)
		return status.Wrap(err, status.InvalidArgument)
	}

	switch pnr.State {
	case entities.RequestStatePendingConfirmed, entities.RequestStateAckConfirmed, entities.RequestStateNackConfirmed:
		err := errors.New("PNR request already confirmed")
		slog.Error(
			err.Error(),
			"id", input.Id,
		)
		return status.Wrap(err, status.InvalidArgument)

	case entities.RequestStateTerminated:
		err := errors.New("Cannot confirm PNR request which has been terminated")
		slog.Error(
			err.Error(),
			"id", input.Id,
		)
		return status.Wrap(err, status.InvalidArgument)

	case entities.RequestStatePending:
		if !u.isThisPIU(pnr.RespondingPIU) {
			err := errors.New("Cannot confirm request in this state")
			slog.Error(
				err.Error(),
				"state", pnr.State,
			)
			return status.Wrap(err, status.InvalidArgument)
		}

	case entities.RequestStateAck, entities.RequestStateNack:
		if !u.isThisPIU(pnr.RequestingPIU) {
			err := errors.New("Cannot confirm request in this state")
			slog.Error(
				err.Error(),
				"state", pnr.State,
			)
			return status.Wrap(err, status.InvalidArgument)
		}
	}

	pnr.State = entities.GetConfirmedState(pnr.State)

	err = u.rep.UpdatePNR(input.Id, pnr)
	if err != nil {
		slog.Error(
			"Could not update PNR request",
			"id", input.Id,
			"error", err,
		)
		return status.Wrap(err, status.Internal)
	}

	if pnr.State == entities.RequestStateAckConfirmed || pnr.State == entities.RequestStateNackConfirmed {
		err := u.rep.PurgePNRData(input.Id)
		if err != nil {
			slog.Error(
				"Could not purge PNR data",
				"id", input.Id,
				"error", err,
			)
			return status.Wrap(err, status.Internal)
		}
	}

	slog.Debug(
		"ConfirmPNR finished",
		"output", output,
	)

	return nil
}

func (u RMTUsecase) TerminatePNRRequest(ctx context.Context, input entities.TerminatePNRRequestInput, output *entities.TerminatePNRRequestOutput) error {
	slog.Debug(
		"TerminatePNRRequest called",
		"input", input,
	)

	pnr, err := u.rep.GetPNR(input.Id)
	if err != nil {
		slog.Error(
			"Could not get PNR request",
			"id", input.Id,
			"error", err,
		)
		return status.Wrap(err, status.InvalidArgument)
	}

	if !u.isThisPIU(pnr.RequestingPIU) && !u.isThisPIU(pnr.RespondingPIU) {
		err := errors.New("Not the requester or responder of this PNR request")
		slog.Error(
			err.Error(),
			"clientId", u.piuId,
			"requestingPIU", pnr.RequestingPIU,
			"respondingPIU", pnr.RespondingPIU,
		)
		return status.Wrap(err, status.InvalidArgument)
	}

	pnr.State = entities.RequestStateTerminated
	err = u.rep.UpdateLocalPNR(input.Id, pnr)

	if err != nil {
		slog.Error(
			"Could not update PNR",
			"error", err,
		)
		return status.Wrap(err, status.Internal)
	}

	err = u.rep.PurgeLocalPNRData(input.Id)

	if err != nil {
		slog.Error(
			"Could not purge PNR data",
			"error", err,
		)
		return status.Wrap(err, status.Internal)
	}

	err = u.rep.DeleteLocalGCMetadata(input.Id)

	if err != nil {
		slog.Error(
			"Could not delete GC metadata",
			"error", err,
		)
		return status.Wrap(err, status.Internal)
	}

	*output = entities.TerminatePNRRequestOutput{}

	slog.Debug(
		"TerminatePNRRequest finished",
		"output", output,
	)

	return nil
}
