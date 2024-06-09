package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// HandleLiveness обработчик для проверки liveness
func (s *Server) HandleLiveness(w http.ResponseWriter, r *http.Request) {
	// Создание конфигурации для доступа к Kubernetes API
	config, err := rest.InClusterConfig()
	if err != nil {
		// Обработка ошибки, если не удается получить конфигурацию из кластера
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Создание клиента для доступа к Kubernetes API
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		// Обработка ошибки, если не удается создать клиент для доступа к Kubernetes API
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверка наличия ресурса Pod с именем вашего сервиса
	podName := os.Getenv("POD_NAME")

	_, err = clientset.CoreV1().Pods("default").Get(context.TODO(), podName, v1.GetOptions{})
	if err != nil {
		// Обработка ошибки, если не удается найти под в Kubernetes
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Если нет ошибки, отправляем успешный статус
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// HandleReadiness обработчик для проверки readiness
func (s *Server) HandleReadiness(w http.ResponseWriter, r *http.Request) {
	// Инициализация клиента Kubernetes
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating Kubernetes configuration: %v", err), http.StatusInternalServerError)
			return
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating Kubernetes client: %v", err), http.StatusInternalServerError)
		return
	}

	// Получение имени пода и пространства имен
	podName := os.Getenv("POD_NAME")
	if podName == "" {
		http.Error(w, "POD_NAME environment variables not set", http.StatusInternalServerError)
		return
	}

	// Получение информации о поде
	pod, err := clientset.CoreV1().Pods("default").Get(context.TODO(), podName, v1.GetOptions{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting pod information: %v", err), http.StatusInternalServerError)
		return
	}
	// Проверка состояния пода
	if pod.Status.Phase != "Running" {
		http.Error(w, "Pod not ready", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
