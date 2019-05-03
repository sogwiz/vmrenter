package models

type Node struct {
	ID              string
	Name            string
	Host            string
	OperatingSystem OperatingSystem
	Username        string
	Password        string
	EsxiIP          string
	EsxiServerID    string
	Esxi            Esxi
}

type NodeDBJson struct {
	ID        string `json:"_id"`
	NodeJson  Node
	ClusterID string
	ExpiresAT string
}

type Esxi struct {
	ID     int
	States []State
}

type State struct {
	Name       string
	SnapshotID int
}

type OperatingSystem struct {
	Name    string
	Version string
}

type Cluster struct {
	ID    string
	Name  string
	Nodes []Node
}

type FileContent struct {
	Clusters []Cluster
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
