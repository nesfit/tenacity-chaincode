package contract_test

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/nesfit/shimtest/pkg/shimtest"

	"github.com/nesfit/tenacity-chaincode/pkg/contract"
	"github.com/nesfit/tenacity-chaincode/pkg/entities"
	"github.com/nesfit/tenacity-chaincode/pkg/repository"
	"github.com/nesfit/tenacity-chaincode/pkg/repository/inmemory"
	"github.com/nesfit/tenacity-chaincode/pkg/testdata"
	"github.com/nesfit/tenacity-chaincode/pkg/usecase"
)

var thisPIUId = testdata.PIUs[0].Id
var peerPIUId = testdata.PIUs[1].Id

var requestData json.RawMessage = lo.Must(json.Marshal("test request data"))
var responseData json.RawMessage = lo.Must(json.Marshal("test request data"))

type testUsecaseFactory struct {
	r repository.Repository
}

func (uf *testUsecaseFactory) New(ctx contractapi.TransactionContextInterface) (usecase.PNRExchangeUsecase, error) {
	piuId, _ := contract.GetClientOrgId(ctx)

	u := usecase.NewRMTUsecase(piuId, uf.r)

	return u, nil
}

func setTransient(ctx contractapi.TransactionContextInterface, data map[string][]byte) error {
	stub := ctx.GetStub().(*shimtest.MockStub)

	txId := uuid.NewString()
	stub.MockTransactionStart(txId)
	err := stub.SetTransient(data)
	stub.MockTransactionEnd(txId)

	return err
}

type ContractTestSuite struct {
	suite.Suite
	c              contract.SmartContract
	thisPIUContext *shimtest.MockTransactionContext
	peerPIUContext *shimtest.MockTransactionContext
}

func (suite *ContractTestSuite) SetupTest() {
	suite.c = contract.NewSmartContract(&testUsecaseFactory{r: inmemory.NewInMemoryRepository()})
	suite.thisPIUContext = shimtest.NewMockTransactionContext("tenacity", "org1", thisPIUId)
	suite.peerPIUContext = shimtest.NewMockTransactionContext("tenacity", "org1", peerPIUId)
}

func (suite *ContractTestSuite) initPIUPair() {
	thisPIUInfo := entities.PIUInfo{
		Name:       "foo",
		AdminEmail: "hello@piu.org",
	}

	peerPIUInfo := entities.PIUInfo{
		Name:       "bar",
		AdminEmail: "bye@piu.org",
	}

	var infoJSON []byte

	infoJSON, _ = json.Marshal(thisPIUInfo)
	suite.c.SetPIUInfo(suite.thisPIUContext, string(infoJSON))

	infoJSON, _ = json.Marshal(peerPIUInfo)
	suite.c.SetPIUInfo(suite.peerPIUContext, string(infoJSON))
}

func TestContractSuite(t *testing.T) {
	suite.Run(t, new(ContractTestSuite))
}

func (suite *ContractTestSuite) TestSetPIUInfo() {
	assert := assert.New(suite.T())

	info := entities.PIUInfo{
		Name:       "foo",
		AdminEmail: "hello@piu.org",
	}

	infoJSON, _ := json.Marshal(info)

	err := suite.c.SetPIUInfo(suite.thisPIUContext, string(infoJSON))
	assert.NoError(err)

	expected := []entities.PIU{entities.NewPIUFromPIUInfo(thisPIUId, info)}

	actual, _ := suite.c.GetPIUs(suite.thisPIUContext)
	assert.ElementsMatch(expected, actual)
}

func (suite *ContractTestSuite) TestGetPIUs() {
	assert := assert.New(suite.T())

	thisPIUInfo := entities.PIUInfo{
		Name:       "foo",
		AdminEmail: "hello@piu.org",
	}

	peerPIUInfo := entities.PIUInfo{
		Name:       "bar",
		AdminEmail: "bye@piu.org",
	}

	var infoJSON []byte

	infoJSON, _ = json.Marshal(thisPIUInfo)
	suite.c.SetPIUInfo(suite.thisPIUContext, string(infoJSON))

	infoJSON, _ = json.Marshal(peerPIUInfo)
	suite.c.SetPIUInfo(suite.peerPIUContext, string(infoJSON))

	expected := []entities.PIU{entities.NewPIUFromPIUInfo(thisPIUId, thisPIUInfo), entities.NewPIUFromPIUInfo(peerPIUId, peerPIUInfo)}

	actual, err := suite.c.GetPIUs(suite.thisPIUContext)
	assert.NoError(err)
	assert.ElementsMatch(expected, actual)
}

