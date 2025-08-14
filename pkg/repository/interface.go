package repository

import (
	"github.com/nesfit/tenacity-chaincode/pkg/entities"
)

type Repository interface {
	PIUExists(id string) (bool, error)
	GetPIU(id string) (entities.PIU, error)
	GetPIUs() ([]entities.PIU, error)
	InsertPIU(id string, piu entities.PIU) error
	UpdatePIU(id string, piu entities.PIU) error
	PNRExists(id string) (bool, error)
	GetPNR(id string) (entities.PNR, error)
	GetPNRs(filter entities.PNRFilter) ([]entities.PNR, error)
	InsertPNR(id string, pnr entities.PNR) error
	UpdatePNR(id string, pnr entities.PNR) error
	UpdateLocalPNR(id string, pnr entities.PNR) error
	PurgePNRData(id string) error
	PurgeLocalPNRData(id string) error
	GCMetadataExists(id string) (bool, error)
	InsertGCMetadata(pnr entities.PNR, gc entities.GCMetadata) error
	UpdateGCMetadata(pnr entities.PNR, gc entities.GCMetadata) error
	DeleteGCMetadata(pnr entities.PNR) error
	DeleteLocalGCMetadata(id string) error
	GetGCMetadata(id string) (entities.GCMetadata, error)
	GetGCMetadatas() ([]entities.GCMetadata, error)
	Close()
}

type RepositoryFactory interface {
	New() (Repository, TransactionManager)
}

type TransactionManager interface {
	Start()
	End()
}
