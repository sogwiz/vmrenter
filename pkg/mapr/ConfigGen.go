package mapr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"vmrenter/pkg/models"
)

func GenerateConfigJson(reservation models.Reservation, generateESXIEntries bool, sourceConfigFilepath string) {
	var configData models.FileContent
	configData = GetConfigObject(sourceConfigFilepath)

	//var nodes [len(reservation.Nodes)]models.Node

	expectedServiceNames := make([]string, 0)
	expectedServiceNames = append(expectedServiceNames, "mapr-core")
	nodes := make([]models.Node, len(reservation.Nodes))
	for i, nodejson := range reservation.Nodes {
		nodes[i] = nodejson.NodeObj
		nodes[i].ExpectedServiceNames = expectedServiceNames
	}

	/*
		var esxiServerID string
		if len(nodes[0].EsxiServerID) > 0 {
			esxiServerID = nodes[0].EsxiServerID
		}
		clusterToReserve := models.Cluster{
			ID:           reservation.ClusterID,
			Name:         reservation.ClusterID,
			Nodes:        nodes,
			EsxiServerID: esxiServerID,
		}
	*/

	clusterToReserve := models.Cluster{
		ID:    reservation.ClusterID,
		Name:  reservation.ClusterID,
		Nodes: nodes,
	}

	//configData.Docker_Nodes = make([]map[string]interface{}, 0)
	//configData.Clusters = make([]models.Cluster, 0)
	configData.Clusters = append(configData.Clusters, clusterToReserve)

	fmt.Println(configData.Clusters)

	if generateESXIEntries {

	}

	outFile, _ := json.MarshalIndent(configData, "", " ")
	outFilepath := "out.json"
	error := ioutil.WriteFile(outFilepath, outFile, 0644)
	if error != nil {
		fmt.Println("Couldn't write to json file", error)
	}
	fmt.Println("Wrote cluster " + reservation.ClusterID + " to " + outFilepath)
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
