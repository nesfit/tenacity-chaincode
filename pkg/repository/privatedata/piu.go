package privatedata

import (
	"encoding/json"

	"github.com/hyperledger/fabric-chaincode-go/shim"

	"github.com/nesfit/tenacity-chaincode/pkg/entities"
)

type piuModel []byte

const piuObjectType = "piu"

func piuEntityToModel(entity entities.PIU) (piuModel, error) {
	model, err := json.Marshal(entity)

	if err != nil {
		return nil, err
	}

	return model, nil
}

func piuModelToEntity(model piuModel) (entities.PIU, error) {
	var entity entities.PIU

	err := json.Unmarshal(model, &entity)

	if err != nil {
		return entities.PIU{}, err
	}

	return entity, nil
}

func getPIUCompositeKey(id string) (string, error) {
	return shim.CreateCompositeKey(piuObjectType, []string{id})
}
