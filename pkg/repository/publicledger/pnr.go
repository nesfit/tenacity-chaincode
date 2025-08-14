package publicledger

import (
	"encoding/json"

	"github.com/hyperledger/fabric-chaincode-go/shim"

	"github.com/nesfit/tenacity-chaincode/pkg/entities"
)

type pnrModel []byte

const pnrObjectType = "pnr"

func pnrEntityToModel(entity entities.PNR) (pnrModel, error) {
	model, err := json.Marshal(entity)

	if err != nil {
		return nil, err
	}

	return model, nil
}

func pnrModelToEntity(model pnrModel) (entities.PNR, error) {
	var entity entities.PNR

	err := json.Unmarshal(model, &entity)

	if err != nil {
		return entities.PNR{}, err
	}

	return entity, nil
}

func getPNRCompositeKey(id string) (string, error) {
	return shim.CreateCompositeKey(pnrObjectType, []string{id})
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
