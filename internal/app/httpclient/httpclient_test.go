package httpclient

import (
	"bytes"
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

//проферка функции GET
func TestHTTPClient_Get(t *testing.T) {
	//запускаей мок сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//записываем в ответ код ответа 503
		w.WriteHeader(http.StatusServiceUnavailable)
		// записываем тело ответа "Hello !!!"
		w.Write([]byte("Hello !!!"))
	}))
	defer ts.Close()
	//определяем URL МОК сервера
	nsqdUrl := ts.URL
	//создаем httpClient
	httpClient := NewHTTPClient()
	//инициализируем httpClient
	httpClient.init()
	// устанавливаем Хидеры
	httpClient.setHeader()
	// устанавливаем Url по которому httpClient должен отправить запрос
	httpClient.url = nsqdUrl
	//вызываем тестируемыею функцию которая отправляет Get запрос
	resp, cod, err := httpClient.Get("", nil)
	if err != nil {
		httpClient.logger.Error(err)
		return
	}
	//httpClient.logger.Info("Code: ",cod, "  Response: ", resp)
	//проверяем плавельность полученого тела ответа и код ответа
	assert.Equal(t, "Hello !!!", resp, "Hello !!!")
	assert.Equal(t, 503, cod, "503")
}

//проверка функции POST
func TestHTTPClient_Post(t *testing.T) {
	type MyStruct struct {
		Name string
		Age  int
	}
	//запускаей мок сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//записываем в ответ код ответа 200
		w.WriteHeader(http.StatusOK)
		// копируем в тело ответа тело запроса
		reqJson, _ := simplejson.NewFromReader(r.Body)
		rr, ee := reqJson.Encode()
		if ee != nil {
			w.Write([]byte(ee.Error()))
		}
		w.Write(rr)
	}))
	defer ts.Close()
	//определяем URL МОК сервера
	nsqdUrl := ts.URL
	//создаем httpClient
	httpClient := NewHTTPClient()
	//инициализируем httpClient
	httpClient.init()
	// устанавливаем Хидеры
	httpClient.setHeader()
	// устанавливаем Url по которому httpClient должен отправить запрос
	httpClient.url = nsqdUrl
	//создаем структуру которую отправляем в тело запроса
	myStrac := MyStruct{
		Name: "Heloo",
		Age:  5,
	}
	//вызываем тестируемыею функцию которая отправляет POST запрос
	resp, cod, err := httpClient.Post("", myStrac)
	if err != nil {
		httpClient.logger.Error(err)
		return
	}
	//конвертируем ответ из interface -> string
	rr := resp.(string)
	//конвертируем ответ из string(JSON) -> MyStruct (в нашу тестовую структуру)
	myStruct := MyStruct{}
	err = json.Unmarshal([]byte(rr), &myStruct)
	if err != nil {
		return
	}
	// Декодируем ответ
	//httpClient.logger.Info("Code: ",cod, "  Response: ", rr)
	//проверяем плавельность полученого тела ответа и код ответа
	assert.Equal(t, myStruct, myStrac, "Hello !!!")
	assert.Equal(t, 200, cod, "503")
}
func TestHTTPClient_PostJson(t *testing.T) {

	session := "123123123"
	payload := `{
        action: 'get',
		controller: 'cache',
		params: {
			names: ['Player:']
		},` +
		session +
		`}`
	//payloadByte := []byte(payload)
	//запускаей мок сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}
		//записываем в ответ код ответа 200
		w.WriteHeader(http.StatusOK)
		// копируем в тело ответа тело запроса
		reqJson, _ := simplejson.NewFromReader(r.Body)
		rr, ee := reqJson.Encode()
		if ee != nil {
			w.Write([]byte(ee.Error()))
		}
		w.Write(rr)
	}))
	defer ts.Close()
	//определяем URL МОК сервера
	nsqdUrl := ts.URL
	//создаем httpClient
	httpClient := NewHTTPClient()
	//инициализируем httpClient
	httpClient.init()
	// устанавливаем Хидеры
	httpClient.setHeader()
	// устанавливаем Url по которому httpClient должен отправить запрос
	httpClient.url = nsqdUrl
	//создаем структуру которую отправляем в тело запроса

	//вызываем тестируемыею функцию которая отправляет POST запрос
	resp, cod, err := httpClient.Post("", payload)
	if err != nil {
		httpClient.logger.Error(err)
		return
	}
	var mm map[string]interface{}
	DumpMap(resp.(string), mm)
	logrus.Error(mm)
	//конвертируем ответ из interface -> string
	rr := resp.(string)
	//проверяем плавельность полученого тела ответа и код ответа
	assert.Equal(t, rr, payload, "Hello !!!")
	assert.Equal(t, 200, cod, "503")
}

//проверка функции POSTFORM
func TestHTTPClient_PostForm(t *testing.T) {
	type MyStruct struct {
		Name string
		Age  int
	}
	//запускаей мок сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ОБЯЗАТЕЛЬНО В САМОМ НАЧАЛЕ ДЕЛАЕМ ПАРСЕР ФОРМ
		//ЕСЛИ СЧИТЫВАЕМ КАКИЕ ЛИБО ДАННЫЕ С ЗАПРОСОВ ПАРСЕР НЕ РАБОТАЕТ
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}
		//записываем в ответ код ответа 200

		w.WriteHeader(http.StatusOK)
		// получаем данные с формы
		v := r.Form
		msg := new(MyStruct)
		decoder := schema.NewDecoder()
		decoder.Decode(msg, v)

		// копируем в тело ответа тело запроса

		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(msg)

		rr := reqBodyBytes.Bytes() // this is the []byte

		w.Write(rr)
	}))
	defer ts.Close()
	//определяем URL МОК сервера
	nsqdUrl := ts.URL
	//создаем httpClient
	httpClient := NewHTTPClient()
	//инициализируем httpClient
	httpClient.init()
	// устанавливаем Хидеры
	httpClient.setHeader()
	// устанавливаем Url по которому httpClient должен отправить запрос
	httpClient.url = nsqdUrl
	//создаем структуру которую отправляем в тело запроса
	myStrac := MyStruct{
		Name: "Heloo",
		Age:  5,
	}
	//вызываем тестируемыею функцию которая отправляет POST запрос
	resp, cod, err := httpClient.PostForm("", myStrac)
	if err != nil {
		httpClient.logger.Error(err)
		return
	}
	//конвертируем ответ из interface -> string
	rr := resp.(string)
	//конвертируем ответ из string(JSON) -> MyStruct (в нашу тестовую структуру)
	myStruct := MyStruct{}
	err = json.Unmarshal([]byte(rr), &myStruct)
	if err != nil {
		return
	}
	// Декодируем ответ
	//httpClient.logger.Info("Code: ",cod, "  Response: ", rr)
	//проверяем плавельность полученого тела ответа и код ответа
	assert.Equal(t, myStruct, myStrac, "Hello !!!")
	assert.Equal(t, 200, cod, "503")
}
