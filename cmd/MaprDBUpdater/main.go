package main

import (
	"log"
	"os"
	"sync"
	zaplogger "vmrenter/logger"
	"vmrenter/pkg/config"
	"vmrenter/pkg/mapr"
	"vmrenter/pkg/models"
	"vmrenter/pkg/utils"

	"go.uber.org/zap"
	"gopkg.in/urfave/cli.v2"
)

var nodesTable = "/user/mapr/nodes"

// Creating NodeDBJsons from nodes
func updateNodeDBJsons(partialNodes []models.PartialReservationForNodesUpdate, listOfMaps []map[string]interface{}) {
	// Updating NodeDBJsons with ClusterID, ExpiresAt
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

	logLevel := c.String("loglevel")
	err := zaplogger.ConfigureLogger(logLevel)
	if err != nil {
		panic(err)
	}

	if !c.IsSet("urldbconn") {
		zap.S().Fatalf("Connection string isn't set, aborting")
		return nil
	}

	if !c.IsSet("nodesfile") {
		zap.S().Fatalf("File path string isn't set, aborting")
		return nil
	}
	config.SetURLDBConn(c.String("urldbconn"))
	dbConn := config.GetURLDBConn()
	var csvFilePath = c.String("nodesfile")

	zap.S().Infof("**** Config **** %v\n", dbConn)

	partialNodes, err := mapr.ExtractPartialNodesData()
	if err != nil {
		zap.S().Panic(err)
	}

	// Getting nodes from csv file
	zap.S().Infof("Starting getting nodes from the csv file...")
	nodes := utils.GetNodesFromCSV(csvFilePath)
	zap.S().Infof("Finished getting nodes from the csv file!")

	listOfMaps := utils.CreateNodeDBJsons(nodes)

	updateNodeDBJsons(partialNodes, listOfMaps)

	// Resetting /user/mapr/node table
	zap.S().Infof("Starting resetting nodes table...")
	resetErr := mapr.Reset(nodesTable)
	if resetErr != nil {
		zap.S().Errorf("Error occurred while resetting /user/mapr/nodes table: %v", err)
		return resetErr
	}
	zap.S().Info("Finished resetting nodes table!")

	mapr.UpdateNodesTable(listOfMaps)
	return nil

}

func main() {

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "urldbconn",
				Aliases: []string{"u"},
				Value:   "",
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