func (suite *ContractTestSuite) TestNewPNRRequest() {
	assert := assert.New(suite.T())

	suite.initPIUPair()

	request := entities.NewPNRRequestInput{
		RespondingPIU:    peerPIUId,
		RequestTimestamp: testdata.MiddleTimestamp,
	}

	transient := map[string][]byte{
		entities.RequestDataTransientKey: requestData,
	}

	err := setTransient(suite.thisPIUContext, transient)
	assert.NoError(err)

	var requestJSON []byte

	requestJSON, _ = json.Marshal(request)
	response, err := suite.c.NewPNRRequest(suite.thisPIUContext, string(requestJSON))

	expected := []entities.PNR{
		{
			Id:               response.Id,
			RequestingPIU:    thisPIUId,
			RespondingPIU:    request.RespondingPIU,
			RequestTimestamp: request.RequestTimestamp,
			State:            entities.RequestStatePending,
			RequestData:      string(requestData),
			PNRHashes:        []string{},
		},
	}

	actual, _ := suite.c.GetPNRs(suite.thisPIUContext, string(lo.Must(json.Marshal(entities.PNRFilter{}))))
	assert.ElementsMatch(expected, actual)
}

func (suite *ContractTestSuite) TestSubmitPNRResponseAck() {
	assert := assert.New(suite.T())

	suite.initPIUPair()

	request := entities.NewPNRRequestInput{
		RespondingPIU:    peerPIUId,
		RequestTimestamp: testdata.MiddleTimestamp,
	}

	transient := map[string][]byte{
		entities.RequestDataTransientKey: requestData,
	}

	err := setTransient(suite.thisPIUContext, transient)
	assert.NoError(err)

	var requestJSON []byte

	requestJSON, _ = json.Marshal(request)
	requestResponse, _ := suite.c.NewPNRRequest(suite.thisPIUContext, string(requestJSON))

	confirm := entities.ConfirmPNRInput{
		Id: requestResponse.Id,
	}

	var confirmJSON []byte
	confirmJSON, _ = json.Marshal(confirm)
	suite.c.ConfirmPNR(suite.peerPIUContext, string(confirmJSON))

	response := entities.SubmitPNRResponseInput{
		Id:                requestResponse.Id,
		ResponseTimestamp: testdata.MiddleTimestamp,
	}

	transient = map[string][]byte{
		entities.ResponseDataTransientKey: responseData,
	}

	err = setTransient(suite.peerPIUContext, transient)
	assert.NoError(err)

	responseJSON, _ := json.Marshal(response)
	err = suite.c.SubmitPNRResponseAck(suite.peerPIUContext, string(responseJSON))
	assert.NoError(err)

	expected := []entities.PNR{
		{
			Id:                requestResponse.Id,
			RequestingPIU:     thisPIUId,
			RespondingPIU:     request.RespondingPIU,
			RequestTimestamp:  request.RequestTimestamp,
			ResponseTimestamp: response.ResponseTimestamp,
			State:             entities.RequestStateAck,
			RequestData:       string(requestData),
			ResponseData:      string(responseData),
			PNRHashes:         []string{},
		},
	}

	actual, _ := suite.c.GetPNRs(suite.thisPIUContext, string(lo.Must(json.Marshal(entities.PNRFilter{}))))
	assert.ElementsMatch(expected, actual)
}

func (suite *ContractTestSuite) TestSubmitPNRResponseNack() {
	assert := assert.New(suite.T())

	suite.initPIUPair()

	request := entities.NewPNRRequestInput{
		RespondingPIU:    peerPIUId,
		RequestTimestamp: testdata.MiddleTimestamp,
	}

	transient := map[string][]byte{
		entities.RequestDataTransientKey: requestData,
	}

	err := setTransient(suite.thisPIUContext, transient)
	assert.NoError(err)

	var requestJSON []byte

	requestJSON, _ = json.Marshal(request)
	requestResponse, _ := suite.c.NewPNRRequest(suite.thisPIUContext, string(requestJSON))

	confirm := entities.ConfirmPNRInput{
		Id: requestResponse.Id,
	}

	var confirmJSON []byte
	confirmJSON, _ = json.Marshal(confirm)
	suite.c.ConfirmPNR(suite.peerPIUContext, string(confirmJSON))

	response := entities.SubmitPNRResponseInput{
		Id:                requestResponse.Id,
		ResponseTimestamp: testdata.MiddleTimestamp,
	}

	transient = map[string][]byte{
		entities.ResponseDataTransientKey: responseData,
	}

	err = setTransient(suite.peerPIUContext, transient)
	assert.NoError(err)

	responseJSON, _ := json.Marshal(response)
	err = suite.c.SubmitPNRResponseNack(suite.peerPIUContext, string(responseJSON))
	assert.NoError(err)

	expected := []entities.PNR{
		{
			Id:                requestResponse.Id,
			RequestingPIU:     thisPIUId,
			RespondingPIU:     request.RespondingPIU,
			RequestTimestamp:  request.RequestTimestamp,
			ResponseTimestamp: response.ResponseTimestamp,
			State:             entities.RequestStateNack,
			RequestData:       string(requestData),
			ResponseData:      string(responseData),
			PNRHashes:         []string{},
		},
	}

	actual, _ := suite.c.GetPNRs(suite.thisPIUContext, string(lo.Must(json.Marshal(entities.PNRFilter{}))))
	assert.ElementsMatch(expected, actual)
}

