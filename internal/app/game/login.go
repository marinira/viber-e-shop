package game

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/marinira/http-rest-api/internal/app/db"
	"github.com/marinira/http-rest-api/internal/app/httpclient"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type Login struct {
	gameWorld        string
	speedTroops      int
	email            string
	sitterType       string
	sitterName       string
	proxyIp          string
	proxyPort        string
	namePlayer       string
	msid             string
	tokenLobby       string
	sessionLobby     string
	cookiesLobby     string
	tokenGameWorld   string
	sessionGameWorld string
	cookiesGameWorld string
	logger           *logrus.Logger
	httpClient       *httpclient.HTTPClient

	//	lobbyEndpoint string//{"https://lobby.kingdoms.com/api/index.php"}
}

//функция задает значение полей структуры
func (l *Login) SetDate(newL Login) {
	l.gameWorld = newL.gameWorld
	l.speedTroops = newL.speedTroops
	l.email = newL.email
	l.sitterType = newL.sitterType
	l.sitterName = newL.sitterName
	l.proxyIp = newL.proxyIp
	l.proxyPort = newL.proxyPort
	l.namePlayer = newL.namePlayer
	l.msid = newL.msid
	l.tokenLobby = newL.tokenLobby
	l.sessionLobby = newL.sessionLobby
	l.cookiesLobby = newL.cookiesLobby
	l.tokenGameWorld = newL.tokenGameWorld
	l.sessionGameWorld = newL.sessionGameWorld
	l.cookiesGameWorld = newL.cookiesGameWorld
}

// функция считывает данные пользователя с базы данных
func (l *Login) ReadCredentials(email string, password string, gameWorld string, sitterType string, sitterName string) (bool, error) {
	//объвляем переменную интерфейсв базы данных
	var dataBase db.DataBaseInterface
	//создаем екземплар MpngoDB
	dataBase = db.NewMongo()
	//объявляем деффер при панеки чтобы закрыла открытое соединение с базою данных
	defer dataBase.Disconnect()
	//соединяемся с базою данных
	err := dataBase.InitConnect()
	if err != nil {
		return false, err
	}
	//создаем тип фильтр для поиска в базе по емейлу
	type MFilter struct {
		email      string
		password   string
		gameWorld  string
		sitterType string
		sitterName string
	}
	//задаем параметры фильтра для поиска в базе данных
	var alArr interface{}
	alArr = MFilter{
		email,
		password,
		gameWorld,
		sitterType,
		sitterName,
	}
	//забираем ОДНУ запись из базе данных согласно фильтру
	var retI interface{}
	retI, err = dataBase.FindOne("accounts", alArr)
	if err != nil {
		return false, err
	}
	//конвертируем ответ в наш тип
	//проверяем если нашел хоть одну запись
	if retI != nil {
		ret := retI.(Login)
		//вставляем полученные значения полей в наш екземпляр
		l.SetDate(ret)
		return true, nil
	}
	return false, nil

}

//функция записывает данные пользователя в базу данных
func (l *Login) WriteCredentials() error {
	//объвляем переменную интерфейсв базы данных
	var dataBase db.DataBaseInterface
	//создаем екземплар MpngoDB
	dataBase = db.NewMongo()
	//объявляем деффер при панеки чтобы закрыла открытое соединение с базою данных
	defer dataBase.Disconnect()
	//соединяемся с базою данных
	err := dataBase.InitConnect()
	if err != nil {
		return err
	}

	//todo if need
	//удаляем все записи из базы которые совпадают с нами
	//err = dataBase.DeleteMany("accounts",l)
	//if err !=nil {
	//	return err
	//}

	//записываем данные в базу данных
	err = dataBase.InsertOne("accounts", l)
	if err != nil {
		return err
	}
	return nil
}

//Функция менеджер подключения к игре к игре
func (l *Login) InitLogin(email string, password string, gameWorld string, sitterType string, sitterName string) error {
	l.logger = logrus.New()
	//создаем httpClient
	l.httpClient = httpclient.NewHTTPClient()
	//прверяем если есть такой клиент в базе и загружаем данные в нашу структуру
	isFind, err := l.ReadCredentials(email, password, gameWorld, sitterType, sitterName)
	if err != nil {
		return err
	}
	//проверяем если нашли клиента то идем дальше

	if isFind {
		l.logger.Info("found lobby session in database...", "Login")
		//записавем кукки для лобби из базы для этого клиента в хидет HTTPClient
		l.httpClient.SetCookie(l.cookiesLobby)
		// проверяем если можем зайти в лобби по сесси из кукки
		//иначе выдаем предупреждение не можем зайти
		if l.testLobbyConnection(l.sessionLobby) {
			l.logger.Info("database lobby connection successful", "Login")
			//записавем кукки для игрового мира из базы для этого клиента в хидет HTTPClient
			l.httpClient.SetCookie(l.cookiesGameWorld)
			if l.testGameWorldConnection(l.gameWorld, l.sessionGameWorld) {
				l.logger.Info("database gameworld connection successful", "Login")
				return nil
			} else {
				l.logger.Warn("database connection to gameworld failed", "Login")
			}
			l.loginToGameWorld(l.gameWorld, l.sitterType, l.sitterName, l.msid, l.sessionLobby)
			return nil
		} else {
			l.logger.Warn("database connection to lobby failed", "Login")
		}

	}
	//todo
	return nil
}

