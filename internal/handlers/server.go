package handlers

import (
	"json-storage-service/internal/storage"
	"json-storage-service/internal/usecase/object"
)

type Server struct {
	storage *storage.Storage
	object  *object.ObjectUsecase
}

func NewServer(storage *storage.Storage) *Server {
	object := object.NewObjectUsecase(storage)

	return &Server{
		storage: storage,
		object:  object,
	}
}
