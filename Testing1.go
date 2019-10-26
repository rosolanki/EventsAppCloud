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

//*****************************
// 2. Assets STRUCTS
//*****************************

//Define the Material structure, with XXX properties.
//Structure tags are used by encoding/json library.
type Material struct {
	Asset_Type           string                        `json:"Asset_Type, omitempty"`
	MaterialID           string                        `json:"MaterialID"`
	OpenPurchaseOrders   []MaterialOpenPurchaseOrder   `json:"OpenPurchaseOrders, omitempty"`
	ClosedPurchaseOrders []MaterialClosedPurchaseOrder `json:"ClosedPurchaseOrders, omitempty"`
	ProductionOrders     []MaterialProductionOrder     `json:"ProductionOrders, omitempty"`
	ActiveBatches        []MaterialBatches             `json:"ActiveBatches, omitempty"`
	Batches              []MaterialBatches             `json:"Batches, omitempty"`
}

type MaterialOpenPurchaseOrder struct {
	PurchaseOrderID       string                         `json:"PurchaseOrderID"`
	Owner                 string                         `json:"Owner"`
	AssociatedSalesOrders []MaterialAssociatedSalesOrder `json:"AssociatedSalesOrders"`
}

type MaterialClosedPurchaseOrder struct {
	PurchaseOrderID       string                         `json:"PurchaseOrderID"`
	Owner                 string                         `json:"Owner"`
	AssociatedSalesOrders []MaterialAssociatedSalesOrder `json:"AssociatedSalesOrders"`
}

type MaterialAssociatedSalesOrder struct {
	SalesOrderID string `json:"SalesOrderID"`
	Owner        string `json:"Owner"`
}

type MaterialProductionOrder struct {
	ProductionOrderID string `json:"ProductionOrderID"`
	Owner             string `json:"Owner"`
}

type MaterialBatches struct {
	BatchNumber string `json:"BatchNumber"`
	Owner       string `json:"Owner"`
}

//Define the Purchase Order structure, with XXX properties.
//Structure tags are used by encoding/json library.
type PurchaseOrder struct {
	Asset_Type      string                  `json:"Asset_Type, omitempty"`
	PurchaseOrderID string                  `json:"PurchaseOrderID"`
	Owner           string                  `json:"Owner"`
	Vendor          string                  `json:"Vendor"`
	LineItems       []PurchaseOrderLineItem `json:"LineItems"`
	Status          string                  `json:"Status"`
	TargetBatch     string                  `json:"TargetBatch, omitempty"`
}

type PurchaseOrderLineItem struct {
	LineItemNumber string `json:"LineItemNumber"`
	MaterialID     string `json:"MaterialID"`
	Quantity       int    `json:"Quantity"`
}

//Define the Sales Order structure, with XXX properties.
//Structure tags are used by encoding/json library.
type SalesOrder struct {
	Asset_Type     string               `json:"Asset_Type, omitempty"`
	SalesOrderID   string               `json:"SalesOrderID"`
	Owner          string               `json:"Owner"`
	POReference    string               `json:"POReference"`
	LineItems      []SalesOrderLineItem `json:"LineItems"`
	DeliveryNumber string               `json:"DeliveryNumber, omitempty"`
	Status         string               `json:"Status"`
}

type SalesOrderLineItem struct {
	LineItemNumber string `json:"LineItemNumber"`
	MaterialID     string `json:"MaterialID"`
	Quantity       int    `json:"Quantity"`
}

//Define the Batch structure, with XXX properties.
//Structure tags are used by encoding/json library.
type Batch struct {
	Asset_Type        string              `json:"Asset_Type, omitempty"`
	BatchNumber       string              `json:"BatchNumber"`
	MaterialID        string              `json:"MaterialID"`
	Owner             string              `json:"Owner"`
	Plant             string              `json:"Plant"`
	StorageLocation   string              `json:"StorageLocation"`
	AvailableQuantity int                 `json:"AvailableQuantity"`
	HandlingUnits     []BatchHandlingUnit `json:"HandlingUnits, omitempty"`
	Status            string              `json:"Status"`
}

type BatchHandlingUnit struct {
	HUID           string `json:"HUID"`
	Quantity       int    `json:"Quantity"`
	DeliveryNumber string `json:"DeliveryNumber, omitempty"`
}

//Define the Production Order structure, with XXX properties.
//Structure tags are used by encoding/json library.
type ProductionOrder struct {
	Asset_Type  string `json:"Asset_Type, omitempty"`
	MaterialID  string `json:"MaterialID"`
	Owner       string `json:"Owner"`
	Quantity    int    `json:"Quantity"`
	TargetBatch string `json:"TargetBatch"`
}

//Define the Delivery Document structure, with XXX properties.
//Structure tags are used by encoding/json library.
type DeliveryDocument struct {
	Asset_Type     string                `json:"Asset_Type, omitempty"`
	DeliveryNumber string                `json:"DeliveryNumber"`
	SalesOrderID   string                `json:"SalesOrderID"`
	Owner          string                `json:"Owner"`
	LineItems      []DeliveryDocLineItem `json:"LineItems"`
	Shipments      []DeliveryDocShipment `json:"Shipments, omitempty"`
}

type DeliveryDocLineItem struct {
	LineItemNumber string `json:"LineItemNumber"`
	MaterialID     string `json:"MaterialID"`
	Quantity       int    `json:"Quantity"`
	SourceBatch    string `json:"SourceBatch"`
}

type DeliveryDocShipment struct {
	ShipmentID string `json:"ShipmentID"`
}

//Define the Shipment structure, with XXX properties.
//Structure tags are used by encoding/json library.
type Shipment struct {
	Asset_Type     string                  `json:"Asset_Type, omitempty"`
	ShipmentID     string                  `json:"ShipmentID"`
	Owner          string                  `json:"Owner"`
	DeliveryNumber string                  `json:"DeliveryNumber"`
	Status         string                  `json:"Status"`
	SensorReadings []ShipmentSensorReading `json:"SensorReadings, omitempty"`
}

type ShipmentSensorReading struct {
	TempCelcius string `json:"TempCelcius"`
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
	if puterr := stub.PutState(strings.ToLower(participantID), jsonBytes); puterr != nil {
		return shim.Error("Invoke Error (Create Participant): Error while storing data into Blockchain")
	}
	return shim.Success(nil)
}

// CASE 02 Get a Participant Info
func (t *Testing1) getParticipant(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 2 {
		return shim.Error("Invoke Error: Incorrect number of arguments - Two Argument expected")
	}
	participantID := string(args[0])

	//Get the Asset from Blockchain
	value, geterr := stub.GetState(strings.ToLower(participantID))
	if geterr != nil || value == nil {
		return shim.Error("Invoke Error (Get Participant): Error while fetching data from Blockchain")
	}
	return shim.Success(value)
}

// CASE 03 Delete a Participant
func (t *Testing1) deleteParticipant(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 2 {
		return shim.Error("Invoke Error: Incorrect number of arguments - Two Argument expected")
	}
	participantID := string(args[0])

	// Check if Exists
	value, geterr := stub.GetState(strings.ToLower(participantID))
	if geterr != nil || value == nil {
		return shim.Error("Invoke Error (Delete Participant): Error while fetching data from Blockchain")
	}

	// Delete if Exists
	if delerr := stub.DelState(strings.ToLower(participantID)); delerr != nil {
		return shim.Error("Invoke Error (Delete Participant): Error while deleting data from Blockchain")
	}
	return shim.Success(nil)
}