//функция логинется с сервером и захрдит в лобби
func (l *Login) loginToLobby(email string, password string) interface{} {
	//todo
	return nil
}

//функция подключается к игровому миру
func (l *Login) loginToGameWorld(gameWorld string, sitterType string, sitterName string, msid string, sessionLobby string) (interface{}, error) {
	speedTroops := "0"
	var mellonURL string
	if sitterType != "" && sitterName != "" {
		// login to sitter or dual account

		avatarID, gameworldID, err := l.getAvatarId(sessionLobby, gameWorld, sitterType, sitterName)
		if err != nil {
			l.logger.Error(err)
			return nil, err
		}
		speedTroops, err = l.getSpeedTroops(sessionLobby, gameworldID)
		if err != nil {
			l.logger.Error(err)
			return nil, err
		}
		mellonURL = "https://mellon-t5.traviangames.com/game-world/join-as-guest/avatarId/" + avatarID + "?msname=msid&msid=" + msid + ""
	} else {
		// login to normal gameworld
		gameworldID, err := l.getGameWorldId(sessionLobby, gameWorld)
		if err != nil {
			l.logger.Error(err)
			return nil, err
		}
		speedTroops, err = l.getSpeedTroops(sessionLobby, gameworldID)
		if err != nil {
			l.logger.Error(err)
			return nil, err
		}
		mellonURL = "https://mellon-t5.traviangames.com/game-world/join/gameWorldId/" + gameworldID + "?msname=msid&msid=" + msid + ""
	}

	ret, code, err := l.httpClient.Get(mellonURL, nil)
	if err != nil || code != 200 {
		l.logger.Error(code, err)
		return nil, err
	}

	//todo
	return nil, nil
}

//функця парсит HTML страничку и ищет токин
func (l *Login) parseToken(rawHtml string) interface{} {
	//todo
	return nil
}

//функция выберает Cookie
func (l *Login) parseCookies(cookieArray []interface{}) string {
	//todo
	return ""
}

//Функция проверяет если код ответа запроса хороший
func (l *Login) validateStatus(status int) bool {
	return status >= 200 && status < 303
}

//функция отпределяет номер игрового мира по его названию
func (l *Login) getGameWorldId(session string, gameWorldString string) (string, error) {

	payload := `{
		action: 'get',
		controller: 'cache',
			params: {
				names: ['Collection:Avatar:']
			},
		` + session + `
	}`
	//устанавливаем URL адрес
	l.httpClient.SetUrl("https://lobby.kingdoms.com/api/index.php")
	//отправляем сам запрос
	resp, _, err := l.httpClient.Post("", payload)
	if err != nil {
		l.logger.Error(err)
		return "", err
	}

	gameWorlds, err := Clash(resp, "data", "cache", "response")
	if err != nil {
		l.logger.Error(err)
		return "", err
	}
	//конвертируем полученный список ишровых миров в слайс
	allgameWorlds, ok := gameWorlds.([]interface{})
	if !ok {
		err := errors.New("Clash: Error convert interface to map ")
		l.logger.Error(err)
		return "", err
	}
	gameWorlds_id, err := Clash(allgameWorlds[0], "data")
	if err != nil {
		l.logger.Error(err)
		return "", err
	}
	gameWorlds_idArr, ok := gameWorlds_id.([]interface{})
	if !ok {
		err := errors.New("Clash: Error convert interface to map ")
		l.logger.Error(err)
		return "", err
	}
	for _, gameWorlds_id := range gameWorlds_idArr {
		gameWorldsName, err := Clash(gameWorlds_id, "data", "worldName")
		if err != nil {
			l.logger.Error(err)
			return "", err
		}
		if strings.ToLower(gameWorldsName.(string)) == gameWorldString {
			gameWorld, err := Clash(gameWorlds_id, "data", "consumersId")
			if err != nil {
				l.logger.Error(err)
				return "", err
			}
			return gameWorld.(string), nil
		}
	}

	err = errors.New("D'ont find world ID!!! ")
	return "", err
}

//функция возвращает скорость передвижения войск игрового мира по его номеру
func (l *Login) getSpeedTroops(session string, gameWorldId string) (string, error) {

	//создаем тело запроса в формате JSON строки
	payload := `{
			action: 'get',
			controller: 'cache',
			params: {
			names: ['GameWorld:'` + gameWorldId + `']
			},
			session:'` + session + `'
			}`
	//устанавливаем URL адрес
	l.httpClient.SetUrl("https://lobby.kingdoms.com/api/index.php")
	//отправляем сам запрос
	resp, _, err := l.httpClient.Post("", payload)
	if err != nil {
		l.logger.Error(err)
		return "", err
	}

	gameWorlds, err := Clash(resp, "data", "cache", "response")
	//конвертируем полученный список ситтеров в слайс
	gameWorldsArr, ok := gameWorlds.([]interface{})
	if !ok {
		err := errors.New("Clash: Error convert interface to map ")
		l.logger.Error(err)
		return "", err
	}
	speedTroops, err := Clash(gameWorldsArr[0].(string), "data", "speedTroops")
	if err != nil {
		l.logger.Error(err)
		return "", err
	}
	return speedTroops.(string), nil
}

