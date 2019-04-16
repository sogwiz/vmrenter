package models

type Node struct {
	Name            string
	Host            string
	OperatingSystem OperatingSystem
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
