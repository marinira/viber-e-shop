package apiserver

import (
	"github.com/gorilla/mux"
	"github.com/marinira/http-rest-api/internal/app/db"
	"github.com/sirupsen/logrus"
	"html/template"
	"io"
	"net/http"
)

// Структура API сервера
type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	mongo  *db.DataBaseInterface
}

//Инициализация API сервера
func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

// Старт API сервера
func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	if err := s.configureDB(); err != nil {
		return err
	}
	s.configureRouter()
	s.logger.Info("starting API server")
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

// конфигурация логера
func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *APIServer) configureDB() error {
	return nil
}

///////////////////////////////////////////////////////////////////////
// *************  Функции приема запроса к серверу
///////////////////////////////////////////////////////////////////////
//конфигурация роутера
func (s *APIServer) configureRouter() {
	// Handle API routes
	api := s.router.PathPrefix("/api/").Subrouter()
	api.HandleFunc("/hello", s.handleHello())

	// Serve static files
	s.router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./public/"))))
	// Serve index page on all unhandled routes
	s.router.PathPrefix("/").HandlerFunc(s.handleIndex())
}

//функция обработки запроса "*"
func (s *APIServer) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("/index.html")
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		type TodoPageData struct {
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}
	}
}

//функция обработки запроса "*/hello"
func (s *APIServer) handleHello() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		io.WriteString(writer, "Hello")
	}
}
