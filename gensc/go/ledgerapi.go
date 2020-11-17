/*  This file has all required wrapped function to interact with Hyperledger Fabric  Ledger*/

package main

import (
	"bytes"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// CreateAsset LedgerAPI to Create Asset into Ledger
// ===========================================================================================
func CreateAsset(stub shim.ChaincodeStubInterface, assetType string, assetkeys []string, assetData []byte) error {

	if len(assetkeys) == 0 {
		return fmt.Errorf(" CreateAsset: Key is not provided for object %s", assetType)
	}
	logger.Infof("Data is %s", string(assetData))
	//creating compositkey for this assetobject
	compositeKey, _ := stub.CreateCompositeKey(assetType, assetkeys)
	data, err := stub.GetState(compositeKey)
	if err != nil {
		logger.Errorf("CreateAsset : Error Querying Object in State Database %s", err)
		return err
	}
	if data != nil {
		return fmt.Errorf(" CreateAsset: Data Exist for Assettype:  %s with Key:  %s", assetType, assetkeys)
	}

	//adding object to ledger
	err = stub.PutState(compositeKey, assetData)
	if err != nil {
		logger.Errorf("CreateAsset : Error inserting Object into State Database %s", err)
		return err
	}

	return nil
}

// ListAllAsset LedgerAPI to ListAllAsset from ledger
// ===========================================================================================
func ListAllAsset(stub shim.ChaincodeStubInterface, assetType string, assetkeys []string) (shim.StateQueryIteratorInterface, error) {

	//creating compositkey for this assetobject
	dataresultIter, err := stub.GetStateByPartialCompositeKey(assetType, assetkeys)
	logger.Infof("Iterator for List data: %+v", dataresultIter)
	if err != nil {
		logger.Errorf("CreateAsset : Error Querying Object in State Database %s", err)
		return nil, err
	}
	return dataresultIter, nil
}

// QueryAsset to Query Asset from Ledger
// ===========================================================================================
func QueryAsset(stub shim.ChaincodeStubInterface, assetType string, assetkeys []string) ([]byte, error) {

	if len(assetkeys) == 0 {
		return nil, fmt.Errorf(" QueryAsset: Key is not provided for object %s", assetType)
	}
	//creating compositkey for this assetobject
	compositeKey, _ := stub.CreateCompositeKey(assetType, assetkeys)
	logger.Infof("QueryAsset : Querying Object in State Database %s", compositeKey)

	data, err := stub.GetState(compositeKey)
	if err != nil {
		logger.Errorf("QueryAsset : Error Querying Object in State Database %s", err)
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf(" QueryAsset: Data does not exist for compositeKey %s ", compositeKey)
	}
	logger.Infof("QueryAsset : data is  %s", string(data))
	return data, nil
}

// GenericQueryAsset to query any kind of asset from Ledger
// This function is only called with HLF ledger is using Couch DB or other reach query ledger like oracle Berkeley DB
// ===========================================================================================
func GenericQueryAsset(stub shim.ChaincodeStubInterface, query string) (shim.StateQueryIteratorInterface, error) {

	dataresultIter, err := stub.GetQueryResult(query)
	logger.Infof("GenericQueryAsset: Iterator for Generic Query list  data: %+v", dataresultIter)
	if err != nil {
		logger.Errorf("GenericQueryAsset : Error Querying Object in State Database %s", err)
		return nil, err
	}
	return dataresultIter, nil

}

// GenericQueryAssetwithPeginations to query any kind of asset from Ledger
//this function is only called with HLF ledger is using Couch DB or other reach query ledger like oracle Berkeley DB
// ===========================================================================================
func GenericQueryAssetwithPeginations(stub shim.ChaincodeStubInterface, query string, pagesize int32, bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {

	dataresultIter, metadata, err := stub.GetQueryResultWithPagination(query, pagesize, bookmark)
	logger.Infof("GenericQueryAsset: Iterator for Generic Query list  data: %+v", dataresultIter)
	if err != nil {
		logger.Errorf("GenericQueryAsset : Error Querying Object in State Database %s", err)
		return nil, nil, err
	}
	return dataresultIter, metadata, nil

}

// ===========================================================================================
// constructQueryResponseFromIterator constructs a JSON array containing query results from
// a given result iterator
// ===========================================================================================
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}

// ===========================================================================================
// addPaginationMetadataToQueryResults adds QueryResponseMetadata, which contains pagination
// info, to the constructed query results
// ===========================================================================================
func addPaginationMetadataToQueryResults(buffer *bytes.Buffer, responseMetadata *pb.QueryResponseMetadata) *bytes.Buffer {

	buffer.WriteString("[{\"ResponseMetadata\":{\"RecordsCount\":")
	buffer.WriteString("\"")
	buffer.WriteString(fmt.Sprintf("%v", responseMetadata.FetchedRecordsCount))
	buffer.WriteString("\"")
	buffer.WriteString(", \"Bookmark\":")
	buffer.WriteString("\"")
	buffer.WriteString(responseMetadata.Bookmark)
	buffer.WriteString("\"}}]")

	return buffer
}
