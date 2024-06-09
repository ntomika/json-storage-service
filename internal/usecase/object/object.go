package object

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"json-storage-service/internal/storage"
)

type ObjectUsecase struct {
	storage *storage.Storage
}

func NewObjectUsecase(s *storage.Storage) *ObjectUsecase {
	return &ObjectUsecase{storage: s}
}

func (u *ObjectUsecase) ProcessObjects(r *http.Request) (map[string]interface{}, error, int) {
	var obj map[string]interface{}
	var found bool

	key := strings.TrimPrefix(r.URL.Path, "/objects/")

	switch r.Method {
	case http.MethodPut:
		if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
			return nil, errors.New("Invalid JSON"), http.StatusBadRequest
		}

		expires := r.Header.Get("Expires")
		ttl := time.Duration(1<<63 - 1)
		if expires != "" {
			expiryTime, err := time.Parse(time.RFC1123, expires)
			if err != nil {
				return nil, errors.New("Invalid Expires header"), http.StatusBadRequest

			}
			ttl = time.Until(expiryTime)
		}

		u.storage.SaveObject(key, obj, ttl)

		log.Println("object added to storage")

	case http.MethodGet:
		obj, found = u.storage.GetObject(key)
		if !found {
			return nil, errors.New("Key not found in storage"), http.StatusNotFound
		}

	default:
		return nil, errors.New("Method not allowed"), http.StatusMethodNotAllowed
	}

	return obj, nil, http.StatusOK
}
