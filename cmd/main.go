package main

import (
	"json-storage-service/internal/file_io"
	"json-storage-service/internal/handlers"
	"json-storage-service/internal/storage"

	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	storage := storage.NewStorage()
	manager := file_io.NewFileManager(storage)

	if err := manager.LoadFromFile(); err != nil {
		log.Println("Failed to load data from file:", err)
	}

	server := handlers.NewServer(storage)

	http.HandleFunc("/objects/", server.HandleObjects)
	http.HandleFunc("/metrics", server.HandleMetrics)
	http.HandleFunc("/probes/liveness", server.HandleLiveness)
	http.HandleFunc("/probes/readiness", server.HandleReadiness)

	go manager.TikerForSave()

	// Завершение работы с сохранением данных
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stop
		log.Println("Shutting down server...")

		// Сохранение данных в файл перед завершением работы
		if err := manager.SaveToFile(); err != nil {
			log.Println("Failed to save data to file:", err)
		}

		os.Exit(0)
	}()

	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
