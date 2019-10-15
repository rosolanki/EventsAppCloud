//Deloitte Consulting LLP.
//**************************** MUST BE USED FOR INTERNAL PURPOSE ONLY ************************************
//****FileName: Blockchain IoT Chaincode
//****Description: Chaincode logic for IoT Enabled Blockchain for Supply Chain of Perishable Goods
//****Author: Rom Solanki
//****Author Email: rosolanki@deloitte.com
//********************************************************************************************************

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

//Import libraries for use

//********************************************************************************************************
//Struct for blockchain
//********************************************************************************************************
type BlockchainIOT struct {
}

//********************************************************************************************************
//Struct for Participants, Assets and Transactions
//********************************************************************************************************

//********************
// ASSETS
//********************

type Product struct {
	Asset_Type         string            `json:"Asset_Type, omitempty"`
	ProductID          string            `json:"ProductID"`
	ProductType        string            `json:"ProductType"`
	TotalQuantity      int               `json:"TotalQuantity"`
	SupplyChainMembers []MaterialDetails `json:"SupplyChainMembers,omitempty"`
	AllMaterials       []string          `json:"AllMaterials, omitempty"`
	Mappings           []Mapping         `json:"Mappings, omitempty"`
	ReverseMappings    []ReverseMapping  `json:"ReverseMappings, omitempty"`
}

type MaterialDetails struct {
	ParticipantID        string `json:"ParticipantID"`
	ParticipantType      string `json:"ParticipantType"`
	MaterialID           string `json:"MaterialID"`
	IsCompromised        bool   `json:"IsCompromised"`
	PotentialCompromised bool   `json:"PotentialCompromised"`
}

type Mapping struct {
	From BatchTradeInfo   `json:"From"`
	To   []BatchTradeInfo `json:"To"`
}

type ReverseMapping struct {
	To   BatchTradeInfo   `json:"To"`
	From []BatchTradeInfo `json:"From"`
}

type BatchTradeInfo struct {
	// Different from BatchInfo in the sence it keeps only the count of exchanged Quantity
	ParticipantID        string   `json:"ParticipantID"`
	MaterialID           string   `json:"MaterialID"`
	BatchNumber          string   `json:"BatchNumber"`
	SerialNumbers        []string `json:"SerialNumbers, omitempty"`
	Quantity             int      `json:"Quantity"`
	IsCompromised        bool     `json:"IsCompromised"`
	PotentialCompromised bool     `json:"PotentialCompromised"`
}

type Material struct {
	Asset_Type          string      `json:"Asset_Type, omitempty"`
	MaterialID          string      `json:"MaterialID"`
	ParticipantID       string      `json:"ParticipantID"`
	MaterialMasterID    string      `json:"MaterialMasterID"`
	ProductBCID         string      `json:"ProductBCID"`
	MaterialDescription string      `json:"MaterialDescription, omitempty"`
	Plant               string      `json:"Plant, omitempty"`
	StorageLocation     string      `json:"StorageLocation, omitempty"`
	UnitOfMeasure       string      `json:"Unit, omitempty"`
	TotalQuantity       int         `json:"TotalQuantity"`
	Batches             []BatchInfo `json:"Batches, omitempty"`
}

type BatchInfo struct {
	ParticipantID        string   `json:"ParticipantID"`
	MaterialID           string   `json:"MaterialID"`
	BatchNumber          string   `json:"BatchNumber"`
	SerialNumbers        []string `json:"SerialNumbers, omitempty"`
	Quantity             int      `json:"Quantity"`
	IsCompromised        bool     `json:"IsCompromised"`
	PotentialCompromised bool     `json:"PotentialCompromised"`
}

type PurchaseOrder struct {
	Asset_Type          string `json:"Asset_Type, omitempty"`
	POID                string `json:"POID"`
	RequestorID         string `json:"RequestorID"`
	RequestorMaterialID string `json:"RequestorMaterialID"`
	VendorID            string `json:"VendorID"`
	VendorMaterialID    string `json:"VendorMaterialID"`
	VendorBatchNumber   string `json:"VendorBatchNumber"`
	Quantity            int    `json:"Quantity"`
	UnitOfMeasure       string `json:"UnitOfMeasure"`
	NetPrice            int    `json:"NetPrice"`
	Currency            string `json:"Currency"`
	DeliveryDate        string `json:"DeliveryDate, omitempty"`
	TimeStamp           string `json:"TimeStamp, omitempty"`
	ShipmentExists      bool   `json:"ShipmentExists, omitempty"`
	ShipmentID          string `json:"ShipmentID, omitempty"`
	Status              string `json:"Status"`
}

type ProductionOrder struct {
	Asset_Type    string `json:"Asset_Type, omitempty"`
	POID          string `json:"POID"`
	ParticipantID string `json:"ParticipantID"`
	MaterialID    string `json:"MaterialID"`
	Quantity      int    `json:"Quantity"`
	UnitOfMeasure string `json:"UnitOfMeasure"`
	TimeStamp     string `json:"TimeStamp, omitempty"`
	Status        string `json:"Status"`
}

type Shipment struct {
	Asset_Type  string          `json:"Asset_Type, omitempty"`
	ShipmentID  string          `json:"ShipmentID"`
	ProductBCID string          `json:"ProductBCID"`
	POID        string          `json:"POID"`
	GPSReading  []GetGPSReading `json:"GPSReading, omitempty"`
	Status      string          `json:"Status, omitempty"`
}

//********************
// PARTICIPANTS
//********************

type Participant struct {
	Asset_Type      string   `json:"Asset_Type, omitempty"`
	ParticipantID   string   `json:"ParticipantID"`
	ParticipantType string   `json:"ParticipantType"` //Valid types are: GROWER, IMPORTER, DISTRIBUTOR, RETAILERS
	Materials       []string `json:"Materials", omitempty`
	CompanyName     string   `json:"CompanyName"`
	ContactEmail    string   `json:"ContactEmail"`
}

//********************
// TRANSACTIONS
//********************

type GetGPSReading struct {
	ShipmentID string  `json:"ShipmentID"`
	Latitude   float64 `json:"Latitude"`
	Longitude  float64 `json:"Longitude"`
	Accuracy   float32 `json:"Accuracy"`
	Timestamp  string  `json:"Timestamp, omitempty"`
}

type GoodsReceipt struct {
	Asset_Type    string   `json:"Asset_Type, omitempty"`
	ReceivedBy    string   `json:"ReceivedBy"`
	GRNumber      string   `json:"GRNumber"`
	Against       string   `json:"Against"`
	POID          string   `json:"POID"`
	BatchNumber   string   `json:"BatchNumber"`
	SerialNumbers []string `json:"SerialNumbers, omitempty"`
}

type BatchContamination struct {
	ParticipantID string `json:"ParticipantID"`
	MaterialID    string `json:"MaterialID"`
	BatchNumber   string `json:"BatchNumber"`
}

//********************************************************************************************************
// JSON Marchal and Unmarshal functions for Assets, Participants and Trasactions
//********************************************************************************************************
func (product *Product) ProductJsonToStruct(input []byte) *Product {
	json.Unmarshal(input, product)
	return product
}

func (product *Product) ProductStructToJson() []byte {
	jsonbytes, _ := json.Marshal(product)
	return jsonbytes
}

//********************************************************************************************************
//REST API Calls Response Format
//********************************************************************************************************
func Success(rc int32, msg string, payload []byte) peer.Response {
	return peer.Response{
		Status:  rc,
		Message: msg,
		Payload: payload,
	}
}

func Error(rc int32, msg string) peer.Response {
	logger.Errorf("Error %d = %s", rc, msg)
	return peer.Response{
		Status:  rc,
		Message: msg,
	}
}

//********************************************************************************************************
// Main Function
//********************************************************************************************************
func main() {
	if err := shim.Start(new(BlockchainIOT)); err != nil {
		fmt.Printf("Error starting BlockchainIOT chaincode: %s", err)
	}
}

//********************************************************************************************************
// INIT, INVOKE And QUERY
//********************************************************************************************************
var logger = shim.NewLogger("chaincode")

func (t *BlockchainIOT) Init(stub shim.ChaincodeStubInterface) peer.Response {
	_, args := stub.GetFunctionAndParameters()
	if len(args) > 0 {
		return Error(http.StatusBadRequest, "Init Error: Incorrect number of arguments - NO ARGUMENT EXPECTED")
	}
	return Success(http.StatusOK, "OK", nil)
}

