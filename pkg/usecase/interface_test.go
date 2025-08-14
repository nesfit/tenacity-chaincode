package usecase_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/samber/lo"

	"github.com/stretchr/testify/assert"

	"github.com/nesfit/tenacity-chaincode/pkg/entities"
	"github.com/nesfit/tenacity-chaincode/pkg/repository"
	"github.com/nesfit/tenacity-chaincode/pkg/repository/inmemory"
	"github.com/nesfit/tenacity-chaincode/pkg/testdata"
	"github.com/nesfit/tenacity-chaincode/pkg/usecase"
)

var testPIUId string = testdata.PIUs[0].Id

func newTestingUsecase() (repository.Repository, usecase.PNRExchangeUsecase) {
	r := inmemory.NewInMemoryRepository()
	return r, usecase.NewRMTUsecase(testPIUId, r)
}

func setupPIUs(r repository.Repository) {
	for _, piu := range testdata.PIUs {
		r.InsertPIU(piu.Id, piu)
	}
}

func TestSetPIUInfoCreate(t *testing.T) {
	assert := assert.New(t)

	r, u := newTestingUsecase()

	expected := []entities.PIU{
		{
			Id:         testPIUId,
			Name:       "testing PIU",
			AdminEmail: "admin@testingPIU.org",
		},
	}

	input := entities.PIUInfo{
		Name:       expected[0].Name,
		AdminEmail: expected[0].AdminEmail,
	}

	var output entities.SetPIUInfoOutput

	err := u.SetPIUInfo(context.TODO(), input, &output)
	assert.NoError(err)

	actual, _ := r.GetPIUs()
	assert.ElementsMatch(expected, actual)
}

func TestSetPIUInfoUpdate(t *testing.T) {
	assert := assert.New(t)

	r, u := newTestingUsecase()

	expected := []entities.PIU{
		{
			Id:         testPIUId,
			Name:       "testing PIU",
			AdminEmail: "admin@testingPIU.org",
		},
	}

	input := entities.PIUInfo{
		Name:       expected[0].Name,
		AdminEmail: expected[0].AdminEmail,
	}

	r.InsertPIU(testPIUId, entities.PIU{
		Id:         testPIUId,
		Name:       "bad PIU",
		AdminEmail: "bad@email.org",
	})

	var output entities.SetPIUInfoOutput

	err := u.SetPIUInfo(context.TODO(), input, &output)
	assert.NoError(err)

	actual, _ := r.GetPIUs()
	assert.ElementsMatch(expected, actual)
}

func TestGetPIUsEmpty(t *testing.T) {
	assert := assert.New(t)

	_, u := newTestingUsecase()

	expected := []entities.PIU{}

	input := entities.GetPIUsInput{}

	var actual []entities.PIU

	err := u.GetPIUs(context.TODO(), input, &actual)
	assert.NoError(err)
	assert.ElementsMatch(expected, actual)
}

func TestGetPIUsNonEmpty(t *testing.T) {
	assert := assert.New(t)

	r, u := newTestingUsecase()

	expected := testdata.PIUs

	for _, piu := range expected {
		r.InsertPIU(piu.Id, piu)
	}

	input := entities.GetPIUsInput{}

	var actual []entities.PIU

	err := u.GetPIUs(context.TODO(), input, &actual)
	assert.NoError(err)
	assert.ElementsMatch(expected, actual)
}

func TestGetPNRsEmpty(t *testing.T) {
	assert := assert.New(t)

	_, u := newTestingUsecase()

	expected := []entities.PNR{}

	input := entities.PNRFilter{}

	var actual []entities.PNR

	err := u.GetPNRs(context.TODO(), input, &actual)
	assert.NoError(err)
	assert.ElementsMatch(expected, actual)
}

func TestGetPNRsNonEmpty(t *testing.T) {
	assert := assert.New(t)

	r, u := newTestingUsecase()

	expected := testdata.PNRs

	for _, pnr := range expected {
		r.InsertPNR(pnr.Id, pnr)
	}

	input := entities.PNRFilter{}

	var actual []entities.PNR

	err := u.GetPNRs(context.TODO(), input, &actual)
	assert.NoError(err)
	assert.ElementsMatch(expected, actual)
}

