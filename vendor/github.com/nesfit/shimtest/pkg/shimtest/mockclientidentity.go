package shimtest

import (
	"crypto/x509"
	"fmt"
)

type MockClientIdentity struct {
	id    string
	mspID string
}

func (ci *MockClientIdentity) GetID() (string, error) {
	return ci.id, nil
}

func (ci *MockClientIdentity) GetMSPID() (string, error) {
	return ci.mspID, nil
}

func (ci *MockClientIdentity) GetAttributeValue(attrName string) (value string, found bool, err error) {
	return "", false, fmt.Errorf("not implemented")
}

func (ci *MockClientIdentity) AssertAttributeValue(attrName, attrValue string) error {
	return fmt.Errorf("not implemented")
}

func (ci *MockClientIdentity) GetX509Certificate() (*x509.Certificate, error) {
	return nil, fmt.Errorf("not implemented")
}

func NewMockClientIdentity(id string, mspID string) *MockClientIdentity {
	return &MockClientIdentity{
		id:    id,
		mspID: mspID,
	}
}