func (t *BlockchainIOT) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case "createParticipant":
		return t.createParticipant(stub, args)
	case "createProduct":
		return t.createProduct(stub, args)
	case "registerMaterial":
		return t.registerMaterial(stub, args)
	case "createProductionOrder":
		return t.createProductionOrder(stub, args)
	case "createPurchaseOrder":
		return t.createPurchaseOrder(stub, args)
	case "createShipment":
		return t.createShipment(stub, args)
	case "trackShipment":
		return t.trackShipment(stub, args)
	case "submitGoodsReceipt":
		return t.submitGoodsReceipt(stub, args)
	case "reportContamination":
		return t.reportContamination(stub, args)
	case "clearContamination":
		return t.clearContamination(stub, args)
	case "getMaterial":
		return t.getMaterial(stub, args)
	case "deleteMaterial":
		return t.deleteMaterial(stub, args)
	case "getAsset":
		return t.getAsset(stub, args)
	case "deleteAsset":
		return t.deleteAsset(stub, args)
	case "getHistory":
		return t.getHistory(stub, args)
	case "customQueries":
		return t.customQueries(stub, args)
	default:
		logger.Warningf("Invalid Function Call - Function '%s' does not exist", function)
		return Error(http.StatusNotImplemented, "Invalid Function Call")
	}
}

//********************************************************************************************************
// Functions Definition
//********************************************************************************************************

// CASE 01 Create a Participant
func (t *BlockchainIOT) createParticipant(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 1 {
		return Error(http.StatusBadRequest, "Invoke Error: Incorrect number of arguments - One Argument expected")
	}

	type QueryData struct {
		ParticipantID   string `json:"ParticipantID"`
		ParticipantType string `json:"ParticipantType"`
		CompanyName     string `json:"CompanyName"`
		ContactEmail    string `json:"ContactEmail"`
	}

	data := string(args[0])
	queryData := QueryData{}
	err := json.Unmarshal([]byte(data), &queryData)
	if err != nil {
		return Error(http.StatusBadRequest, "Invoke Error: Invalid Data - Check Payload")
	}

	participant := Participant{}
	participant.Asset_Type = "PARTICIPANT"

	participant.ParticipantID = queryData.ParticipantID
	participant.ParticipantType = queryData.ParticipantType
	participant.CompanyName = queryData.CompanyName
	participant.ContactEmail = queryData.ContactEmail

	// Check If Exists
	participantID := strings.ToLower(participant.ParticipantID)
	if value, geterr := stub.GetState(participantID); !(geterr == nil && value == nil) {
		return Error(http.StatusConflict, "Participant Already Exists! \n Please Specify Another ID")
	}

	// Check Participant Type
	// Valid Participant Types are:
	// 1. GROWER
	// 2. IMPORTER
	// 3. DISTRIBUTOR
	// 4. RETAILER
	participantType := strings.ToUpper(participant.ParticipantType)
	if participantType != "GROWER" && participantType != "IMPORTER" && participantType != "DISTRIBUTOR" && participantType != "RETAILER" {
		return Error(http.StatusBadRequest, "Invoke Error: Invalid Data - Participant Type must be one of the following: \n 1) GROWER \n 2) IMPORTER \n 3) DISTRIBUTOR \n 4) RETAILER")
	}

	// Store in Blockchain
	jsonBytes, _ := json.Marshal(participant)
	if puterr := stub.PutState(participantID, jsonBytes); puterr != nil {
		return Error(http.StatusInternalServerError, puterr.Error())
	}
	return Success(http.StatusCreated, "Participant Created", nil)
}

// CASE 02 Create a Product
func (t *BlockchainIOT) createProduct(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 1 {
		return Error(http.StatusBadRequest, "Invoke Error: Incorrect number of arguments - One Argument expected")
	}

	type QueryData struct {
		ProductID   string `json:"ProductID"`
		ProductType string `json:"ProductType"`
	}

	data := string(args[0])
	queryData := QueryData{}
	err := json.Unmarshal([]byte(data), &queryData)
	if err != nil {
		return Error(http.StatusBadRequest, "Invoke Error: Invalid Data - Check Payload")
	}

	product := Product{}
	product.Asset_Type = "PRODUCT"
	product.TotalQuantity = 0

	product.ProductID = queryData.ProductID
	product.ProductType = queryData.ProductType

	// Check If Exists
	productID := strings.ToLower(product.ProductID)
	if value, geterr := stub.GetState(productID); !(geterr == nil && value == nil) {
		return Error(http.StatusConflict, "Product Already Exists! \n Please Specify Another ID")
	}

	// Store in Blockchain
	jsonBytes, _ := json.Marshal(product)
	if puterr := stub.PutState(productID, jsonBytes); puterr != nil {
		return Error(http.StatusInternalServerError, puterr.Error())
	}
	return Success(http.StatusCreated, "Product Created", nil)
}

// CASE 03 Register a Material
func (t *BlockchainIOT) registerMaterial(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 1 {
		return Error(http.StatusBadRequest, "Invoke Error: Incorrect number of arguments - One Argument expected")
	}

	type QueryData struct {
		ParticipantID       string `json:"ParticipantID"`
		MaterialMasterID    string `json:"MaterialMasterID"`
		ProductBCID         string `json:"ProductBCID"`
		MaterialDescription string `json:"MaterialDescription"`
		Plant               string `json:"Plant"`
		StorageLocation     string `json:"StorageLocation"`
		UnitOfMeasure       string `json:"UnitOfMeasure"`
	}

	data := string(args[0])
	queryData := QueryData{}
	err := json.Unmarshal([]byte(data), &queryData)
	if err != nil {
		return Error(http.StatusBadRequest, "Invoke Error: Invalid Data - Check Payload")
	}

	material := Material{}
	material.Asset_Type = "MATERIAL"
	material.TotalQuantity = 0

	material.ParticipantID = queryData.ParticipantID
	material.MaterialMasterID = queryData.MaterialMasterID
	material.ProductBCID = queryData.ProductBCID
	material.MaterialDescription = queryData.MaterialDescription
	material.Plant = queryData.Plant
	material.StorageLocation = queryData.StorageLocation
	material.UnitOfMeasure = queryData.UnitOfMeasure
	material.MaterialID = material.ParticipantID + "-" + material.MaterialMasterID

	// Check If Exists
	materialID := strings.ToLower(material.MaterialID)
	if value, geterr := stub.GetState(materialID); !(geterr == nil && value == nil) {
		return Error(http.StatusConflict, "Material Already Exists! \n Please Specify Another ID")
	}

	// Check If Product Exists and Get the Product
	productValue, productGetErr := stub.GetState(strings.ToLower(material.ProductBCID))
	if productGetErr != nil || productValue == nil {
		return Error(http.StatusNotFound, "Product Does Not Exists! \n Please Specify Another Product ID")
	}
	product := Product{}
	json.Unmarshal(productValue, &product)

	// Check If Participant Exists and Get the Participant
	participantValue, participantGetErr := stub.GetState(strings.ToLower(material.ParticipantID))
	if participantGetErr != nil || participantValue == nil {
		return Error(http.StatusNotFound, "Participant Does Not Exists! \n Please Specify Another Participant ID")
	}
	participant := Participant{}
	json.Unmarshal(participantValue, &participant)

	for _, element := range participant.Materials {
		if strings.ToLower(element) == strings.ToLower(materialID) {
			return Error(http.StatusConflict, "Material Already Present with Participant!")
		}
	}

	for _, element := range product.AllMaterials {
		if strings.ToLower(element) == strings.ToLower(materialID) {
			return Error(http.StatusConflict, "Material Already Present with Product!")
		}
	}

	// Register the Material to Product and Participant
	product.AllMaterials = append(product.AllMaterials, material.MaterialID)
	participant.Materials = append(participant.Materials, material.MaterialID)

	// Store Product and Material to Blockchain
	productJsonBytes, _ := json.Marshal(product)
	if puterr := stub.PutState(strings.ToLower(product.ProductID), productJsonBytes); puterr != nil {
		return Error(http.StatusInternalServerError, puterr.Error())
	}
	participantJsonBytes, _ := json.Marshal(participant)
	if puterr := stub.PutState(strings.ToLower(participant.ParticipantID), participantJsonBytes); puterr != nil {
		return Error(http.StatusInternalServerError, puterr.Error())
	}
	materialJsonBytes, _ := json.Marshal(material)
	if puterr := stub.PutState(strings.ToLower(materialID), materialJsonBytes); puterr != nil {
		return Error(http.StatusInternalServerError, puterr.Error())
	}
	return Success(http.StatusCreated, "Material Registered", nil)
}

