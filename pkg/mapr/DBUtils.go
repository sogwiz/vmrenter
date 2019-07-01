package mapr

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
	"vmrenter/pkg/config"
	"vmrenter/pkg/models"

	client "github.com/mapr/maprdb-go-client"
)

const tableNodes = "/user/mapr/nodes"
const tableReservations = "/user/mapr/reservations"

func GetConnection() (*client.Connection, error) {

	connection, err := client.MakeConnection(config.GetURLDBConn())

	if err != nil {
		panic(err)
	}

	return connection, err
}

func IsRequestDoable(numNodes int, osName string, osVersion string) bool {

	fmt.Println("Request: ", numNodes, osName)

	nodeDBJsons := GetAvailableNodes("", "centos")

	fmt.Println("Available nodes ", len(nodeDBJsons))
	if len(nodeDBJsons) < numNodes {
		fmt.Println("Can't fulfill this request unfortunately")
		return false
	}

	fmt.Println("Can fullfill request")
	return true
}

func GetAvailableNodes(clusterID string, operatingSystem string) []models.NodeDBJson {
	connection, err := GetConnection()
	if err != nil {
		fmt.Println("error getting connection", err)
	}

	// this will get called when function exits after this point, irregardless of returning a value or error
	defer connection.Close()

	// Options for find request
	// THIS IS CRITICAL - must be FALSE
	options := &client.FindOptions{ResultAsDocument: false}

	store, err := connection.GetStore(tableNodes)
	if err != nil {
		panic(err)
	}

	// query for nodes where the ExpiresAt field has already passed
	//now.Add(3*24*time.Hour)
	queryStr := fmt.Sprintf(`{"$where":{"$and":[{"$matches":{"NodeObj.OperatingSystem.Name":"(?i)%s"}},{"$lt":{"ExpiresAT": "%s"}}] }}`, operatingSystem, time.Now().Format(time.RFC3339))
	fmt.Println(queryStr)

	findResult, err := store.FindQueryString(queryStr, options)
	if err != nil {
		panic(err)
	}

	nodeDBJsons := make([]models.NodeDBJson, 0)

	// Print OJAI Documents from document stream
	for _, doc := range findResult.DocumentList() {
		tmpNode := models.NodeDBJson{}
		tmp, _ := json.Marshal(doc)
		//fmt.Println(err)
		//fmt.Println(string(tmp))
		err = json.Unmarshal(tmp, &tmpNode)
		nodeDBJsons = append(nodeDBJsons, tmpNode)
		//fmt.Println(err)
		//fmt.Println(tmpNode)
	}

	return nodeDBJsons
}

func getAllNodes() []models.NodeDBJson {
	connection, err := GetConnection()
	if err != nil {
		fmt.Println("error getting connection", err)
	}

	// this will get called when function exits after this point, irregardless of returning a value or error
	defer connection.Close()

	// Options for find request
	options := &client.FindOptions{ResultAsDocument: false}

	store, err := connection.GetStore(tableNodes)
	if err != nil {
		panic(err)
	}

	findResult, err := store.FindAll(options)
	if err != nil {
		panic(err)
	}

	nodeDBJsons := make([]models.NodeDBJson, 0)

	// Print OJAI Documents from document stream
	for _, doc := range findResult.DocumentList() {
		tmpNode := models.NodeDBJson{}
		tmp, _ := json.Marshal(doc)
		//fmt.Println(err)
		//fmt.Println(string(tmp))
		err = json.Unmarshal(tmp, &tmpNode)
		nodeDBJsons = append(nodeDBJsons, tmpNode)
		//fmt.Println(err)
		//fmt.Println(tmpNode)
	}

	return nodeDBJsons
}

func getUnavailableNodes(clusterID string, operatingSystem string) []models.Node {
	connection, err := GetConnection()
	if err != nil {
		fmt.Println("error getting connection", err)
	}

	// this will get called when function exits after this point, irregardless of returning a value or error
	defer connection.Close()

	// Options for find request
	options := &client.FindOptions{ResultAsDocument: false}

	store, err := connection.GetStore(tableNodes)
	if err != nil {
		panic(err)
	}

	queryStr := fmt.Sprintf(`{"$where":{"$and":[{"$eq":{"Node.OperatingSystem.Name":"Ubuntu"}},{"$gt":{"ExpiresAt": "%s"}}] }}`, time.Now().Add(3 * 24 * time.Hour).Format(time.RFC3339))
	fmt.Println(queryStr)

	findResult, err := store.FindQueryString(queryStr, options)
	if err != nil {
		panic(err)
	}

	// Print OJAI Documents from document stream
	for _, doc := range findResult.DocumentList() {
		fmt.Println(doc)
	}

	return nil

}

