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
	"bytes"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"net/http"
	"io/ioutil"
	)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// Payment structure
type Payment struct {
	Destination string `json:"destination"`
	SourceAmount string `json:"sourceAmount"`
	DestinationAmount string `json:"destinationAmount"`
	Message	string `json:"message"`
}

//==============================================================================================================================
//	Account - Defines the structure for a account object. JSON on right tells it what JSON fields to map to
//			  that element when reading a JSON object into the struct e.g. JSON make -> Struct Make.
//==============================================================================================================================
type Account struct {
	AccountId            string `json:"accountId"`
	AccountName           string `json:"accountName"`
	Balance                string `json:"balance"`
	TimeStamp             string    `json:"timeStamp"`
	
}
// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("############Error starting  Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
fmt.Println("Init is running " )

	if len(args) != 3 {
		return nil, errors.New("############Incorrect number  of arguments. Expecting 1")
	}
	stub.PutState("Default_Open_Balance", []byte(args[0]))
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
    }else if function == "createAccount"{
	return t.createAccount(stub, args[0],args[1],args[2]);
    }	
	fmt.Println("############invoke did not find  func: " + function)					//error

	return nil, errors.New("############Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "dummy_query" {											//read a variable
		fmt.Println("########################hi Nagmani " + function)						//error
		return nil, nil;
	} else if  function == "fetchAccountDetails"   {
	 return t.fetchAccountDetails(stub, args);
	}
	fmt.Println("query  did not find func: " + function)						//error

	return nil, errors.New("############Received unknown function query: " + function)
}
//Check Balance
func (t *SimpleChaincode) fetchAccountDetails(stub shim.ChaincodeStubInterface,args []string) ([]byte, error) {
    var a Account
    var temp []byte
        a, err :=t.retrieve_Account(stub, args[0]);
        temp, err = t.get_account_details(stub,a)
		if err != nil {
	return nil, errors.New("fetchAccountDetails: Invalid Account")
	}
	    return temp, nil
}
//==============================================================================================================================
//	 retrieve_Account - Gets the state of the data at accountId in the ledger then converts it from the stored
//					JSON into the Account struct for use in the contract. Returns the Account struct.
//					Returns empty a if it errors.
//==============================================================================================================================
func (t *SimpleChaincode) retrieve_Account(stub shim.ChaincodeStubInterface, accountId string) (Account, error) {

	var a Account

	bytes, err := stub.GetState(accountId);

	if err != nil {	
	fmt.Printf("RETRIEVE_Account: Failed to invoke Account_Id: %s", err); return a, errors.New("RETRIEVE_Account: Error retrieving account with Account Id = " + accountId) 
	}

	err = json.Unmarshal(bytes, &a);

    if err != nil {	fmt.Printf("RETRIEVE_account: Corrupt account record "+string(bytes)+": %s", err); return a, errors.New("RETRIEVE_account: Corrupt account record"+string(bytes))	}

	return a, nil
}
func (t *SimpleChaincode) get_account_details(stub shim.ChaincodeStubInterface, a Account) ([]byte, error) {

	bytes, err := json.Marshal(a)

	if err != nil {
	return nil, errors.New("GET_Account: Invalid Account")
	}

return bytes, nil
}
//transfer money
func (t *SimpleChaincode) sendMoney(stub shim.ChaincodeStubInterface,args []string) ([]byte  , error) {
	amount, err := stub.GetState("Default_Open_Balance");
	payStatusCd,uUID, err := t.makePayment(args);
	if err != nil {
		fmt.Printf("Make Payment Status: Error storing payment record: %s", payStatusCd[:]); 
		fmt.Printf("Make Payment Status ID : Error storing payment record: %s", uUID[:]); 
	}
	var balAmt, transferAmt int64;
	var newBalance []byte;
	balAmt, err = strconv.ParseInt(string(amount[:]),0,64);
	transferAmt, err = strconv.ParseInt(args[0],0, 64);
	newBalance = []byte(strconv.Itoa( int(balAmt) - int(transferAmt)))
	err = stub.PutState("Initial_Amount", newBalance);

	//err = stub.PutState("Initial_Amount", []byte(strconv.Itoa( balAmt- transferAmt)));
	
	if err != nil { 
		fmt.Printf("SAVE_CHANGES: Error storing payment record: %s", err); 
		return nil, errors.New("Error storing payment record") 
	}
	//return nil, errors.New("############Received unknown function query: "+string(amount[:]))
	return nil, nil
}