// CASE 04 Create Production Order
func (t *BlockchainIOT) createProductionOrder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 1 {
		return Error(http.StatusBadRequest, "Invoke Error: Incorrect number of arguments - One Argument expected")
	}

	type QueryData struct {
		POID          string `json:"POID"`
		ParticipantID string `json:"ParticipantID"`
		MaterialID    string `json:"MaterialID"`
		Quantity      int    `json:"Quantity"`
		UnitOfMeasure string `json:"UnitOfMeasure"`
	}

	data := string(args[0])
	queryData := QueryData{}
	err := json.Unmarshal([]byte(data), &queryData)
	if err != nil {
		return Error(http.StatusBadRequest, "Invoke Error: Invalid Data - Check Payload")
	}

	productionOrder := ProductionOrder{}
	productionOrder.Asset_Type = "PRODUCTION ORDER"
	productionOrder.Status = "OPEN"

	productionOrder.POID = queryData.POID
	productionOrder.ParticipantID = queryData.ParticipantID
	productionOrder.MaterialID = queryData.MaterialID
	productionOrder.Quantity = queryData.Quantity
	productionOrder.UnitOfMeasure = queryData.UnitOfMeasure

	// Check If Exists
	productionOrderID := strings.ToLower(productionOrder.POID)
	if value, geterr := stub.GetState(productionOrderID); !(geterr == nil && value == nil) {
		return Error(http.StatusConflict, "Production Order Already Exists! \n Please Specify Another ID")
	}

	// Check If Participant Exists
	participantValue, participantGetErr := stub.GetState(strings.ToLower(productionOrder.ParticipantID))
	if participantGetErr != nil || participantValue == nil {
		return Error(http.StatusNotFound, "Participant Does Not Exists! \n Please Specify Another Participant ID")
	}

	// Check If Material Exists
	materialID := productionOrder.ParticipantID + "-" + productionOrder.MaterialID
	materialValue, materialGetErr := stub.GetState(strings.ToLower(materialID))
	if materialGetErr != nil || materialValue == nil {
		return Error(http.StatusNotFound, "Material Does Not Exists! \n Please Specify Another Material ID")
	}

	// Store in Blockchain
	jsonBytes, _ := json.Marshal(productionOrder)
	if puterr := stub.PutState(productionOrderID, jsonBytes); puterr != nil {
		return Error(http.StatusInternalServerError, puterr.Error())
	}
	return Success(http.StatusCreated, "Production Order Created", nil)
}

// CASE 05 Create Purchase Order
func (t *BlockchainIOT) createPurchaseOrder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 1 {
		return Error(http.StatusBadRequest, "Invoke Error: Incorrect number of arguments - One Argument expected")
	}

	type QueryData struct {
		POID                string `json:"POID"`
		RequestorID         string `json:"RequestorID"`
		RequestorMaterialID string `json:"RequestorMaterialID"`
		VendorID            string `json:"VendorID"`
		VendorMaterialID    string `json:"VendorMaterialID"`
		VendorBatchNumber   string `json:"VendorBatchNumber"`
		Quantity            int    `json:"Quantity"`
		UnitOfMeasure       string `json:"UnitOfMeasure"`
		NetPrice            int    `json:"NetPrice"`
		Currency            string `json:"Currency"`
	}

	data := string(args[0])
	queryData := QueryData{}
	err := json.Unmarshal([]byte(data), &queryData)
	if err != nil {
		return Error(http.StatusBadRequest, "Invoke Error: Invalid Data - Check Payload")
	}

	purchaseOrder := PurchaseOrder{}
	purchaseOrder.Asset_Type = "PURCHASE ORDER"
	purchaseOrder.Status = "OPEN"

	purchaseOrder.POID = queryData.POID
	purchaseOrder.RequestorID = queryData.RequestorID
	purchaseOrder.RequestorMaterialID = queryData.RequestorMaterialID
	purchaseOrder.VendorID = queryData.VendorID
	purchaseOrder.VendorMaterialID = queryData.VendorMaterialID
	purchaseOrder.VendorBatchNumber = queryData.VendorBatchNumber
	purchaseOrder.Quantity = queryData.Quantity
	purchaseOrder.UnitOfMeasure = queryData.UnitOfMeasure
	purchaseOrder.NetPrice = queryData.NetPrice
	purchaseOrder.Currency = queryData.Currency

	// Check If Exists
	purchaseOrderID := strings.ToLower(purchaseOrder.POID)
	if value, geterr := stub.GetState(purchaseOrderID); !(geterr == nil && value == nil) {
		return Error(http.StatusConflict, "Purchase Order Already Exists! \n Please Specify Another ID")
	}

	// Check If Vendor and Requestor Exists
	vendorValue, vendorGetErr := stub.GetState(strings.ToLower(purchaseOrder.VendorID))
	if vendorGetErr != nil || vendorValue == nil {
		return Error(http.StatusNotFound, "Vendor Does Not Exists! \n Please Specify Another Vendor ID")
	}

	requestorValue, requestorGetErr := stub.GetState(strings.ToLower(purchaseOrder.RequestorID))
	if requestorGetErr != nil || requestorValue == nil {
		return Error(http.StatusNotFound, "Requestor Does Not Exists! \n Please Specify Another Requestor ID")
	}

	// Check If Material Exists for Vendor and Requestor
	vendorMaterialID := purchaseOrder.VendorID + "-" + purchaseOrder.VendorMaterialID
	vendorMaterialValue, vendorMaterialGetErr := stub.GetState(strings.ToLower(vendorMaterialID))
	if vendorMaterialGetErr != nil || vendorMaterialValue == nil {
		return Error(http.StatusNotFound, "Vendor Material Does Not Exists! \n Please Specify Another Vendor Material ID")
	}

	requestorMaterialID := purchaseOrder.RequestorID + "-" + purchaseOrder.RequestorMaterialID
	requestorMaterialValue, requestorMaterialGetErr := stub.GetState(strings.ToLower(requestorMaterialID))
	if requestorMaterialGetErr != nil || requestorMaterialValue == nil {
		return Error(http.StatusNotFound, "Requestor Material Does Not Exists! \n Please Specify Another Requestor Material ID")
	}

	// Check if Vendor Batch Exists
	vendorMaterial := Material{}
	json.Unmarshal(vendorMaterialValue, &vendorMaterial)

	batchExists := false
	for _, element := range vendorMaterial.Batches {
		if strings.ToLower(element.BatchNumber) == strings.ToLower(purchaseOrder.VendorBatchNumber) {
			element.Quantity -= purchaseOrder.Quantity
			if element.Quantity < 0 {
				return Error(http.StatusBadRequest, "Not Enough Quantity Available in this Batch!")
			}
			batchExists = true
			break
		}
	}

	if batchExists == false {
		return Error(http.StatusNotFound, "Vendor Batch Not Found!")
	}

	// Store in Blockchain
	jsonBytes, _ := json.Marshal(purchaseOrder)
	if puterr := stub.PutState(purchaseOrderID, jsonBytes); puterr != nil {
		return Error(http.StatusInternalServerError, puterr.Error())
	}
	return Success(http.StatusCreated, "Production Order Created", nil)
}

// CASE 06 Create Shipment
func (t *BlockchainIOT) createShipment(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 1 {
		return Error(http.StatusBadRequest, "Invoke Error: Incorrect number of arguments - One Argument expected")
	}

	type QueryData struct {
		ShipmentID  string `json:"ShipmentID"`
		ProductBCID string `json:"ProductBCID"`
		POID        string `json:"POID"`
	}

	data := string(args[0])
	queryData := QueryData{}
	err := json.Unmarshal([]byte(data), &queryData)
	if err != nil {
		return Error(http.StatusBadRequest, "Invoke Error: Invalid Data - Check Payload")
	}

	shipment := Shipment{}
	shipment.Asset_Type = "SHIPMENT"
	shipment.Status = "SHIPPING"

	shipment.ShipmentID = queryData.ShipmentID
	shipment.ProductBCID = queryData.ProductBCID
	shipment.POID = queryData.POID

	// Check If Exists
	shipmentID := strings.ToLower(shipment.ShipmentID)
	if value, geterr := stub.GetState(shipmentID); !(geterr == nil && value == nil) {
		return Error(http.StatusConflict, "Shipment Already Exists! \n Please Specify Another ID")
	}

	// Check if Purchase Order Exists and get the Purchase Order
	POValue, POGetErr := stub.GetState(strings.ToLower(shipment.POID))
	if POGetErr != nil || POValue == nil {
		return Error(http.StatusNotFound, "Purchase Order Does Not Exists! \n Please Specify Another POID")
	}
	purchaseOrder := PurchaseOrder{}
	json.Unmarshal(POValue, &purchaseOrder)

	// Check if Purchase Order is Completed
	if purchaseOrder.Status == "COMPLETED" {
		return Error(http.StatusBadRequest, "Goods are Already Delivered for this Purchase Order")
	}

	// Check if Shipment Exists for this Purchase Order
	if purchaseOrder.ShipmentExists == true {
		return Error(http.StatusBadRequest, "Shipment Already Exists for this Purchase Order")
	}

	purchaseOrder.ShipmentExists = true
	purchaseOrder.ShipmentID = shipment.ShipmentID

	// Get Vendor Material
	vendorMaterialID := purchaseOrder.VendorID + "-" + purchaseOrder.VendorMaterialID
	vendorMaterialValue, _ := stub.GetState(strings.ToLower(vendorMaterialID))
	vendorMaterial := Material{}
	json.Unmarshal(vendorMaterialValue, &vendorMaterial)

	// Update Vendor Material
	// Reduce the Available Quantity from Vendor Material and its Batch
	vendorMaterial.TotalQuantity -= purchaseOrder.Quantity
	for index, element := range vendorMaterial.Batches {
		if strings.ToLower(element.BatchNumber) == strings.ToLower(purchaseOrder.VendorBatchNumber) {
			element.Quantity -= purchaseOrder.Quantity
			if element.Quantity < 0 {
				return Error(http.StatusBadRequest, "Not Enough Quantity Present in this Batch!")
			}
			vendorMaterial.Batches[index] = element
			break
		}
	}

	// Store in Blockchain
	shipmentJsonBytes, _ := json.Marshal(shipment)
	if puterr := stub.PutState(strings.ToLower(shipment.ShipmentID), shipmentJsonBytes); puterr != nil {
		return Error(http.StatusInternalServerError, puterr.Error())
	}
	purchaseOrderJsonBytes, _ := json.Marshal(purchaseOrder)
	if puterr := stub.PutState(strings.ToLower(purchaseOrder.POID), purchaseOrderJsonBytes); puterr != nil {
		return Error(http.StatusInternalServerError, puterr.Error())
	}
	vendorMaterialJsonBytes, _ := json.Marshal(vendorMaterial)
	if puterr := stub.PutState(strings.ToLower(vendorMaterialID), vendorMaterialJsonBytes); puterr != nil {
		return Error(http.StatusInternalServerError, puterr.Error())
	}
	return Success(http.StatusCreated, "Shipment Created", nil)
}

