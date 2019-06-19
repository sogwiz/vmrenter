package utils

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"vmrenter/pkg/models"

	"github.com/fatih/structs"
)

/*
type SSHCommander struct {
	User string
	IP   string
}

func (s *SSHCommander) Command(cmd ...string) *exec.Cmd {
	arg := append(
		[]string{
			fmt.Sprintf("%s@%s", s.User, s.IP),
		},
		cmd...,
	)

	return exec.Command("/bin/ssh", arg...)
}

func main() {
	commander := SSHCommander{"root", "10.10.99.165"}

	cmd := []string{
		"cat /etc/*elease",
	}

	// am I doing this automation thing right?
	if err := commander.Command(cmd...); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running SSH command: ", err)
		os.Exit(1)
	}
}
*/

func getNodeOperatingSystems(ips []string) []models.Node {

	nodes := make([]models.Node, 0)

	for _, ip := range ips {
		cmd := exec.Command("ssh", "root@"+ip, "cat /etc/*elease")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		outstr := out.String()
		parts := strings.Split(outstr, "NAME=")
		nameParts := strings.Split(parts[1], "\"")

		versionArr := strings.Split(outstr, "VERSION_ID=")
		versionParts := strings.Split(versionArr[1], "\"")

		osName := strings.Split(nameParts[1], " ")[0]

		osFullVersion := versionParts[1]

		// if centos, then we need additional parsing
		if strings.Contains(strings.ToLower(outstr), strings.ToLower("centos")) {
			osFullVersion = strings.Split(strings.Split(parts[0], " release ")[1], " ")[0]
			majorMinorPatch := strings.Split(osFullVersion, ".")
			if len(majorMinorPatch) == 3 {
				osFullVersion = majorMinorPatch[0] + "." + majorMinorPatch[1]
			}
		}
		//fmt.Println(outstr)
		//fmt.Println(osName, osFullVersion)
		node := models.Node{
			Host: ip,
			Name: "node." + ip,
			OperatingSystem: models.OperatingSystem{
				Name:    osName,
				Version: osFullVersion,
			},
		}
		nodes = append(nodes, node)

	}

	return nodes
}

func GetNodesFromCSV(csvFilename string) []models.Node {
	csvFile, _ := os.Open(csvFilename)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("error reading all lines: %v", err)
	}

	nodes := make([]models.Node, 0)

	for i, line := range lines {
		if i == 0 {
			// skip header line
			continue
		}

		esxiID, _ := strconv.Atoi(line[5])
		snapshotID, _ := strconv.Atoi(line[6])
		ram, _ := strconv.Atoi(line[13])
		node := models.Node{
			ID:           line[0],
			Host:         line[1],
			Name:         line[2],
			EsxiIP:       line[3],
			EsxiServerID: line[4],
			Esxi: models.Esxi{
				ID:     esxiID,
				States: []models.State{{Name: line[7], SnapshotID: snapshotID}},
			},
			OperatingSystem: models.OperatingSystem{
				Name:    line[8],
				Version: line[9],
			},
			RAM:      ram,
			Username: os.Getenv("DEFAULT_USERNAME"),
			Password: os.Getenv("DEFAULT_PASSWORD"),
		}

		nodes = append(nodes, node)
	}

	return nodes
}

func GetNodeJsonDocMap(node models.Node) map[string]interface{} {
	nodeDbJson := models.NodeDBJson{
		NodeObj: node,
		ID:      node.ID,
	}

	return structs.Map(nodeDbJson)

}

func getNodeJsonDocString(node models.Node) string {
	nodeDbJson := models.NodeDBJson{
		NodeObj: node,
		ID:      node.ID,
	}

	nodeJson, err := json.Marshal(nodeDbJson)
	if err != nil {
		log.Fatalf("couldn't marshal obj to json: %v", err)
	}

	return string(nodeJson)
}

func getNodeOS() {

}

func nodeCSVToJSonFile() {

}
