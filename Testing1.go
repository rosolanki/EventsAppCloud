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
	AssociatedSalesOrders []MaterialAssociatedSalesOrder `json:"AssociatedSalesOrders, omitempty"`
	Deleted               bool                           `json:"Deleted"`
}

type MaterialClosedPurchaseOrder struct {
	PurchaseOrderID       string                         `json:"PurchaseOrderID"`
	Owner                 string                         `json:"Owner"`
	AssociatedSalesOrders []MaterialAssociatedSalesOrder `json:"AssociatedSalesOrders, omitempty"`
	Deleted               bool                           `json:"Deleted"`
}

type MaterialAssociatedSalesOrder struct {
	SalesOrderID string `json:"SalesOrderID"`
	Owner        string `json:"Owner"`
	Deleted      bool   `json:"Deleted"`
}

type MaterialProductionOrder struct {
	ProductionOrderID string `json:"ProductionOrderID"`
	Owner             string `json:"Owner"`
	Deleted           bool   `json:"Deleted"`
}

type MaterialBatches struct {
	BatchNumber string `json:"BatchNumber"`
	Owner       string `json:"Owner"`
	Deleted     bool   `json:"Deleted"`
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
	POOwner        string               `json:"POOwner, omitempty"`
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
	Asset_Type        string `json:"Asset_Type, omitempty"`
	ProductionOrderID string `json:"ProductionOrderID"`
	MaterialID        string `json:"MaterialID"`
	Owner             string `json:"Owner"`
	Quantity          int    `json:"Quantity"`
	TargetBatch       string `json:"TargetBatch"`
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
	case "createPurchaseOrder":
		return t.createPurchaseOrder(stub, args)
	case "getPurchaseOrder":
		return t.getPurchaseOrder(stub, args)
	case "deletePurchaseOrder":
		return t.deletePurchaseOrder(stub, args)
	case "reportProductionOrderGR":
		return t.reportProductionOrderGR(stub, args)
	case "getProductionOrder":
		return t.getProductionOrder(stub, args)
	case "deleteProductionOrder":
		return t.deleteProductionOrder(stub, args)
	case "getBatch":
		return t.getBatch(stub, args)
	case "deleteBatch":
		return t.deleteBatch(stub, args)
	case "createSalesOrder":
		return t.createSalesOrder(stub, args)
	case "getSalesOrder":
		return t.getSalesOrder(stub, args)
	case "deleteSalesOrder":
		return t.deleteSalesOrder(stub, args)
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

	//Get Data
	data := string(args[0])
	queryData := QueryData{}
	err := json.Unmarshal([]byte(data), &queryData)
	if err != nil {
		return shim.Error("Invoke Error (Create Participant):  Invalid Data - Check Payload")
	}
	//Get Invoking Participant
	participantID := string(args[1])
	//Define Namespace
	namespace := "PARTICIPANT"

	//Define a new Participant
	participant := Participant{}
	participant.Asset_Type = namespace
	participant.ParticipantID = participantID
	participant.ParticipantType = queryData.ParticipantType
	participant.OrgName = queryData.OrgName
	participant.Email = queryData.Email

	//Key for fetching/storing the Asset
	keystring := namespace + "-" + participantID

	//Check if Participant already exists.
	if value, geterr := stub.GetState(strings.ToLower(keystring)); !(geterr == nil && value == nil) {
		return shim.Error("Invoke Error (Create Participant): Participant Already Exists! Please Specify Another ID")
	}

	//Store Participant in Blockchain
	jsonBytes, _ := json.Marshal(participant) //Get Bytes from struct
	if puterr := stub.PutState(strings.ToLower(keystring), jsonBytes); puterr != nil {
		return shim.Error("Invoke Error (Create Participant): Error while storing data into Blockchain")
	}
	return shim.Success(nil)
}

// CASE 02 Get a Participant Info
func (t *Testing1) getParticipant(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//Checks appropriate number of arguments in incoming invoke request
	if len(args) < 2 {
		return shim.Error("Invoke Error: Incorrect number of arguments - Two Argument expected")
	}
	//Get Data
	data := string(args[0])
	//Get Invoking Participant
	participantID := string(args[1])
	//Define Namespace
	namespace := "PARTICIPANT"

	//Key for fetching/storing the Asset
	keystring := namespace + "-" + data

	//Check if Invoking Participant already exists, return error if not.
	if value, geterr := stub.GetState(strings.ToLower(participantID)); geterr != nil || value == nil {
		return shim.Error("Invoke Error (Get Participant): Invoking Participant Does Not Exists! Please Enroll Participant")
	}

	//Get the Asset from Blockchain
	value, geterr := stub.GetState(strings.ToLower(keystring))
	if geterr != nil || value == nil {
		return shim.Error("Invoke Error (Get Participant): Error while fetching data from Blockchain")
	}
	return shim.Success(value)
}

