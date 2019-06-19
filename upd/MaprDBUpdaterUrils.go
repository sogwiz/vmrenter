package main

import (
	"fmt"
	"sync"
	"vmrenter/pkg/mapr"
	"vmrenter/pkg/utils"
)

func SeedNodesTable(csvFilePath string) {
	nodes := utils.GetNodesFromCSV(csvFilePath)

	listOfMaps := make([]map[string]interface{}, 0)

	var wg sync.WaitGroup
	for _, node := range nodes {
		mapIntface := utils.GetNodeJsonDocMap(node)
		mapIntface["_id"] = node.ID
		listOfMaps = append(listOfMaps, mapIntface)

		wg.Add(1)
		go func(mapIntface map[string]interface{}) {
			defer wg.Done()
			error := mapr.WriteToDBWithTableMap(mapIntface, "/user/mapr/nodes")
			if error != nil {
				fmt.Println("Error writing to table", error)
			}
		}(mapIntface)
	}
	wg.Wait()

}

