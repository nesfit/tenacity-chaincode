package publicledger

import (
	"errors"
	"log/slog"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"

	"github.com/nesfit/tenacity-chaincode/pkg/entities"
)

type PublicLedgerRepository struct {
	ctx contractapi.TransactionContextInterface
}

func NewPublicLedgerRepository(ctx contractapi.TransactionContextInterface) *PublicLedgerRepository {
	return &PublicLedgerRepository{
		ctx: ctx,
	}
}

func (r *PublicLedgerRepository) PIUExists(id string) (bool, error) {
	key, err := getPIUCompositeKey(id)

	if err != nil {
		slog.Error(
			"could not create PIU composite key",
			"id", id,
			"error", err,
		)
		return false, err
	}

	piuModel, err := r.ctx.GetStub().GetState(key)

	if err != nil {
		slog.Error(
			"could not get PIU",
			"id", id,
			"error", err,
		)
		return false, err
	}

	exists := piuModel != nil

	return exists, nil
}

func (r *PublicLedgerRepository) GetPIU(id string) (entities.PIU, error) {
	key, err := getPIUCompositeKey(id)

	if err != nil {
		slog.Error(
			"could not create PIU composite key",
			"id", id,
			"error", err,
		)
		return entities.PIU{}, err
	}

	piuModel, err := r.ctx.GetStub().GetState(key)

	if err != nil {
		slog.Error(
			"could not get PIU",
			"id", id,
			"error", err,
		)
		return entities.PIU{}, err
	}

	exists := piuModel != nil

	if !exists {
		err = errors.New("PIU not found")
		slog.Error(
			err.Error(),
			"id", id,
		)
		return entities.PIU{}, err
	}

	return piuModelToEntity(piuModel)
}

func (r *PublicLedgerRepository) GetPIUs() ([]entities.PIU, error) {
	var result []entities.PIU

	iterator, err := r.ctx.GetStub().GetStateByPartialCompositeKey(piuObjectType, []string{})
	if err != nil {
		slog.Error(
			err.Error(),
		)
		return []entities.PIU{}, err
	}
	defer iterator.Close()

	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		piu, err := piuModelToEntity(queryResponse.Value)
		if err != nil {
			return nil, err
		}
		result = append(result, piu)
	}

	return result, nil
}

func (r *PublicLedgerRepository) InsertPIU(id string, piu entities.PIU) error {
	exists, _ := r.PIUExists(id)

	if exists {
		return errors.New("PIU already exists")
	}

	key, err := getPIUCompositeKey(id)

	if err != nil {
		slog.Error(
			"could not create PIU composite key",
			"id", id,
			"error", err,
		)
		return err
	}

	piuModel, err := piuEntityToModel(piu)

	if err != nil {
		slog.Error(
			"could not map PIU entity to model",
			"id", id,
			"error", err,
		)
		return err
	}

	err = r.ctx.GetStub().PutState(key, piuModel)

	if err != nil {
		slog.Error(
			"could not put model into ledger",
			"id", id,
			"error", err,
		)
		return err
	}

	return nil
}

func (r *PublicLedgerRepository) UpdatePIU(id string, piu entities.PIU) error {
	exists, _ := r.PIUExists(id)

	if !exists {
		return errors.New("PIU does not exist")
	}

	key, err := getPIUCompositeKey(id)

	if err != nil {
		slog.Error(
			"could not create PIU composite key",
			"id", id,
			"error", err,
		)
		return err
	}

	piuModel, err := piuEntityToModel(piu)

	if err != nil {
		slog.Error(
			"could not map PIU entity to model",
			"id", id,
			"error", err,
		)
		return err
	}

	err = r.ctx.GetStub().PutState(key, piuModel)

	if err != nil {
		slog.Error(
			"could not put model into ledger",
			"id", id,
			"error", err,
		)
		return err
	}

	return nil
}

func (r *PublicLedgerRepository) PNRExists(id string) (bool, error) {
	key, err := getPNRCompositeKey(id)

	if err != nil {
		slog.Error(
			"could not create PNR composite key",
			"id", id,
			"error", err,
		)
		return false, err
	}

	pnrModel, err := r.ctx.GetStub().GetState(key)

	if err != nil {
		slog.Error(
			"could not get PNR",
			"id", id,
			"error", err,
		)
		return false, err
	}

	exists := pnrModel != nil

	return exists, nil
}