// CASE 03 Delete a Participant
func (t *Testing1) deleteParticipant(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//Checks appropriate number of arguments in incoming invoke request
	if len(args) < 2 {
		return shim.Error("Invoke Error: Incorrect number of arguments - Two Argument expected")
	}
	//Get Data
	data := string(args[0])
	//Get Invoking Participant
	participantID := string(args[1])
	//Define Namespace
	namespace := "PARTICIPANT"

	//Key for fetching/storing the Asset
	keystring := namespace + "-" + data

	//Check if Invoking Participant already exists, return error if not.
	if value, geterr := stub.GetState(strings.ToLower(participantID)); geterr != nil || value == nil {
		return shim.Error("Invoke Error (Delete Participant): Invoking Participant Does Not Exists! Please Enroll Participant")
	}

	//Check if Asset exists and get the Asset.
	value, geterr := stub.GetState(strings.ToLower(keystring))
	if geterr != nil || value == nil {
		return shim.Error("Invoke Error (Delete Participant): Error while fetching data from Blockchain")
	}
	participant := Participant{}
	json.Unmarshal(value, &participant)

	//Check if Invoking Participant is authorised for Delete
	if strings.ToLower(participantID) == strings.ToLower(participant.ParticipantID) {
		// Delete if Exists
		if delerr := stub.DelState(strings.ToLower(keystring)); delerr != nil {
			return shim.Error("Invoke Error (Delete Participant): Error while deleting data from Blockchain")
		}
		return shim.Success(nil)
	} else {
		return shim.Error("Invoke Error (Delete Participant): Not Authorized to Delete Participant")
	}
}

// CASE 04 Create a Purchase Order
func (t *Testing1) createPurchaseOrder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//Checks appropriate number of arguments in incoming invoke request
	if len(args) < 2 {
		return shim.Error("Invoke Error: Incorrect number of arguments - Two Argument expected")
	}

	//Define the structure for expected incoming JSON as argument
	type QueryData struct {
		PurchaseOrderID string `json:"PurchaseOrderID"`
		Vendor          string `json:"Vendor"`
		LineItemNumber  string `json:"LineItemNumber"`
		MaterialID      string `json:"MaterialID"`
		Quantity        int    `json:"Quantity"`
	}

	//Get Data
	data := string(args[0])
	queryData := QueryData{}
	err := json.Unmarshal([]byte(data), &queryData)
	if err != nil {
		return shim.Error("Invoke Error (Create Purchase Order):  Invalid Data - Check Payload")
	}
	//Get Invoking Participant
	participantID := string(args[1])
	//Define Namespace
	namespace := "PURCHASEORDER"

	//One Assets will be created:
	//(1) Purchase Order
	//One Asset will be updated:
	//(1) Material

	//Define a new Purchase Order and Line Item.
	//Line Item
	poLineItem := PurchaseOrderLineItem{}
	poLineItem.LineItemNumber = queryData.LineItemNumber
	poLineItem.MaterialID = queryData.MaterialID
	poLineItem.Quantity = queryData.Quantity
	//Purchase Order
	purchaseOrder := PurchaseOrder{}
	purchaseOrder.Asset_Type = namespace
	purchaseOrder.PurchaseOrderID = queryData.PurchaseOrderID
	purchaseOrder.Owner = participantID
	purchaseOrder.Vendor = queryData.Vendor
	purchaseOrder.Status = "OPEN"
	purchaseOrder.LineItems = append(purchaseOrder.LineItems, poLineItem)

	//Key for fetching/storing the Asset
	keystring := namespace + "-" + participantID + "-" + queryData.PurchaseOrderID

	//Check if Invoking Participant already exists, return error if not.
	if value, geterr := stub.GetState(strings.ToLower(participantID)); geterr != nil || value == nil {
		return shim.Error("Invoke Error (Create Purchase Order): Invoking Participant Does Not Exists! Please Enroll Participant")
	}

	//Check If PO already exists.
	if value, geterr := stub.GetState(strings.ToLower(keystring)); !(geterr == nil && value == nil) {
		return shim.Error("Invoke Error (Create Purchase Order): PO Already Exists! Please Specify Another ID")
	}

	//**************************************************
	//Updating new PO information inside Material Asset
	//**************************************************
	//Check If Material Exists and get Material, else create new Material
	matNamespace := "MATERIAL"
	matkeystring := matNamespace + "-" + queryData.MaterialID
	matValue, matGetErr := stub.GetState(strings.ToLower(matkeystring))
	if matGetErr != nil || matValue == nil {
		//If Material does not exist,
		//Create New Material
		material := Material{}
		material.Asset_Type = matNamespace
		material.MaterialID = queryData.MaterialID

		//Define new Open Purchase Order information for Material
		matOpenPO := MaterialOpenPurchaseOrder{}
		matOpenPO.PurchaseOrderID = queryData.PurchaseOrderID
		matOpenPO.Owner = participantID
		matOpenPO.Deleted = false
		//Update Material with new Open PO information
		material.OpenPurchaseOrders = append(material.OpenPurchaseOrders, matOpenPO)

		// Store Material in Blockchain
		matJsonBytes, _ := json.Marshal(material) //Get Bytes from struct
		if puterr := stub.PutState(strings.ToLower(matkeystring), matJsonBytes); puterr != nil {
			return shim.Error("Invoke Error (Create PO - Create Material): Error while storing data into Blockchain")
		}
		// Store Purchase Order in Blockchain
		jsonBytes, _ := json.Marshal(purchaseOrder) //Get Bytes from struct
		if puterr := stub.PutState(strings.ToLower(keystring), jsonBytes); puterr != nil {
			return shim.Error("Invoke Error (Create Purchase Order): Error while storing data into Blockchain")
		}
		return shim.Success(nil)
	} else {
		//If Material already exists,
		//Fetch Material from Blockchain
		material := Material{}
		json.Unmarshal(matValue, &material)

		//Check if PO is Present in Material
		// presentFlag := false
		// for _, element := range material.OpenPurchaseOrders {
		// 	if (strings.ToLower(element.PurchaseOrderID) == strings.ToLower(queryData.PurchaseOrderID)) && (strings.ToLower(element.Owner) == strings.ToLower(participantID)) && (element.Deleted == false) {
		// 		presentFlag = true
		// 		break
		// 	}
		// }
		// if presentFlag == true {
		// 	// Store Purchase Order in Blockchain
		// 	jsonBytes, _ := json.Marshal(purchaseOrder) //Get Bytes from struct
		// 	if puterr := stub.PutState(strings.ToLower(keystring), jsonBytes); puterr != nil {
		// 		return shim.Error("Invoke Error (Create Purchase Order): Error while storing data into Blockchain")
		// 	}
		// 	return shim.Success(nil)
		// }

		//Define new Open Purchase Order information for Material
		matOpenPO := MaterialOpenPurchaseOrder{}
		matOpenPO.PurchaseOrderID = queryData.PurchaseOrderID
		matOpenPO.Owner = participantID
		matOpenPO.Deleted = false
		//Update Material with new Open PO information
		material.OpenPurchaseOrders = append(material.OpenPurchaseOrders, matOpenPO)

		// Store Material in Blockchain
		matJsonBytes, _ := json.Marshal(material) //Get Bytes from struct
		if puterr := stub.PutState(strings.ToLower(matkeystring), matJsonBytes); puterr != nil {
			return shim.Error("Invoke Error (Create PO - Update Material): Error while storing data into Blockchain")
		}
		// Store Purchase Order in Blockchain
		jsonBytes, _ := json.Marshal(purchaseOrder) //Get Bytes from struct
		if puterr := stub.PutState(strings.ToLower(keystring), jsonBytes); puterr != nil {
			return shim.Error("Invoke Error (Create Purchase Order): Error while storing data into Blockchain")
		}
		return shim.Success(nil)
	}
}

