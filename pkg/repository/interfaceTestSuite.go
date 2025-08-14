package repository

import (
	"slices"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/nesfit/tenacity-chaincode/pkg/entities"
	"github.com/nesfit/tenacity-chaincode/pkg/testdata"
)

type RepositoryTestSuite struct {
	suite.Suite
	r   Repository
	rf  RepositoryFactory
	txm TransactionManager
}

func NewRepositoryTestSuite(rf RepositoryFactory) suite.TestingSuite {
	s := new(RepositoryTestSuite)
	s.rf = rf
	return s
}

func (s *RepositoryTestSuite) SetupTest() {
	s.r, s.txm = s.rf.New()
}

func (s *RepositoryTestSuite) SetupSubTest() {
	s.r, s.txm = s.rf.New()
}

func (s *RepositoryTestSuite) TestPIUExistsEmpty() {
	assert := assert.New(s.T())

	exists, err := s.r.PIUExists("missing")
	assert.NoError(err)
	assert.False(exists)
}

func (s *RepositoryTestSuite) TestPIUExistsNonEmptyNotMatching() {
	assert := assert.New(s.T())

	for _, piu := range testdata.PIUs {
		s.r.InsertPIU(piu.Id, piu)
	}

	exists, err := s.r.PIUExists("missing")
	assert.NoError(err)
	assert.False(exists)
}

func (s *RepositoryTestSuite) TestPIUExistsNonEmptyMatching() {
	assert := assert.New(s.T())

	s.txm.Start()
	for _, piu := range testdata.PIUs {
		s.r.InsertPIU(piu.Id, piu)
	}
	s.txm.End()

	exists, err := s.r.PIUExists(testdata.PIUs[1].Id)
	assert.NoError(err)
	assert.True(exists)
}

func (s *RepositoryTestSuite) TestGetPIUEmpty() {
	assert := assert.New(s.T())

	_, err := s.r.GetPIU("missing")
	assert.Error(err)
}

func (s *RepositoryTestSuite) TestGetPIUNotMatching() {
	assert := assert.New(s.T())

	s.txm.Start()
	for _, piu := range testdata.PIUs {
		s.r.InsertPIU(piu.Id, piu)
	}
	s.txm.End()

	_, err := s.r.GetPIU("missing")
	assert.Error(err)
}

func (s *RepositoryTestSuite) TestGetPIUMatching() {
	assert := assert.New(s.T())

	s.txm.Start()
	for _, piu := range testdata.PIUs {
		s.r.InsertPIU(piu.Id, piu)
	}
	s.txm.End()

	expected := testdata.PIUs[1]
	actual, err := s.r.GetPIU(expected.Id)

	assert.NoError(err)
	assert.Equal(expected, actual)
}

func (s *RepositoryTestSuite) TestGetPIUsEmpty() {
	assert := assert.New(s.T())

	expected := []entities.PIU{}

	actual, err := s.r.GetPIUs()
	assert.NoError(err)
	assert.ElementsMatch(expected, actual)
}

func (s *RepositoryTestSuite) TestGetPIUsNonEmpty() {
	assert := assert.New(s.T())

	s.txm.Start()
	for _, piu := range testdata.PIUs {
		s.r.InsertPIU(piu.Id, piu)
	}
	s.txm.End()

	expected := testdata.PIUs
	actual, err := s.r.GetPIUs()
	assert.NoError(err)
	assert.ElementsMatch(expected, actual)
}

func (s *RepositoryTestSuite) TestInsertPIU() {
	assert := assert.New(s.T())

	insertedPIU := testdata.PIUs[1]
	expected := []entities.PIU{insertedPIU}

	s.txm.Start()
	err := s.r.InsertPIU(insertedPIU.Id, insertedPIU)
	s.txm.End()
	assert.NoError(err)

	actual, _ := s.r.GetPIUs()
	assert.ElementsMatch(expected, actual)
}

