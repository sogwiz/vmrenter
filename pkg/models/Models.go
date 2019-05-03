package models

type Node struct {
	ID                   string
	Name                 string          `json:"name"`
	Host                 string          `json:"host"`
	OperatingSystem      OperatingSystem `json:"operatingSystem"`
	Username             string          `json:"username"`
	Password             string          `json:"password"`
	EsxiIP               string
	EsxiServerID         string
	Esxi                 Esxi     `json:"esxi"`
	ExpectedServiceNames []string `json:"expectedServiceNames"`
}

type NodeDBJson struct {
	ID        string `json:"_id"`
	NodeObj   Node
	ClusterID string
	ExpiresAT string
}

type Esxi struct {
	ID     int `json:"id"`
	States []State
}

type State struct {
	Name       string `json:"name"`
	SnapshotID int    `json:"snapshotId"`
}

type OperatingSystem struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Cluster struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Nodes []Node `json:"nodes"`
}

type FileContent struct {
	Clusters       []Cluster                `json:"clusters"`
	Rest           map[string]interface{}   `json:"rest"`
	ClusterTesting map[string]interface{}   `json:"clusterTesting"`
	Cucumber       map[string]interface{}   `json:"cucumber"`
	Esxi           map[string]interface{}   `json:"esxi"`
	Docker_Nodes   []map[string]interface{} `json:"docker_nodes"`
}

type Reservation struct {
	ID              string `json:"_id"`
	CreatedAt       string
	ExpiresAt       string
	JenkinsJobURL   string
	Nodes           []NodeDBJson
	ClusterID       string
	ReservationType string
}
