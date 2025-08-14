package privatedata

import (
	"encoding/json"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"

	"github.com/nesfit/tenacity-chaincode/pkg/entities"
)

type pnrMeta struct {
	Id                string                `json:"id" required:"true" format:"uuid" description:"Id of PNR request"`
	RequestingPIU     string                `json:"requestingPIU" required:"true" description:"Id of requesting PIU"`
	RespondingPIU     string                `json:"respondingPIU" required:"true" description:"Id of responding PIU"`
	RequestTimestamp  time.Time             `json:"requestTimestamp" required:"false" description:"Timestamp of request"`
	ResponseTimestamp time.Time             `json:"responseTimestamp" required:"false" description:"Timestamp of response"`
	State             entities.RequestState `json:"state" required:"true" enum:"Pending,PendingConfirmed,Ack,AckConfirmed,Nack,NackConfirmed,Terminated" description:"State of the PNR request"`
	PNRHashes         []string              `json:"pnrHashes" required:"true" description:"Hashes of PNRs included in response"`
}

type pnrData struct {
	RequestData  string `json:"requestData" required:"true" description:"PNR request data"`
	ResponseData string `json:"responseData" required:"true" description:"PNR response data"`
}

type pnrModel []byte

const pnrMetaObjectType = "pnrMeta"
const pnrDataObjectType = "pnrData"

func pnrEntityToMetaEntity(entity entities.PNR) pnrMeta {
	return pnrMeta{
		Id:                entity.Id,
		RequestingPIU:     entity.RequestingPIU,
		RespondingPIU:     entity.RespondingPIU,
		RequestTimestamp:  entity.RequestTimestamp,
		ResponseTimestamp: entity.ResponseTimestamp,
		State:             entity.State,
		PNRHashes:         entity.PNRHashes,
	}
}

func pnrEntityToDataEntity(entity entities.PNR) pnrData {
	return pnrData{
		RequestData:  entity.RequestData,
		ResponseData: entity.ResponseData,
	}
}

func pnrEntitiesToEntity(metaEntity pnrMeta, dataEntity pnrData) entities.PNR {
	return entities.PNR{
		Id:                metaEntity.Id,
		RequestingPIU:     metaEntity.RequestingPIU,
		RespondingPIU:     metaEntity.RespondingPIU,
		RequestTimestamp:  metaEntity.RequestTimestamp,
		ResponseTimestamp: metaEntity.ResponseTimestamp,
		State:             metaEntity.State,
		PNRHashes:         metaEntity.PNRHashes,
		RequestData:       dataEntity.RequestData,
		ResponseData:      dataEntity.ResponseData,
	}
}

func pnrEntityToMetaModel(entity entities.PNR) (pnrModel, error) {
	metaEntity := pnrEntityToMetaEntity(entity)

	model, err := json.Marshal(metaEntity)

	if err != nil {
		return nil, err
	}

	return model, nil
}

func pnrEntityToDataModel(entity entities.PNR) (pnrModel, error) {
	dataEntity := pnrEntityToDataEntity(entity)

	model, err := json.Marshal(dataEntity)

	if err != nil {
		return nil, err
	}

	return model, nil
}

func metaModelToMetaEntity(model pnrModel) (pnrMeta, error) {
	var entity pnrMeta

	err := json.Unmarshal(model, &entity)

	if err != nil {
		return pnrMeta{}, err
	}

	return entity, nil
}

func dataModelToDataEntity(model pnrModel) (pnrData, error) {
	var entity pnrData

	err := json.Unmarshal(model, &entity)

	if err != nil {
		return pnrData{}, err
	}

	return entity, nil
}

func pnrModelsToEntity(metaModel pnrModel, dataModel pnrModel) (entities.PNR, error) {
	metaEntity, err := metaModelToMetaEntity(metaModel)

	if err != nil {
		return entities.PNR{}, err
	}

	dataEntity, err := dataModelToDataEntity(dataModel)

	if err != nil {
		return entities.PNR{}, err
	}

	return pnrEntitiesToEntity(metaEntity, dataEntity), nil
}

func getPNRMetaCompositeKey(id string) (string, error) {
	return shim.CreateCompositeKey(pnrMetaObjectType, []string{id})
}

func getPNRDataCompositeKey(id string) (string, error) {
	return shim.CreateCompositeKey(pnrDataObjectType, []string{id})
}

func getRemotePIU(pnr entities.PNR, localPIU string) string {
	if pnr.RequestingPIU == localPIU {
		return pnr.RespondingPIU
	} else {
		return pnr.RequestingPIU
	}
}

type gcMetadataModel []byte

const gcMetadataObjectType = "gc"

func gcMetadataEntityToModel(entity entities.GCMetadata) (gcMetadataModel, error) {
	model, err := json.Marshal(entity)

	if err != nil {
		return nil, err
	}

	return model, nil
}

func gcMetadataModelToEntity(model pnrModel) (entities.GCMetadata, error) {
	var entity entities.GCMetadata

	err := json.Unmarshal(model, &entity)

	if err != nil {
		return entities.GCMetadata{}, err
	}

	return entity, nil
}

func getGCMetatadaCompositeKey(id string) (string, error) {
	return shim.CreateCompositeKey(gcMetadataObjectType, []string{id})
}
