package main

import (
	"go_grpc_realtime/lib/core/server"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Ltime)

	server.RunServer()
}