func TestGetPNRsFilter(t *testing.T) {
	var timeOffset = 20 * time.Minute

	testCases := map[string]struct {
		Filter   entities.PNRFilter
		Inputs   []entities.PNR
		Expected []entities.PNR
	}{
		"empty": {
			Filter:   entities.PNRFilter{},
			Expected: testdata.PNRs,
		},
		"start": {
			Filter: entities.PNRFilter{
				Start: testdata.EarliestTimestamp.Add(timeOffset),
			},
			Expected: lo.Filter(testdata.PNRs, func(v entities.PNR, i int) bool {
				return !v.RequestTimestamp.Before(testdata.EarliestTimestamp.Add(timeOffset))
			}),
		},
		"end": {
			Filter: entities.PNRFilter{
				End: testdata.LatestTimestamp.Add(-1 * timeOffset),
			},
			Expected: lo.Filter(testdata.PNRs, func(v entities.PNR, i int) bool {
				return !v.RequestTimestamp.After(testdata.LatestTimestamp.Add(-1 * timeOffset))
			}),
		},
		"startAndEnd": {
			Filter: entities.PNRFilter{
				Start: testdata.EarliestTimestamp.Add(timeOffset),
				End:   testdata.LatestTimestamp.Add(-1 * timeOffset),
			},
			Expected: lo.Filter(testdata.PNRs, func(v entities.PNR, i int) bool {
				return !v.RequestTimestamp.Before(testdata.EarliestTimestamp.Add(timeOffset)) && !v.RequestTimestamp.After(testdata.LatestTimestamp.Add(-1*timeOffset))
			}),
		},
		"state": {
			Filter: entities.PNRFilter{
				State: entities.RequestStateAck,
			},
			Expected: lo.Filter(testdata.PNRs, func(v entities.PNR, i int) bool {
				return v.State == entities.RequestStateAck
			}),
		},
		"requestingPIU": {
			Filter: entities.PNRFilter{
				RequestingPIU: testdata.PNRs[1].RequestingPIU,
			},
			Expected: lo.Filter(testdata.PNRs, func(v entities.PNR, i int) bool {
				return v.RequestingPIU == testdata.PNRs[1].RequestingPIU
			}),
		},
		"respondingPIU": {
			Filter: entities.PNRFilter{
				RespondingPIU: testdata.PNRs[1].RespondingPIU,
			},
			Expected: lo.Filter(testdata.PNRs, func(v entities.PNR, i int) bool {
				return v.RespondingPIU == testdata.PNRs[1].RespondingPIU
			}),
		},
		"exact": {
			Filter: entities.PNRFilter{
				Start:         testdata.PNRs[1].RequestTimestamp.Add(-1 * time.Microsecond),
				End:           testdata.PNRs[1].RequestTimestamp.Add(time.Microsecond),
				State:         testdata.PNRs[1].State,
				RequestingPIU: testdata.PNRs[1].RequestingPIU,
				RespondingPIU: testdata.PNRs[1].RespondingPIU,
			},
			Expected: []entities.PNR{testdata.PNRs[1]},
		},
		"exactMismatch": {
			Filter: entities.PNRFilter{
				Start:         testdata.PNRs[1].RequestTimestamp.Add(-1 * time.Microsecond),
				End:           testdata.PNRs[1].RequestTimestamp.Add(time.Microsecond),
				State:         entities.RequestStateTerminated,
				RequestingPIU: testdata.PNRs[1].RequestingPIU,
				RespondingPIU: testdata.PNRs[1].RespondingPIU,
			},
			Expected: []entities.PNR{},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			r, u := newTestingUsecase()

			for _, pnr := range testdata.PNRs {
				r.InsertPNR(pnr.Id, pnr)
			}

			var actual []entities.PNR

			err := u.GetPNRs(context.TODO(), testCase.Filter, &actual)
			assert.NoError(err)
			assert.ElementsMatch(actual, testCase.Expected)
		})
	}
}

func TestNewPNRRequest(t *testing.T) {
	assert := assert.New(t)

	r, u := newTestingUsecase()
	setupPIUs(r)

	var requestData json.RawMessage = lo.Must(json.Marshal("test request data"))

	input := entities.NewPNRRequestInput{
		RespondingPIU:    testdata.PIUs[1].Id,
		RequestTimestamp: testdata.MiddleTimestamp,
		RequestData:      &requestData,
	}

	var output entities.NewPNRRequestOutput

	err := u.NewPNRRequest(context.TODO(), input, &output)
	assert.NoError(err)

	expected := entities.PNR{
		Id:               output.Id,
		RequestingPIU:    testPIUId,
		RespondingPIU:    input.RespondingPIU,
		RequestTimestamp: input.RequestTimestamp,
		State:            entities.RequestStatePending,
		RequestData:      string(*input.RequestData),
		PNRHashes:        []string{},
	}

	actual, _ := r.GetPNR(output.Id)
	assert.Equal(expected, actual)
}

