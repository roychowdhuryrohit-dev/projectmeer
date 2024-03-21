package config

import (
	"log"
	"net"
	"os"
	"sync"
)

const (
	ReplicaID = "REPLICA_ID"
	NodeListValue  = "NODE_LIST"
	Port ="PORT"
)

var ConfigMap sync.Map
var NodeList []string

func Config() {

	if envar, ok := os.LookupEnv(ReplicaID); ok && envar != "" {
		switch envar {
		case "LOCAL_IP":
			conn, err := net.Dial("udp", "8.8.8.8:80")
			if err != nil {
				log.Panicln(err.Error())
			}
			defer conn.Close()
			ConfigMap.Store(ReplicaID, conn.LocalAddr().String())
		case "PUBLIC_IP":
		default:
			ConfigMap.Store(ReplicaID, envar)
		}
	} else {
		log.Panicf("%s not set", ReplicaID)
	}

	if envar, ok := os.LookupEnv(NodeListValue); ok && envar != "" {
		ConfigMap.Store(NodeListValue, envar)
	} else {
		log.Panicf("%s not set", NodeListValue)
	}

	if envar, ok := os.LookupEnv(Port); ok && envar != "" {
		ConfigMap.Store(Port, envar)
	} else {
		log.Panicf("%s not set", Port)
	}

}
