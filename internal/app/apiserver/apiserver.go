package apiserver

import (
	"github.com/gorilla/mux"
	"github.com/marinira/viber-e-shop/internal/app/db"
	"github.com/marinira/viber-e-shop/internal/app/viberapi"
	"github.com/sirupsen/logrus"
	"html/template"
	"io"
	"log"
	"net/http"
)

// Структура API сервера
type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	mongo  *db.DataBaseInterface
	viber  *viberapi.Viber
}

//Инициализация API сервера
func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
		viber:  viberapi.NewViber(),
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
	//s.viberapi.InitViberChat()

	//return nil
	err := http.ListenAndServe(s.config.BindAddr, s.router)

	return err

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

	// Handle Viber routes
	//s.viberChat.Uri = "https://083ccc95.ngrok.io"

	s.router.Handle("/viber/webhook/", s.viber)
	//s.router.HandleFunc("/",s.handleViber())
	api := s.router.PathPrefix("/api/").Subrouter()
	api.HandleFunc("/hello/", s.handleHello())
	api.HandleFunc("/startbot/", s.handleStartBot())

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

//функция обработки запроса "*/hello"
func (s *APIServer) handleViber() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		s.logger.Error("Viber")
		//s.viberapi.ServeHTTP(writer,request)
	}
}

//функция обработки запроса "*/hello"
func (s *APIServer) handleStartBot() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		ak, e := s.viber.SetWebhook("", nil)
		ak, e = s.viber.SetWebhook(s.viber.Uri, nil)
		log.Println(ak.Status, ak.StatusMessage)

		if e != nil {
			io.WriteString(writer, e.Error())
		} else {

			a, err := s.viber.AccountInfo()
			if err != nil {
				log.Println("AccountInfo Viber error:", err)
				return
			}
			adminID := ""
			userID := ""
			// print all admministrators
			//mes := s.viberapi.NewTextMessage("Hello to everyone")

			for _, m := range a.Members {

				if m.Role == "admin" {
					adminID = m.ID
				} else {
					userID = m.ID
				}
				ms := s.viber.NewTextMessage("Hello, world!")
				ms.Sender = viberapi.Sender{
					Name:   "SomeOtherName",
					Avatar: "https://mysite.com/img/other_avatar.jpg",
				}
				ea1, ee1 := s.viber.SendMessage(m.ID, ms)
				log.Print(m.ID, m.Name, m.Role, a.ID, adminID, ea1, ee1)
			}

			mes1 := s.viber.NewTextMessage("Hello to everyone1")
			mes1.Sender = viberapi.Sender{
				Name:   "SomeOtherName1",
				Avatar: "https://mysite.com/img/other_avatar.jpg",
			}
			//ea, ee := s.viberapi.SendPublicMessage(a.ID, mes)
			//log.Print(ea , ee)
			ea, ee := s.viber.SendPublicMessage(adminID, mes1)
			log.Print(ea, ee)

			mes2 := s.viber.NewTextMessage("Hello to everyone1")
			mes2.Sender = viberapi.Sender{
				Name:   "SomeOtherName2",
				Avatar: "https://mysite.com/img/other_avatar.jpg",
			}
			//ea, ee := s.viberapi.SendPublicMessage(a.ID, mes)
			//log.Print(ea , ee)
			ea2, ee2 := s.viber.SendPublicMessage(a.ID, mes2)
			log.Print(ea2, ee2)

			mes3 := s.viber.NewTextMessage("Hello to everyone1")
			mes3.Sender = s.viber.Sender
			//ea, ee := s.viberapi.SendPublicMessage(a.ID, mes)
			//log.Print(ea , ee)
			ea3, ee3 := s.viber.SendPublicMessage(a.ID, mes3)
			log.Print(ea3, ee3)

			m := s.viber.NewTextMessage("Hello, world!")
			m.Sender = viberapi.Sender{
				Name:   "SomeOtherName",
				Avatar: "https://mysite.com/img/other_avatar.jpg",
			}
			ea1, ee1 := s.viber.SendMessage(userID, m)
			log.Print(ea1, ee1)
			io.WriteString(writer, "200")
		}

	}
}