func (s *RepositoryTestSuite) TestInsertPIUAlreadyExists() {
	assert := assert.New(s.T())

	insertedPIU := testdata.PIUs[1]

	s.txm.Start()
	s.r.InsertPIU(insertedPIU.Id, insertedPIU)
	s.txm.End()

	s.txm.Start()
	err := s.r.InsertPIU(insertedPIU.Id, insertedPIU)
	s.txm.End()
	assert.Error(err)
}

func (s *RepositoryTestSuite) TestUpdatePIU() {
	assert := assert.New(s.T())

	updatedPIU := entities.PIU{
		Id:         testdata.PIUs[1].Id,
		Name:       "New PIU 1",
		AdminEmail: "admin@new.piu1.org",
	}
	expected := slices.Concat(testdata.PIUs[0:1], []entities.PIU{updatedPIU}, testdata.PIUs[2:])

	s.txm.Start()
	for _, piu := range testdata.PIUs {
		s.r.InsertPIU(piu.Id, piu)
	}
	s.txm.End()

	s.txm.Start()
	err := s.r.UpdatePIU(updatedPIU.Id, updatedPIU)
	s.txm.End()
	assert.NoError(err)

	actual, _ := s.r.GetPIUs()
	assert.ElementsMatch(expected, actual)
}

func (s *RepositoryTestSuite) TestUpdatePIUDoesNotExist() {
	assert := assert.New(s.T())

	s.txm.Start()
	for _, piu := range testdata.PIUs {
		s.r.InsertPIU(piu.Id, piu)
	}
	s.txm.End()

	s.txm.Start()
	err := s.r.UpdatePIU("missing", entities.PIU{})
	s.txm.End()
	assert.Error(err)
}

func (s *RepositoryTestSuite) TestPNRExistsEmpty() {
	assert := assert.New(s.T())

	exists, err := s.r.PNRExists("missing")
	assert.NoError(err)
	assert.False(exists)
}

func (s *RepositoryTestSuite) TestPNRExistsNonEmptyNotMatching() {
	assert := assert.New(s.T())

	s.txm.Start()
	for _, pnr := range testdata.PNRs {
		s.r.InsertPNR(pnr.Id, pnr)
	}
	s.txm.End()

	exists, err := s.r.PNRExists("missing")
	assert.NoError(err)
	assert.False(exists)
}

func (s *RepositoryTestSuite) TestPNRExistsNonEmptyMatching() {
	assert := assert.New(s.T())

	s.txm.Start()
	for _, pnr := range testdata.PNRs {
		s.r.InsertPNR(pnr.Id, pnr)
	}
	s.txm.End()

	exists, err := s.r.PNRExists(testdata.PNRs[1].Id)
	assert.NoError(err)
	assert.True(exists)
}

func (s *RepositoryTestSuite) TestGetPNREmpty() {
	assert := assert.New(s.T())

	_, err := s.r.GetPNR("missing")
	assert.Error(err)
}

func (s *RepositoryTestSuite) TestGetPNRNotMatching() {
	assert := assert.New(s.T())

	s.txm.Start()
	for _, pnr := range testdata.PNRs {
		s.r.InsertPNR(pnr.Id, pnr)
	}
	s.txm.End()

	_, err := s.r.GetPNR("missing")
	assert.Error(err)
}

func (s *RepositoryTestSuite) TestGetPNRMatching() {
	assert := assert.New(s.T())

	s.txm.Start()
	for _, pnr := range testdata.PNRs {
		s.r.InsertPNR(pnr.Id, pnr)
	}
	s.txm.End()

	expected := testdata.PNRs[1]
	actual, err := s.r.GetPNR(expected.Id)

	assert.NoError(err)
	assert.Equal(expected, actual)
}

func (s *RepositoryTestSuite) TestGetPNRsEmpty() {
	assert := assert.New(s.T())

	expected := []entities.PNR{}

	actual, err := s.r.GetPNRs(entities.PNRFilter{})
	assert.NoError(err)
	assert.ElementsMatch(expected, actual)
}