func TestNewPNRRequestMissingRequestingPIU(t *testing.T) {
	assert := assert.New(t)

	_, u := newTestingUsecase()

	input := entities.NewPNRRequestInput{
		RespondingPIU: testdata.PIUs[1].Id,
	}

	var output entities.NewPNRRequestOutput

	err := u.NewPNRRequest(context.TODO(), input, &output)
	assert.Error(err)
}

func TestNewPNRRequestMissingRespondingPIU(t *testing.T) {
	assert := assert.New(t)

	_, u := newTestingUsecase()

	input := entities.NewPNRRequestInput{
		RespondingPIU: "missingPIU",
	}

	var output entities.NewPNRRequestOutput

	err := u.NewPNRRequest(context.TODO(), input, &output)
	assert.Error(err)
}

func TestNewPNRRequestFromSelf(t *testing.T) {
	assert := assert.New(t)

	_, u := newTestingUsecase()

	input := entities.NewPNRRequestInput{
		RespondingPIU: testPIUId,
	}

	var output entities.NewPNRRequestOutput

	err := u.NewPNRRequest(context.TODO(), input, &output)
	assert.Error(err)
}

func TestSubmitPNRResponse(t *testing.T) {
	testCases := []entities.RequestState{
		entities.RequestStateAck,
		entities.RequestStateNack,
	}

	for _, state := range testCases {
		t.Run(string(state), func(t *testing.T) {
			assert := assert.New(t)

			r, u := newTestingUsecase()
			setupPIUs(r)

			var requestData json.RawMessage = lo.Must(json.Marshal("test request data"))
			var responseData json.RawMessage = lo.Must(json.Marshal("test response data"))

			originalRequest := entities.PNR{
				Id:               "someId",
				RequestingPIU:    testdata.PIUs[1].Id,
				RespondingPIU:    testPIUId,
				RequestTimestamp: testdata.MiddleTimestamp,
				State:            entities.RequestStatePendingConfirmed,
				RequestData:      string(requestData),
				PNRHashes:        []string{},
			}

			r.InsertPNR(originalRequest.Id, originalRequest)
			r.InsertGCMetadata(originalRequest, entities.GCMetadata{Id: originalRequest.Id, CreationTimestamp: originalRequest.RequestTimestamp})

			input := entities.SubmitPNRResponseInput{
				Id:                originalRequest.Id,
				ResponseTimestamp: testdata.LatestTimestamp,
				ResponseData:      &responseData,
			}

			var output entities.SubmitPNRResponseOutput
			var err error

			switch state {
			case entities.RequestStateAck:
				err = u.SubmitPNRResponseAck(context.TODO(), input, &output)
				break
			case entities.RequestStateNack:
				err = u.SubmitPNRResponseNack(context.TODO(), input, &output)
				break
			default:
				assert.Fail("wrong response type")
			}

			assert.NoError(err)

			expected := originalRequest
			expected.ResponseTimestamp = input.ResponseTimestamp
			expected.ResponseData = string(responseData)
			expected.State = state

			actual, _ := r.GetPNR(originalRequest.Id)
			assert.Equal(expected, actual)
		})
	}
}

