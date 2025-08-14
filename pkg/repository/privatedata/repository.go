package privatedata

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"

	"github.com/nesfit/tenacity-chaincode/pkg/entities"
)

func getCollectionName(clientOrgID string) string {
	return fmt.Sprintf("%sCollection", clientOrgID)
}

type PrivateDataRepository struct {
	ctx       contractapi.TransactionContextInterface
	piuId     string
	localData string
}

func NewPrivateDataRepository(ctx contractapi.TransactionContextInterface, piuId string) *PrivateDataRepository {
	return &PrivateDataRepository{
		ctx:       ctx,
		piuId:     piuId,
		localData: getCollectionName(piuId),
	}
}

func (r *PrivateDataRepository) PIUExists(id string) (bool, error) {
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

func (r *PrivateDataRepository) GetPIU(id string) (entities.PIU, error) {
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

func (r *PrivateDataRepository) GetPIUs() ([]entities.PIU, error) {
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

func (r *PrivateDataRepository) InsertPIU(id string, piu entities.PIU) error {
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

func (r *PrivateDataRepository) UpdatePIU(id string, piu entities.PIU) error {
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

func (r *PrivateDataRepository) getPNRMeta(id string) (string, pnrMeta, error) {
	metaKey, err := getPNRMetaCompositeKey(id)

	if err != nil {
		slog.Error(
			"could not create PNR metadata composite key",
			"id", id,
			"error", err,
		)
		return "", pnrMeta{}, err
	}

	metaModel, err := r.ctx.GetStub().GetPrivateData(r.localData, metaKey)

	if err != nil {
		slog.Error(
			"could not get PNR metadata",
			"id", id,
			"error", err,
		)
		return "", pnrMeta{}, err
	}

	metaEntity, err := metaModelToMetaEntity(metaModel)

	if err != nil {
		slog.Error(
			"could not convert PNR metadata model to entity",
			"id", id,
			"error", err,
		)
		return "", pnrMeta{}, err
	}

	return metaKey, metaEntity, nil
}

func (r *PrivateDataRepository) getPNRData(id string) (string, pnrData, error) {
	dataKey, err := getPNRDataCompositeKey(id)

	if err != nil {
		slog.Error(
			"could not create PNR datadata composite key",
			"id", id,
			"error", err,
		)
		return "", pnrData{}, err
	}

	dataModel, err := r.ctx.GetStub().GetPrivateData(r.localData, dataKey)

	if err != nil {
		slog.Error(
			"could not get PNR datadata",
			"id", id,
			"error", err,
		)
		return "", pnrData{}, err
	}

	dataEntity, err := dataModelToDataEntity(dataModel)

	if err != nil {
		slog.Error(
			"could not convert PNR datadata model to entity",
			"id", id,
			"error", err,
		)
		return "", pnrData{}, err
	}

	return dataKey, dataEntity, nil
}

func (r *PrivateDataRepository) putToBothPrivateCollections(remoteCollection string, key string, value []byte) error {
	var err error

	err = r.ctx.GetStub().PutPrivateData(remoteCollection, key, value)

	if err != nil {
		slog.Error(
			"could not put data into remote private collection",
			"key", key,
			"error", err,
		)
		return err
	}

	err = r.ctx.GetStub().PutPrivateData(r.localData, key, value)

	if err != nil {
		slog.Error(
			"could not put data into local private collection",
			"key", key,
			"error", err,
		)
		return err
	}

	return nil
}

func (r *PrivateDataRepository) PNRExists(id string) (bool, error) {
	key, err := getPNRMetaCompositeKey(id)

	if err != nil {
		slog.Error(
			"could not create PNR composite key",
			"id", id,
			"error", err,
		)
		return false, err
	}

	pnrModel, err := r.ctx.GetStub().GetPrivateData(r.localData, key)

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

func (r *PrivateDataRepository) GetPNR(id string) (entities.PNR, error) {
	_, metaEntity, err := r.getPNRMeta(id)

	if err != nil {
		return entities.PNR{}, err
	}

	if !entities.HasData(metaEntity.State) {
		return pnrEntitiesToEntity(metaEntity, pnrData{}), nil
	}

	_, dataEntity, err := r.getPNRData(id)

	if err != nil {
		return entities.PNR{}, err
	}

	return pnrEntitiesToEntity(metaEntity, dataEntity), nil
}

func (r *PrivateDataRepository) GetPNRs(filter entities.PNRFilter) ([]entities.PNR, error) {
	var result []entities.PNR

	iterator, err := r.ctx.GetStub().GetPrivateDataByPartialCompositeKey(r.localData, pnrMetaObjectType, []string{})
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
			slog.Error(
				"failed calling iterator.Next()",
				"error", err,
			)
			continue
		}

		meta, err := metaModelToMetaEntity(queryResponse.Value)
		if err != nil {
			slog.Error(
				"failed to map PNR metadata model to entity",
				"model", queryResponse.Value,
				"error", err,
			)
			continue
		}

		pnr := pnrEntitiesToEntity(meta, pnrData{})

		if entities.IsMatchingPNR(filter, pnr) {
			id := meta.Id

			switch pnr.State {
			case entities.RequestStateAckConfirmed, entities.RequestStateNackConfirmed, entities.RequestStateTerminated:
				result = append(result, pnr)
			default:
				_, dataEntity, err := r.getPNRData(id)

				if err != nil {
					continue
				}

				pnrWithData := pnrEntitiesToEntity(meta, dataEntity)

				result = append(result, pnrWithData)
			}
		}
	}

	return result, nil
}

func (r *PrivateDataRepository) InsertPNR(id string, pnr entities.PNR) error {
	exists, _ := r.PNRExists(id)

	if exists {
		return errors.New("PNR already exists")
	}

	metaKey, err := getPNRMetaCompositeKey(id)

	if err != nil {
		slog.Error(
			"could not create PNR metadata composite key",
			"id", id,
			"error", err,
		)
		return err
	}

	dataKey, err := getPNRDataCompositeKey(id)

	if err != nil {
		slog.Error(
			"could not create PNR data composite key",
			"id", id,
			"error", err,
		)
		return err
	}

	metaModel, err := pnrEntityToMetaModel(pnr)
	if err != nil {
		slog.Error(
			"could not map PNR entity to metadata model",
			"id", id,
			"error", err,
		)
		return err
	}

	dataModel, err := pnrEntityToDataModel(pnr)
	if err != nil {
		slog.Error(
			"could not map PNR entity to data model",
			"id", id,
			"error", err,
		)
		return err
	}

	remotePIU := getRemotePIU(pnr, r.piuId)
	remoteData := getCollectionName(remotePIU)

	err = r.putToBothPrivateCollections(remoteData, metaKey, metaModel)

	if err != nil {
		return err
	}

	return r.putToBothPrivateCollections(remoteData, dataKey, dataModel)
}

func (r *PrivateDataRepository) UpdatePNR(id string, pnr entities.PNR) error {
	exists, _ := r.PNRExists(id)

	if !exists {
		return errors.New("PNR does not exist")
	}

	metaKey, metaEntity, err := r.getPNRMeta(id)

	if err != nil {
		return err
	}

	metaModel, err := pnrEntityToMetaModel(pnr)
	if err != nil {
		slog.Error(
			"could not map PNR entity to metadata model",
			"id", id,
			"error", err,
		)
		return err
	}

	remotePIU := getRemotePIU(pnr, r.piuId)
	remoteData := getCollectionName(remotePIU)

	err = r.putToBothPrivateCollections(remoteData, metaKey, metaModel)

	if err != nil {
		return err
	}

	if !entities.HasData(metaEntity.State) {
		return nil
	}

	dataKey, _, err := r.getPNRData(id)

	if err != nil {
		return err
	}

	dataModel, err := pnrEntityToDataModel(pnr)
	if err != nil {
		slog.Error(
			"could not map PNR entity to data model",
			"id", id,
			"error", err,
		)
		return err
	}

	return r.putToBothPrivateCollections(remoteData, dataKey, dataModel)
}

func (r *PrivateDataRepository) UpdateLocalPNR(id string, pnr entities.PNR) error {
	exists, _ := r.PNRExists(id)

	if !exists {
		return errors.New("PNR does not exist")
	}

	metaKey, metaEntity, err := r.getPNRMeta(id)

	if err != nil {
		return err
	}

	metaModel, err := pnrEntityToMetaModel(pnr)
	if err != nil {
		slog.Error(
			"could not map PNR entity to metadata model",
			"id", id,
			"error", err,
		)
		return err
	}

	err = r.ctx.GetStub().PutPrivateData(r.localData, metaKey, metaModel)

	if err != nil {
		return err
	}

	if !entities.HasData(metaEntity.State) {
		return nil
	}

	dataKey, _, err := r.getPNRData(id)

	if err != nil {
		return err
	}

	dataModel, err := pnrEntityToDataModel(pnr)
	if err != nil {
		slog.Error(
			"could not map PNR entity to data model",
			"id", id,
			"error", err,
		)
		return err
	}

	return r.ctx.GetStub().PutPrivateData(r.localData, dataKey, dataModel)
}

func (r *PrivateDataRepository) PurgePNRData(id string) error {
	exists, _ := r.PNRExists(id)

	if !exists {
		return errors.New("PNR does not exist")
	}

	_, metaEntity, err := r.getPNRMeta(id)

	if err != nil {
		return err
	}

	dataKey, err := getPNRDataCompositeKey(id)

	if err != nil {
		slog.Error(
			"could not create PNR data composite key",
			"id", id,
			"error", err,
		)
		return err
	}

	pnr := pnrEntitiesToEntity(metaEntity, pnrData{})

	remotePIU := getRemotePIU(pnr, r.piuId)
	remoteData := getCollectionName(remotePIU)

	err = r.ctx.GetStub().PurgePrivateData(remoteData, dataKey)

	if err != nil {
		slog.Error(
			"could not purge PNR data from remote collection",
			"id", id,
			"error", err,
		)
		return err
	}

	err = r.ctx.GetStub().PurgePrivateData(r.localData, dataKey)

	if err != nil {
		slog.Error(
			"could not purge PNR data from local collection",
			"id", id,
			"error", err,
		)
		return err
	}

	return nil
}

func (r *PrivateDataRepository) PurgeLocalPNRData(id string) error {
	exists, _ := r.PNRExists(id)

	if !exists {
		return errors.New("PNR does not exist")
	}

	dataKey, err := getPNRDataCompositeKey(id)

	if err != nil {
		slog.Error(
			"could not create PNR data composite key",
			"id", id,
			"error", err,
		)
		return err
	}

	err = r.ctx.GetStub().PurgePrivateData(r.localData, dataKey)

	if err != nil {
		slog.Error(
			"could not purge PNR data from local collection",
			"id", id,
			"error", err,
		)
		return err
	}

	return nil
}

func (r *PrivateDataRepository) GCMetadataExists(id string) (bool, error) {
	key, err := getGCMetatadaCompositeKey(id)

	if err != nil {
		slog.Error(
			"could not create GC metadata composite key",
			"id", id,
			"error", err,
		)
		return false, err
	}

	gcMetadataModel, err := r.ctx.GetStub().GetPrivateData(r.localData, key)

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

func (r *PrivateDataRepository) InsertGCMetadata(pnr entities.PNR, gc entities.GCMetadata) error {
	exists, _ := r.GCMetadataExists(pnr.Id)

	if exists {
		return errors.New("GC metadata already exists")
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

	remotePIU := getRemotePIU(pnr, r.piuId)
	remoteData := getCollectionName(remotePIU)

	err = r.putToBothPrivateCollections(remoteData, key, gcMetadataModel)

	if err != nil {
		slog.Error(
			"could not put model into private collections",
			"id", pnr.Id,
			"error", err,
		)
		return err
	}

	return nil
}

func (r *PrivateDataRepository) UpdateGCMetadata(pnr entities.PNR, gc entities.GCMetadata) error {
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

	remotePIU := getRemotePIU(pnr, r.piuId)
	remoteData := getCollectionName(remotePIU)

	err = r.putToBothPrivateCollections(remoteData, key, gcMetadataModel)

	if err != nil {
		slog.Error(
			"could not put model into private collections",
			"id", pnr.Id,
			"error", err,
		)
		return err
	}

	return nil
}

func (r *PrivateDataRepository) DeleteGCMetadata(pnr entities.PNR) error {
	exists, _ := r.GCMetadataExists(pnr.Id)

	if !exists {
		return errors.New("GC metadata does not exist")
	}

	remotePIU := getRemotePIU(pnr, r.piuId)
	remoteData := getCollectionName(remotePIU)

	key, err := getGCMetatadaCompositeKey(pnr.Id)

	if err != nil {
		slog.Error(
			"could not create GC metadata composite key",
			"id", pnr.Id,
			"error", err,
		)
		return err
	}

	err = r.ctx.GetStub().DelPrivateData(remoteData, key)

	if err != nil {
		slog.Error(
			"could not delete GC metadata from remote collection",
			"id", pnr.Id,
			"error", err,
		)
		return err
	}

	err = r.ctx.GetStub().DelPrivateData(r.localData, key)

	if err != nil {
		slog.Error(
			"could not delete GC metadata from local collection",
			"id", pnr.Id,
			"error", err,
		)
		return err
	}

	return nil
}

func (r *PrivateDataRepository) DeleteLocalGCMetadata(id string) error {
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

	err = r.ctx.GetStub().DelPrivateData(r.localData, key)

	if err != nil {
		slog.Error(
			"could not delete GC metadata from local collection",
			"id", id,
			"error", err,
		)
		return err
	}

	return nil
}

func (r *PrivateDataRepository) GetGCMetadata(id string) (entities.GCMetadata, error) {
	key, err := getGCMetatadaCompositeKey(id)

	if err != nil {
		slog.Error(
			"could not create GC metadata composite key",
			"id", id,
			"error", err,
		)
		return entities.GCMetadata{}, err
	}

	gcMetadataModel, err := r.ctx.GetStub().GetPrivateData(r.localData, key)

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

func (r *PrivateDataRepository) GetGCMetadatas() ([]entities.GCMetadata, error) {
	var result []entities.GCMetadata

	iterator, err := r.ctx.GetStub().GetPrivateDataByPartialCompositeKey(r.localData, gcMetadataObjectType, []string{})
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

func (r *PrivateDataRepository) Close() {
}
