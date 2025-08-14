package publicledger_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/suite"
	"github.com/nesfit/shimtest/pkg/shimtest"

	"github.com/nesfit/tenacity-chaincode/pkg/repository"
	"github.com/nesfit/tenacity-chaincode/pkg/repository/publicledger"
)

func newMockTransactionContext() *shimtest.MockTransactionContext {
	return shimtest.NewMockTransactionContext("tenacity", "org1", "org1MSP")
}

type publicLedgerRepositoryFactory struct {
}

func (f publicLedgerRepositoryFactory) New() (repository.Repository, repository.TransactionManager) {
	ctx := newMockTransactionContext()
	return publicledger.NewPublicLedgerRepository(ctx), newTransactionManager(ctx)
}

type publicledgerTransactionManager struct {
	stub *shimtest.MockStub
	txId string
}

func newTransactionManager(ctx contractapi.TransactionContextInterface) repository.TransactionManager {
	return publicledgerTransactionManager{stub: ctx.GetStub().(*shimtest.MockStub)}
}

func (txm publicledgerTransactionManager) Start() {
	txm.txId = uuid.NewString()
	txm.stub.MockTransactionStart(txm.txId)
}

func (txm publicledgerTransactionManager) End() {
	txm.stub.MockTransactionEnd(txm.txId)
	txm.txId = ""
}

func TestRepositorySuite(t *testing.T) {
	s := repository.NewRepositoryTestSuite(publicLedgerRepositoryFactory{})
	suite.Run(t, s)
}
