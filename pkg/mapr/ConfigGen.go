package mapr

import (
	"encoding/json"
	"go.uber.org/zap"
	"io/ioutil"
	"path"
	"path/filepath"
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

	zap.S().Info(configData.Clusters)

	if generateESXIEntries {

	}

	outFile, _ := json.MarshalIndent(configData, "", " ")
	outDir, _ := path.Split(sourceConfigFilepath)
	outFilepath := filepath.FromSlash(outDir + "/out.json")
	error := ioutil.WriteFile(outFilepath, outFile, 0644)
	if error != nil {
		zap.S().Error("Couldn't write to json file", error)
	}
	zap.S().Infof("Wrote cluster %v to %v", reservation.ClusterID, outFilepath)
}

func GetConfigObject(filePath string) models.FileContent {

	theFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		zap.S().Fatal(err)
	}

	var configData models.FileContent

	err = json.Unmarshal(theFile, &configData)

	if err != nil {
		zap.S().Error("Error during unmarshal", err)
	}

	return configData
}
