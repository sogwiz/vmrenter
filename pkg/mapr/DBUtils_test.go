package mapr

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
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
	fmt.Println(IsRequestDoable(1, "centos", "7"))
}

func TestMakeReservation(t *testing.T) {
	numNodes := 1
	nodes := GetAvailableNodes("sharedpool", "centos", "7.3")
	if len(nodes) < numNodes {
		panic("Can't fulfill request. exiting")
	}

	reservation, error := MakeReservation("sharedcluster", "sbenjamin@mapr.com", nodes[0:numNodes], "jenkinsurl", "vmsonly", 24)

	if error != nil {
		fmt.Println("Error making reservation", error)
	}

	fmt.Println(reservation)
	GenerateConfigJson(reservation, false, "/home/vlad/Work/sample-config.json")
}

func TestReset(t *testing.T) {
	Reset(tableReservations)
	Reset(tableNodes)
}

func workerUnreserver(id int, jobs <-chan models.NodeDBJson, results chan<- int) {
	for node := range jobs {
		fmt.Println("worker", id, "started  job", node)
		err := ReserveNode(node.ID, "", "")
		if err != nil {
			fmt.Println("Found error for node ", node, err)
		}
		time.Sleep(time.Second)
		results <- id
	}
}

func TestUnreserveNodes(t *testing.T) {
	nodes := getAllNodes()
	fmt.Println("nodes length: ", len(nodes))

	//only unreserve 5 at a time, concurrently
	jobs := make(chan models.NodeDBJson, len(nodes))
	results := make(chan int, len(nodes))

	for w := 1; w <= 5; w++ {
		//fmt.Println("iteration ", node.ID)
		go workerUnreserver(w, jobs, results)
	}

	for _, node := range nodes {
		jobs <- node
	}
	close(jobs)

	for i := 0; i < len(nodes); i++ {
		fmt.Println("Finished with result", <-results)
	}
}

func TestReserveNodes(t *testing.T) {
	now := time.Now()
	expiry := now.Add(108 * time.Hour)

	nodeIds := [2]string{"node.10.10.108.231", "node.10.10.108.233"}

	for _, nodeId := range nodeIds {
		err := ReserveNode(nodeId, expiry.Format(time.RFC3339), "sargon")
		if err != nil {
			fmt.Println("Error reserving node", err)
		}
	}

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
