package file_io

import (
	"encoding/json"
	"json-storage-service/internal/storage"
	"log"
	"os"
	"time"
)

const filePath = "internal/storage/saved_datas.json" // Путь к файлу хранения данных

type FileManager struct {
	storage *storage.Storage
}

func NewFileManager(s *storage.Storage) *FileManager {
	return &FileManager{storage: s}
}

// LoadFromFile загружает все данные из файла в хранилище
func (m *FileManager) LoadFromFile() error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var data map[string]storage.Object
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return err
	}

	for key, obj := range data {
		ttl := time.Until(obj.ExpiryTime)
		if ttl > 0 {
			m.storage.SaveObject(key, obj.Data, ttl)
		}
	}
	return nil
}

// SaveToFile сохраняет все данные из хранилища в файл
func (m *FileManager) SaveToFile() error {
	data := m.storage.GetAllData()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(data)
}

// Тикер на периодическое сохранение данных в файл
func (m *FileManager) TikerForSave() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		if err := m.SaveToFile(); err != nil {
			log.Println("Failed to save data to file:", err)
		}
	}
}
