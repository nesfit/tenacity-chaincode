package inmemory_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/nesfit/tenacity-chaincode/pkg/repository"
	"github.com/nesfit/tenacity-chaincode/pkg/repository/inmemory"
)

type inmemoryRepositoryFactory struct {
}

func (f inmemoryRepositoryFactory) New() (repository.Repository, repository.TransactionManager) {
	return inmemory.NewInMemoryRepository(), inmemoryTransactionManager{}
}

type inmemoryTransactionManager struct {
}

func (txm inmemoryTransactionManager) Start() {
}

func (txm inmemoryTransactionManager) End() {
}

func TestRepositorySuite(t *testing.T) {
	s := repository.NewRepositoryTestSuite(inmemoryRepositoryFactory{})
	suite.Run(t, s)
}
