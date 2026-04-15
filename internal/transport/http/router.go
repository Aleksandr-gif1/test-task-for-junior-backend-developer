package transporthttp

import (
	"net/http"

	"github.com/gorilla/mux"

	swaggerdocs "example.com/taskservice/internal/transport/http/docs"
	httphandlers "example.com/taskservice/internal/transport/http/handlers"
)

func NewRouter(taskHandler *httphandlers.TaskHandler, docsHandler *swaggerdocs.Handler) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/swagger/openapi.json", docsHandler.ServeSpec).Methods(http.MethodGet)
	router.HandleFunc("/swagger/", docsHandler.ServeUI).Methods(http.MethodGet)
	router.HandleFunc("/swagger", docsHandler.RedirectToUI).Methods(http.MethodGet)

	api := router.PathPrefix("/api/v1").Subrouter()

	// Задачи (шаблоны)
	api.HandleFunc("/tasks", taskHandler.Create).Methods(http.MethodPost)
	api.HandleFunc("/tasks", taskHandler.List).Methods(http.MethodGet)
	api.HandleFunc("/tasks/{id:[0-9]+}", taskHandler.GetByID).Methods(http.MethodGet)
	api.HandleFunc("/tasks/{id:[0-9]+}", taskHandler.Update).Methods(http.MethodPut)
	api.HandleFunc("/tasks/{id:[0-9]+}", taskHandler.Delete).Methods(http.MethodDelete)

	// Генерация экземпляров
	api.HandleFunc("/tasks/generate", taskHandler.GenerateTasks).Methods(http.MethodPost)

	// Экземпляры задач
	api.HandleFunc("/tasks/instances", taskHandler.ListInstances).Methods(http.MethodGet)
	api.HandleFunc("/tasks/instances/{id:[0-9]+}", taskHandler.GetInstanceByID).Methods(http.MethodGet)
	api.HandleFunc("/tasks/instances/{id:[0-9]+}/status", taskHandler.UpdateInstanceStatus).Methods(http.MethodPut)

	return router
}