// CASE 05 Get a Purchase Order Info
func (t *Testing1) getPurchaseOrder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//Checks appropriate number of arguments in incoming invoke request
	if len(args) < 2 {
		return shim.Error("Invoke Error: Incorrect number of arguments - Two Argument expected")
	}

	//Define the structure for expected incoming JSON as argument
	type QueryData struct {
		Owner           string `json:"Owner"`
		PurchaseOrderID string `json:"PurchaseOrderID"`
	}

	//Get Data
	data := string(args[0])
	queryData := QueryData{}
	err := json.Unmarshal([]byte(data), &queryData)
	if err != nil {
		return shim.Error("Invoke Error (Get Purchase Order):  Invalid Data - Check Payload")
	}
	//Get Invoking Participant
	participantID := string(args[1])
	//Define Namespace
	namespace := "PURCHASEORDER"

	//Key for fetching/storing the Asset
	keystring := namespace + "-" + queryData.Owner + "-" + queryData.PurchaseOrderID

	//Check if Invoking Participant already exists, return error if not.
	if value, geterr := stub.GetState(strings.ToLower(participantID)); geterr != nil || value == nil {
		return shim.Error("Invoke Error (Get Purchase Order): Invoking Participant Does Not Exists! Please Enroll Participant")
	}

	//Get the Asset from Blockchain
	value, geterr := stub.GetState(strings.ToLower(keystring))
	if geterr != nil || value == nil {
		return shim.Error("Invoke Error (Get Purchase Order): Error while fetching data from Blockchain")
	}
	return shim.Success(value)
}

