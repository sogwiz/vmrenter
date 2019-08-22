package main

import (
	"go.uber.org/zap"
	"log"
	"os"
	"strconv"
	zaplogger "vmrenter/logger"

	"vmrenter/pkg/config"

	"vmrenter/pkg/mapr"
	"vmrenter/pkg/models"

	"gopkg.in/urfave/cli.v2"
)

var configData models.FileContent

const maximumRentingTime = 168 // in hours

func getNodesInCluster(clusterID string) []models.Node {
	res := make([]models.Node, 0)
	for _, cluster := range configData.Clusters {
		if cluster.ID == clusterID {
			res = append(res, cluster.Nodes...)
		}
	}
	return res
}

func getAvailableSharedNodes(operatingSystem string) (n []models.Node) {
	//sharedNodes := getNodesInCluster("sharedClusterId")
	return
}

func start(c *cli.Context) error {

	logLevel := c.String("loglevel")
	err := zaplogger.ConfigureLogger(logLevel)
	if err != nil {
		panic(err)
	}

	filePath := c.String("file")
	clusterID := c.String("cluster")
	config.SetURLDBConn(c.String("urldbconn"))
	dbConn := config.GetURLDBConn()
	requestedNumNodes := c.Int("nodes")
	requestedOperatingSystem := c.String("os")
	emailAddr := c.String("email")
	ram := c.Int("ram")
	hoursToReserve := c.Int("hourstoreserve")
	osVersion := c.String("nodeosversion")

	if !c.IsSet("hourstoreserve") {
		hoursToReserve = 24
		zap.S().Warnf("Time for renting is not specified. The nodes will be reserved for %v hours from now", hoursToReserve)
	} else {
		switch {
		case hoursToReserve <= 0:
			zap.S().Fatalf("You have specified invalid numer of hours - %v. Aborting reservation.\n", hoursToReserve)
		case hoursToReserve > maximumRentingTime:
			zap.S().Fatalf("Your time for renting is %v hours - exceeds maximum amount of hours - %v (1 week). "+
				"Aborting reservation.\n", hoursToReserve, maximumRentingTime)
		}
	}

	zap.S().Infof("**** Config **** \nfilePath=", filePath, "\nclusterId=", clusterID,
		"\nNodes Requested=", requestedNumNodes, "\nURL_DB_CONN=", dbConn, "\nTime for rent (in hours):", hoursToReserve,
		"\nRequested operating system:", requestedOperatingSystem, "\nOs version:", osVersion)

	//configData := mapr.GetConfigObject(filePath)

	nodes := mapr.GetAvailableNodes("sharedpool", requestedOperatingSystem, osVersion)
	if len(nodes) < requestedNumNodes {
		zap.S().Fatalf("Can't fulfill request. Only %v nodes available matching your request requirements", strconv.Itoa(len(nodes)))
	}

	if len(nodes) == 0 {
		zap.S().Fatal("Must submit a non-zero number of nodes. eg -n 1")
	}

	// Check if a constraint for RAM is posed at all
	switch {
	case !c.IsSet("ram"): // Check if the flag is set at all
		zap.S().Warn("You have not provided RAM constraint for VMs in a cluster , thus RAM limitations will be neglected")
	case ram <= 0:
		zap.S().Fatalf("You have provided invalid value for RAM constraint - %v. Aborting reserving.", ram)
	case ram > 0: // Constraint is present
		// Check if there are enough nodes with RAM equal or more than needed
		numOfRAMPassingNodes := 0 // Number of nodes that adhere to RAM constraints
		for i := range nodes {
			if nodes[i].NodeObj.RAM >= ram {
				numOfRAMPassingNodes += 1
			}
		}
		if numOfRAMPassingNodes < requestedNumNodes {
			zap.S().Fatalf("You are trying to reserve %d nodes. There are %d matching nodes but only %d nodes have %d or more RAM", requestedNumNodes, len(nodes), numOfRAMPassingNodes, ram)
		}
	}

	reservation, err := mapr.MakeReservation(clusterID, emailAddr, nodes[0:requestedNumNodes], "http://jenkinshost:jenkinsport/view/VIEW_NAME/job/JOB_NAME/5607/", "vmsonly", hoursToReserve)
	if err != nil {
		zap.S().Fatalf("error calling MakeReservation", err)
	}

	//fmt.Println("Reservation=", reservation)
	mapr.GenerateConfigJson(reservation, false, filePath)

	return nil
}

func main() {

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Usage:   "The location of the config.json file",
			},
			&cli.StringFlag{
				Name:    "cluster",
				Aliases: []string{"c"},
				Value:   "sharedcluster",
				Usage:   "The cluster id",
			},
			&cli.StringFlag{
				Name:    "os",
				Aliases: []string{"o"},
				Value:   "centos",
				Usage:   "The requested operating system",
			},
			&cli.StringFlag{
				Name:    "urldbconn",
				Aliases: []string{"u"},
				Value:   "DBHOST:DBPORT?auth=basic;user=USERNAME;password=PASSWORD;ssl=false",
				Usage:   "DB Connection URL",
				EnvVars: []string{"URL_DB_CONN"},
			},
			&cli.StringFlag{
				Name:    "email",
				Aliases: []string{"e"},
				Value:   "default@mapr.com",
				Usage:   "Your Email Address",
			},
			&cli.IntFlag{
				Name:    "nodes",
				Aliases: []string{"n"},
				Value:   0,
				Usage:   "number of nodes requested to reserve",
			},
			&cli.IntFlag{
				Name:    "ram",
				Aliases: []string{"m"},
				Usage:   "VMs RAM in gigabytes. All vms in the cluster should have equal or more than specified RAM",
			},
			&cli.IntFlag{
				Name:    "hourstoreserve",
				Aliases: []string{"hours"},
				Value:   24,
				Usage:   "Number of hours for reservation",
			},
			&cli.StringFlag{
				Name:    "osversion",
				Aliases: []string{"nodeosversion"},
				Usage:   "Specific version of OS of the node",
			},
			&cli.StringFlag{
				Name:    "loglevel",
				Aliases: []string{"l"},
				Usage:   "Log level",
			},
		},
		Name:   "vmrenter",
		Usage:  "Parameters Usage",
		Action: start,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("\t*** FATAL ERROR *** \n\tError occured while running the app: %v", err)
	}

	//for k, v := range m {
	//	switch vv := v.(type) {
	//	case string:
	//		fmt.Println(k, "is string", vv)
	//	case float64:
	//		fmt.Println(k, "is float64", vv)
	//	case []interface{}:
	//		fmt.Println(k, "is an array:")
	//		for i, u := range vv {
	//			fmt.Println(i, u)
	//		}
	//	case map[string]interface{}:
	//		fmt.Println(k, "is dictionary", vv)
	//
	//	default:
	//		fmt.Println(k, "is of a type I don't know how to handle")
	//	}

}
