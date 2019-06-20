package main

import (
	"fmt"
	"sync"
	"vmrenter/pkg/mapr"
	"vmrenter/pkg/models"
	"vmrenter/pkg/utils"
)

var nodesTable = "/user/mapr/nodes"
var csvFilePath = "/home/user6bb0/Work/vm-renter/my_nodes1.csv"

func main() {

	fmt.Println("Starting getting nodes id, ExpiresAT, ClusterID...")
	partialNodes := mapr.GetPartialReservationsForNodesUpdate() // ID doesn't work, but _id does - investigate why
	err := mapr.Reset(nodesTable)
	if err != nil {
		fmt.Printf("Error occured while resetting /user/mapr/nodes table: %v", err)
		return
	}
	fmt.Println("Finished getting nodes id, ExpiresAT, ClusterID!")

	// Getting nodes from csv file
	fmt.Println("Starting getting nodes from the csv file...")
	nodes := utils.GetNodesFromCSV(csvFilePath)
	fmt.Println("Finished getting nodes from the csv file!")

	listOfMaps := make([]map[string]interface{}, 0) // List of NodeDBJsons

	var wg sync.WaitGroup

	// Creating NodeDBJsons from nodes
	fmt.Println("Starting creating NodeDBJsons from nodes...")
	for _, node := range nodes {
		wg.Add(1)
		go func(node models.Node) {
			defer wg.Done()
			mapIntface := utils.GetNodeJsonDocMap(node)
			mapIntface["_id"] = node.ID
			listOfMaps = append(listOfMaps, mapIntface)
		}(node)
	}

	wg.Wait()
	fmt.Println("Finished creating NodeDBJsons from nodes!")

	// Update NodeDBJsons with ClusterID, ExpiresAT
	for _, partialNode := range partialNodes {
	comparingOnePartialNodeToListOfMaps:
		for i := range listOfMaps {
			nodeDBJson := listOfMaps[i]
			if partialNode.ID == nodeDBJson["_id"] {
				nodeDBJson["ClusterID"] = partialNode.ClusterID
				nodeDBJson["ExpiresAT"] = partialNode.ExpiresAt
				break comparingOnePartialNodeToListOfMaps
			}
		}
	}

	fmt.Println("Starting resetting nodes table...")
	resetErr := mapr.Reset(nodesTable)
	if resetErr != nil {
		fmt.Printf("Error occurred while resetting /user/mapr/nodes table: %v", err)
		return
	}
	fmt.Println("Finished resetting nodes table!")

	// Updating the nodes table
	//var wg2 sync.WaitGroup
	fmt.Println("Starting writing to nodes table...")
	for _, mapIntface := range listOfMaps {

		//wg2.Add(1)
		//go func(mapIntface map[string]interface{}) {
		//defer wg2.Done()
		error := mapr.WriteToDBWithTableMap(mapIntface, "/user/mapr/nodes")
		if error != nil {
			fmt.Println("Error writing to table", error)
		}
		//}(mapIntface)
	}
	//wg2.Wait()
	fmt.Println("Finished writing to nodes table!")

}

//for
//SeedNodesTable(csvFilePath)