// CASE 06 Delete a Purchase Order
func (t *Testing1) deletePurchaseOrder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//Checks appropriate number of arguments in incoming invoke request
	if len(args) < 2 {
		return shim.Error("Invoke Error: Incorrect number of arguments - Two Argument expected")
	}

	//Define the structure for expected incoming JSON as argument
	type QueryData struct {
		Owner           string `json:"Owner"`
		PurchaseOrderID string `json:"PurchaseOrderID"`
	}

	//Get Data
	data := string(args[0])
	queryData := QueryData{}
	err := json.Unmarshal([]byte(data), &queryData)
	if err != nil {
		return shim.Error("Invoke Error (Delete Purchase Order):  Invalid Data - Check Payload")
	}
	//Get Invoking Participant
	participantID := string(args[1])
	//Define Namespace
	namespace := "PURCHASEORDER"

	//Key for fetching/storing the Asset
	keystring := namespace + "-" + queryData.Owner + "-" + queryData.PurchaseOrderID

	//Check if Invoking Participant already exists, return error if not.
	if value, geterr := stub.GetState(strings.ToLower(participantID)); geterr != nil || value == nil {
		return shim.Error("Invoke Error (Delete Purchase Order): Invoking Participant Does Not Exists! Please Enroll Participant")
	}

	//Check if Asset exists and get the Asset.
	value, geterr := stub.GetState(strings.ToLower(keystring))
	if geterr != nil || value == nil {
		return shim.Error("Invoke Error (Delete Purchase Order): Error while fetching data from Blockchain")
	}
	purchaseOrder := PurchaseOrder{}
	json.Unmarshal(value, &purchaseOrder)

	//Check if Invoking Participant is authorised for Delete
	if strings.ToLower(participantID) == strings.ToLower(purchaseOrder.Owner) {
		// Delete if Exists
		if delerr := stub.DelState(strings.ToLower(keystring)); delerr != nil {
			return shim.Error("Invoke Error (Delete Purchase Order): Error while deleting data from Blockchain")
		}
		return shim.Success(nil)
	} else {
		return shim.Error("Invoke Error (Delete Purchase Order): Not Authorized to Delete Purchase Order")
	}
}

