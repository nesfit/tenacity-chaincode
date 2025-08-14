package entities

import (
	"encoding/json"
	"time"
)

type PIUInfo struct {
	Name       string `json:"name" required:"false" description:"Name of PIU"`
	AdminEmail string `json:"adminEmail" required:"false" description:"Email of Administrator"`
}

type PIU struct {
	Id         string `json:"id" required:"true" description:"ID is a unique uuid string that identifies a PIU."`
	Name       string `json:"name" required:"false" description:"Name of PIU"`
	AdminEmail string `json:"adminEmail" required:"false" description:"Email of Administrator"`
}

func NewPIUFromPIUInfo(id string, info PIUInfo) PIU {
	return PIU{
		Id:         id,
		Name:       info.Name,
		AdminEmail: info.AdminEmail,
	}
}

func UpdatePIUFromPIUInfo(piu PIU, info PIUInfo) PIU {
	if info.Name != "" {
		piu.Name = info.Name
	}
	if info.AdminEmail != "" {
		piu.AdminEmail = info.AdminEmail
	}
	return piu
}

type SetPIUInfoOutput struct {
}

type GetPIUsInput struct {
}

type RequestState string

const (
	RequestStatePending          RequestState = "Pending"
	RequestStatePendingConfirmed RequestState = "PendingConfirmed"
	RequestStateAck              RequestState = "Ack"
	RequestStateAckConfirmed     RequestState = "AckConfirmed"
	RequestStateNack             RequestState = "Nack"
	RequestStateNackConfirmed    RequestState = "NackConfirmed"
	RequestStateTerminated       RequestState = "Terminated"
)

const RequestDataTransientKey string = "requestData"
const ResponseDataTransientKey string = "responseData"

func GetConfirmedState(state RequestState) RequestState {
	switch state {
	case RequestStatePending:
		return RequestStatePendingConfirmed
	case RequestStateAck:
		return RequestStateAckConfirmed
	case RequestStateNack:
		return RequestStateNackConfirmed
	default:
		return state
	}
}

func HasData(state RequestState) bool {
	switch state {
	case RequestStateAckConfirmed, RequestStateNackConfirmed, RequestStateTerminated:
		return false
	default:
		return true
	}
}

func OptionalMessage(msg *json.RawMessage) string {
	if msg != nil {
		return string(*msg)
	} else {
		return ""
	}
}

func OptionalJSONMessage(msg string) *json.RawMessage {
	if msg != "" {
		json := json.RawMessage([]byte(msg))
		return &json
	} else {
		return nil
	}
}

func IsMatchingPNR(filter PNRFilter, pnr PNR) bool {
	{
		if (filter.Start != time.Time{}) {
			if pnr.RequestTimestamp.Before(filter.Start) {
				return false
			}
		}

		if (filter.End != time.Time{}) {
			if pnr.RequestTimestamp.After(filter.End) {
				return false
			}
		}

		if filter.State != "" {
			if pnr.State != filter.State {
				return false
			}
		}

		if filter.RequestingPIU != "" {
			if pnr.RequestingPIU != filter.RequestingPIU {
				return false
			}
		}

		if filter.RespondingPIU != "" {
			if pnr.RespondingPIU != filter.RespondingPIU {
				return false
			}
		}

		return true
	}
}

type PNR struct {
	Id                string       `json:"id" required:"true" format:"uuid" description:"Id of PNR request"`
	RequestingPIU     string       `json:"requestingPIU" required:"true" description:"Id of requesting PIU"`
	RespondingPIU     string       `json:"respondingPIU" required:"true" description:"Id of responding PIU"`
	RequestTimestamp  time.Time    `json:"requestTimestamp" required:"false" description:"Timestamp of request"`
	ResponseTimestamp time.Time    `json:"responseTimestamp" required:"false" description:"Timestamp of response"`
	State             RequestState `json:"state" required:"true" enum:"Pending,PendingConfirmed,Ack,AckConfirmed,Nack,NackConfirmed,Terminated" description:"State of the PNR request"`
	RequestData       string       `json:"requestData" required:"true" description:"PNR request data"`
	ResponseData      string       `json:"responseData" required:"true" description:"PNR response data"`
	PNRHashes         []string     `json:"pnrHashes" required:"true" description:"Hashes of PNRs included in response"`
}

type PNRFilter struct {
	Start         time.Time    `query:"start" required:"false" description:"Start of time period"`
	End           time.Time    `query:"end" required:"false" description:"End of time period"`
	State         RequestState `query:"state" required:"false" enum:"Pending,PendingConfirmed,Ack,AckConfirmed,Nack,NackConfirmed,Terminated" description:"State of the PNR request"`
	RequestingPIU string       `query:"requestingPIU" required:"false" description:"Id of requesting PIU"`
	RespondingPIU string       `query:"respondingPIU" required:"false" description:"Id of responding PIU"`
}

type NewPNRRequestInput struct {
	Id               string           `json:"id" required:"true" format:"uuid" description:"Id of PNR request"`
	RespondingPIU    string           `query:"respondingPIU" required:"true" description:"Id of responding PIU"`
	RequestTimestamp time.Time        `json:"requestTimestamp" required:"true" description:"Timestamp of request"`
	RequestData      *json.RawMessage `json:"requestData"`
}

type NewPNRRequestOutput struct {
	Id string `json:"id" required:"true" format:"uuid" description:"Id of the PNR request"`
}

type SubmitPNRResponseInput struct {
	Id                string           `query:"id" required:"true" format:"uuid"`
	ResponseTimestamp time.Time        `json:"responseTimestamp" required:"true" description:"Timestamp of response"`
	ResponseData      *json.RawMessage `json:"responseData"`
}

type SubmitPNRResponseOutput struct {
}

type ConfirmPNRInput struct {
	Id string `query:"id" required:"true" format:"uuid"`
}

type ConfirmPNROutput struct {
}

type TerminatePNRRequestInput struct {
	Id string `query:"id" required:"true" format:"uuid"`
}

type TerminatePNRRequestOutput struct {
}

type GCMetadata struct {
	Id                string    `json:"id" required:"true" format:"uuid" description:"Id of PNR request"`
	CreationTimestamp time.Time `json:"creationTimestamp" required:"true" description:"Creation timestamp of the PNR record"`
}