func TestSubmitPNRResponsePNRHash(t *testing.T) {
	assert := assert.New(t)

	r, u := newTestingUsecase()
	setupPIUs(r)

	var requestData json.RawMessage = lo.Must(json.Marshal("test request data"))

	responseData, _ := os.ReadFile("../testdata/response.json")

	originalRequest := entities.PNR{
		Id:               "someId",
		RequestingPIU:    testdata.PIUs[1].Id,
		RespondingPIU:    testPIUId,
		RequestTimestamp: testdata.MiddleTimestamp,
		State:            entities.RequestStatePendingConfirmed,
		RequestData:      string(requestData),
	}

	r.InsertPNR(originalRequest.Id, originalRequest)
	r.InsertGCMetadata(originalRequest, entities.GCMetadata{Id: originalRequest.Id, CreationTimestamp: originalRequest.RequestTimestamp})

	input := entities.SubmitPNRResponseInput{
		Id:                originalRequest.Id,
		ResponseTimestamp: testdata.LatestTimestamp,
		ResponseData:      (*json.RawMessage)(&responseData),
	}

	var output entities.SubmitPNRResponseOutput
	var err error

	err = u.SubmitPNRResponseAck(context.TODO(), input, &output)

	assert.NoError(err)

	expected := []string{
		"a66a1ca0c0972304bd92d49e30a51d7bcb3488b52013e5b31f0af03420c91999",
		"734e2185105c32f43d48d47452b2cb8a3503f3738d665dc3e768e29eba08c1c0",
		"54b514ed2a93d95d03245feb7094699c8f362e1a896bcd9080f4f09a099f6f2f",
		"69250d9f0025e8826189081319e00bd494fcbbf3c7ce51083b9ac37587fd6b51",
		"9cf0206ad9857d5932625ee87007763a07d0ead58394b53e201109790e9e1df0",
		"b4b208aab68dff185c557911f0b5dd02d8cefcac8dbcc66a6adf75a3543667cb",
		"f2d9ee6f011faf23538a3bcb072cc7538d752f7b8e89e4c3382ff5fbd2718b8c",
		"3ae092497253d021fbc7e02a00bc00296564bd349aab2efee117eb7df193e1ce",
		"bb587c725bbcbf78a0c2a404978a3485ed3b785a87ec882742f862dae86e317b",
		"94b9e52e8372125ee4aa137cbb746df732cb13c3a8c9378c9906ee65d133d9e6",
	}

	actual, _ := r.GetPNR(originalRequest.Id)
	assert.Equal(expected, actual.PNRHashes)
}

func TestSubmitPNRResponseUpdatedGCMetadataTimestamp(t *testing.T) {
	assert := assert.New(t)

	r, u := newTestingUsecase()
	setupPIUs(r)

	var requestData json.RawMessage = lo.Must(json.Marshal("test request data"))

	responseData, _ := os.ReadFile("../testdata/response.json")

	originalRequest := entities.PNR{
		Id:               "someId",
		RequestingPIU:    testdata.PIUs[1].Id,
		RespondingPIU:    testPIUId,
		RequestTimestamp: testdata.MiddleTimestamp,
		State:            entities.RequestStatePendingConfirmed,
		RequestData:      string(requestData),
	}

	r.InsertPNR(originalRequest.Id, originalRequest)
	r.InsertGCMetadata(originalRequest, entities.GCMetadata{Id: originalRequest.Id, CreationTimestamp: originalRequest.RequestTimestamp})

	input := entities.SubmitPNRResponseInput{
		Id:                originalRequest.Id,
		ResponseTimestamp: testdata.LatestTimestamp,
		ResponseData:      (*json.RawMessage)(&responseData),
	}

	var output entities.SubmitPNRResponseOutput
	var err error

	err = u.SubmitPNRResponseAck(context.TODO(), input, &output)

	assert.NoError(err)

	expected := entities.GCMetadata{Id: originalRequest.Id, CreationTimestamp: time.Date(2025, time.January, 17, 12, 27, 18, 0, time.UTC)}

	actual, _ := r.GetGCMetadata(originalRequest.Id)
	assert.Equal(expected, actual)
}

func TestSubmitPNRResponseWrongPNRId(t *testing.T) {
	testCases := []entities.RequestState{
		entities.RequestStateAck,
		entities.RequestStateNack,
	}

	for _, state := range testCases {
		t.Run(string(state), func(t *testing.T) {
			assert := assert.New(t)

			r, u := newTestingUsecase()
			setupPIUs(r)

			var responseData json.RawMessage = lo.Must(json.Marshal("test response data"))

			input := entities.SubmitPNRResponseInput{
				Id:           "missing",
				ResponseData: &responseData,
			}

			var output entities.SubmitPNRResponseOutput
			var err error

			switch state {
			case entities.RequestStateAck:
				err = u.SubmitPNRResponseAck(context.TODO(), input, &output)
				break
			case entities.RequestStateNack:
				err = u.SubmitPNRResponseNack(context.TODO(), input, &output)
				break
			default:
				assert.Fail("wrong response type")
			}

			assert.Error(err)
		})
	}
}