// CASE 07 Report a Production Order Goods Receipt
func (t *Testing1) reportProductionOrderGR(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//Checks appropriate number of arguments in incoming invoke request
	if len(args) < 2 {
		return shim.Error("Invoke Error: Incorrect number of arguments - Two Argument expected")
	}

	//Define the structure for expected incoming JSON as argument
	type QueryData struct {
		ProductionOrderID string `json:"ProductionOrderID"`
		MaterialID        string `json:"MaterialID"`
		Quantity          int    `json:"Quantity"`
		Plant             string `json:"Plant"`
		StorageLocation   string `json:"StorageLocation"`
		BatchNumber       string `json:"BatchNumber"`
	}

	//Get Data
	data := string(args[0])
	queryData := QueryData{}
	err := json.Unmarshal([]byte(data), &queryData)
	if err != nil {
		return shim.Error("Invoke Error (GR Production Order):  Invalid Data - Check Payload")
	}
	//Get Invoking Participant
	participantID := string(args[1])
	//Define Namespace
	productionOrderNamespace := "PRODUCTIONORDER"
	batchNamespace := "BATCH"

	//One Assets will be created:
	//(1) Production Order
	//Two Asset will be updated:
	//(1) Batch
	//(2) Material

	//Define a new Production Order.
	productionOrder := ProductionOrder{}
	productionOrder.Asset_Type = productionOrderNamespace
	productionOrder.ProductionOrderID = queryData.ProductionOrderID
	productionOrder.MaterialID = queryData.MaterialID
	productionOrder.Owner = participantID
	productionOrder.Quantity = queryData.Quantity
	productionOrder.TargetBatch = queryData.BatchNumber

	//Key for fetching/storing the Asset
	productionorderkeystring := productionOrderNamespace + "-" + participantID + "-" + queryData.ProductionOrderID
	batchkeystring := batchNamespace + "-" + participantID + "-" + queryData.MaterialID + "-" + queryData.BatchNumber

	//Check if Invoking Participant already exists, return error if not.
	if value, geterr := stub.GetState(strings.ToLower(participantID)); geterr != nil || value == nil {
		return shim.Error("Invoke Error (GR Production Order): Invoking Participant Does Not Exists! Please Enroll Participant")
	}

	//Check If Production Order already exists.
	if value, geterr := stub.GetState(strings.ToLower(productionorderkeystring)); !(geterr == nil && value == nil) {
		return shim.Error("Invoke Error (GR Production Order): Production Order Already Exists! Please Specify Another ID")
	}

	// Check If Batch already Exists and get the Batch to update quantity, else create a new Batch.
	batchValue, batchGetErr := stub.GetState(strings.ToLower(batchkeystring))
	batch := Batch{}
	if batchGetErr != nil || batchValue == nil {
		//If Batch does not exist,
		//Create New Batch
		batch.Asset_Type = batchNamespace
		batch.BatchNumber = queryData.BatchNumber
		batch.MaterialID = queryData.MaterialID
		batch.Owner = participantID
		batch.Plant = queryData.Plant
		batch.StorageLocation = queryData.StorageLocation
		batch.AvailableQuantity = queryData.Quantity
		batch.Status = "OK"
	} else {
		//If Batch exists,
		//Get the Batch
		json.Unmarshal(batchValue, &batch)
		batch.AvailableQuantity += queryData.Quantity
		batch.Plant = queryData.Plant
		batch.StorageLocation = queryData.StorageLocation
	}

	//****************************************************************
	//Updating new Production Order information inside Material Asset
	//****************************************************************
	// Check If Material Exists and get Material, else create new Material
	matNamespace := "MATERIAL"
	matkeystring := matNamespace + "-" + queryData.MaterialID
	matValue, matGetErr := stub.GetState(strings.ToLower(matkeystring))
	if matGetErr != nil || matValue == nil {
		//If Material does not exist,
		//Create New Material
		material := Material{}
		material.Asset_Type = matNamespace
		material.MaterialID = queryData.MaterialID

		//Define new Material Production Order information for Material
		matProdOrder := MaterialProductionOrder{}
		matProdOrder.ProductionOrderID = queryData.ProductionOrderID
		matProdOrder.Owner = participantID
		matProdOrder.Deleted = false
		//Define new Material Batch information for Material
		matBatch := MaterialBatches{}
		matBatch.BatchNumber = queryData.BatchNumber
		matBatch.Owner = participantID
		matBatch.Deleted = false
		//Update Material with new Material Production Order and Batch information
		material.ProductionOrders = append(material.ProductionOrders, matProdOrder)
		material.ActiveBatches = append(material.ActiveBatches, matBatch)
		material.Batches = append(material.Batches, matBatch)

		// Store Material in Blockchain
		matJsonBytes, _ := json.Marshal(material) //Get Bytes from struct
		if puterr := stub.PutState(strings.ToLower(matkeystring), matJsonBytes); puterr != nil {
			return shim.Error("Invoke Error (GR Production Order - Create Material): Error while storing data into Blockchain")
		}
		// Store Production Order in Blockchain
		jsonBytes, _ := json.Marshal(productionOrder) //Get Bytes from struct
		if puterr := stub.PutState(strings.ToLower(productionorderkeystring), jsonBytes); puterr != nil {
			return shim.Error("Invoke Error (GR Production Order): Error while storing data into Blockchain")
		}
		// Store Batch in Blockchain
		batchjsonBytes, _ := json.Marshal(batch) //Get Bytes from struct
		if puterr := stub.PutState(strings.ToLower(batchkeystring), batchjsonBytes); puterr != nil {
			return shim.Error("Invoke Error (GR Production Order - Batch): Error while storing data into Blockchain")
		}
		return shim.Success(nil)
	} else {
		//If Material already exists,
		//Fetch Material from Blockchain
		material := Material{}
		json.Unmarshal(matValue, &material)

		//Check if PO is Present in Material
		// presentFlag := false
		// for _, element := range material.ProductionOrders {
		// 	if (strings.ToLower(element.ProductionOrderID) == strings.ToLower(queryData.ProductionOrderID)) && (strings.ToLower(element.Owner) == strings.ToLower(participantID)) && (element.Deleted == false) {
		// 		presentFlag = true
		// 		break
		// 	}
		// }
		// if presentFlag == true {
		// 	// Store Production Order in Blockchain
		// 	jsonBytes, _ := json.Marshal(productionOrder) //Get Bytes from struct
		// 	if puterr := stub.PutState(strings.ToLower(productionorderkeystring), jsonBytes); puterr != nil {
		// 		return shim.Error("Invoke Error (GR Production Order): Error while storing data into Blockchain")
		// 	}
		// 	// Store Batch in Blockchain
		// 	batchjsonBytes, _ := json.Marshal(batch) //Get Bytes from struct
		// 	if puterr := stub.PutState(strings.ToLower(batchkeystring), batchjsonBytes); puterr != nil {
		// 		return shim.Error("Invoke Error (GR Production Order - Batch): Error while storing data into Blockchain")
		// 	}
		// 	return shim.Success(nil)
		// }

		//Define new Material Production Order information for Material
		matProdOrder := MaterialProductionOrder{}
		matProdOrder.ProductionOrderID = queryData.ProductionOrderID
		matProdOrder.Owner = participantID
		matProdOrder.Deleted = false
		//Define new Material Batch information for Material
		matBatch := MaterialBatches{}
		matBatch.BatchNumber = queryData.BatchNumber
		matBatch.Owner = participantID
		matBatch.Deleted = false
		//Update Material with new Material Production Order and Batch information
		material.ProductionOrders = append(material.ProductionOrders, matProdOrder)
		material.ActiveBatches = append(material.ActiveBatches, matBatch)
		material.Batches = append(material.Batches, matBatch)

		// Store Material in Blockchain
		matJsonBytes, _ := json.Marshal(material) //Get Bytes from struct
		if puterr := stub.PutState(strings.ToLower(matkeystring), matJsonBytes); puterr != nil {
			return shim.Error("Invoke Error (GR Production Order - Update Material): Error while storing data into Blockchain")
		}
		// Store Production Order in Blockchain
		jsonBytes, _ := json.Marshal(productionOrder) //Get Bytes from struct
		if puterr := stub.PutState(strings.ToLower(productionorderkeystring), jsonBytes); puterr != nil {
			return shim.Error("Invoke Error (GR Production Order): Error while storing data into Blockchain")
		}
		// Store Batch in Blockchain
		batchjsonBytes, _ := json.Marshal(batch) //Get Bytes from struct
		if puterr := stub.PutState(strings.ToLower(batchkeystring), batchjsonBytes); puterr != nil {
			return shim.Error("Invoke Error (GR Production Order - Batch): Error while storing data into Blockchain")
		}
		return shim.Success(nil)
	}
}

