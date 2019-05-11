package main

import (
	"fmt"
	"os"
	"strconv"

	"vmrenter/pkg/config"

	"vmrenter/pkg/mapr"
	"vmrenter/pkg/models"

	"gopkg.in/urfave/cli.v2"
)

var configData models.FileContent

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
	filePath := c.String("file")
	clusterID := c.String("cluster")
	config.SetURLDBConn(c.String("urldbconn"))
	dbConn := config.GetURLDBConn()
	requestedNumNodes := c.Int("nodes")
	requestedOperatingSystem := c.String("os")
	emailAddr := c.String("email")

	fmt.Println("**** Config **** \nfilePath=", filePath, "\nclusterId=", clusterID, "\nNodes Requested=", requestedNumNodes, "\nURL_DB_CONN=", dbConn)

	//configData := mapr.GetConfigObject(filePath)

	nodes := mapr.GetAvailableNodes("sharedpool", requestedOperatingSystem)
	if len(nodes) < requestedNumNodes {
		errorStr := "Can't full request. Only " + strconv.Itoa(len(nodes)) + " nodes available matching your request requirements"
		panic(errorStr)
	}

	if len(nodes) == 0 {
		panic("Must submit a non-zero number of nodes. eg -n 1")
	}

	reservation, err := mapr.MakeReservation(clusterID, emailAddr, nodes[0:requestedNumNodes], "http://jenkinshost:jenkinsport/view/VIEW_NAME/job/JOB_NAME/5607/", "vmsonly")
	if err != nil {
		fmt.Println("error calling MakeReservation", err)
	}

	fmt.Println("Reservation=", reservation)
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
		},
		Name:   "vmrenter",
		Usage:  "Parameters Usage",
		Action: start,
	}

	app.Run(os.Args)

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