// CASE 07 Track Shipment
func (t *BlockchainIOT) trackShipment(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 1 {
		return Error(http.StatusBadRequest, "Invoke Error: Incorrect number of arguments - One Argument expected")
	}

	type QueryData struct {
		ShipmentID string  `json:"ShipmentID"`
		Latitude   float64 `json:"Latitude"`
		Longitude  float64 `json:"Longitude"`
		Accuracy   float32 `json:"Accuracy"`
		Timestamp  string  `json:"Timestamp, omitempty"`
	}

	data := string(args[0])
	queryData := QueryData{}
	err := json.Unmarshal([]byte(data), &queryData)
	if err != nil {
		return Error(http.StatusBadRequest, "Invoke Error: Invalid Data - Check Payload")
	}

	gpsReading := GetGPSReading{}

	gpsReading.ShipmentID = queryData.ShipmentID
	gpsReading.Latitude = queryData.Latitude
	gpsReading.Longitude = queryData.Longitude
	gpsReading.Accuracy = queryData.Accuracy
	gpsReading.Timestamp = queryData.Timestamp

	// Check if Shipment Exists and Get Shipment
	shipmentValue, shipmentGetErr := stub.GetState(strings.ToLower(gpsReading.ShipmentID))
	if shipmentGetErr != nil || shipmentValue == nil {
		return Error(http.StatusNotFound, "Shipment Does Not Exists! \n Please Specify Another Shipment ID")
	}
	shipment := Shipment{}
	json.Unmarshal(shipmentValue, &shipment)

	// Check if Shipment is Completed
	if shipment.Status == "COMPLETED" {
		return Error(http.StatusBadRequest, "Shipment is already Completed")
	}

	// Enter Location for Shipment
	shipment.GPSReading = append(shipment.GPSReading, gpsReading)

	// Store Updated Shipment in Blockchain
	shipmentJsonBytes, _ := json.Marshal(shipment)
	if puterr := stub.PutState(strings.ToLower(shipment.ShipmentID), shipmentJsonBytes); puterr != nil {
		return Error(http.StatusInternalServerError, puterr.Error())
	}
	return Success(http.StatusCreated, "Shipment Location Updated", nil)
}

