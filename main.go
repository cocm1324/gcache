package main

import (
	"log"
	"net"
	"time"

	"github.com/cocm1324/cstorage"
	"github.com/cocm1324/gcache/internal/storageserver"
	storagepb "github.com/cocm1324/gcache/protos/storage"
	"google.golang.org/grpc"
)

const port string = ":9000"

func main() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("main: failed to open listing on port %s: %v", port, err)
	}

	cstorageConfig := cstorage.CStorageConfig{
		Ttl:      time.Hour * 24 * 7,
		Capacity: 1000,
	}

	storage := cstorage.New(cstorageConfig)
	storageServer := storageserver.Init(storage)
	grpcServer := grpc.NewServer()
	storagepb.RegisterStorageServer(grpcServer, storageServer)

	log.Printf("main: start grpc server on %s\n", port)
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("main: failed to start grpc server: %v", err)
	}
}
