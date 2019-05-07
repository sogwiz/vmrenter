package mapr

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"testing"
	"vmrenter/pkg/models"

	"github.com/stretchr/testify/assert"
)

func TestGetUnavailableNodes(t *testing.T) {
	getUnavailableNodes("", "")
	assert.Equal(t, 1, 1)
}

func TestGetAvailableNodes(t *testing.T) {
	//getAvailableNodes("", "centos")
	//assert.Equal(t, 1, 1)
	fmt.Println(IsRequestDoable(20, "centos", "7"))
}

func TestMakeReservation(t *testing.T) {
	numNodes := 1
	nodes := GetAvailableNodes("sharedpool", "centos")
	if len(nodes) < numNodes {
		panic("Can't fulllfill request. exiting")
	}

	reservation, error := MakeReservation("sharedcluster", "sargon", nodes[0:numNodes], "jenkinsurl", "vmsonly")

	if error != nil {
		fmt.Println("Error making reservation", error)
	}

	fmt.Println(reservation)
	GenerateConfigJson(reservation, false, "/Users/sargonbenjamin/dev/src/private-installer/testing/configuration/config.json")
}

func TestReset(t *testing.T) {
	reset(tableReservations)
	reset(tableNodes)
}

func TestUnreserveNodes(t *testing.T) {
	nodes := getAllNodes()
	fmt.Println("nodes length: ", len(nodes))

	var wg sync.WaitGroup
	for _, node := range nodes {
		wg.Add(1)
		//fmt.Println("iteration ", node.ID)
		go func(node models.NodeDBJson) {
			defer wg.Done()

			err := ReserveNode(node.ID, "", "")
			if err != nil {
				fmt.Println("Found error for node ", node, err)
			}

		}(node)
	}
	wg.Wait()

	/*
		c := make(chan struct{}) // We don't need any data to be passed, so use an empty struct

		for _, node := range nodes {

			//fmt.Println("iteration ", node.ID)
			go func(node models.NodeDBJson) {

				err := ReserveNode(node.ID, "", "")
				if err != nil {
					fmt.Println("Found error for node ", node, err)
				}
				c <- struct{}{} // signal that the routine has completed
			}(node)
		}
		for i := 0; i < len(nodes); i++ {
			<-c
		}
	*/

}

func TestDeleteAndCreateTable(t *testing.T) {

	url := "https://mapr:mapr@10.10.99.165:8443/rest/table/delete"

	payload := strings.NewReader("path=/user/mapr/nodes")

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "*/*")
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}
