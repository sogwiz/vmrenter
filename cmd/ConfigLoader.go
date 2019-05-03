package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"vmrenter/pkg/config"

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

	fmt.Println("Using config", "filePath=", filePath, "clusterId=", clusterID)

	theFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error during readFile", err)
	}

	//var configData map[string]interface{}

	err = json.Unmarshal(theFile, &configData)

	if err != nil {
		fmt.Println("Error during unmarshal", err)
	}

	/*=
	nodes := getNodesInCluster(clusterID)
	for _, node := range nodes {
		fmt.Println(node.Host)
	}

	_, err := mapr.MakeReservation(clusterID, nodes, "http://jenkinshost:jenkinsport/view/VIEW_NAME/job/JOB_NAME/5607/", "vmsonly")
	if err != nil {
		fmt.Println("error calling MakeReservation", err)
	}
	*/

	return nil
}

func main() {

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Value:   "/Users/sargonbenjamin/dev/src/private-installer/testing/configuration/config.json",
				Usage:   "The location of the config.json file",
			},
			&cli.StringFlag{
				Name:    "cluster",
				Aliases: []string{"c"},
				Value:   "CLUSTER_ID",
				Usage:   "The cluster id",
			},
			&cli.StringFlag{
				Name:    "urldbconn",
				Aliases: []string{"u"},
				Value:   "DBHOST:DBPORT?auth=basic;user=USERNAME;password=PASSWORD;ssl=false",
				Usage:   "DB Connection URL",
				EnvVars: []string{"URL_DB_CONN"},
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
