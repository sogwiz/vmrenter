package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"vmrenter/pkg/mapr"
	"vmrenter/pkg/models"
)

func TestGetNodesFromCSV(t *testing.T) {
	csvFilename := "/home/vlad/Work/nodes.csv"
	nodes := GetNodesFromCSV(csvFilename)
	assert.True(t, len(nodes) > 0, "couldn't load nodes from csv to memory data model")
}

func TestGetNodeJsonDocString(t *testing.T) {
	csvFilename := "/home/vlad/Work/nodes.csv"
	nodes := GetNodesFromCSV(csvFilename)

	nodeStrings := make([]string, 0)

	for _, node := range nodes {
		nodeJsonStr := getNodeJsonDocString(node)
		fmt.Println(nodeJsonStr)
		nodeStrings = append(nodeStrings, nodeJsonStr)
	}

	assert.True(t, len(nodeStrings) > 0, "couldn't load nodes to json strings")
}

func writeNodeToDb(goroutineId int, jobs <-chan models.Node, results chan<- map[string]interface{}) {
	for node := range jobs {
		fmt.Printf("Worker %d starts job %v", goroutineId, node)
		mapIntface := GetNodeJsonDocMap(node)
		mapIntface["_id"] = node.ID
		err := mapr.WriteToDBWithTableMap(mapIntface, "/user/mapr/nodes")
		if err != nil {
			fmt.Println("Error writing to table", err)
		}
		results <- mapIntface
	}
}

func TestDataSeed(t *testing.T) {
	csvFilename := "/home/vlad/Work/nodes.csv"
	nodes := GetNodesFromCSV(csvFilename)

	listOfMaps := make([]map[string]interface{}, 0)

	// *** Usual for loop ***
	//for _, node := range nodes {
	//	mapIntface := GetNodeJsonDocMap(node)
	//	mapIntface["_id"] = node.ID
	//	err := mapr.WriteToDBWithTableMap(mapIntface, "/user/mapr/nodes")
	//	if err != nil {
	//		fmt.Println("Error writing to table", err)
	//	}
	//	listOfMaps = append(listOfMaps, mapIntface)
	//}

	// *** Fully async way with goroutine for every node ***
	//var wg sync.WaitGroup
	//wg.Add(len(nodes))
	//
	//queue := make(chan map[string]interface{}, 1)
	//
	//for _, node := range nodes {
	//	go func(node models.Node) {
	//		mapIntface := GetNodeJsonDocMap(node)
	//		mapIntface["_id"] = node.ID
	//		err := mapr.WriteToDBWithTableMap(mapIntface, "/user/mapr/nodes")
	//		if err != nil {
	//			fmt.Println("Error writing to table", err)
	//		}
	//		queue <- mapIntface
	//	}(node)
	//}
	//
	//go func() {
	//	for n := range queue {
	//		listOfMaps = append(listOfMaps, n)
	//		wg.Done()
	//	}
	//}()
	//wg.Wait()

	// Thread pool analog
	var jobs = make(chan models.Node, len(nodes))
	var results = make(chan map[string]interface{}, len(nodes))

	var goroutineCount = 5

	for i := 1; i <= goroutineCount; i++ {
		go writeNodeToDb(i, jobs, results)
	}

	for _, node := range nodes {
		jobs <- node
	}
	close(jobs)

	for i := 0; i < len(nodes); i++ {
		res := <-results
		listOfMaps = append(listOfMaps, res)
		fmt.Println("Finished with result", res["_id"])
	}

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