func TestSubmitPNRResponseToSelf(t *testing.T) {
	testCases := map[string]struct {
		ResponseType entities.RequestState
		UsecaseName  string
	}{
		"ack": {
			ResponseType: entities.RequestStateAck,
		},
		"nack": {
			ResponseType: entities.RequestStateNack,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			r, u := newTestingUsecase()
			setupPIUs(r)

			var requestData json.RawMessage = lo.Must(json.Marshal("test request data"))
			var responseData json.RawMessage = lo.Must(json.Marshal("test response data"))

			originalRequest := entities.PNR{
				Id:               "someId",
				RequestingPIU:    testPIUId,
				RespondingPIU:    testdata.PIUs[1].Id,
				RequestTimestamp: testdata.MiddleTimestamp,
				State:            entities.RequestStatePendingConfirmed,
				RequestData:      string(requestData),
			}

			r.InsertPNR(originalRequest.Id, originalRequest)

			input := entities.SubmitPNRResponseInput{
				Id:                originalRequest.Id,
				ResponseTimestamp: testdata.LatestTimestamp,
				ResponseData:      &responseData,
			}

			var output entities.SubmitPNRResponseOutput
			var err error

			switch testCase.ResponseType {
			case entities.RequestStateAck:
				err = u.SubmitPNRResponseAck(context.TODO(), input, &output)
				break
			case entities.RequestStateNack:
				err = u.SubmitPNRResponseNack(context.TODO(), input, &output)
				break
			default:
				assert.Fail("wrong response type")
			}

			assert.Error(err)
		})
	}
}

func TestSubmitPNRResponseWrongState(t *testing.T) {
	testCases := []entities.RequestState{
		entities.RequestStateAck,
		entities.RequestStateNack,
	}

	for _, responseState := range testCases {
		for _, requestState := range testCases {
			t.Run(string(responseState+"/"+requestState), func(t *testing.T) {
				assert := assert.New(t)

				r, u := newTestingUsecase()
				setupPIUs(r)

				var requestData json.RawMessage = lo.Must(json.Marshal("test request data"))
				var responseData json.RawMessage = lo.Must(json.Marshal("test response data"))

				originalRequest := entities.PNR{
					Id:               "someId",
					RequestingPIU:    testdata.PIUs[1].Id,
					RespondingPIU:    testPIUId,
					RequestTimestamp: testdata.MiddleTimestamp,
					State:            requestState,
					RequestData:      string(requestData),
				}

				r.InsertPNR(originalRequest.Id, originalRequest)

				input := entities.SubmitPNRResponseInput{
					Id:                originalRequest.Id,
					ResponseTimestamp: testdata.LatestTimestamp,
					ResponseData:      &responseData,
				}

				var output entities.SubmitPNRResponseOutput
				var err error

				switch responseState {
				case entities.RequestStateAck:
					err = u.SubmitPNRResponseAck(context.TODO(), input, &output)
					break
				case entities.RequestStateNack:
					err = u.SubmitPNRResponseNack(context.TODO(), input, &output)
					break
				default:
					assert.Fail("wrong response type")
				}

				assert.Error(err)

			})
		}
	}
}

func TestConfirmPNR(t *testing.T) {
	testCases := map[entities.RequestState]struct {
		RequestingPIU string
		RespondingPIU string
	}{
		entities.RequestStatePending: {
			RequestingPIU: testdata.PIUs[1].Id,
			RespondingPIU: testPIUId,
		},
		entities.RequestStateAck: {
			RequestingPIU: testPIUId,
			RespondingPIU: testdata.PIUs[1].Id,
		},
		entities.RequestStateNack: {
			RequestingPIU: testPIUId,
			RespondingPIU: testdata.PIUs[1].Id,
		},
	}

	for state, testCase := range testCases {
		t.Run(string(state), func(t *testing.T) {
			assert := assert.New(t)

			r, u := newTestingUsecase()
			setupPIUs(r)

			var requestData json.RawMessage = lo.Must(json.Marshal("test request data"))
			var responseData json.RawMessage = lo.Must(json.Marshal("test response data"))

			originalRequest := entities.PNR{
				Id:               "someId",
				RequestingPIU:    testCase.RequestingPIU,
				RespondingPIU:    testCase.RespondingPIU,
				RequestTimestamp: testdata.MiddleTimestamp,
				State:            state,
				RequestData:      string(requestData),
				ResponseData:     string(responseData),
			}

			r.InsertPNR(originalRequest.Id, originalRequest)

			input := entities.ConfirmPNRInput{
				Id: originalRequest.Id,
			}

			var output entities.ConfirmPNROutput

			err := u.ConfirmPNR(context.TODO(), input, &output)
			assert.NoError(err)

			expected := originalRequest
			expected.State = entities.GetConfirmedState(originalRequest.State)

			if state == entities.RequestStateAck || state == entities.RequestStateNack {
				expected.RequestData = ""
				expected.ResponseData = ""
			}

			actual, _ := r.GetPNR(originalRequest.Id)
			assert.Equal(expected, actual)
		})
	}
}

