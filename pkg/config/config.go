package config

import "os"

var URLDBConn string

func GetURLDBConn() string {
	if URLDBConn == "" {
		URLDBConn = os.Getenv("URL_DB_CONN")
	}
	return URLDBConn
}
