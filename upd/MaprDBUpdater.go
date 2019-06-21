package main

import (
	"fmt"
	"sync"
	"vmrenter/pkg/mapr"
	"vmrenter/pkg/models"
	"vmrenter/pkg/utils"
)

var nodesTable = "/user/mapr/nodes"
var csvFilePath = "/home/user6bb0/Work/vm-renter/my_nodes.csv"

func main() {

	// Getting nodes id, ExpiresAt and ClusterID from /user/mapr/nodes table
	fmt.Println("Starting getting nodes id, ExpiresAT, ClusterID...")
	partialNodes := mapr.GetPartialReservationsForNodesUpdate()
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

	// Creating NodeDBJsons from nodes
	var wg sync.WaitGroup
	nodeDbJsonQueue := make(chan map[string]interface{}, 1)
	wg.Add(len(nodes))

	fmt.Println("Starting creating NodeDBJsons from nodes...")
	for _, node := range nodes {
		go func(node models.Node) {
			mapIntface := utils.GetNodeJsonDocMap(node)
			mapIntface["_id"] = node.ID
			nodeDbJsonQueue <- mapIntface
		}(node)
	}

	go func() {
		for n := range nodeDbJsonQueue {
			listOfMaps = append(listOfMaps, n)
			wg.Done()
		}
	}()

	wg.Wait()
	fmt.Println("Finished creating NodeDBJsons from nodes!")

	// Updating NodeDBJsons with ClusterID, ExpiresAT
	wg1 := sync.WaitGroup{}
	for _, partialNode := range partialNodes {
		wg1.Add(1)
		go func(partialNode models.PartialReservationForNodesUpdate) {
			defer wg1.Done()
			for i := range listOfMaps {
				nodeDBJson := listOfMaps[i]

				if partialNode.ID == nodeDBJson["_id"] {
					nodeDBJson["ClusterID"] = partialNode.ClusterID
					nodeDBJson["ExpiresAT"] = partialNode.ExpiresAt
				}
			}
		}(partialNode)
	}

	wg1.Wait()

	// Resetting /user/mapr/node table
	fmt.Println("Starting resetting nodes table...")
	resetErr := mapr.Reset(nodesTable)
	if resetErr != nil {
		fmt.Printf("Error occurred while resetting /user/mapr/nodes table: %v", err)
		return
	}
	fmt.Println("Finished resetting nodes table!")

	// Updating the nodes table
	fmt.Println("Starting writing to nodes table...")

	for _, mapIntface := range listOfMaps {
		writeErr := mapr.WriteToDBWithTableMap(mapIntface, "/user/mapr/nodes")
		if writeErr != nil {
			fmt.Println("Error writing to table", writeErr)
		}
	}
	fmt.Println("Finished writing to nodes table!")

	fmt.Println("Starting writing to nodes table...")

	// Asynchronous writing to the table - fails because of syncPut() Investigate further
	//var wg2 = sync.WaitGroup{}
	//for _, mapIntface := range listOfMaps {
	//	wg2.Add(1)
	//	go func(mapIntface map[string]interface{}) {
	//		defer wg2.Done()
	//		writeErr := mapr.WriteToDBWithTableMap(mapIntface, "/user/mapr/nodes")
	//		if writeErr != nil {
	//			fmt.Println("Error writing to table", writeErr)
	//		}
	//	}(mapIntface)
	//}
	//wg2.Wait()
	//fmt.Println("Finished writing to nodes table!")

}
