package config

import (
	"os"
	"sync"
)

var urldbConn string
var m = &sync.Mutex{}

func SetURLDBConn(urlconnection string) {
	urldbConn = urlconnection
}

func GetURLDBConn() string {
	m.Lock()
	if urldbConn == "" {
		urldbConn = os.Getenv("URL_DB_CONN")
	}
	m.Unlock()
	return urldbConn
}