// CASE 08 Submit Goods Receipt
func (t *BlockchainIOT) submitGoodsReceipt(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 1 {
		return Error(http.StatusBadRequest, "Invoke Error: Incorrect number of arguments - One Argument expected")
	}

	type QueryData struct {
		GRNumber      string   `json:"GRNumber"`
		ReceivedBy    string   `json:"ReceivedBy"`
		Against       string   `json:"Against"`
		POID          string   `json:"POID"`
		BatchNumber   string   `json:"BatchNumber"`
		SerialNumbers []string `json:"SerialNumbers, omitempty"`
	}

	data := string(args[0])
	queryData := QueryData{}
	err := json.Unmarshal([]byte(data), &queryData)
	if err != nil {
		return Error(http.StatusBadRequest, "Invoke Error: Invalid Data - Check Payload")
	}

	goodsReceipt := GoodsReceipt{}
	goodsReceipt.Asset_Type = "Goods Receipt"

	goodsReceipt.GRNumber = queryData.GRNumber
	goodsReceipt.ReceivedBy = queryData.ReceivedBy
	goodsReceipt.Against = queryData.Against
	goodsReceipt.POID = queryData.POID
	goodsReceipt.BatchNumber = queryData.BatchNumber
	goodsReceipt.SerialNumbers = queryData.SerialNumbers

	if strings.ToUpper(goodsReceipt.Against) == "PRODUCTION ORDER" {
		// Check If Production Order Exists and Get the Order
		POValue, POGetErr := stub.GetState(strings.ToLower(goodsReceipt.POID))
		if POGetErr != nil || POValue == nil {
			return Error(http.StatusNotFound, "Production Order Does Not Exists! \n Please Specify Another POID")
		}
		productionOrder := ProductionOrder{}
		json.Unmarshal(POValue, &productionOrder)

		// Check the Status of Order
		if productionOrder.Status == "COMPLETED" {
			return Error(http.StatusBadRequest, "Goods Already Received for this Production Order")
		}

		// Check for Valid Receiver
		if strings.ToLower(goodsReceipt.ReceivedBy) != strings.ToLower(productionOrder.ParticipantID) {
			return Error(http.StatusBadRequest, "Not a Valid Receiver for this Production Order")
		}

		// Get Material
		materialID := productionOrder.ParticipantID + "-" + productionOrder.MaterialID
		materialValue, _ := stub.GetState(strings.ToLower(materialID))
		material := Material{}
		json.Unmarshal(materialValue, &material)

		// Get Participant
		participantValue, _ := stub.GetState(strings.ToLower(material.ParticipantID))
		participant := Participant{}
		json.Unmarshal(participantValue, &participant)

		// Get Product
		productValue, _ := stub.GetState(strings.ToLower(material.ProductBCID))
		product := Product{}
		json.Unmarshal(productValue, &product)

		// Update Product
		// Step 1: Add Quantity to Total Product Quantity
		product.TotalQuantity += productionOrder.Quantity
		// Step 2: Add Participant, BatchInfo and Material Details
		batchInfo := BatchInfo{}
		batchInfo.ParticipantID = productionOrder.ParticipantID
		batchInfo.MaterialID = productionOrder.MaterialID
		batchInfo.BatchNumber = goodsReceipt.BatchNumber
		batchInfo.SerialNumbers = goodsReceipt.SerialNumbers
		batchInfo.Quantity = productionOrder.Quantity
		batchInfo.IsCompromised = false
		batchInfo.PotentialCompromised = false

		materialDetails := MaterialDetails{}
		materialDetails.ParticipantID = productionOrder.ParticipantID
		materialDetails.ParticipantType = participant.ParticipantType
		materialDetails.MaterialID = productionOrder.MaterialID
		materialDetails.IsCompromised = false
		materialDetails.PotentialCompromised = false

		participantExists := false
		for _, element := range product.SupplyChainMembers {
			if strings.ToLower(element.ParticipantID) == strings.ToLower(productionOrder.ParticipantID) {
				participantExists = true
				break
			}
		}
		if participantExists == false {
			product.SupplyChainMembers = append(product.SupplyChainMembers, materialDetails)
		}

		// Update Material
		// Step 1: Add Quantity to Material
		material.TotalQuantity += productionOrder.Quantity
		// Step 2: Update Batch
		materialBatchExists := false

		for index, element := range material.Batches {
			if strings.ToLower(element.BatchNumber) == strings.ToLower(goodsReceipt.BatchNumber) {
				element.Quantity += productionOrder.Quantity
				material.Batches[index] = element
				materialBatchExists = true
				break
			}
		}
		if materialBatchExists == false {
			material.Batches = append(material.Batches, batchInfo)
		}

		// Update Production Order Status
		productionOrder.Status = "COMPLETED"

		// Store Data into Blockchain (Update Product and Material)
		POjsonBytes, _ := json.Marshal(productionOrder)
		if puterr := stub.PutState(strings.ToLower(productionOrder.POID), POjsonBytes); puterr != nil {
			return Error(http.StatusInternalServerError, puterr.Error())
		}

		ProductjsonBytes, _ := json.Marshal(product)
		if puterr := stub.PutState(strings.ToLower(product.ProductID), ProductjsonBytes); puterr != nil {
			return Error(http.StatusInternalServerError, puterr.Error())
		}

		MaterialjsonBytes, _ := json.Marshal(material)
		if puterr := stub.PutState(strings.ToLower(materialID), MaterialjsonBytes); puterr != nil {
			return Error(http.StatusInternalServerError, puterr.Error())
		}

		GRjsonBytes, _ := json.Marshal(goodsReceipt)
		if puterr := stub.PutState(strings.ToLower(goodsReceipt.GRNumber), GRjsonBytes); puterr != nil {
			return Error(http.StatusInternalServerError, puterr.Error())
		}

		return Success(http.StatusCreated, "Goods Received Against Production Order", nil)

	} else if strings.ToUpper(goodsReceipt.Against) == "PURCHASE ORDER" {
		// Check If Purchase Order Exists and Get the Order
		POValue, POGetErr := stub.GetState(strings.ToLower(goodsReceipt.POID))
		if POGetErr != nil || POValue == nil {
			return Error(http.StatusNotFound, "Purchase Order Does Not Exists! \n Please Specify Another POID")
		}
		purchaseOrder := PurchaseOrder{}
		json.Unmarshal(POValue, &purchaseOrder)

		// Check the Status of Order
		if purchaseOrder.Status == "COMPLETED" {
			return Error(http.StatusBadRequest, "Goods Already Received for this Purchase Order")
		}

		// Check for Valid Receiver
		if strings.ToLower(goodsReceipt.ReceivedBy) != strings.ToLower(purchaseOrder.RequestorID) {
			return Error(http.StatusBadRequest, "Not a Valid Receiver for this Production Order")
		}

		// Get Materials
		vendorMaterialID := purchaseOrder.VendorID + "-" + purchaseOrder.VendorMaterialID
		vendorMaterialValue, _ := stub.GetState(strings.ToLower(vendorMaterialID))
		vendorMaterial := Material{}
		json.Unmarshal(vendorMaterialValue, &vendorMaterial)

		receiverMaterialID := purchaseOrder.RequestorID + "-" + purchaseOrder.RequestorMaterialID
		receiverMaterialValue, _ := stub.GetState(strings.ToLower(receiverMaterialID))
		receiverMaterial := Material{}
		json.Unmarshal(receiverMaterialValue, &receiverMaterial)

		// Get Participants
		vendorParticipantValue, _ := stub.GetState(strings.ToLower(purchaseOrder.VendorID))
		vendor := Participant{}
		json.Unmarshal(vendorParticipantValue, &vendor)

		receiverParticipantValue, _ := stub.GetState(strings.ToLower(purchaseOrder.RequestorID))
		receiver := Participant{}
		json.Unmarshal(receiverParticipantValue, &receiver)

		// Get Shipment
		shipmentValue, _ := stub.GetState(strings.ToLower(purchaseOrder.ShipmentID))
		shipment := Shipment{}
		json.Unmarshal(shipmentValue, &shipment)

		// Get Product
		productValue, _ := stub.GetState(strings.ToLower(receiverMaterial.ProductBCID))
		product := Product{}
		json.Unmarshal(productValue, &product)

		// Update Materials
		// Step 1: Update Receiver Material
		// Get Vendor Batch Info
		vendorbatchInfo := BatchInfo{}
		for _, element := range vendorMaterial.Batches {
			if strings.ToLower(element.BatchNumber) == strings.ToLower(purchaseOrder.VendorBatchNumber) {
				vendorbatchInfo = element
			}
		}

		// Create Receiver Batch Info
		receiverbatchInfo := BatchInfo{}
		receiverbatchInfo.ParticipantID = purchaseOrder.RequestorID
		receiverbatchInfo.MaterialID = purchaseOrder.RequestorMaterialID
		receiverbatchInfo.BatchNumber = goodsReceipt.BatchNumber
		receiverbatchInfo.SerialNumbers = goodsReceipt.SerialNumbers
		receiverbatchInfo.Quantity = purchaseOrder.Quantity
		receiverbatchInfo.IsCompromised = vendorbatchInfo.IsCompromised
		receiverbatchInfo.PotentialCompromised = vendorbatchInfo.PotentialCompromised

		if receiverbatchInfo.IsCompromised == true {
			receiverbatchInfo.PotentialCompromised = false
		}

		// Step 2: Update Total Quantity of Receiver Material
		receiverMaterial.TotalQuantity += purchaseOrder.Quantity

		// Step 3: Check if Batch Exists or Add New Batch with Updated Quantity
		receiverBatchExists := false
		for index, element := range receiverMaterial.Batches {
			if strings.ToLower(element.BatchNumber) == strings.ToLower(goodsReceipt.BatchNumber) {
				element.Quantity += purchaseOrder.Quantity
				receiverBatchExists = true
				receiverMaterial.Batches[index] = element
				break
			}
		}
		if receiverBatchExists == false {
			receiverMaterial.Batches = append(receiverMaterial.Batches, receiverbatchInfo)
		}

		// Update Product
		// Step 1: Create BatchTradeInfo for Mapping
		batchTradeInfoFROM := BatchTradeInfo{}
		batchTradeInfoFROM.ParticipantID = vendorbatchInfo.ParticipantID
		batchTradeInfoFROM.MaterialID = vendorbatchInfo.MaterialID
		batchTradeInfoFROM.BatchNumber = vendorbatchInfo.BatchNumber
		batchTradeInfoFROM.SerialNumbers = vendorbatchInfo.SerialNumbers
		batchTradeInfoFROM.IsCompromised = vendorbatchInfo.IsCompromised
		batchTradeInfoFROM.PotentialCompromised = vendorbatchInfo.PotentialCompromised
		batchTradeInfoFROM.Quantity += purchaseOrder.Quantity

		batchTradeInfoTO := BatchTradeInfo{}
		batchTradeInfoTO.ParticipantID = receiverbatchInfo.ParticipantID
		batchTradeInfoTO.MaterialID = receiverbatchInfo.MaterialID
		batchTradeInfoTO.BatchNumber = receiverbatchInfo.BatchNumber
		batchTradeInfoTO.SerialNumbers = receiverbatchInfo.SerialNumbers
		batchTradeInfoTO.IsCompromised = receiverbatchInfo.IsCompromised
		batchTradeInfoTO.PotentialCompromised = receiverbatchInfo.PotentialCompromised
		batchTradeInfoTO.Quantity += purchaseOrder.Quantity

		productMapping := Mapping{}
		productMapping.From = batchTradeInfoFROM
		productMapping.To = append(productMapping.To, batchTradeInfoTO)

		productRevMapping := ReverseMapping{}
		productRevMapping.To = batchTradeInfoTO
		productRevMapping.From = append(productRevMapping.From, batchTradeInfoFROM)

		// Step 2: Check if Mapping Exists
		fromMapExists := false
		ToMapExists := false
		for index, element := range product.Mappings {
			if (strings.ToLower(element.From.ParticipantID) == strings.ToLower(batchTradeInfoFROM.ParticipantID)) && (strings.ToLower(element.From.BatchNumber) == strings.ToLower(batchTradeInfoFROM.BatchNumber)) {
				for index1, element1 := range element.To {
					if strings.ToLower(element1.ParticipantID) == strings.ToLower(batchTradeInfoTO.ParticipantID) && strings.ToLower(element1.BatchNumber) == strings.ToLower(batchTradeInfoTO.BatchNumber) {
						element1.Quantity += purchaseOrder.Quantity
						element.To[index1] = element1
						ToMapExists = true
						break
					}
				}
				if ToMapExists == false {
					element.To = append(element.To, batchTradeInfoTO)
				}
				element.From.Quantity += purchaseOrder.Quantity
				product.Mappings[index] = element
				fromMapExists = true
				break
			}
		}
		if fromMapExists == false {
			product.Mappings = append(product.Mappings, productMapping)
		}

		// Check if Reverse Mapping Exists
		revfromMapExists := false
		revToMapExists := false

		for index, element := range product.ReverseMappings {
			if (strings.ToLower(element.To.ParticipantID) == strings.ToLower(batchTradeInfoTO.ParticipantID)) && (strings.ToLower(element.To.BatchNumber) == strings.ToLower(batchTradeInfoTO.BatchNumber)) {
				for index1, element1 := range element.From {
					if strings.ToLower(element1.ParticipantID) == strings.ToLower(batchTradeInfoFROM.ParticipantID) && strings.ToLower(element1.BatchNumber) == strings.ToLower(batchTradeInfoFROM.BatchNumber) {
						element1.Quantity += purchaseOrder.Quantity
						element.From[index1] = element1
						revfromMapExists = true
						break
					}
				}
				if revfromMapExists == false {
					element.From = append(element.From, batchTradeInfoFROM)
				}
				element.To.Quantity += purchaseOrder.Quantity
				product.ReverseMappings[index] = element
				revToMapExists = true
				break
			}
		}
		if revToMapExists == false {
			product.ReverseMappings = append(product.ReverseMappings, productRevMapping)
		}

		// Update Participant in Product
		// See if Vendor Exists
		vendorMaterialExists := false
		vendorMaterialDetails := MaterialDetails{}
		for _, element := range product.SupplyChainMembers {
			if strings.ToLower(element.ParticipantID) == strings.ToLower(purchaseOrder.VendorID) {
				vendorMaterialDetails = element
				vendorMaterialExists = true
			}
		}
		if vendorMaterialExists == false {
			return Error(http.StatusBadRequest, "Vendor Material Does Not Exist in Product Information! Register a Material First and Produce some Quantity")
		}

		receiverMaterialDetail := MaterialDetails{}
		receiverMaterialDetail.ParticipantID = receiverMaterial.ParticipantID
		receiverMaterialDetail.MaterialID = receiverMaterial.MaterialMasterID
		receiverMaterialDetail.ParticipantType = receiver.ParticipantType
		receiverMaterialDetail.IsCompromised = vendorMaterialDetails.IsCompromised
		receiverMaterialDetail.PotentialCompromised = vendorMaterialDetails.PotentialCompromised

		participantExists := false
		for _, element := range product.SupplyChainMembers {
			if strings.ToLower(element.ParticipantID) == strings.ToLower(receiver.ParticipantID) {
				participantExists = true
				break
			}
		}
		if participantExists == false {
			product.SupplyChainMembers = append(product.SupplyChainMembers, receiverMaterialDetail)
		}

		// Update Shipment
		shipment.Status = "COMPLETED"

		//Complete Purchase Order
		purchaseOrder.Status = "COMPLETED"

		// Store Information in Blockchain
		POjsonBytes, _ := json.Marshal(purchaseOrder)
		if puterr := stub.PutState(strings.ToLower(purchaseOrder.POID), POjsonBytes); puterr != nil {
			return Error(http.StatusInternalServerError, puterr.Error())
		}

		ShipmentjsonBytes, _ := json.Marshal(shipment)
		if puterr := stub.PutState(strings.ToLower(shipment.ShipmentID), ShipmentjsonBytes); puterr != nil {
			return Error(http.StatusInternalServerError, puterr.Error())
		}

		ProductjsonBytes, _ := json.Marshal(product)
		if puterr := stub.PutState(strings.ToLower(product.ProductID), ProductjsonBytes); puterr != nil {
			return Error(http.StatusInternalServerError, puterr.Error())
		}

		vendorMaterialjsonBytes, _ := json.Marshal(vendorMaterial)
		if puterr := stub.PutState(strings.ToLower(vendorMaterialID), vendorMaterialjsonBytes); puterr != nil {
			return Error(http.StatusInternalServerError, puterr.Error())
		}

		receiverMaterialjsonBytes, _ := json.Marshal(receiverMaterial)
		if puterr := stub.PutState(strings.ToLower(receiverMaterialID), receiverMaterialjsonBytes); puterr != nil {
			return Error(http.StatusInternalServerError, puterr.Error())
		}

		GRjsonBytes, _ := json.Marshal(goodsReceipt)
		if puterr := stub.PutState(strings.ToLower(goodsReceipt.GRNumber), GRjsonBytes); puterr != nil {
			return Error(http.StatusInternalServerError, puterr.Error())
		}
		return Success(http.StatusCreated, "Goods Received Against Production Order", nil)
	} else {
		return Error(http.StatusBadRequest, "Invoke Error: Invalid Data. Currently, Valid GR Types are Against: \n 1) PRODUCTION ORDER \n 2) PURCHASE ORDER ")
	}
}

