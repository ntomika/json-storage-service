package storage

import (
	"sync"
	"time"
)

// Основная структура хранилища
type Storage struct {
	data map[string]Object
	sync.RWMutex
}

// Структура для хранения данных объекта с таймером
type Object struct {
	Data       map[string]interface{}
	ExpiryTime time.Time
}

// Инициализация хранилища
func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]Object),
	}
}

// SaveObject сохраняет объект в хранилище с указанным временем жизни (ttl)
func (s *Storage) SaveObject(key string, data map[string]interface{}, ttl time.Duration) {
	s.Lock()
	defer s.Unlock()

	expiryTime := time.Now().Add(ttl)
	s.data[key] = Object{
		Data:       data,
		ExpiryTime: expiryTime,
	}

	// Таймер для удаления объекта по истечении времени жизни
	go func() {
		<-time.After(ttl)
		s.Lock()
		defer s.Unlock()
		delete(s.data, key)
	}()
}

// GetObject возвращает объект по ключу, если он существует и еще не истек
func (s *Storage) GetObject(key string) (map[string]interface{}, bool) {
	s.RLock()
	defer s.RUnlock()

	obj, exists := s.data[key]
	if !exists || time.Now().After(obj.ExpiryTime) {
		return nil, false
	}

	return obj.Data, true
}

// GetAllData возвращает все данные из хранилища, игнорируя объекты, срок действия которых истек
func (s *Storage) GetAllData() map[string]Object {
	s.RLock()
	defer s.RUnlock()

	data := make(map[string]Object)
	for key, obj := range s.data {
		if time.Now().Before(obj.ExpiryTime) {
			data[key] = obj
		}
	}
	return data
}
