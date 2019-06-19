package utils

import (
	"fmt"
	"sync"
	"testing"
	"vmrenter/pkg/mapr"

	"github.com/stretchr/testify/assert"
)

func TestGetNodesFromCSV(t *testing.T) {
	csvFilename := "/Users/sargonbenjamin/Downloads/nodes.csv"
	nodes := GetNodesFromCSV(csvFilename)
	assert.True(t, len(nodes) > 0, "couldn't load nodes from csv to memory data model")
}

func TestGetNodeJsonDocString(t *testing.T) {
	csvFilename := "/Users/sargonbenjamin/Downloads/nodes.csv"
	nodes := GetNodesFromCSV(csvFilename)

	nodeStrings := make([]string, 0)

	for _, node := range nodes {
		nodeJsonStr := getNodeJsonDocString(node)
		fmt.Println(nodeJsonStr)
		nodeStrings = append(nodeStrings, nodeJsonStr)
	}

	assert.True(t, len(nodeStrings) > 0, "couldn't load nodes to json strings")
}


func TestDataSeed(t *testing.T) {
	//csvFilename := "/home/user6bb0/Work/vm-renter/nodes.csv"
	csvFilename := "/home/user6bb0/Work/vm-renter/my_nodes.csv"
	nodes := GetNodesFromCSV(csvFilename)

	listOfMaps := make([]map[string]interface{}, 0)

	var wg sync.WaitGroup
	for _, node := range nodes {
		//nodeJsonStr := getNodeJsonDocString(node)
		mapIntface := GetNodeJsonDocMap(node)
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

	assert.True(t, len(listOfMaps) > 0, "couldn't load nodes to map")

}

func TestGetNodeOperatingSystems(t *testing.T) {
	/*ips := []string{"10.10.99.165",
		"10.10.99.171",
		"10.10.99.172",
		"10.10.99.173",
		"10.10.99.174",
		"10.10.99.176",
		"10.10.108.241",
		"10.10.111.21",
		"10.10.111.22",
		"10.10.111.23",
		"10.10.111.24",
		"10.10.111.26",
		"10.10.111.27",
		"10.10.111.28",
		"10.10.111.29",
		"10.10.111.30",
		"10.10.111.32",
		"10.10.111.33",
		"10.10.111.34",
		"10.10.111.35",
		"10.10.111.36",
		"10.10.111.37",
		"10.10.111.38",
		"10.10.111.39",
		"10.10.111.40",
	}*/
	ips := []string{"10.10.99.178",
		"10.10.99.179", "10.10.99.181", "10.10.99.182", "10.10.99.183",
		"10.10.30.71", "10.10.30.72",
	}
	nodes := getNodeOperatingSystems(ips)
	fmt.Println(nodes)
}
