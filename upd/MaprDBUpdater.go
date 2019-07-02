package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v2"
	"log"
	"os"
	"sync"
	zaplogger "vmrenter/logger"
	"vmrenter/pkg/config"
	"vmrenter/pkg/mapr"
	"vmrenter/pkg/models"
	"vmrenter/pkg/utils"
)

var nodesTable = "/user/mapr/nodes"

func update(c *cli.Context) error {

	if !c.IsSet("urldbconn") {
		fmt.Println("Connection string isn't set, aborting")
		return nil
	}

	if !c.IsSet("nodesfile") {
		fmt.Println("File path string isn't set, aborting")
		return nil
	}

	logLevel := c.String("loglevel")
	lgr, err := zaplogger.BuildLogger(logLevel)
	if err != nil {
		panic(err)
	}

	lgr.Info("Good!")

	config.SetURLDBConn(c.String("urldbconn"))
	dbConn := config.GetURLDBConn()
	var csvFilePath = c.String("nodesfile")

	fmt.Printf("**** Config **** %v\n", dbConn)

	// Getting nodes id, ExpiresAt and ClusterID from /user/mapr/nodes table
	fmt.Println("Starting getting nodes id, ExpiresAT, ClusterID...")
	partialNodes := mapr.GetPartialReservationsForNodesUpdate()
	err = mapr.Reset(nodesTable)
	if err != nil {
		fmt.Printf("Error occured while resetting /user/mapr/nodes table: %v", err)
		return err
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

	// Resetting /user/mapr/node table
	fmt.Println("Starting resetting nodes table...")
	resetErr := mapr.Reset(nodesTable)
	if resetErr != nil {
		fmt.Printf("Error occurred while resetting /user/mapr/nodes table: %v", err)
		return resetErr
	}
	fmt.Println("Finished resetting nodes table!")

	// Updating the nodes table
	fmt.Println("Starting writing to nodes table...")

	// Synchronous way to update table until the error with goroutines is fixed
	for _, mapIntface := range listOfMaps {
		writeErr := mapr.WriteToDBWithTableMap(mapIntface, "/user/mapr/nodes")
		if writeErr != nil {
			fmt.Println("Error writing to table", writeErr)
		}
	}

	// Asynchronous writing to the table - fails because of syncPut(). Uncomment when the bug is fixed.
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
	fmt.Println("Finished writing to nodes table!")

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
			&cli.StringFlag{
				Name:    "loglevel",
				Aliases: []string{"l"},
				Usage:   "Log level",
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