// CASE 08 Get Production Order Info
func (t *Testing1) getProductionOrder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//Checks appropriate number of arguments in incoming invoke request
	if len(args) < 2 {
		return shim.Error("Invoke Error: Incorrect number of arguments - Two Argument expected")
	}

	//Define the structure for expected incoming JSON as argument
	type QueryData struct {
		Owner             string `json:"Owner"`
		ProductionOrderID string `json:"ProductionOrderID"`
	}

	//Get Data
	data := string(args[0])
	queryData := QueryData{}
	err := json.Unmarshal([]byte(data), &queryData)
	if err != nil {
		return shim.Error("Invoke Error (Get Production Order):  Invalid Data - Check Payload")
	}
	//Get Invoking Participant
	participantID := string(args[1])
	//Define Namespace
	namespace := "PRODUCTIONORDER"

	//Key for fetching/storing the Asset
	keystring := namespace + "-" + queryData.Owner + "-" + queryData.ProductionOrderID

	//Check if Invoking Participant already exists, return error if not.
	if value, geterr := stub.GetState(strings.ToLower(participantID)); geterr != nil || value == nil {
		return shim.Error("Invoke Error (Get Production Order): Invoking Participant Does Not Exists! Please Enroll Participant")
	}

	//Get the Asset from Blockchain
	value, geterr := stub.GetState(strings.ToLower(keystring))
	if geterr != nil || value == nil {
		return shim.Error("Invoke Error (Get Production Order): Error while fetching data from Blockchain")
	}
	return shim.Success(value)
}

// CASE 09 Delete a Production Order
func (t *Testing1) deleteProductionOrder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//Checks appropriate number of arguments in incoming invoke request
	if len(args) < 2 {
		return shim.Error("Invoke Error: Incorrect number of arguments - Two Argument expected")
	}

	//Define the structure for expected incoming JSON as argument
	type QueryData struct {
		Owner             string `json:"Owner"`
		ProductionOrderID string `json:"ProductionOrderID"`
	}

	//Get Data
	data := string(args[0])
	queryData := QueryData{}
	err := json.Unmarshal([]byte(data), &queryData)
	if err != nil {
		return shim.Error("Invoke Error (Delete Production Order):  Invalid Data - Check Payload")
	}
	//Get Invoking Participant
	participantID := string(args[1])
	//Define Namespace
	namespace := "PRODUCTIONORDER"

	//Key for fetching/storing the Asset
	keystring := namespace + "-" + queryData.Owner + "-" + queryData.ProductionOrderID

	//Check if Invoking Participant already exists, return error if not.
	if value, geterr := stub.GetState(strings.ToLower(participantID)); geterr != nil || value == nil {
		return shim.Error("Invoke Error (Delete Production Order): Invoking Participant Does Not Exists! Please Enroll Participant")
	}

	//Check if Asset exists and get the Asset.
	value, geterr := stub.GetState(strings.ToLower(keystring))
	if geterr != nil || value == nil {
		return shim.Error("Invoke Error (Delete Production Order): Error while fetching data from Blockchain")
	}
	productionOrder := ProductionOrder{}
	json.Unmarshal(value, &productionOrder)

	//Check if Invoking Participant is authorised for Delete
	if strings.ToLower(participantID) == strings.ToLower(productionOrder.Owner) {
		// Delete if Exists
		if delerr := stub.DelState(strings.ToLower(keystring)); delerr != nil {
			return shim.Error("Invoke Error (Delete Production Order): Error while deleting data from Blockchain")
		}
		return shim.Success(nil)
	} else {
		return shim.Error("Invoke Error (Delete Production Order): Not Authorized to Delete Production Order")
	}
}

// CASE 10 Get Batch Info
func (t *Testing1) getBatch(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//Checks appropriate number of arguments in incoming invoke request
	if len(args) < 2 {
		return shim.Error("Invoke Error: Incorrect number of arguments - Two Argument expected")
	}

	//Define the structure for expected incoming JSON as argument
	type QueryData struct {
		Owner       string `json:"Owner"`
		MaterialID  string `json:"MaterialID"`
		BatchNumber string `json:"BatchNumber"`
	}

	//Get Data
	data := string(args[0])
	queryData := QueryData{}
	err := json.Unmarshal([]byte(data), &queryData)
	if err != nil {
		return shim.Error("Invoke Error (Get Batch):  Invalid Data - Check Payload")
	}
	//Get Invoking Participant
	participantID := string(args[1])
	//Define Namespace
	namespace := "BATCH"

	//Key for fetching/storing the Asset
	keystring := namespace + "-" + queryData.Owner + "-" + queryData.MaterialID + "-" + queryData.BatchNumber

	//Check if Invoking Participant already exists, return error if not.
	if value, geterr := stub.GetState(strings.ToLower(participantID)); geterr != nil || value == nil {
		return shim.Error("Invoke Error (Get Batch): Invoking Participant Does Not Exists! Please Enroll Participant")
	}

	//Get the Asset from Blockchain
	value, geterr := stub.GetState(strings.ToLower(keystring))
	if geterr != nil || value == nil {
		return shim.Error("Invoke Error (Get Batch): Error while fetching data from Blockchain")
	}
	return shim.Success(value)
}