func WriteToDBWithTableMap(inputMap map[string]interface{}, table string) error {
	fmt.Println("The time is", time.Now())

	connectionString := config.GetURLDBConn()

	//options := &client.ConnectionOptions{MaxAttempt:3, WaitBetweenSeconds:10, CallTimeoutSeconds:60}
	storeName := table

	//fmt.Println("Connection string is ", connectionString)
	if connectionString == "" {
		panic("Connection string must not be empty")
	}
	connection, err := client.MakeConnection(connectionString)

	if err != nil {
		panic(err)
	}

	// this will get called when function exits after this point, irregardless of returning a value or error
	defer connection.Close()

	store, err := connection.GetStore(storeName)

	if err != nil {
		panic(err)
	}

	newDocument := connection.CreateDocumentFromMap(inputMap)

	err = store.InsertOrReplaceDocument(newDocument)

	if err != nil {
		fmt.Println("Error calling InsertOrReplaceDocument", err)
	}
	// Options for find request
	//options := &client.FindOptions{ResultAsDocument: true}

	//query := map[string]interface{}{"$select": []interface{}{"CourseName", "RaceID"},
	//	"$where": map[string]interface{}{
	//		"$like": map[string]interface{}{"DisplayName": "sargon%benjamin"}}}
	//
	//findResult, err := store.FindQueryMap(query, options)
	//
	//iterations := 0
	//for _, doc := range findResult.DocumentList() {
	//	fmt.Println(doc)
	//	iterations+=1
	//}
	//
	//fmt.Println(iterations)

	return nil

}

func ReserveNode(nodeID string, expiresAT string, clusterID string) error {
	fmt.Println(nodeID)
	connectionString := config.GetURLDBConn()

	storeName := tableNodes

	if connectionString == "" {
		panic("Connection string must not be empty")
	}
	connection, err := client.MakeConnection(connectionString)

	if err != nil {
		panic(err)
	}

	// this will get called when function exits after this point, irregardless of returning a value or error
	defer connection.Close()

	store, err := connection.GetStore(storeName)

	if err != nil {
		panic(err)
	}

	mutation := map[string]interface{}{"$set": []interface{}{
		map[string]interface{}{"ExpiresAT": expiresAT},
		map[string]interface{}{"ClusterID": clusterID},
	},
	}
	//mutation := map[string]interface{}{"$set": map[string]interface{}{"ExpiresAT": expiresAT}}
	//mutation := map[string]interface{}{"$set": []interface{}{{"ExpiresAT": expiresAT},{"ClusterID": clusterID}}}
	//mutationStr := "{\"$set\":[{\"ExpiresAT\":" + expiresAT + "},{\"ClusterID\":" + clusterID + "}]}"

	//mutation := map[string]interface{}{"$set": map[string]interface{}{"ExpiresAT": expiresAT}}
	docID := client.BosiFromString(nodeID)
	//docMutation := client.
	docMutation := client.MosmFromMap(mutation)

	err = store.Update(docID, docMutation)

	if err != nil {
		fmt.Println("Error updating node", err)
		return err
	}
	return nil
}

func WriteToDBWithTable(inputStr string, table string) (*client.Document, error) {
	fmt.Println("The time is", time.Now())

	connectionString := config.GetURLDBConn()

	//options := &client.ConnectionOptions{MaxAttempt:3, WaitBetweenSeconds:10, CallTimeoutSeconds:60}
	storeName := table

	if connectionString == "" {
		panic("Connection string must not be empty")
	}
	connection, err := client.MakeConnection(connectionString)

	if err != nil {
		panic(err)
	}

	// this will get called when function exits after this point, irregardless of returning a value or error
	defer connection.Close()

	store, err := connection.GetStore(storeName)

	if err != nil {
		panic(err)
	}

	newDocument, err := connection.CreateDocumentFromString(inputStr)

	if err != nil {
		return nil, err
	}

	err = store.InsertOrReplaceDocument(newDocument)

	if err != nil {
		fmt.Println("Error calling InsertOrReplaceDocument", err)
	}
	// Options for find request
	//options := &client.FindOptions{ResultAsDocument: true}

	//query := map[string]interface{}{"$select": []interface{}{"CourseName", "RaceID"},
	//	"$where": map[string]interface{}{
	//		"$like": map[string]interface{}{"DisplayName": "sargon%benjamin"}}}
	//
	//findResult, err := store.FindQueryMap(query, options)
	//
	//iterations := 0
	//for _, doc := range findResult.DocumentList() {
	//	fmt.Println(doc)
	//	iterations+=1
	//}
	//
	//fmt.Println(iterations)

	return newDocument, err
}