func (suite *ContractTestSuite) TestConfirmPNR() {
	assert := assert.New(suite.T())

	suite.initPIUPair()

	request := entities.NewPNRRequestInput{
		RespondingPIU:    peerPIUId,
		RequestTimestamp: testdata.MiddleTimestamp,
	}

	transient := map[string][]byte{
		entities.RequestDataTransientKey: requestData,
	}

	err := setTransient(suite.thisPIUContext, transient)
	assert.NoError(err)

	var requestJSON []byte

	requestJSON, _ = json.Marshal(request)
	requestResponse, _ := suite.c.NewPNRRequest(suite.thisPIUContext, string(requestJSON))

	confirmation := entities.ConfirmPNRInput{
		Id: requestResponse.Id,
	}

	confirmationJSON, _ := json.Marshal(confirmation)
	err = suite.c.ConfirmPNR(suite.peerPIUContext, string(confirmationJSON))
	assert.NoError(err)

	expected := []entities.PNR{
		{
			Id:               requestResponse.Id,
			RequestingPIU:    thisPIUId,
			RespondingPIU:    request.RespondingPIU,
			RequestTimestamp: request.RequestTimestamp,
			State:            entities.RequestStatePendingConfirmed,
			RequestData:      string(requestData),
			PNRHashes:        []string{},
		},
	}

	actual, _ := suite.c.GetPNRs(suite.thisPIUContext, string(lo.Must(json.Marshal(entities.PNRFilter{}))))
	assert.ElementsMatch(expected, actual)
}

func (suite *ContractTestSuite) TestTerminatePNRRequest() {
	assert := assert.New(suite.T())

	suite.initPIUPair()

	request := entities.NewPNRRequestInput{
		RespondingPIU:    peerPIUId,
		RequestTimestamp: testdata.MiddleTimestamp,
	}

	transient := map[string][]byte{
		entities.RequestDataTransientKey: requestData,
	}

	err := setTransient(suite.thisPIUContext, transient)
	assert.NoError(err)

	var requestJSON []byte

	requestJSON, _ = json.Marshal(request)
	requestResponse, _ := suite.c.NewPNRRequest(suite.thisPIUContext, string(requestJSON))

	confirm := entities.ConfirmPNRInput{
		Id: requestResponse.Id,
	}

	var confirmJSON []byte
	confirmJSON, _ = json.Marshal(confirm)
	suite.c.ConfirmPNR(suite.peerPIUContext, string(confirmJSON))

	response := entities.SubmitPNRResponseInput{
		Id:                requestResponse.Id,
		ResponseTimestamp: testdata.MiddleTimestamp,
	}

	transient = map[string][]byte{
		entities.ResponseDataTransientKey: responseData,
	}

	err = setTransient(suite.peerPIUContext, transient)
	assert.NoError(err)

	responseJSON, _ := json.Marshal(response)
	suite.c.SubmitPNRResponseNack(suite.peerPIUContext, string(responseJSON))

	purge := entities.TerminatePNRRequestInput{
		Id: requestResponse.Id,
	}
	purgeJSON, _ := json.Marshal(purge)
	err = suite.c.TerminatePNRRequest(suite.thisPIUContext, string(purgeJSON))
	assert.NoError(err)

	expected := []entities.PNR{
		{
			Id:                requestResponse.Id,
			RequestingPIU:     thisPIUId,
			RespondingPIU:     request.RespondingPIU,
			RequestTimestamp:  request.RequestTimestamp,
			ResponseTimestamp: response.ResponseTimestamp,
			State:             entities.RequestStateTerminated,
			PNRHashes:         []string{},
		},
	}

	actual, _ := suite.c.GetPNRs(suite.thisPIUContext, string(lo.Must(json.Marshal(entities.PNRFilter{}))))
	assert.ElementsMatch(expected, actual)
}