//	 Create Account - Creates the initial JSON for the Account and then saves it to the ledger.
//==============================================================================================================================
func (t *SimpleChaincode) createAccount(stub shim.ChaincodeStubInterface,accountId string, accountName string,timestamp_val string) ([]byte  , error) {
amount, err := stub.GetState("Default_Open_Balance");
var acc Account
	acountId         := "\"AccountId\":\""+accountId+"\", "							// Variables to define the JSON
	acountName         := "\"AccountName\":\""+accountName+"\", "	
	balance           := "\"Balance\":\""+string(amount[:])+"\", "	
	timestamp          := "\"TimeStamp\":\""+timestamp_val+"\""
	

	account_json := "{"+acountId+acountName+balance+timestamp+"}" 	// Concatenates the variables to create the total JSON 
    err = json.Unmarshal([]byte(account_json), &acc)	
	// If not an error then a account exists so cant create a new account with this acountId as it must be unique
	record, err := stub.GetState(acc.AccountId) 
    if record != nil { 
	return nil, errors.New("Account already exists") 
	}
    _, err  = t.openAccount(stub, acc)

	if err != nil { 
	fmt.Printf("CREATE_ACCOUNT: Error saving changes: %s", err); 
	return nil, errors.New("Error saving changes") 
	}

	return nil, nil
}
//==============================================================================================================================
// openAccount - Writes to the ledger the Account struct passed in a JSON format. Uses the shim file's
//				  method 'PutState'.
//==============================================================================================================================
func (t *SimpleChaincode) openAccount(stub shim.ChaincodeStubInterface, a Account) (bool, error) {

	bytes, err := json.Marshal(a)

	if err != nil {
	fmt.Printf("OPEN_ACCOUNT: Error converting Account record: %s", err);
	return false, errors.New("Error converting Account record") 
	}

	err = stub.PutState(a.AccountId, bytes)

	if err != nil { 
	fmt.Printf("SAVE_CHANGES: Error storing Account record: %s", err); return false, errors.New("Error storing Account record") 
	}

	return true, nil
}

//==============================================================================================================================
// makePayment- Sends money to Wallet Account using REST service exposed by Wallet
//==============================================================================================================================
func (t *SimpleChaincode) makePayment(args []string) (string, string, error){
	//w http.ResponseWriter, r *http.Request
	//r.ParseForm()
	url := "http://services-uscentral.skytap.com:10504/payments/9efa70ec-08b9-11e6-b512-3e1d05defe78"

	var payment Payment
	var resp *http.Response 
	// res := r.FormValue("<your param name>")
	//payment.Destination = r.FormValue("destination");
	//payment.SourceAmount = r.FormValue("sourceAmount");
	//payment.DestinationAmount = r.FormValue("destinationAmount");
	//payment.Message = r.FormValue("message");

	
	payment.Destination = args[0];
	payment.SourceAmount = args[1];
	payment.DestinationAmount = args[2];
	payment.Message = args[3];
	
	// try Option 1
	//bufObj := new(bytes.Buffer)
	//json.NewEncoder(bufObj).Encode(payment)
	// res, _ := http.Post(url, "application/json; charset=utf-8", bufObj)
	
	// try Option 2
	paymentByteAry, err := json.Marshal(payment)
	resp, err = http.Post(url, "application/json", bytes.NewBuffer(paymentByteAry))
	
	if err != nil { 
		fmt.Printf("MAKEPAYMENTS: Error making payment : %s", err); 
		return "", "", errors.New("Error while maying payment") 
	}
	// try Option 3
	//req, err := http.NewRequest("POST", url, bytes.NewBuffer(payment))
    	//req.Header.Set("Content-Type", "application/json")
	//client := &http.Client{}
	//resp, err := client.Do(req)
	//if err != nil {
	//    panic(err)
	//}
	//defer resp.Body.Close()

	
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return resp.Status,resp.Status, errors.New("Error while maying payment") 
}