func TestConfirmPNRWrongState(t *testing.T) {
	testCases := map[entities.RequestState]struct {
		RequestingPIU string
		RespondingPIU string
	}{
		entities.RequestStatePending: {
			RequestingPIU: testPIUId,
			RespondingPIU: testdata.PIUs[1].Id,
		},
		entities.RequestStateAck: {
			RequestingPIU: testdata.PIUs[1].Id,
			RespondingPIU: testPIUId,
		},
		entities.RequestStateNack: {
			RequestingPIU: testdata.PIUs[1].Id,
			RespondingPIU: testPIUId,
		},
	}

	for state, testCase := range testCases {
		t.Run(string(state), func(t *testing.T) {
			assert := assert.New(t)

			r, u := newTestingUsecase()
			setupPIUs(r)

			var requestData json.RawMessage = lo.Must(json.Marshal("test request data"))

			originalRequest := entities.PNR{
				Id:               "someId",
				RequestingPIU:    testCase.RequestingPIU,
				RespondingPIU:    testCase.RespondingPIU,
				RequestTimestamp: testdata.MiddleTimestamp,
				State:            state,
				RequestData:      string(requestData),
			}

			r.InsertPNR(originalRequest.Id, originalRequest)

			input := entities.ConfirmPNRInput{
				Id: originalRequest.Id,
			}

			var output entities.ConfirmPNROutput

			err := u.ConfirmPNR(context.TODO(), input, &output)
			assert.Error(err)
		})
	}
}

func TestConfirmPNRAlreadyConfirmed(t *testing.T) {
	testCases := map[entities.RequestState]struct {
		RequestingPIU string
		RespondingPIU string
	}{
		entities.RequestStatePendingConfirmed: {
			RequestingPIU: testdata.PIUs[1].Id,
			RespondingPIU: testPIUId,
		},
		entities.RequestStateAckConfirmed: {
			RequestingPIU: testPIUId,
			RespondingPIU: testdata.PIUs[1].Id,
		},
		entities.RequestStateNackConfirmed: {
			RequestingPIU: testPIUId,
			RespondingPIU: testdata.PIUs[1].Id,
		},
	}

	for state, testCase := range testCases {
		t.Run(string(state), func(t *testing.T) {
			assert := assert.New(t)

			r, u := newTestingUsecase()
			setupPIUs(r)

			var requestData json.RawMessage = lo.Must(json.Marshal("test request data"))

			originalRequest := entities.PNR{
				Id:               "someId",
				RequestingPIU:    testCase.RequestingPIU,
				RespondingPIU:    testCase.RespondingPIU,
				RequestTimestamp: testdata.MiddleTimestamp,
				State:            state,
				RequestData:      string(requestData),
			}

			r.InsertPNR(originalRequest.Id, originalRequest)

			input := entities.ConfirmPNRInput{
				Id: originalRequest.Id,
			}

			var output entities.ConfirmPNROutput

			err := u.ConfirmPNR(context.TODO(), input, &output)
			assert.Error(err)
		})
	}
}