// CASE 09 Report Contamination
func (t *BlockchainIOT) reportContamination(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 1 {
		return Error(http.StatusBadRequest, "Invoke Error: Incorrect number of arguments - One Argument expected")
	}

	type QueryData struct {
		ParticipantID string `json:"ParticipantID"`
		MaterialID    string `json:"MaterialID"`
		BatchNumber   string `json:"BatchNumber"`
	}

	data := string(args[0])
	queryData := QueryData{}
	err := json.Unmarshal([]byte(data), &queryData)
	if err != nil {
		return Error(http.StatusBadRequest, "Invoke Error: Invalid Data - Check Payload")
	}

	contaminatedBatch := BatchContamination{}
	contaminatedBatch.ParticipantID = queryData.ParticipantID
	contaminatedBatch.MaterialID = queryData.MaterialID
	contaminatedBatch.BatchNumber = queryData.BatchNumber

	// Get Material
	materialID := contaminatedBatch.ParticipantID + "-" + contaminatedBatch.MaterialID
	materialValue, materialGetErr := stub.GetState(strings.ToLower(materialID))
	if materialGetErr != nil || materialValue == nil {
		return Error(http.StatusNotFound, "Material Does Not Exist! Please Check Participant ID and Material ID!")
	}
	material := Material{}
	json.Unmarshal(materialValue, &material)

	// Get the Associated Product for Material
	productValue, productGetErr := stub.GetState(strings.ToLower(material.ProductBCID))
	if productGetErr != nil || productValue == nil {
		return Error(http.StatusNotFound, "Product Does Not Exist for this Material!")
	}
	myproduct := Product{}
	json.Unmarshal(productValue, &myproduct)

	newproduct := setContamination(stub, contaminatedBatch.ParticipantID, contaminatedBatch.MaterialID, contaminatedBatch.BatchNumber, myproduct)
	finalproduct := setPotentialContamination(stub, contaminatedBatch.ParticipantID, contaminatedBatch.MaterialID, contaminatedBatch.BatchNumber, newproduct)

	// Update the Material
	for index, element := range material.Batches {
		if strings.ToLower(element.BatchNumber) == strings.ToLower(contaminatedBatch.BatchNumber) {
			element.IsCompromised = true
			element.PotentialCompromised = false
			material.Batches[index] = element
			break
		}
	}

	// Update initial SCM in Product
	for index, element := range finalproduct.SupplyChainMembers {
		if strings.ToLower(element.ParticipantID) == strings.ToLower(contaminatedBatch.ParticipantID) {
			element.IsCompromised = true
			element.PotentialCompromised = false
			finalproduct.SupplyChainMembers[index] = element
			break
		}
	}

	// Store Updated Product
	materialJsonBytes, _ := json.Marshal(material)
	if puterr := stub.PutState(strings.ToLower(materialID), materialJsonBytes); puterr != nil {
		return Error(http.StatusInternalServerError, puterr.Error())
	}

	JsonBytes, _ := json.Marshal(finalproduct)
	if puterr := stub.PutState(strings.ToLower(finalproduct.ProductID), JsonBytes); puterr != nil {
		return Error(http.StatusInternalServerError, puterr.Error())
	}
	return Success(http.StatusCreated, "Product Updated", nil)
}