func (r *PublicLedgerRepository) GetPNR(id string) (entities.PNR, error) {
	key, err := getPNRCompositeKey(id)

	if err != nil {
		slog.Error(
			"could not create PNR composite key",
			"id", id,
			"error", err,
		)
		return entities.PNR{}, err
	}

	pnrModel, err := r.ctx.GetStub().GetState(key)

	if err != nil {
		slog.Error(
			"could not get PNR",
			"id", id,
			"error", err,
		)
		return entities.PNR{}, err
	}

	exists := pnrModel != nil

	if !exists {
		err = errors.New("PNR not found")
		slog.Error(
			err.Error(),
			"id", id,
		)
		return entities.PNR{}, err
	}

	return pnrModelToEntity(pnrModel)
}

func (r *PublicLedgerRepository) GetPNRs(filter entities.PNRFilter) ([]entities.PNR, error) {
	var result []entities.PNR

	iterator, err := r.ctx.GetStub().GetStateByPartialCompositeKey(pnrObjectType, []string{})
	if err != nil {
		slog.Error(
			err.Error(),
		)
		return []entities.PNR{}, err
	}
	defer iterator.Close()

	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		pnr, err := pnrModelToEntity(queryResponse.Value)
		if err != nil {
			return nil, err
		}
		if entities.IsMatchingPNR(filter, pnr) {
			result = append(result, pnr)
		}
	}

	return result, nil
}

func (r *PublicLedgerRepository) InsertPNR(id string, pnr entities.PNR) error {
	exists, _ := r.PNRExists(id)

	if exists {
		return errors.New("PNR already exists")
	}

	key, err := getPNRCompositeKey(id)

	if err != nil {
		slog.Error(
			"could not create PNR composite key",
			"id", id,
			"error", err,
		)
		return err
	}

	pnrModel, err := pnrEntityToModel(pnr)

	if err != nil {
		slog.Error(
			"could not map PNR entity to model",
			"id", id,
			"error", err,
		)
		return err
	}

	err = r.ctx.GetStub().PutState(key, pnrModel)

	if err != nil {
		slog.Error(
			"could not put model into ledger",
			"id", id,
			"error", err,
		)
		return err
	}

	return nil
}

func (r *PublicLedgerRepository) UpdatePNR(id string, pnr entities.PNR) error {
	exists, _ := r.PNRExists(id)

	if !exists {
		return errors.New("PNR does not exist")
	}

	key, err := getPNRCompositeKey(id)

	if err != nil {
		slog.Error(
			"could not create PNR composite key",
			"id", id,
			"error", err,
		)
		return err
	}

	pnrModel, err := pnrEntityToModel(pnr)

	if err != nil {
		slog.Error(
			"could not map PNR entity to model",
			"id", id,
			"error", err,
		)
		return err
	}

	err = r.ctx.GetStub().PutState(key, pnrModel)

	if err != nil {
		slog.Error(
			"could not put model into ledger",
			"id", id,
			"error", err,
		)
		return err
	}

	return nil
}

func (r *PublicLedgerRepository) UpdateLocalPNR(id string, pnr entities.PNR) error {
	return r.UpdatePNR(id, pnr)
}

func (r *PublicLedgerRepository) PurgePNRData(id string) error {
	pnr, err := r.GetPNR(id)

	if err != nil {
		return err
	}

	pnr.RequestData = ""
	pnr.ResponseData = ""

	r.UpdatePNR(id, pnr)

	return nil
}

func (r *PublicLedgerRepository) PurgeLocalPNRData(id string) error {
	return r.PurgePNRData(id)
}

func (r *PublicLedgerRepository) GCMetadataExists(id string) (bool, error) {
	key, err := getGCMetatadaCompositeKey(id)

	if err != nil {
		slog.Error(
			"could not create GC metadata composite key",
			"id", id,
			"error", err,
		)
		return false, err
	}

	gcMetadataModel, err := r.ctx.GetStub().GetState(key)

	if err != nil {
		slog.Error(
			"could not get GC metadata",
			"id", id,
			"error", err,
		)
		return false, err
	}

	exists := gcMetadataModel != nil

	return exists, nil
}

