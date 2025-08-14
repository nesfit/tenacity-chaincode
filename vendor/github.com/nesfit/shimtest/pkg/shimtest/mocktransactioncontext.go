package shimtest

import (
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/nesfit/shimtest/pkg/shimtest/mock"
)

type MockTransactionContext struct {
	stub     *MockStub
	identity *MockClientIdentity
}

func (m *MockTransactionContext) GetStub() shim.ChaincodeStubInterface {
	return m.stub
}

func (m *MockTransactionContext) GetClientIdentity() cid.ClientIdentity {
	return m.identity
}

func NewMockTransactionContext(name string, id string, mspID string) *MockTransactionContext {
	stub := NewMockStub(name, &mock.Chaincode{})
	identity := NewMockClientIdentity(id, mspID)
	return &MockTransactionContext{
		stub:     stub,
		identity: identity,
	}
}