//функция возвращает скорость передвижения войск игрового мира по его номеру для Заместителя или дуала
func (l *Login) getSitterSpeedTroops(session string, gameWorldString string, sitterType string, sitterName string) string {
	//todo
	return ""
}

//функция возвращает ID играка
//возвращает avatarId  и Id-игрового мира
func (l *Login) getAvatarId(session string, gameWorldString string, sitterType string, sitterName string) (string, string, error) {
	// ignore sitter type for now

	// there are only 4 sitter spots right now, but just to be safe
	var sitterArray string
	for i := 0; i < 10; i++ {
		sitterArray += "{Collection:Sitter:" + string(i) + "},"
	}
	//удаляем последнюю запятую
	sitterArray = sitterArray[:len(sitterArray)-1]
	payload := `{
		action: 'get',
		controller: 'cache',
			params: {
				names: [` + sitterArray + `]
			},
		` + session + `
	}`
	//устанавливаем URL адрес
	l.httpClient.SetUrl("https://lobby.kingdoms.com/api/index.php")
	//отправляем сам запрос
	resp, _, err := l.httpClient.Post("", payload)
	if err != nil {
		l.logger.Error(err)
		return "", "", err
	}
	// получаем список ситтеров с их данными
	sitters, err := Clash(resp, "data", "cache", "response")
	if err != nil {
		l.logger.Error(err)
		return "", "", err
	}
	//конвертируем полученный список ситтеров в слайс
	allSitters, ok := sitters.([]interface{})
	if !ok {
		err := errors.New("Clash: Error convert interface to map ")
		l.logger.Error(err)
		return "", "", err
	}
	//проходим по слайсу ситтеров в поисках нужного ситтера
	for _, sitter := range allSitters {
		// получаем название игрогого мира
		worldName, err := Clash(sitter, "data", "worldName")
		if err != nil {
			l.logger.Error(err)
			return "", "", err
		}
		// получаем ник ситтера
		avatarName, err := Clash(sitter, "data", "avatarName")
		if err != nil {
			l.logger.Error(err)
			return "", "", err
		}
		//проверяем если игровой мир и ситтер совпадает
		//и выводим его айди
		if strings.ToLower(worldName.(string)) == strings.ToLower(gameWorldString) && strings.ToLower(avatarName.(string)) == strings.ToLower(sitterName) {
			avatarIdentifier, err := Clash(sitter, "data", "avatarIdentifier")
			if err != nil {
				l.logger.Error(err)
				return "", "", err
			}
			consumersId, err := Clash(sitter, "data", "consumersId")
			if err != nil {
				l.logger.Error(err)
				return "", "", err
			}
			return avatarIdentifier.(string), consumersId.(string), nil
		}
	}

	err = errors.New("Not finder sitter ")
	return "", "", err
}

//функция проверяет если есть соедине в лобби по сессии
func (l *Login) testLobbyConnection(arsession string) bool {
	//создаем тело запроса в формате JSON строки
	payload := `{
			action: 'getPossibleNewGameworlds',
			controller: 'gameworld',
			params: {},
			session:'` + arsession + `'
			}`
	//устанавливаем URL адрес
	l.httpClient.SetUrl("https://lobby.kingdoms.com/api/index.php")
	//отправляем сам запрос
	resp, _, err := l.httpClient.Post("", payload)
	if err != nil {
		l.logger.Error(err)
		return false
	}

	//возвращаем ответ
	//конвертируем полученный интерфейс в мап и выбираем по очередно с проверкой на ошибку нужный нам елемент

	node3, err := Clash(resp, "data", "error")
	if err != nil {
		l.logger.Error(err)
		return false
	}
	return !(node3.(bool))
}

//функция проверяет если есть соединение с игровым миром по сессии
func (l *Login) testGameWorldConnection(gameWorld string, session string) bool {

	payload := `{
			action: 'get',
			controller: 'cache',
			params: {
				names: ['Player:']
			},
			  session:'` + session + `'
		  }`

	l.httpClient.SetUrl("")
	now := time.Now().UTC().Unix() * 1000
	resp, _, err := l.httpClient.Post("https://"+gameWorld+"}.kingdoms.com/api/?c=cache&a=get&t"+string(now), payload)
	if err != nil {
		return false
	}
	//возвращаем ответ

	//возвращаем ответ
	//конвертируем полученный интерфейс в мап и выбираем по очередно с проверкой на ошибку нужный нам елемент

	node3, err := Clash(resp, "data", "error")
	if err != nil {
		l.logger.Error(err)
		return false
	}
	return !(node3.(bool))
}

//функция определяет скорость игрового мира для ситтера
