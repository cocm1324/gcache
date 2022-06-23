package main

import (
	"log"
	"net"
	"time"

	"github.com/cocm1324/gcache/internal/storage"
	"github.com/cocm1324/gcache/internal/storageserver"
	storagepb "github.com/cocm1324/gcache/protos/storage"
	"google.golang.org/grpc"
)

const PORT string = ":9000"

func main() {
	listen, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("main: failed to open listing on port %s: %v", PORT, err)
	}
	storageConfig := storage.StorageConfig{
		Ttl:      time.Hour * 24,
		Capacity: 1000,
	}
	storage := storage.New(storageConfig)
	storageServer := storageserver.Init(storage)
	grpcServer := grpc.NewServer()
	storagepb.RegisterStorageServer(grpcServer, storageServer)

	log.Printf("main: start grpc server on %s\n", PORT)
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("main: failed to start grpc server: %v", err)
	}
}
