package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	//"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	//"net"
	"net/http"
	//"net/url"
	"reflect"
	"time"
)

type HTTPClient struct {
	client      *http.Client
	logger      *logrus.Logger
	header      *http.Header
	handler_str string
	url         string
	port        string
	proxy       string
}

//Конструктор - создает екземпляр класса HTTPClient и
//инициализирует поля по умолчанию
func NewHTTPClient() *HTTPClient {
	//создаем новый екземпляр класса HTTPClient
	ret := new(HTTPClient)
	//создаем екземпляр logrus (логирование) для екземпляра  класса HTTPClient
	ret.logger = logrus.New()
	ret.client = &http.Client{}
	ret.header = &http.Header{}
	// инициализаруем занчения полей по умолчанию
	ret.init()
	ret.setHeader()
	return ret
}

//функция инициализации значения полей по умолчанию.
//данные берутся из констант из mongo_config
func (m *HTTPClient) init() {
	m.handler_str = ""
	m.url = ""
	m.port = ""
	m.proxy = ""
	//инициализация клиента

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	m.client.Transport = tr

}

// setHeader функция устанавливает узерагент для Header в запросах
func (m *HTTPClient) setHeader() {
	// добавляем заголовки
	m.header.Add("Accept", "text/html")     // добавляем заголовок Accept
	m.header.Add("User-Agent", "MSIE/15.0") // добавляем заголовок User-Agent
	return
}
func (m *HTTPClient) SetUrl(u string) {
	m.url = u
}
func (m *HTTPClient) Get(url string, parameter interface{}) (response interface{}, code int, err error) {
	//Создаем параметры запроса и метод запроса
	req, err := http.NewRequest(
		"GET", m.url+url, nil,
	)
	if err != nil {
		m.logger.Error(err)
		return
	}
	// задаем хидеры запроса
	req.Header = *m.header
	//отправляем сам запрос
	resp, err := m.client.Do(req)
	if err != nil {
		m.logger.Error(err, resp.StatusCode, resp.Status)
		return
	}
	//коныертируем ответ в формат utf8
	utf8, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		m.logger.Error("Encoding error:", err)
		return
	}
	//считываем ответ
	body, err := ioutil.ReadAll(utf8)
	if err != nil {
		m.logger.Error(err)
	}
	retcode := resp.StatusCode
	defer resp.Body.Close()
	return string(body), retcode, nil
}

func (m *HTTPClient) SetCookie(cookie string) {
	m.header.Set("Cookie", cookie)
}

func (m *HTTPClient) Post(url string, parameter interface{}) (response interface{}, code int, err error) {

	//конвертируем тело запроса в тип []byte
	b, _ := json.Marshal(parameter)
	reque := bytes.NewBuffer(b)
	//Создаем параметры запроса и метод запроса
	req, err := http.NewRequest("POST", m.url+url, reque)
	if err != nil {
		m.logger.Error(err)
		return
	}
	// задаем хидеры запроса
	req.Header = *m.header
	//отправляем сам запрос
	resp, err := m.client.Do(req)
	if err != nil {
		m.logger.Error(err, resp.StatusCode, resp.Status)
		return
	}
	//коныертируем ответ в формат json
	reqJson, _ := simplejson.NewFromReader(resp.Body)
	rr, ee := reqJson.Encode()
	if ee != nil {
		m.logger.Error("Encoding error:", err)
		return
	}
	//считываем ответ
	if err != nil {
		m.logger.Error(err)
	}

	defer resp.Body.Close()
	return string(rr), resp.StatusCode, nil
}

func (m *HTTPClient) PostForm(myurl string, parameter interface{}) (response interface{}, code int, err error) {
	//Создаем параметры запроса и метод запроса

	//переводим из стурктуры в мап
	m1 := StructToMap(parameter)
	//записываем данные в реквест
	form := make(url.Values)
	for i, v := range m1 {
		form.Add(i, fmt.Sprint(v))
	}
	//формируем запрос
	req, err := http.NewRequest("POST", m.url+myurl, strings.NewReader(form.Encode()))
	if err != nil {
		m.logger.Error(err)
	}
	// задаем хидеры запроса
	req.Header = *m.header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//отправляем сам запрос
	resp, err := m.client.Do(req)
	//resp, err := http.PostForm(m.url+myurl, form)

	if err != nil {
		m.logger.Error(err, resp.StatusCode, resp.Status)
		return
	}
	//коныертируем ответ в формат json
	reqJson, _ := simplejson.NewFromReader(resp.Body)
	rr, ee := reqJson.Encode()
	if ee != nil {
		m.logger.Error("Encoding error:", err)
		return
	}

	defer resp.Body.Close()
	return string(rr), resp.StatusCode, nil
}

//recursive struct to map
func StructToMap(item interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	nn := v.NumField()
	for i := 0; i < nn; i++ {
		tag := v.Field(i).Tag.Get("json")

		if tag != "" && tag != "-" {
			tag = v.Field(i).Tag.Get("json")
		} else {
			tag = v.Field(i).Name
		}
		if v.Field(i).Type.Kind() == reflect.Struct {
			field := reflectValue.Field(i).Interface()
			res = StructToMap(field)
		} else {
			field := reflectValue.Field(i)
			res[tag] = field.Interface()
		}
	}
	return res
}
