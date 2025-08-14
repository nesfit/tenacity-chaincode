package testdata

import (
	"time"

	"github.com/samber/lo"
	"github.com/nesfit/tenacity-chaincode/pkg/entities"
)

var PIUs = []entities.PIU{
	{
		Id:         "piu1",
		Name:       "PIU 1",
		AdminEmail: "admin@piu1.org",
	},
	{
		Id:         "piu2",
		Name:       "PIU 2",
		AdminEmail: "admin@piu2.org",
	},
	{
		Id:         "piu3",
		Name:       "PIU 3",
		AdminEmail: "admin@piu3.org",
	},
}

var EarliestTimestamp = lo.Must(time.Parse(time.RFC3339, "2025-11-19T12:00:00Z"))
var MiddleTimestamp = lo.Must(time.Parse(time.RFC3339, "2025-11-19T13:00:00Z"))
var LatestTimestamp = lo.Must(time.Parse(time.RFC3339, "2025-11-19T14:59:00Z"))

var PNRs = []entities.PNR{
	{
		Id:               "pnr1",
		RequestingPIU:    "piu1",
		RespondingPIU:    "piu2",
		RequestTimestamp: EarliestTimestamp,
		State:            entities.RequestStatePending,
		RequestData:      "\"requestData\"",
		PNRHashes:        []string{},
	},
	{
		Id:                "pnr2",
		RequestingPIU:     "piu2",
		RespondingPIU:     "piu1",
		RequestTimestamp:  MiddleTimestamp,
		ResponseTimestamp: MiddleTimestamp.Add(time.Minute),
		State:             entities.RequestStateAck,
		RequestData:       "\"requestData\"",
		ResponseData:      "\"responseData\"",
		PNRHashes:         []string{},
	},
	{
		Id:                "pnr3",
		RequestingPIU:     "piu2",
		RespondingPIU:     "piu1",
		RequestTimestamp:  MiddleTimestamp,
		ResponseTimestamp: MiddleTimestamp.Add(time.Minute),
		State:             entities.RequestStateAckConfirmed,
		RequestData:       "",
		ResponseData:      "",
		PNRHashes:         []string{},
	},
	{
		Id:                "pnr4",
		RequestingPIU:     "piu1",
		RespondingPIU:     "piu2",
		RequestTimestamp:  LatestTimestamp,
		ResponseTimestamp: LatestTimestamp.Add(time.Minute),
		State:             entities.RequestStateNack,
		RequestData:       "\"requestData\"",
		ResponseData:      "\"responseData\"",
		PNRHashes:         []string{},
	},
}
