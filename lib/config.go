package config

import (
	"log"
	"net"
	"os"
	"sync"
)

const (
	ReplicaID = "REPLICA_ID"
	NodeList = "NODE_LIST"
)

var ConfigMap sync.Map

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
	}

	if envar, ok := os.LookupEnv(NodeList); ok && envar != "" {
		ConfigMap.Store(NodeList, envar)
	} else {
		log.Fatalf("%s not set", NodeList)
	}


	


}
