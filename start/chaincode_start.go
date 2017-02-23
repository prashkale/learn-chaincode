/*
Copyright IBM Corp 2016 All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
		 http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
    "strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("############Error starting Nagmani Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
fmt.Println("Init is running " )

	if len(args) != 3 {
		return nil, errors.New("############Incorrect number  of arguments. Expecting 1")
	}
 stub.PutState("Initial_Amount", []byte(args[0]))
  stub.PutState("account_Namet", []byte(args[1]))
   stub.PutState("timeStamp", []byte(args[2]))
  fmt.Println(" Data writing done " )
	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	}else if function == "sendMoney"{
	return t.sendMoney(stub, args);
    }	
	fmt.Println("############invoke did not find Nagmani func: " + function)					//error

	return nil, errors.New("############Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "dummy_query" {											//read a variable
		fmt.Println("########################hi Nagmani " + function)						//error
		return nil, nil;
	} else if  function == "checkBalance"   {
	 return t.checkBalance(stub, args);
	}
	fmt.Println("query  did not find func: " + function)						//error

	return nil, errors.New("############Received unknown function query: " + function)
}
//transfer money
func (t *SimpleChaincode) checkBalance(stub shim.ChaincodeStubInterface,args []string) ([]byte, error) {
 //amount, err := stub.GetState(args[0]);
 amount, err := stub.GetState("Initial_Amount"); 
	if err != nil { return nil, errors.New("Couldn't get attribute 'amount'. Error: "+string(amount[:]) + err.Error()) }
	return amount, nil
}
//transfer money
func (t *SimpleChaincode) sendMoney(stub shim.ChaincodeStubInterface,args []string) ([]byte  , error) {
	amount, err := stub.GetState("Initial_Amount");
	var balAmt, transferAmt int;
	balAmt, err = strconv.ParseInt(string(amount[:]),0,32);
	transferAmt, err = strconv.ParseInt(args[0],0, 32);
    err = stub.PutState("Initial_Amount", []byte(strconv.Itoa( balAmt- transferAmt)));

	if err != nil { 
		fmt.Printf("SAVE_CHANGES: Error storing payment record: %s", err); 
		return nil, errors.New("Error storing payment record") 
	}
	return nil, errors.New("############Received unknown function query: "+string(amount[:]))
}