func TestTerminatePNRRequest(t *testing.T) {
	testCases := map[entities.RequestState]struct {
		RequestingPIU string
		RespondingPIU string
	}{
		"requester": {
			RequestingPIU: testPIUId,
			RespondingPIU: testdata.PIUs[1].Id,
		},
		"responder": {
			RequestingPIU: testdata.PIUs[1].Id,
			RespondingPIU: testPIUId,
		},
	}

	for name, testCase := range testCases {
		t.Run(string(name), func(t *testing.T) {
			assert := assert.New(t)

			r, u := newTestingUsecase()
			setupPIUs(r)

			var requestData json.RawMessage = lo.Must(json.Marshal("test request data"))
			var responseData json.RawMessage = lo.Must(json.Marshal("test response data"))

			originalRequest := entities.PNR{
				Id:                "someId",
				RequestingPIU:     testCase.RequestingPIU,
				RespondingPIU:     testCase.RespondingPIU,
				RequestTimestamp:  testdata.MiddleTimestamp,
				ResponseTimestamp: testdata.LatestTimestamp,
				State:             entities.RequestStateAck,
				RequestData:       string(requestData),
				ResponseData:      string(responseData),
				PNRHashes:         []string{},
			}

			r.InsertPNR(originalRequest.Id, originalRequest)

			r.InsertGCMetadata(originalRequest, entities.GCMetadata{Id: originalRequest.Id, CreationTimestamp: originalRequest.RequestTimestamp})

			input := entities.TerminatePNRRequestInput{
				Id: originalRequest.Id,
			}

			var output entities.TerminatePNRRequestOutput

			err := u.TerminatePNRRequest(context.TODO(), input, &output)
			assert.NoError(err)

			expected := originalRequest
			expected.State = entities.RequestStateTerminated
			expected.RequestData = ""
			expected.ResponseData = ""

			actual, _ := r.GetPNR(originalRequest.Id)
			assert.Equal(expected, actual)
		})
	}
}

func TestTerminatePNRRequestMissingPNR(t *testing.T) {
	assert := assert.New(t)

	r, u := newTestingUsecase()
	setupPIUs(r)

	input := entities.TerminatePNRRequestInput{
		Id: "missing",
	}

	var output entities.TerminatePNRRequestOutput

	err := u.TerminatePNRRequest(context.TODO(), input, &output)
	assert.Error(err)
}

func TestTerminatePNRRequestMissingGCMetadata(t *testing.T) {
	testCases := map[entities.RequestState]struct {
		RequestingPIU string
		RespondingPIU string
	}{
		"requester": {
			RequestingPIU: testPIUId,
			RespondingPIU: testdata.PIUs[1].Id,
		},
		"responder": {
			RequestingPIU: testdata.PIUs[1].Id,
			RespondingPIU: testPIUId,
		},
	}

	for name, testCase := range testCases {
		t.Run(string(name), func(t *testing.T) {
			assert := assert.New(t)

			r, u := newTestingUsecase()
			setupPIUs(r)

			var requestData json.RawMessage = lo.Must(json.Marshal("test request data"))
			var responseData json.RawMessage = lo.Must(json.Marshal("test response data"))

			originalRequest := entities.PNR{
				Id:                "someId",
				RequestingPIU:     testCase.RequestingPIU,
				RespondingPIU:     testCase.RespondingPIU,
				RequestTimestamp:  testdata.MiddleTimestamp,
				ResponseTimestamp: testdata.LatestTimestamp,
				State:             entities.RequestStateAck,
				RequestData:       string(requestData),
				ResponseData:      string(responseData),
				PNRHashes:         []string{},
			}

			r.InsertPNR(originalRequest.Id, originalRequest)

			input := entities.TerminatePNRRequestInput{
				Id: originalRequest.Id,
			}

			var output entities.TerminatePNRRequestOutput

			err := u.TerminatePNRRequest(context.TODO(), input, &output)
			assert.Error(err)
		})
	}
}

func TestTerminatePNRRequestUnrelatedPIU(t *testing.T) {
	assert := assert.New(t)

	r, u := newTestingUsecase()
	setupPIUs(r)

	var requestData json.RawMessage = lo.Must(json.Marshal("test request data"))
	var responseData json.RawMessage = lo.Must(json.Marshal("test response data"))

	originalRequest := entities.PNR{
		Id:                "someId",
		RequestingPIU:     testdata.PIUs[1].Id,
		RespondingPIU:     testdata.PIUs[2].Id,
		RequestTimestamp:  testdata.MiddleTimestamp,
		ResponseTimestamp: testdata.LatestTimestamp,
		State:             entities.RequestStateAck,
		RequestData:       string(requestData),
		ResponseData:      string(responseData),
		PNRHashes:         []string{},
	}

	r.InsertPNR(originalRequest.Id, originalRequest)

	input := entities.TerminatePNRRequestInput{
		Id: originalRequest.Id,
	}

	var output entities.TerminatePNRRequestOutput

	err := u.TerminatePNRRequest(context.TODO(), input, &output)
	assert.Error(err)
}