// CASE 11 Delete a Batch
func (t *Testing1) deleteBatch(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//Checks appropriate number of arguments in incoming invoke request
	if len(args) < 2 {
		return shim.Error("Invoke Error: Incorrect number of arguments - Two Argument expected")
	}

	//Define the structure for expected incoming JSON as argument
	type QueryData struct {
		Owner       string `json:"Owner"`
		MaterialID  string `json:"MaterialID"`
		BatchNumber string `json:"BatchNumber"`
	}

	//Get Data
	data := string(args[0])
	queryData := QueryData{}
	err := json.Unmarshal([]byte(data), &queryData)
	if err != nil {
		return shim.Error("Invoke Error (Delete Batch):  Invalid Data - Check Payload")
	}
	//Get Invoking Participant
	participantID := string(args[1])
	//Define Namespace
	namespace := "BATCH"

	//Key for fetching/storing the Asset
	keystring := namespace + "-" + queryData.Owner + "-" + queryData.MaterialID + "-" + queryData.BatchNumber

	//Check if Invoking Participant already exists, return error if not.
	if value, geterr := stub.GetState(strings.ToLower(participantID)); geterr != nil || value == nil {
		return shim.Error("Invoke Error (Delete Batch): Invoking Participant Does Not Exists! Please Enroll Participant")
	}

	//Check if Asset exists and get the Asset.
	value, geterr := stub.GetState(strings.ToLower(keystring))
	if geterr != nil || value == nil {
		return shim.Error("Invoke Error (Delete Batch): Error while fetching data from Blockchain")
	}
	batch := Batch{}
	json.Unmarshal(value, &batch)

	//Check if Invoking Participant is authorised for Delete
	if strings.ToLower(participantID) == strings.ToLower(batch.Owner) {
		// Delete if Exists
		if delerr := stub.DelState(strings.ToLower(keystring)); delerr != nil {
			return shim.Error("Invoke Error (Delete Batch): Error while deleting data from Blockchain")
		}
		return shim.Success(nil)
	} else {
		return shim.Error("Invoke Error (Delete Batch): Not Authorized to Delete Batch")
	}
}

// CASE 12 Delete a Batch
func (t *Testing1) createSalesOrder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//Checks appropriate number of arguments in incoming invoke request
	if len(args) < 2 {
		return shim.Error("Invoke Error: Incorrect number of arguments - Two Argument expected")
	}

	//Define the structure for expected incoming JSON as argument
	type QueryData struct {
		SalesOrderID   string `json:"SalesOrderID"`
		POReference    string `json:"POReference"`
		LineItemNumber string `json:"LineItemNumber"`
		MaterialID     string `json:"MaterialID"`
		Quantity       int    `json:"Quantity"`
	}

	//Get Data
	data := string(args[0])
	queryData := QueryData{}
	err := json.Unmarshal([]byte(data), &queryData)
	if err != nil {
		return shim.Error("Invoke Error (Create Sales Order):  Invalid Data - Check Payload")
	}
	//Get Invoking Participant
	participantID := string(args[1])
	//Define Namespace
	namespace := "SALESORDER"

	//One Assets will be created:
	//(1) Sales Order
	//One Asset will be updated:
	//(1) Material

	//Define a new Purchase Order and Line Item.
	//Line Item
	salesOrderLineItem := SalesOrderLineItem{}
	salesOrderLineItem.LineItemNumber = queryData.LineItemNumber
	salesOrderLineItem.MaterialID = queryData.MaterialID
	salesOrderLineItem.Quantity = queryData.Quantity
	//Purchase Order
	salesOrder := SalesOrder{}
	salesOrder.Asset_Type = namespace
	salesOrder.SalesOrderID = queryData.SalesOrderID
	salesOrder.Owner = participantID
	salesOrder.POReference = queryData.POReference
	salesOrder.Status = "OPEN"
	salesOrder.LineItems = append(salesOrder.LineItems, salesOrderLineItem)

	//Key for fetching/storing the Asset
	keystring := namespace + "-" + participantID + "-" + queryData.SalesOrderID

	//Check if Invoking Participant already exists, return error if not.
	if value, geterr := stub.GetState(strings.ToLower(participantID)); geterr != nil || value == nil {
		return shim.Error("Invoke Error (Create Sales Order): Invoking Participant Does Not Exists! Please Enroll Participant")
	}

	//Check If Sales Order already exists.
	if value, geterr := stub.GetState(strings.ToLower(keystring)); !(geterr == nil && value == nil) {
		return shim.Error("Invoke Error (Create Sales Order): Sales Order Already Exists! Please Specify Another ID")
	}

	//***********************************************************
	//Updating new Sales Order information inside Material Asset
	//***********************************************************
	//Get Material
	matNamespace := "MATERIAL"
	matkeystring := matNamespace + "-" + queryData.MaterialID
	matValue, matGetErr := stub.GetState(strings.ToLower(matkeystring))
	if matGetErr != nil || matValue == nil {
		return shim.Error("Invoke Error (Create Sales Order): Material Does Not Exists! Please Check Payload")
	}
	material := Material{}
	json.Unmarshal(matValue, &material)

	//Add Sales Order information to Open PO in Material
	for index, element := range material.OpenPurchaseOrders {
		if strings.ToLower(element.PurchaseOrderID) == strings.ToLower(queryData.POReference) {
			//Define Material Sales Order Information
			materialAssociatedSalesOrder := MaterialAssociatedSalesOrder{}
			materialAssociatedSalesOrder.SalesOrderID = queryData.SalesOrderID
			materialAssociatedSalesOrder.Owner = participantID
			materialAssociatedSalesOrder.Deleted = false

			//Update Sales Order POOwner information
			salesOrder.POOwner = element.Owner

			//Update the OpenPO in Material with Sales Order Information
			element.AssociatedSalesOrders = append(element.AssociatedSalesOrders, materialAssociatedSalesOrder)
			material.OpenPurchaseOrders[index] = element
		} else {
			return shim.Error("Invoke Error (Create Sales Order): PO Does not exists in Material! Please Specify Open PO Reference!")
		}
	}

	// Store Material in Blockchain
	matJsonBytes, _ := json.Marshal(material) //Get Bytes from struct
	if puterr := stub.PutState(strings.ToLower(matkeystring), matJsonBytes); puterr != nil {
		return shim.Error("Invoke Error (Create Sales Order - Material Update): Error while storing data into Blockchain")
	}
	// Store Sales Order in Blockchain
	jsonBytes, _ := json.Marshal(salesOrder) //Get Bytes from struct
	if puterr := stub.PutState(strings.ToLower(keystring), jsonBytes); puterr != nil {
		return shim.Error("Invoke Error (Create Sales Order): Error while storing data into Blockchain")
	}
	return shim.Success(nil)
}

