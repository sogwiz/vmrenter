package main

import (
	"fmt"
	"vmrenter/pkg/mapr"
	"vmrenter/pkg/utils"
)

var nodesTable = "/user/mapr/nodes"
var csvFilePath = "/home/user6bb0/Work/vm-renter/my_nodes.csv"

func main() {

	partialNodes := mapr.GetPartialReservationsForNodesUpdate()
	fmt.Println(partialNodes)
	err := mapr.Reset(nodesTable)
	if err != nil {
		fmt.Printf("Error occured: %v", err)
	}

	nodes := utils.GetNodesFromCSV(csvFilePath)
	fmt.Println(nodes)
	//SeedNodesTable(csvFilePath)
}
