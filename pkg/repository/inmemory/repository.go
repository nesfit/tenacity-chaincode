package inmemory

import (
	"errors"

	"github.com/nesfit/tenacity-chaincode/pkg/entities"
)

type InMemoryRepository struct {
	pius        map[string]entities.PIU
	pnrs        map[string]entities.PNR
	gcMetadatas map[string]entities.GCMetadata
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		pius:        make(map[string]entities.PIU),
		pnrs:        make(map[string]entities.PNR),
		gcMetadatas: make(map[string]entities.GCMetadata),
	}
}

func (r *InMemoryRepository) PIUExists(id string) (bool, error) {
	_, ok := r.pius[id]

	return ok, nil
}

func (r *InMemoryRepository) GetPIU(id string) (entities.PIU, error) {
	entity, ok := r.pius[id]

	if !ok {
		return entities.PIU{}, errors.New("PIU not found")
	}

	return entity, nil
}

func (r *InMemoryRepository) GetPIUs() ([]entities.PIU, error) {
	result := make([]entities.PIU, 0, len(r.pius))

	for _, entity := range r.pius {
		result = append(result, entity)
	}

	return result, nil
}

func (r *InMemoryRepository) InsertPIU(id string, piu entities.PIU) error {
	exists, _ := r.PIUExists(id)

	if exists {
		return errors.New("PIU already exists")
	}

	r.pius[id] = piu

	return nil
}

func (r *InMemoryRepository) UpdatePIU(id string, piu entities.PIU) error {
	exists, _ := r.PIUExists(id)

	if !exists {
		return errors.New("PIU does not exist")
	}

	r.pius[id] = piu

	return nil
}

func (r *InMemoryRepository) PNRExists(id string) (bool, error) {
	_, ok := r.pnrs[id]

	return ok, nil
}

func (r *InMemoryRepository) GetPNR(id string) (entities.PNR, error) {
	entity, ok := r.pnrs[id]

	if !ok {
		return entities.PNR{}, errors.New("PNR not found")
	}

	return entity, nil
}

func (r *InMemoryRepository) GetPNRs(filter entities.PNRFilter) ([]entities.PNR, error) {
	result := make([]entities.PNR, 0, len(r.pnrs))

	for _, entity := range r.pnrs {
		if entities.IsMatchingPNR(filter, entity) {
			result = append(result, entity)
		}
	}

	return result, nil
}

func (r *InMemoryRepository) InsertPNR(id string, pnr entities.PNR) error {
	exists, _ := r.PNRExists(id)

	if exists {
		return errors.New("PNR already exists")
	}

	r.pnrs[id] = pnr

	return nil
}

func (r *InMemoryRepository) UpdatePNR(id string, pnr entities.PNR) error {
	exists, _ := r.PNRExists(id)

	if !exists {
		return errors.New("PNR does not exist")
	}

	r.pnrs[id] = pnr

	return nil
}

func (r *InMemoryRepository) UpdateLocalPNR(id string, pnr entities.PNR) error {
	return r.UpdatePNR(id, pnr)
}

func (r *InMemoryRepository) PurgePNRData(id string) error {
	pnr, err := r.GetPNR(id)

	if err != nil {
		return err
	}

	pnr.RequestData = ""
	pnr.ResponseData = ""

	r.UpdatePNR(id, pnr)

	return nil
}

func (r *InMemoryRepository) PurgeLocalPNRData(id string) error {
	return r.PurgePNRData(id)
}

func (r *InMemoryRepository) GCMetadataExists(id string) (bool, error) {
	_, ok := r.gcMetadatas[id]

	return ok, nil
}

func (r *InMemoryRepository) InsertGCMetadata(pnr entities.PNR, gc entities.GCMetadata) error {
	exists, _ := r.GCMetadataExists(pnr.Id)

	if exists {
		return errors.New("PNR GC metadata already exists")
	}

	r.gcMetadatas[pnr.Id] = gc

	return nil
}

func (r *InMemoryRepository) UpdateGCMetadata(pnr entities.PNR, gc entities.GCMetadata) error {
	exists, _ := r.GCMetadataExists(pnr.Id)

	if !exists {
		return errors.New("PNR GC metadata does not exist")
	}

	r.gcMetadatas[pnr.Id] = gc

	return nil
}

func (r *InMemoryRepository) DeleteGCMetadata(pnr entities.PNR) error {
	return r.DeleteLocalGCMetadata(pnr.Id)
}

func (r *InMemoryRepository) DeleteLocalGCMetadata(id string) error {
	_, err := r.GetGCMetadata(id)

	if err != nil {
		return err
	}

	delete(r.gcMetadatas, id)

	return nil
}

func (r *InMemoryRepository) GetGCMetadata(id string) (entities.GCMetadata, error) {
	entity, ok := r.gcMetadatas[id]

	if !ok {
		return entities.GCMetadata{}, errors.New("PNR GC metadata not found")
	}

	return entity, nil
}

func (r *InMemoryRepository) GetGCMetadatas() ([]entities.GCMetadata, error) {
	result := make([]entities.GCMetadata, 0, len(r.gcMetadatas))

	for _, entity := range r.gcMetadatas {
		result = append(result, entity)
	}

	return result, nil
}

func (r *InMemoryRepository) Close() {
}
