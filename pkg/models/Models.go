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
	Available bool
}

type Esxi struct {
	ID     string
	States []State
}

type State struct {
	Name       string
	SnapshotID string
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
