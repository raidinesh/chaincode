/*This is main program to start Generic Smart Contract for adding Asset to Fabric Ledger*/

package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ============================================================================================================================
// Main function for FebSC call
// ============================================================================================================================
var logger = shim.NewLogger("FabSc")

func main() {
	logger.SetLevel(shim.LogDebug)
	err := shim.Start(new(FabSc))
	if err != nil {
		logger.Errorf("Error starting FabSc Smart Contract - %s", err)
	}
}
