package mapr

import (
	"encoding/json"
	"fmt"
	"time"
	"vmrenter/pkg/config"
	"vmrenter/pkg/models"

	client "github.com/mapr/maprdb-go-client"
)

const tableName = "/user/mapr/nodes"

func getConnection() (*client.Connection, error) {

	connection, err := client.MakeConnection(config.GetURLDBConn())

	if err != nil {
		panic(err)
	}

	return connection, err
}

func getUnavailableNodes(clusterID string, operatingSystem string) []models.Node {
	connection, err := getConnection()
	if err != nil {
		fmt.Println("error getting connection", err)
	}

	// this will get called when function exits after this point, irregardless of returning a value or error
	defer connection.Close()

	// Options for find request
	options := &client.FindOptions{ResultAsDocument: true}

	store, err := connection.GetStore(tableName)
	if err != nil {
		panic(err)
	}

	queryStr := fmt.Sprintf(`{"$where":{"$and":[{"$eq":{"Node.OperatingSystem.Name":"Ubuntu"}},{"$gt":{"ExpiresAt": "%s"}}] }}`, time.Now().Add(3*24*time.Hour).Format(time.RFC3339))
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

func writeToDB(inputStr string) error {
	fmt.Println("The time is", time.Now())

	connectionString := config.URLDBConn

	//options := &client.ConnectionOptions{MaxAttempt:3, WaitBetweenSeconds:10, CallTimeoutSeconds:60}
	storeName := tableName

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
		return err
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

	return nil
}

func MakeReservation(clusterID string, nodes []models.Node, jenkinsJobURL string) error {
	var s struct {
		ID            string `json:"_id"`
		ClusterID     string `json:"ClusterID"`
		Node          models.Node
		JenkinsJobURL string
		ExpiresAt     time.Time
	}
	s.ClusterID = clusterID
	s.Node = nodes[0]
	now := time.Now()
	s.ID = s.Node.Host + now.Format("2006-01-02_150405")
	s.JenkinsJobURL = jenkinsJobURL
	s.ExpiresAt = now.Add(24 * time.Hour)

	by, _ := json.Marshal(&s)
	a := string(by)

	return writeToDB(a)
}