func (r *PublicLedgerRepository) InsertGCMetadata(pnr entities.PNR, gc entities.GCMetadata) error {
	exists, _ := r.GCMetadataExists(pnr.Id)

	if exists {
		return errors.New("GC metadata already exists")
	}

	key, err := getGCMetatadaCompositeKey(pnr.Id)

	if err != nil {
		slog.Error(
			"could not create PNR composite key",
			"id", pnr.Id,
			"error", err,
		)
		return err
	}

	gcMetadataModel, err := gcMetadataEntityToModel(gc)

	if err != nil {
		slog.Error(
			"could not map GC metadata entity to model",
			"id", pnr.Id,
			"error", err,
		)
		return err
	}

	err = r.ctx.GetStub().PutState(key, gcMetadataModel)

	if err != nil {
		slog.Error(
			"could not put model into ledger",
			"id", pnr.Id,
			"error", err,
		)
		return err
	}

	return nil
}

func (r *PublicLedgerRepository) UpdateGCMetadata(pnr entities.PNR, gc entities.GCMetadata) error {
	exists, _ := r.GCMetadataExists(pnr.Id)

	if !exists {
		return errors.New("GC metadata does not exist")
	}

	key, err := getGCMetatadaCompositeKey(pnr.Id)

	if err != nil {
		slog.Error(
			"could not create GC metadata composite key",
			"id", pnr.Id,
			"error", err,
		)
		return err
	}

	gcMetadataModel, err := gcMetadataEntityToModel(gc)

	if err != nil {
		slog.Error(
			"could not map GC metadata entity to model",
			"id", pnr.Id,
			"error", err,
		)
		return err
	}

	err = r.ctx.GetStub().PutState(key, gcMetadataModel)

	if err != nil {
		slog.Error(
			"could not put model into ledger",
			"id", pnr.Id,
			"error", err,
		)
		return err
	}

	return nil
}

func (r *PublicLedgerRepository) DeleteGCMetadata(pnr entities.PNR) error {
	return r.DeleteLocalGCMetadata(pnr.Id)
}

func (r *PublicLedgerRepository) DeleteLocalGCMetadata(id string) error {
	exists, _ := r.GCMetadataExists(id)

	if !exists {
		return errors.New("GC metadata does not exist")
	}

	key, err := getGCMetatadaCompositeKey(id)

	if err != nil {
		slog.Error(
			"could not create GC metadata composite key",
			"id", id,
			"error", err,
		)
		return err
	}

	err = r.ctx.GetStub().DelState(key)

	if err != nil {
		slog.Error(
			"could not delete GC metadata from ledger",
			"id", id,
			"error", err,
		)
		return err
	}

	return nil
}

func (r *PublicLedgerRepository) GetGCMetadata(id string) (entities.GCMetadata, error) {
	key, err := getGCMetatadaCompositeKey(id)

	if err != nil {
		slog.Error(
			"could not create GC metadata composite key",
			"id", id,
			"error", err,
		)
		return entities.GCMetadata{}, err
	}

	gcMetadataModel, err := r.ctx.GetStub().GetState(key)

	if err != nil {
		slog.Error(
			"could not get GC metadata",
			"id", id,
			"error", err,
		)
		return entities.GCMetadata{}, err
	}

	exists := gcMetadataModel != nil

	if !exists {
		err = errors.New("GC metadata not found")
		slog.Error(
			err.Error(),
			"id", id,
		)
		return entities.GCMetadata{}, err
	}

	return gcMetadataModelToEntity(gcMetadataModel)
}

func (r *PublicLedgerRepository) GetGCMetadatas() ([]entities.GCMetadata, error) {
	var result []entities.GCMetadata

	iterator, err := r.ctx.GetStub().GetStateByPartialCompositeKey(gcMetadataObjectType, []string{})
	if err != nil {
		slog.Error(
			err.Error(),
		)
		return []entities.GCMetadata{}, err
	}
	defer iterator.Close()

	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		pnr, err := gcMetadataModelToEntity(queryResponse.Value)
		if err != nil {
			return nil, err
		}
		result = append(result, pnr)
	}

	return result, nil
}

func (r *PublicLedgerRepository) Close() {
}
