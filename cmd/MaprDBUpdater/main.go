package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v2"
	"log"
	"os"
	"sync"
	"vmrenter/pkg/config"
	"vmrenter/pkg/mapr"
	"vmrenter/pkg/models"
	"vmrenter/pkg/utils"
)

var nodesTable = "/user/mapr/nodes"

// Creating NodeDBJsons from nodes
func updateNodeDBJsons(partialNodes []models.PartialReservationForNodesUpdate, listOfMaps []map[string]interface{}) {
	// Updating NodeDBJsons with ClusterID, ExpiresAT
	wg1 := sync.WaitGroup{}
	m := &sync.Mutex{}
	for _, partialNode := range partialNodes {
		wg1.Add(1)
		go func(partialNode models.PartialReservationForNodesUpdate) {
			defer wg1.Done()
			for i := range listOfMaps {
				nodeDBJson := listOfMaps[i]

				m.Lock()
				if partialNode.ID == nodeDBJson["_id"] {
					nodeDBJson["ClusterID"] = partialNode.ClusterID
					nodeDBJson["ExpiresAT"] = partialNode.ExpiresAt
				}
				m.Unlock()
			}
		}(partialNode)
	}

	wg1.Wait()
}


func update(c *cli.Context) error {

	if !c.IsSet("urldbconn") {
		fmt.Println("Connection string isn't set, aborting")
		return nil
	}

	if !c.IsSet("nodesfile") {
		fmt.Println("File path string isn't set, aborting")
		return nil
	}
	config.SetURLDBConn(c.String("urldbconn"))
	dbConn := config.GetURLDBConn()
	var csvFilePath = c.String("nodesfile")

	fmt.Printf("**** Config **** %v\n", dbConn)

	partialNodes, err := mapr.ExtractPartialNodesData()
	if err != nil {
		panic(err)
	}

	// Getting nodes from csv file
	fmt.Println("Starting getting nodes from the csv file...")
	nodes := utils.GetNodesFromCSV(csvFilePath)
	fmt.Println("Finished getting nodes from the csv file!")

	listOfMaps := utils.CreateNodeDBJsons(nodes)

	updateNodeDBJsons(partialNodes, listOfMaps)

	// Resetting /user/mapr/node table
	fmt.Println("Starting resetting nodes table...")
	resetErr := mapr.Reset(nodesTable)
	if resetErr != nil {
		fmt.Printf("Error occurred while resetting /user/mapr/nodes table: %v", err)
		return resetErr
	}
	fmt.Println("Finished resetting nodes table!")

	mapr.UpdateNodesTable(listOfMaps)
	return nil

}

func main() {

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "urldbconn",
				Aliases: []string{"u"},
				Value:   "DBHOST:DBPORT?auth=basic;user=USERNAME;password=PASSWORD;ssl=false",
				Usage:   "DB Connection URL",
				EnvVars: []string{"URL_DB_CONN"},
			},
			&cli.StringFlag{
				Name:    "nodesfile",
				Aliases: []string{"f"},
				Usage:   "Location of 'nodes' file",
			},
		},
		Name:   "vmrenter",
		Usage:  "Parameters Usage",
		Action: update,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("\t*** FATAL ERROR *** \n\tError occured while running the app: %v", err)
	}
}