func (s *RepositoryTestSuite) TestGetPNRsNonEmpty() {
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
		s.Run(name, func() {
			assert := assert.New(s.T())

			s.txm.Start()
			for _, pnr := range testdata.PNRs {
				s.r.InsertPNR(pnr.Id, pnr)
			}
			s.txm.End()

			actual, err := s.r.GetPNRs(testCase.Filter)
			assert.NoError(err)
			assert.ElementsMatch(actual, testCase.Expected)
		})
	}
}

func (s *RepositoryTestSuite) TestInsertPNR() {
	assert := assert.New(s.T())

	insertedPNR := testdata.PNRs[1]
	expected := []entities.PNR{insertedPNR}

	s.txm.Start()
	err := s.r.InsertPNR(insertedPNR.Id, insertedPNR)
	s.txm.End()
	assert.NoError(err)

	actual, _ := s.r.GetPNRs(entities.PNRFilter{})
	assert.ElementsMatch(expected, actual)
}

func (s *RepositoryTestSuite) TestInsertPNRAlreadyExists() {
	assert := assert.New(s.T())

	insertedPNR := testdata.PNRs[1]
	expected := []entities.PNR{insertedPNR}

	s.txm.Start()
	s.r.InsertPNR(insertedPNR.Id, insertedPNR)
	s.txm.End()

	s.txm.Start()
	err := s.r.InsertPNR(insertedPNR.Id, insertedPNR)
	s.txm.End()
	assert.Error(err)

	actual, _ := s.r.GetPNRs(entities.PNRFilter{})
	assert.ElementsMatch(expected, actual)
}

func (s *RepositoryTestSuite) TestUpdatePNR() {
	assert := assert.New(s.T())

	s.txm.Start()
	for _, pnr := range testdata.PNRs {
		s.r.InsertPNR(pnr.Id, pnr)
	}
	s.txm.End()

	updatedPNR := entities.PNR{
		Id:                testdata.PNRs[1].Id,
		RequestingPIU:     "new-piu",
		RespondingPIU:     "other-new-piu",
		RequestTimestamp:  testdata.PNRs[1].RequestTimestamp.Add(time.Minute),
		ResponseTimestamp: testdata.PNRs[1].ResponseTimestamp.Add(time.Minute),
		State:             entities.RequestStateNack,
		RequestData:       "\"new data\"",
		ResponseData:      "\"other new data\"",
	}
	expected := slices.Concat(testdata.PNRs[0:1], []entities.PNR{updatedPNR}, testdata.PNRs[2:])

	s.txm.Start()
	err := s.r.UpdatePNR(updatedPNR.Id, updatedPNR)
	s.txm.End()
	assert.NoError(err)

	actual, _ := s.r.GetPNRs(entities.PNRFilter{})
	assert.ElementsMatch(expected, actual)
}

func (s *RepositoryTestSuite) TestUpdatePNRDoesNotExist() {
	assert := assert.New(s.T())

	s.txm.Start()
	for _, pnr := range testdata.PNRs {
		s.r.InsertPNR(pnr.Id, pnr)
	}
	s.txm.End()

	s.txm.Start()
	err := s.r.UpdatePNR("missing", entities.PNR{})
	s.txm.End()
	assert.Error(err)
}

func (s *RepositoryTestSuite) TestPurgePNRData() {
	assert := assert.New(s.T())

	insertedPNR := testdata.PNRs[2]
	expected := insertedPNR
	expected.RequestData = ""
	expected.ResponseData = ""

	s.txm.Start()
	s.r.InsertPNR(insertedPNR.Id, insertedPNR)
	s.txm.End()

	s.txm.Start()
	err := s.r.PurgePNRData(insertedPNR.Id)
	s.txm.End()
	assert.NoError(err)

	actual, _ := s.r.GetPNR(insertedPNR.Id)
	assert.Equal(expected, actual)
}