// CASE 13 Delete a Batch
func (t *Testing1) getSalesOrder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//Checks appropriate number of arguments in incoming invoke request
	if len(args) < 2 {
		return shim.Error("Invoke Error: Incorrect number of arguments - Two Argument expected")
	}

	//Define the structure for expected incoming JSON as argument
	type QueryData struct {
		Owner        string `json:"Owner"`
		SalesOrderID string `json:"SalesOrderID"`
	}

	//Get Data
	data := string(args[0])
	queryData := QueryData{}
	err := json.Unmarshal([]byte(data), &queryData)
	if err != nil {
		return shim.Error("Invoke Error (Get Sales Order):  Invalid Data - Check Payload")
	}
	//Get Invoking Participant
	participantID := string(args[1])
	//Define Namespace
	namespace := "SALESORDER"

	//Key for fetching/storing the Asset
	keystring := namespace + "-" + queryData.Owner + "-" + queryData.SalesOrderID

	//Check if Invoking Participant already exists, return error if not.
	if value, geterr := stub.GetState(strings.ToLower(participantID)); geterr != nil || value == nil {
		return shim.Error("Invoke Error (Get Sales Order): Invoking Participant Does Not Exists! Please Enroll Participant")
	}

	//Get the Asset from Blockchain
	value, geterr := stub.GetState(strings.ToLower(keystring))
	if geterr != nil || value == nil {
		return shim.Error("Invoke Error (Get Sales Order): Error while fetching data from Blockchain")
	}
	return shim.Success(value)
}

// CASE 14 Delete a Batch
func (t *Testing1) deleteSalesOrder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//Checks appropriate number of arguments in incoming invoke request
	if len(args) < 2 {
		return shim.Error("Invoke Error: Incorrect number of arguments - Two Argument expected")
	}

	//Define the structure for expected incoming JSON as argument
	type QueryData struct {
		Owner        string `json:"Owner"`
		SalesOrderID string `json:"SalesOrderID"`
	}

	//Get Data
	data := string(args[0])
	queryData := QueryData{}
	err := json.Unmarshal([]byte(data), &queryData)
	if err != nil {
		return shim.Error("Invoke Error (Delete Sales Order):  Invalid Data - Check Payload")
	}
	//Get Invoking Participant
	participantID := string(args[1])
	//Define Namespace
	namespace := "SALESORDER"

	//Key for fetching/storing the Asset
	keystring := namespace + "-" + queryData.Owner + "-" + queryData.SalesOrderID

	//Check if Invoking Participant already exists, return error if not.
	if value, geterr := stub.GetState(strings.ToLower(participantID)); geterr != nil || value == nil {
		return shim.Error("Invoke Error (Delete Sales Order): Invoking Participant Does Not Exists! Please Enroll Participant")
	}

	//Check if Asset exists and get the Asset.
	value, geterr := stub.GetState(strings.ToLower(keystring))
	if geterr != nil || value == nil {
		return shim.Error("Invoke Error (Delete Sales Order): Error while fetching data from Blockchain")
	}
	salesOrder := SalesOrder{}
	json.Unmarshal(value, &salesOrder)

	//Check if Invoking Participant is authorised for Delete
	if strings.ToLower(participantID) == strings.ToLower(salesOrder.Owner) {
		// Delete if Exists
		if delerr := stub.DelState(strings.ToLower(keystring)); delerr != nil {
			return shim.Error("Invoke Error (Delete Sales Order): Error while deleting data from Blockchain")
		}
		return shim.Success(nil)
	} else {
		return shim.Error("Invoke Error (Delete Sales Order): Not Authorized to Delete Sales Order")
	}
}
