package config

import "os"

var urldbConn string

func SetURLDBConn(urlconnection string) {
	urldbConn = urlconnection
}

func GetURLDBConn() string {
	if urldbConn == "" {
		urldbConn = os.Getenv("URL_DB_CONN")
	}
	return urldbConn
}