func setContamination(stub shim.ChaincodeStubInterface, participant string, material string, batch string, product Product) *Product {
	// Gets All Mapping and Checks for the Contaminated Batch in FROM
	for index, element := range product.Mappings {
		if (strings.ToLower(element.From.ParticipantID) == strings.ToLower(participant)) && (strings.ToLower(element.From.BatchNumber) == strings.ToLower(batch)) {
			// Set FROM to Compromised
			element.From.IsCompromised = true
			element.From.PotentialCompromised = false
			// Get All TOs for the FROM Compromised
			for index1, element1 := range element.To {
				// Set TOs to Compromised
				element1.IsCompromised = true
				element1.PotentialCompromised = false
				element.To[index1] = element1
				// Get all TOs in Reverse Mapping
				for index4, element4 := range product.ReverseMappings {
					if (strings.ToLower(element4.To.ParticipantID) == strings.ToLower(element1.ParticipantID)) && (strings.ToLower(element4.To.BatchNumber) == strings.ToLower(element1.BatchNumber)) {
						// Set reverse mapping TOs as Compromised
						element4.To.IsCompromised = true
						element4.To.PotentialCompromised = false
						product.ReverseMappings[index4] = element4
						/*
							for index5, element5 := range element4.From {
								if (strings.ToLower(element5.ParticipantID) == strings.ToLower(element1.ParticipantID)) && (strings.ToLower(element5.BatchNumber) == strings.ToLower(element1.BatchNumber)) {
									element5.IsCompromised = element4.To.IsCompromised
									element5.PotentialCompromised = element4.To.PotentialCompromised
									element4.From[index5] = element5
									break
								}
							}
						*/
					}
				}
				for index2, element2 := range product.SupplyChainMembers {
					if strings.ToLower(element2.ParticipantID) == strings.ToLower(element1.ParticipantID) {
						element2.IsCompromised = true
						element2.PotentialCompromised = false
						product.SupplyChainMembers[index2] = element2
						break
					}
				}
				// Get the Material for Batch in TO
				newMaterialID := element1.ParticipantID + "-" + element1.MaterialID
				newMaterialValue, _ := stub.GetState(strings.ToLower(newMaterialID))
				newMaterial := Material{}
				json.Unmarshal(newMaterialValue, &newMaterial)
				for index3, element3 := range newMaterial.Batches {
					if strings.ToLower(element3.BatchNumber) == strings.ToLower(element1.BatchNumber) {
						element3.IsCompromised = true
						element3.PotentialCompromised = false
						newMaterial.Batches[index3] = element3
					}
				}
				// Store the Material
				JsonBytes, _ := json.Marshal(newMaterial)
				if puterr := stub.PutState(strings.ToLower(newMaterialID), JsonBytes); puterr != nil {
					newMaterial.Plant = "FFFF"
				}
			}
			product.Mappings[index] = element
			for _, element5 := range element.To {
				setContamination(stub, element5.ParticipantID, element5.MaterialID, element5.BatchNumber, product)
			}
		}
	}
	return &product
}

func setPotentialContamination(stub shim.ChaincodeStubInterface, participant string, material string, batch string, product *Product) *Product {
	for index, element := range product.ReverseMappings {
		if (strings.ToLower(element.To.ParticipantID) == strings.ToLower(participant)) && (strings.ToLower(element.To.BatchNumber) == strings.ToLower(batch)) {
			for index1, element1 := range element.From {
				element1.IsCompromised = false
				element1.PotentialCompromised = true
				element.From[index1] = element1
				for index4, element4 := range product.Mappings {
					if (strings.ToLower(element4.From.ParticipantID) == strings.ToLower(element1.ParticipantID)) && (strings.ToLower(element4.From.BatchNumber) == strings.ToLower(element1.BatchNumber)) {
						element4.From.IsCompromised = false
						element4.From.PotentialCompromised = true
						product.Mappings[index4] = element4
						/*
							for index5, element5 := range element4.To {
								if (strings.ToLower(element5.ParticipantID) == strings.ToLower(element1.ParticipantID)) && (strings.ToLower(element5.BatchNumber) == strings.ToLower(element1.BatchNumber)) {
									element5.IsCompromised = element4.From.IsCompromised
									element5.PotentialCompromised = element4.From.PotentialCompromised
									element4.To[index5] = element5
									break
								}
							}
						*/
					}
				}
				for index2, element2 := range product.SupplyChainMembers {
					if strings.ToLower(element2.ParticipantID) == strings.ToLower(element1.ParticipantID) {
						element2.IsCompromised = false
						element2.PotentialCompromised = true
						product.SupplyChainMembers[index2] = element2
						break
					}
				}
				// Get the Material for Batch in FROM
				newMaterialID := element1.ParticipantID + "-" + element1.MaterialID
				newMaterialValue, _ := stub.GetState(strings.ToLower(newMaterialID))
				newMaterial := Material{}
				json.Unmarshal(newMaterialValue, &newMaterial)
				for index3, element3 := range newMaterial.Batches {
					if strings.ToLower(element3.BatchNumber) == strings.ToLower(element1.BatchNumber) {
						element3.IsCompromised = false
						element3.PotentialCompromised = true
						newMaterial.Batches[index3] = element3
					}
				}
				// Store the Material
				JsonBytes, _ := json.Marshal(newMaterial)
				if puterr := stub.PutState(strings.ToLower(newMaterialID), JsonBytes); puterr != nil {
					newMaterial.Plant = "FFFF"
				}
			}
			product.ReverseMappings[index] = element
			for _, element5 := range element.From {
				setPotentialContamination(stub, element5.ParticipantID, element5.MaterialID, element5.BatchNumber, product)
			}
		}
	}
	return product
}

// CASE 10 Clear Contamination
func (t *BlockchainIOT) clearContamination(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 1 {
		return Error(http.StatusBadRequest, "Invoke Error: Incorrect number of arguments - One Argument expected")
	}

	type QueryData struct {
		ParticipantID string `json:"ParticipantID"`
		MaterialID    string `json:"MaterialID"`
		BatchNumber   string `json:"BatchNumber"`
	}

	data := string(args[0])
	queryData := QueryData{}
	err := json.Unmarshal([]byte(data), &queryData)
	if err != nil {
		return Error(http.StatusBadRequest, "Invoke Error: Invalid Data - Check Payload")
	}

	contaminatedBatch := BatchContamination{}
	contaminatedBatch.ParticipantID = queryData.ParticipantID
	contaminatedBatch.MaterialID = queryData.MaterialID
	contaminatedBatch.BatchNumber = queryData.BatchNumber

	// Get Material
	materialID := contaminatedBatch.ParticipantID + "-" + contaminatedBatch.MaterialID
	materialValue, materialGetErr := stub.GetState(strings.ToLower(materialID))
	if materialGetErr != nil || materialValue == nil {
		return Error(http.StatusNotFound, "Material Does Not Exist! Please Check Participant ID and Material ID!")
	}
	material := Material{}
	json.Unmarshal(materialValue, &material)

	// Get the Associated Product for Material
	productValue, productGetErr := stub.GetState(strings.ToLower(material.ProductBCID))
	if productGetErr != nil || productValue == nil {
		return Error(http.StatusNotFound, "Product Does Not Exist for this Material!")
	}
	myproduct := Product{}
	json.Unmarshal(productValue, &myproduct)

	newproduct := clearBatchContamination(stub, contaminatedBatch.ParticipantID, contaminatedBatch.MaterialID, contaminatedBatch.BatchNumber, myproduct)
	finalproduct := clearPotentialContamination(stub, contaminatedBatch.ParticipantID, contaminatedBatch.MaterialID, contaminatedBatch.BatchNumber, newproduct)

	// Update the Material
	for index, element := range material.Batches {
		if strings.ToLower(element.BatchNumber) == strings.ToLower(contaminatedBatch.BatchNumber) {
			element.IsCompromised = false
			element.PotentialCompromised = false
			material.Batches[index] = element
			break
		}
	}

	// Update initial SCM in Product
	for index, element := range finalproduct.SupplyChainMembers {
		if strings.ToLower(element.ParticipantID) == strings.ToLower(contaminatedBatch.ParticipantID) {
			element.IsCompromised = false
			element.PotentialCompromised = false
			finalproduct.SupplyChainMembers[index] = element
			break
		}
	}

	// Store Updated Product
	materialJsonBytes, _ := json.Marshal(material)
	if puterr := stub.PutState(strings.ToLower(materialID), materialJsonBytes); puterr != nil {
		return Error(http.StatusInternalServerError, puterr.Error())
	}

	JsonBytes, _ := json.Marshal(finalproduct)
	if puterr := stub.PutState(strings.ToLower(finalproduct.ProductID), JsonBytes); puterr != nil {
		return Error(http.StatusInternalServerError, puterr.Error())
	}
	return Success(http.StatusCreated, "Product Updated", nil)
}

