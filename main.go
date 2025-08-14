package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/nesfit/tenacity-chaincode/pkg/contract"
)

func main() {
	var handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})
	var logger = slog.New(handler)

	slog.SetDefault(logger)

	c := contract.NewSmartContract(&contract.LedgerUsecaseFactory{})
	chaincode, err := contractapi.NewChaincode(&c)
	if err != nil {
		log.Panicf("Error creating chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting  chaincode: %v", err)
	}
}
