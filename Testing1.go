//**************************** MUST BE USED FOR INTERNAL PURPOSE ONLY ************************************
//****FileName: Blockchain IoT Chaincode
//****Description: Chaincode logic for IoT Enabled Blockchain for Supply Chain of Perishable Goods
//****Author: Rom Solanki
//****Author Email: rosolanki@deloitte.com
//********************************************************************************************************

package main

//Importing 6 libraries for handling bytes, encoding, reading and writing JSON and string manipulation and formatting.
//Importing 2 Hyperledger Specific Libraries for Smart Contract.
import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

//Define the Smart Contract structure.
type Testing1 struct {
}

//Dividing the keys stored inside Blockchain into 3 main categories:
//1. Participants
//2. Assets
//3. Transactions

//*****************************
// 1. PARTICIPANTS STRUCTS
//*****************************

//Define the Participant structure, with 5 properties.
//Structure tags are used by encoding/json library.
type Participant struct {
	Asset_Type      string `json:"Asset_Type, omitempty"`
	ParticipantID   string `json:"ParticipantID"`
	ParticipantType string `json:"ParticipantType"`
	OrgName         string `json:"OrgName"`
	Email           string `json:"Email"`
}

// Main function (only used for Unit Testing)
func main() {
	if err := shim.Start(new(Testing1)); err != nil {
		fmt.Printf("Error starting Testing1 chaincode: %s", err)
	}
}

//Initializing new logger for logging objects used by chaincode
var logger = shim.NewLogger("Testing1")

//The Init method is called when the Smart Contract "Testing1" is instantiated by the blockchain network
func (t *Testing1) Init(stub shim.ChaincodeStubInterface) peer.Response {
	//Retrieves the arguments when instantiating chaincode
	_, args := stub.GetFunctionAndParameters()
	//Checks and returns error if any argument exists at the time of chaincode instantiation
	if len(args) > 0 {
		return shim.Error("Init Error: Incorrect number of arguments - NO ARGUMENT EXPECTED")
	}
	return shim.Success(nil)
}

//The Invoke method is called as a result of an application request to run the Smart Contract "Testing1".
//The calling application program must specify the particular smart contract function to be called, with arguments.
func (t *Testing1) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	//Get the requested Smart Contract function and arguments
	function, args := stub.GetFunctionAndParameters()

	// Route to the appropriate handler function to interact with the ledger appropriately
	switch function {
	case "createParticipant":
		return t.createParticipant(stub, args)
	case "getParticipant":
		return t.getParticipant(stub, args)
	case "deleteParticipant":
		return t.deleteParticipant(stub, args)
	default:
		logger.Warningf("Invalid Function Call - Function '%s' does not exist", function)
		return shim.Error("Invoke Error: Invalid Function Call - Function does not exist")
	}
}

//Function Definitions

// CASE 01 Create a Participant
func (t *Testing1) createParticipant(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//Checks appropriate number of arguments in incoming invoke request
	if len(args) < 2 {
		return shim.Error("Invoke Error: Incorrect number of arguments - Two Argument expected")
	}

	//Define the structure for expected incoming JSON as argument
	type QueryData struct {
		ParticipantType string `json:"ParticipantType"`
		OrgName         string `json:"OrgName"`
		Email           string `json:"Email"`
	}

	data := string(args[0])
	participantID := string(args[1])

	queryData := QueryData{}
	err := json.Unmarshal([]byte(data), &queryData)
	if err != nil {
		return shim.Error("Invoke Error (Create Participant):  Invalid Data - Check Payload")
	}
	participant := Participant{}
	participant.Asset_Type = "PARTICIPANT"
	participant.ParticipantID = participantID
	participant.ParticipantType = queryData.ParticipantType
	participant.OrgName = queryData.OrgName
	participant.Email = queryData.Email

	// Check If Participant already Exists and returns error if it does.
	if value, geterr := stub.GetState(strings.ToLower(participantID)); !(geterr == nil && value == nil) {
		return shim.Error("Invoke Error (Create Participant): Participant Already Exists! Please Specify Another ID")
	}

	// Store Participant in Blockchain
	jsonBytes, _ := json.Marshal(participant) //Get Bytes from struct
	if puterr := stub.PutState(participantID, jsonBytes); puterr != nil {
		return shim.Error(puterr.Error())
	}
	return shim.Success(nil)
}

// CASE 02 Get a Participant Info
func (t *Testing1) getParticipant(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 1 {
		return shim.Error("Invoke Error: Incorrect number of arguments - One Argument expected")
	}
	participantID := string(args[0])

	//Get the Asset from Blockchain
	value, geterr := stub.GetState(strings.ToLower(participantID))
	if geterr != nil || value == nil {
		return shim.Error(geterr.Error())
	}
	return shim.Success(value)
}

// CASE 03 Delete a Participant
func (t *Testing1) deleteParticipant(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 1 {
		return shim.Error("Invoke Error: Incorrect number of arguments - One Argument expected")
	}
	participantID := string(args[0])

	// Check if Exists
	value, geterr := stub.GetState(strings.ToLower(participantID))
	if geterr != nil || value == nil {
		return shim.Error(geterr.Error())
	}

	// Delete if Exists
	if delerr := stub.DelState(strings.ToLower(participantID)); delerr != nil {
		return shim.Error(delerr.Error())

	return shim.Success(nil)
}