func clearBatchContamination(stub shim.ChaincodeStubInterface, participant string, material string, batch string, product Product) *Product {
	for index, element := range product.Mappings {
		if (strings.ToLower(element.From.ParticipantID) == strings.ToLower(participant)) && (strings.ToLower(element.From.BatchNumber) == strings.ToLower(batch)) {
			element.From.IsCompromised = false
			element.From.PotentialCompromised = false
			for index1, element1 := range element.To {
				element1.IsCompromised = false
				element1.PotentialCompromised = false
				element.To[index1] = element1
				for index4, element4 := range product.ReverseMappings {
					if (strings.ToLower(element4.To.ParticipantID) == strings.ToLower(element1.ParticipantID)) && (strings.ToLower(element4.To.BatchNumber) == strings.ToLower(element1.BatchNumber)) {
						element4.To.IsCompromised = false
						element4.To.PotentialCompromised = false
						product.ReverseMappings[index4] = element4
						/*
							for index5, element5 := range element4.From {
								if (strings.ToLower(element5.ParticipantID) == strings.ToLower(element1.ParticipantID)) && (strings.ToLower(element5.BatchNumber) == strings.ToLower(element1.BatchNumber)) {
									element5.IsCompromised = false
									element5.PotentialCompromised = false
									element4.From[index5] = element5
									break
								}
							}
						*/
					}
				}
				for index2, element2 := range product.SupplyChainMembers {
					if strings.ToLower(element2.ParticipantID) == strings.ToLower(element1.ParticipantID) {
						element2.IsCompromised = false
						element2.PotentialCompromised = false
						product.SupplyChainMembers[index2] = element2
						break
					}
				}
				// Get the Material for Batch in TO
				newMaterialID := element1.ParticipantID + "-" + element1.MaterialID
				newMaterialValue, _ := stub.GetState(strings.ToLower(newMaterialID))
				newMaterial := Material{}
				json.Unmarshal(newMaterialValue, &newMaterial)
				for index3, element3 := range newMaterial.Batches {
					if strings.ToLower(element3.BatchNumber) == strings.ToLower(element1.BatchNumber) {
						element3.IsCompromised = false
						element3.PotentialCompromised = false
						newMaterial.Batches[index3] = element3
					}
				}
				// Store the Material
				JsonBytes, _ := json.Marshal(newMaterial)
				if puterr := stub.PutState(strings.ToLower(newMaterialID), JsonBytes); puterr != nil {
					newMaterial.Plant = "FFFF"
				}
			}
			product.Mappings[index] = element
			for _, element5 := range element.To {
				clearBatchContamination(stub, element5.ParticipantID, element5.MaterialID, element5.BatchNumber, product)
			}
		}
	}
	return &product
}

func clearPotentialContamination(stub shim.ChaincodeStubInterface, participant string, material string, batch string, product *Product) *Product {
	for index, element := range product.ReverseMappings {
		if (strings.ToLower(element.To.ParticipantID) == strings.ToLower(participant)) && (strings.ToLower(element.To.BatchNumber) == strings.ToLower(batch)) {
			for index1, element1 := range element.From {
				element1.IsCompromised = false
				element1.PotentialCompromised = false
				element.From[index1] = element1
				for index4, element4 := range product.Mappings {
					if (strings.ToLower(element4.From.ParticipantID) == strings.ToLower(element1.ParticipantID)) && (strings.ToLower(element4.From.BatchNumber) == strings.ToLower(element1.BatchNumber)) {
						element4.From.IsCompromised = false
						element4.From.PotentialCompromised = false
						product.Mappings[index4] = element4
						/*
							for index5, element5 := range element4.To {
								if (strings.ToLower(element5.ParticipantID) == strings.ToLower(element1.ParticipantID)) && (strings.ToLower(element5.BatchNumber) == strings.ToLower(element1.BatchNumber)) {
									element5.IsCompromised = false
									element5.PotentialCompromised = false
									element4.To[index5] = element5
									break
								}
							}
						*/
					}
				}
				for index2, element2 := range product.SupplyChainMembers {
					if strings.ToLower(element2.ParticipantID) == strings.ToLower(element1.ParticipantID) {
						element2.IsCompromised = false
						element2.PotentialCompromised = false
						product.SupplyChainMembers[index2] = element2
						break
					}
				}
				// Get the Material for Batch in FROM
				newMaterialID := element1.ParticipantID + "-" + element1.MaterialID
				newMaterialValue, _ := stub.GetState(strings.ToLower(newMaterialID))
				newMaterial := Material{}
				json.Unmarshal(newMaterialValue, &newMaterial)
				for index3, element3 := range newMaterial.Batches {
					if strings.ToLower(element3.BatchNumber) == strings.ToLower(element1.BatchNumber) {
						element3.IsCompromised = false
						element3.PotentialCompromised = false
						newMaterial.Batches[index3] = element3
					}
				}
				// Store the Material
				JsonBytes, _ := json.Marshal(newMaterial)
				if puterr := stub.PutState(strings.ToLower(newMaterialID), JsonBytes); puterr != nil {
					newMaterial.Plant = "FFFF"
				}
			}
			product.ReverseMappings[index] = element
			for _, element3 := range element.From {
				clearPotentialContamination(stub, element3.ParticipantID, element3.MaterialID, element3.BatchNumber, product)
			}
		}
	}
	return product
}

// CASE 11 Get Materials
func (t *BlockchainIOT) getMaterial(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 1 {
		return Error(http.StatusBadRequest, "Invoke Error: Incorrect number of arguments - One Argument expected")
	}
	data := string(args[0])
	data1 := string(args[1])

	materialID := data + "-" + data1

	//Get the Material from Blockchain
	value, geterr := stub.GetState(strings.ToLower(materialID))
	if geterr != nil || value == nil {
		return Error(http.StatusNotFound, "Not Found")
	}
	return Success(http.StatusOK, "OK", value)
}

// CASE 12 Delete Material Asset
func (t *BlockchainIOT) deleteMaterial(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 1 {
		return Error(http.StatusBadRequest, "Invoke Error: Incorrect number of arguments - One Argument expected")
	}
	data := string(args[0])
	data1 := string(args[1])

	materialID := data + "-" + data1

	// Check if Exists
	value, geterr := stub.GetState(strings.ToLower(materialID))
	if geterr != nil || value == nil {
		return Error(http.StatusNotFound, "Not Found")
	}

	// Delete if Exists
	if delerr := stub.DelState(strings.ToLower(materialID)); delerr != nil {
		return Error(http.StatusInternalServerError, delerr.Error())
	}
	return Success(http.StatusNoContent, "Material Deleted", nil)
}

// CASE 13 Get Any Asset
func (t *BlockchainIOT) getAsset(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 1 {
		return Error(http.StatusBadRequest, "Invoke Error: Incorrect number of arguments - One Argument expected")
	}
	data := strings.ToLower(string(args[0]))

	//Get the Asset from Blockchain
	value, geterr := stub.GetState(data)
	if geterr != nil || value == nil {
		return Error(http.StatusNotFound, "Not Found")
	}
	return Success(http.StatusOK, "OK", value)
}

// CASE 14 Delete Any Asset
func (t *BlockchainIOT) deleteAsset(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 1 {
		return Error(http.StatusBadRequest, "Invoke Error: Incorrect number of arguments - One Argument expected")
	}
	data := strings.ToLower(string(args[0]))

	// Check if Exists
	value, geterr := stub.GetState(data)
	if geterr != nil || value == nil {
		return Error(http.StatusNotFound, "Not Found")
	}

	// Delete if Exists
	if delerr := stub.DelState(data); delerr != nil {
		return Error(http.StatusInternalServerError, delerr.Error())
	}
	return Success(http.StatusNoContent, "Asset Deleted", nil)
}

//********************************************************************************************************
// Micellanious Functions
//********************************************************************************************************

// Get Transactions History From Blockchain
func (t *BlockchainIOT) getHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	key := string(args[0])
	historyResult, err := getHistoryExecution(stub, key)
	if err != nil {
		return Error(http.StatusInternalServerError, err.Error())
	}
	return Success(http.StatusOK, "OK", historyResult)
}

func getHistoryExecution(stub shim.ChaincodeStubInterface, key string) ([]byte, error) {
	historyKey := strings.ToLower(key)
	resultsIterator, err := stub.GetHistoryForKey(historyKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	//JSON Array Buffer
	var buffer bytes.Buffer
	buffer.WriteString("[")

	alreadyFetched := false
	for resultsIterator.HasNext() {
		output, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		if alreadyFetched == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TransactionId\":")
		buffer.WriteString("\"")
		buffer.WriteString(output.TxId)
		buffer.WriteString("\"")
		buffer.WriteString(", \"Value\":")
		if output.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(output.Value))
		}
		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(output.Timestamp.Seconds, int64(output.Timestamp.Nanos)).String())
		buffer.WriteString("\"")
		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(output.IsDelete))
		buffer.WriteString("\"")
		buffer.WriteString("}")
		alreadyFetched = true
	}
	buffer.WriteString("]")
	return buffer.Bytes(), nil
}

// Custom Queries
func (t *BlockchainIOT) customQueries(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	searchString := string(args[0])
	queryResults, err := queryexecution(stub, searchString)
	if err != nil {
		return Error(http.StatusInternalServerError, err.Error())
	}
	return Success(http.StatusOK, "OK", queryResults)
}

func queryexecution(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	//JSON Array Buffer
	var buffer bytes.Buffer
	buffer.WriteString("[")

	alreadyFetched := false
	for resultsIterator.HasNext() {
		output, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		if alreadyFetched == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(output.Key)
		buffer.WriteString("\"")
		buffer.WriteString(", \"Record\":")
		buffer.WriteString(string(output.Value))
		buffer.WriteString("}")
		alreadyFetched = true
	}
	buffer.WriteString("]")
	return buffer.Bytes(), nil
}
