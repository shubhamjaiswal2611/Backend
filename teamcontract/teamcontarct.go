package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Contract represents the contract details
type SmartContract struct {
	ContractID   string `json:"contractID"`
	ContractFile string `json:"contractFile"`
	ContractType string `json:"contractType"`
	TeamName     string `json:"teamName"`
	ContractWith string `json:"contractWith"`
	UploadedBy   string `json:"uploadedBy"`
	RecordDate   string `json:"recordDate"`
	Status       string `json:"status"`
	ActionBy     string `json:"actionBy"`
	ActionDate   string `json:"actionDate"`
	Comment      string `json:"comment"`
}

// ContractChaincode represents the smart contract
type SmartContractChaincode struct {
	contractapi.Contract
}

func (cc *SmartContractChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// Initialize any initial data if required
	return nil
}

// CreateContract creates a new contract
func (cc *SmartContractChaincode) AddTeamContract(ctx contractapi.TransactionContextInterface, contractID string, contractDetailsJSON string) error {
	// Check if the contract ID already exists
	existing, err := ctx.GetStub().GetState(contractID) //check---
	if err != nil {
		return fmt.Errorf("failed to read from world state: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("the contract ID already exists")
	}

	// Parse the contract details JSON
	contractDetails := &SmartContract{}
	err = json.Unmarshal([]byte(contractDetailsJSON), contractDetails)
	if err != nil {
		return fmt.Errorf("failed to unmarshal contract details: %w", err)
	}

	// Set the contract details in the world state
	contractDetailsBytes, err := json.Marshal(contractDetails)
	if err != nil {
		return fmt.Errorf("failed to marshal contract details: %w", err)
	}
	err = ctx.GetStub().PutState(contractID, contractDetailsBytes)
	if err != nil {
		return fmt.Errorf("failed to put contract details in world state: %w", err)
	}

	return nil
}

// ApproveContract approves a contract by updating the status, action by, and action date
func (cc *SmartContractChaincode) ApproveTeamContract(ctx contractapi.TransactionContextInterface, contractID string, status string, actionBy string, actionDate string, comment string) error {
	// Get the contract details from the world state
	contractDetailsBytes, err := ctx.GetStub().GetState(contractID)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %w", err)
	}
	if contractDetailsBytes == nil {
		return fmt.Errorf("the contract ID does not exist")
	}

	// Parse the contract details JSON
	contractDetails := &SmartContract{}
	err = json.Unmarshal(contractDetailsBytes, contractDetails)
	if err != nil {
		return fmt.Errorf("failed to unmarshal contract details: %w", err)
	}

	// Update the contract details
	contractDetails.Status = status
	contractDetails.ActionBy = actionBy
	contractDetails.ActionDate = actionDate
	contractDetails.Comment = comment

	// Set the updated contract details in the world state
	updatedContractDetailsBytes, err := json.Marshal(contractDetails)
	if err != nil {
		return fmt.Errorf("failed to marshal updated contract details: %w", err)
	}
	err = ctx.GetStub().PutState(contractID, updatedContractDetailsBytes)
	if err != nil {
		return fmt.Errorf("failed to put updated contract details in world state: %w", err)
	}

	return nil
}

func (s *SmartContractChaincode) GetAllContract(ctx contractapi.TransactionContextInterface) ([]*SmartContract, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")

	if err != nil {
		return nil, err
	}

	defer resultsIterator.Close()
	var assets []*SmartContract
	for resultsIterator.HasNext() {

		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var asset SmartContract
		err = json.Unmarshal(queryResponse.Value, &asset)

		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}
	return assets, nil

}

//Query the status of a contract
//history contract

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContractChaincode{})
	if err != nil {
		fmt.Printf("Error creating contract chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting contract chaincode: %s", err.Error())
	}
}
