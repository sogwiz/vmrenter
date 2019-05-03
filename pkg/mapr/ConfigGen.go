package mapr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"vmrenter/pkg/models"
)

func GenerateConfigJson(reservation models.Reservation, generateESXIEntries bool) {
	var configData models.FileContent
	configData = GetConfigObject("/Users/sargonbenjamin/dev/src/private-installer/testing/configuration/config.json")

	//var nodes [len(reservation.Nodes)]models.Node

	expectedServiceNames := make([]string, 0)
	expectedServiceNames = append(expectedServiceNames, "mapr-core")
	nodes := make([]models.Node, len(reservation.Nodes))
	for i, nodejson := range reservation.Nodes {
		nodes[i] = nodejson.NodeObj
		nodes[i].ExpectedServiceNames = expectedServiceNames
	}

	clusterToReserve := models.Cluster{
		ID:    reservation.ClusterID,
		Name:  reservation.ClusterID,
		Nodes: nodes,
	}
	configData.Clusters = append(configData.Clusters, clusterToReserve)
	fmt.Println(configData.Clusters)

	if generateESXIEntries {

	}

	outFile, _ := json.MarshalIndent(configData, "", " ")
	error := ioutil.WriteFile("/Users/sargonbenjamin/Downloads/out1.json", outFile, 0644)
	if error != nil {
		fmt.Println("Couldn't write to json file", error)
	}
}

func GetConfigObject(filePath string) models.FileContent {

	theFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var configData models.FileContent

	err = json.Unmarshal(theFile, &configData)

	if err != nil {
		fmt.Println("Error during unmarshal", err)
	}

	return configData
}