func WriteToDB(inputStr string) (*client.Document, error) {
	return WriteToDBWithTable(inputStr, tableNodes)
}

/**
	this has to update 2 tables
 - 1. nodes table : update the clusterID and expiresAT for each node
 - 2. reservations table:
*/
func MakeReservation(clusterID string, requestor string, nodes []models.NodeDBJson, jenkinsJobURL string,
	reservationType string, hoursToReserve int) (models.Reservation, error) {

	if len(nodes) > 5 && requestor != "sbenjamin@mapr.com" {
		panic("Can't request more than 5 nodes")
	}

	now := time.Now()
	expiry := now.Add(time.Duration(hoursToReserve) * time.Hour)

	var wg sync.WaitGroup
	for i, _ := range nodes {
		nodes[i].ExpiresAT = expiry.Format(time.RFC3339)
		nodes[i].ClusterID = clusterID

		node := nodes[i]

		wg.Add(1)
		go func(node models.NodeDBJson) {
			defer wg.Done()
			err := ReserveNode(node.ID, node.ExpiresAT, node.ClusterID)
			if err != nil {
				fmt.Println("Found error for node", node, err)
			}
		}(node)
	}
	wg.Wait()

	reservationID := clusterID + "_" + requestor + "_" + now.Format("2006-01-02_150405")

	reservation := models.Reservation{
		ID:              reservationID,
		CreatedAt:       now.Format("2006-01-02_150405"),
		ExpiresAt:       expiry.Format("2006-01-02_150405"),
		JenkinsJobURL:   jenkinsJobURL,
		Nodes:           nodes,
		ReservationType: strings.ToLower(reservationType),
		ClusterID:       clusterID,
	}

	by, _ := json.Marshal(&reservation)
	a := string(by)

	_, error := WriteToDBWithTable(a, tableReservations)
	return reservation, error
}

func Reset(tableName string) error {
	connectionString := config.GetURLDBConn()

	//options := &client.ConnectionOptions{MaxAttempt:3, WaitBetweenSeconds:10, CallTimeoutSeconds:60}
	fmt.Println("Connection string is: ", connectionString)

	if connectionString == "" {
		panic("Connection string must not be empty")
	}

	connection, err := client.MakeConnection(connectionString)

	if err != nil {
		panic(err)
	}

	// this will get called when function exits after this point, irregardless of returning a value or error
	defer connection.Close()

	error := connection.DeleteStore(tableName)
	if error != nil {
		fmt.Println("Couldn't delete table ", tableName, error)
	}
	_, error = connection.CreateStore(tableName)
	if error != nil {
		fmt.Println("Couldn't create table ", tableName, error)
	}

	return error
}

func GetPartialReservationsForNodesUpdate() []models.PartialReservationForNodesUpdate {
	connection, err := GetConnection()
	if err != nil {
		fmt.Println("error getting connection", err)
	}

	// this will get called when function exits after this point, irregardless of returning a value or error
	defer connection.Close()

	// Options for find request
	// THIS IS CRITICAL - must be FALSE
	options := &client.FindOptions{ResultAsDocument: false}

	store, err := connection.GetStore(tableNodes)
	if err != nil {
		panic(err)
	}

	// query for nodes where the ExpiresAt field has not passed yet
	query := fmt.Sprintf(`{"$select":["_id","ExpiresAT","ClusterID"],"$where":{"$gt":{"ExpiresAT":"%v"}}}`, time.Now().Format(time.RFC3339))

	findResult, err := store.FindQueryString(query, options)
	if err != nil {
		panic(err)
	}

	partialReservationsForNodesUpdate := make([]models.PartialReservationForNodesUpdate, 0)

	// Print OJAI Documents from document stream
	for _, doc := range findResult.DocumentList() {
		tmpPartialReservation := models.PartialReservationForNodesUpdate{}
		tmp, _ := json.Marshal(doc)
		err = json.Unmarshal(tmp, &tmpPartialReservation)
		partialReservationsForNodesUpdate = append(partialReservationsForNodesUpdate, tmpPartialReservation)
	}

	return partialReservationsForNodesUpdate
}
