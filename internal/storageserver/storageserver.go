package storageserver

import (
	"context"
	"log"

	"github.com/cocm1324/gcache/internal/storage"
	storagepb "github.com/cocm1324/gcache/protos/storage"
)

type StorageServer struct {
	storagepb.StorageServer
	storage *storage.Storage
}

func Init(s *storage.Storage) *StorageServer {
	return &StorageServer{
		storage: s,
	}
}

func (s *StorageServer) Get(ctx context.Context, req *storagepb.GetRequest) (*storagepb.GetResponse, error) {
	log.Println("storageserver: get invoked")
	value, hit := s.storage.Get(req.Key)
	return &storagepb.GetResponse{
		Hit:   hit,
		Value: value,
	}, nil
}

func (s *StorageServer) Put(ctx context.Context, req *storagepb.PutRequest) (*storagepb.PutResponse, error) {
	log.Println("storageserver: put invoked")
	hit := s.storage.Put(req.Key, req.Value)
	return &storagepb.PutResponse{
		Hit: hit,
	}, nil
}

func (s *StorageServer) Delete(ctx context.Context, req *storagepb.DeleteRequest) (*storagepb.DeleteResponse, error) {
	log.Println("storageserver: delete invoked")
	hit := s.storage.Delete(req.Key)
	return &storagepb.DeleteResponse{
		Hit: hit,
	}, nil
}

func (s *StorageServer) Clear(ctx context.Context, req *storagepb.ClearRequest) (*storagepb.ClearResponse, error) {
	log.Println("storageserver: clear invoked")
	s.storage.Clear()
	return &storagepb.ClearResponse{}, nil
}