func (s *RepositoryTestSuite) TestPurgePNRDataDoesNotExist() {
	assert := assert.New(s.T())

	s.txm.Start()
	err := s.r.PurgePNRData("missing")
	s.txm.End()
	assert.Error(err)
}

func (s *RepositoryTestSuite) TestGetGCMetadatasEmpty() {
	assert := assert.New(s.T())

	actual, err := s.r.GetGCMetadatas()
	assert.NoError(err)
	assert.Empty(actual)
}

func (s *RepositoryTestSuite) TestInsertGCMetadata() {
	assert := assert.New(s.T())

	pnr := testdata.PNRs[0]
	gc := entities.GCMetadata{Id: pnr.Id, CreationTimestamp: pnr.RequestTimestamp}

	s.txm.Start()
	err := s.r.InsertGCMetadata(pnr, gc)
	s.txm.End()

	expected := []entities.GCMetadata{gc}

	actual, _ := s.r.GetGCMetadatas()
	assert.NoError(err)
	assert.ElementsMatch(expected, actual)
}

func (s *RepositoryTestSuite) TestInsertGCMetadataAlreadyExists() {
	assert := assert.New(s.T())

	pnr := testdata.PNRs[0]
	gc := entities.GCMetadata{Id: pnr.Id, CreationTimestamp: pnr.RequestTimestamp}

	s.txm.Start()
	s.r.InsertGCMetadata(pnr, gc)
	s.txm.End()

	expected := []entities.GCMetadata{gc}
	gc = entities.GCMetadata{Id: pnr.Id, CreationTimestamp: pnr.ResponseTimestamp}

	s.txm.Start()
	err := s.r.InsertGCMetadata(pnr, gc)
	s.txm.End()

	actual, _ := s.r.GetGCMetadatas()
	assert.Error(err)
	assert.ElementsMatch(expected, actual)
}

func (s *RepositoryTestSuite) TestUpdateGCMetadata() {
	assert := assert.New(s.T())

	pnr := testdata.PNRs[0]
	gc := entities.GCMetadata{Id: pnr.Id, CreationTimestamp: pnr.RequestTimestamp}

	s.txm.Start()
	s.r.InsertGCMetadata(pnr, gc)
	s.txm.End()

	gc.CreationTimestamp = pnr.ResponseTimestamp

	s.txm.Start()
	err := s.r.UpdateGCMetadata(pnr, gc)
	s.txm.End()

	expected := []entities.GCMetadata{gc}

	actual, _ := s.r.GetGCMetadatas()
	assert.NoError(err)
	assert.ElementsMatch(expected, actual)
}

func (s *RepositoryTestSuite) TestUpdateGCMetadataDoesNotExist() {
	assert := assert.New(s.T())

	pnr := testdata.PNRs[0]
	gc := entities.GCMetadata{Id: pnr.Id, CreationTimestamp: pnr.RequestTimestamp}

	s.txm.Start()
	err := s.r.UpdateGCMetadata(pnr, gc)
	s.txm.End()

	actual, _ := s.r.GetGCMetadatas()
	assert.Error(err)
	assert.Empty(actual)
}

func (s *RepositoryTestSuite) TestDeleteGCMetadata() {
	assert := assert.New(s.T())

	var gcs []entities.GCMetadata

	s.txm.Start()
	for _, pnr := range testdata.PNRs {
		gc := entities.GCMetadata{Id: pnr.Id, CreationTimestamp: pnr.RequestTimestamp}
		gcs = append(gcs, gc)
		s.r.InsertGCMetadata(pnr, gc)
	}
	s.txm.End()

	s.txm.Start()
	err := s.r.DeleteGCMetadata(testdata.PNRs[1])
	s.txm.End()

	expected := slices.Delete(gcs, 1, 2)

	actual, _ := s.r.GetGCMetadatas()
	assert.NoError(err)
	assert.ElementsMatch(expected, actual)
}
