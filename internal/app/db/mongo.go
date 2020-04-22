package db

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type Mongo struct {
	nameDataBase string
	hostName     string
	portName     string
	userName     string
	pass         string
	client       *mongo.Client
	logger       *logrus.Logger
	sessions     map[string]interface{}
}

//Конструктор - создает екземпляр класса Mongo и инициализирует поля по умолчанию
func NewMongo() *Mongo {
	//создаем новый екземпляр класса Mongo
	ret := new(Mongo)
	//создаем екземпляр logrus (логирование) для екземпляра  класса Mongo
	ret.logger = logrus.New()
	// инициализаруем занчения полей по умолчанию
	ret.init()
	return ret
}

//функция инициализации значения полей по умолчанию.
//данные берутся из констант из mongo_config
func (m *Mongo) init() {
	m.nameDataBase = NameDataBaseDefault
	m.hostName = HostNameDefault
	m.portName = PortNameDefault
	m.userName = UserNameDefault
	m.pass = PassDefault
}

//функция инициализации соединения с базою данных. возвращает ошибку при возникновении проблемм
func (m *Mongo) InitConnect() (err error) {

	// Create client
	m.client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://" + m.hostName + ":" + m.portName))
	if err != nil {
		m.logger.Error(err)
		return err
	}

	// Create context
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// Create connect
	err = m.client.Connect(ctx)
	if err != nil {
		m.logger.Error(err)
		return err
	}
	return nil
}

// функция Disconnect с Базаю данных. возвращает NIL или ошибку
func (m *Mongo) Disconnect() (err error) {
	//проверяем если открыто соединение с базою данных
	if err_ping := m.ping(); err_ping != nil {
		m.logger.Error(err_ping)
		return err_ping
	}
	//создаем контекст с ожиданием в 10 сек
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//отсоединяемся от базы данных
	err = m.client.Disconnect(ctx)
	if err != nil {
		m.logger.Error(err)
		return err
	}
	return nil
}

// функция для проверки если есть соединения с базаю данных
func (m *Mongo) ping() (err error) {
	// Check the connection
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	if m.client == nil {
		return errors.New("Error: SetConfig Mongo, not correct count of arguments ")
	}
	err = m.client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}
	return nil
}

//установить названия базы данных
func (m *Mongo) SetNameDataBase(n string) {
	m.nameDataBase = n
}

//получить название базы данных
func (m *Mongo) GetNameDataBase() string {
	return m.nameDataBase
}

// установить все поля для структуры Mongo
func (m *Mongo) SetConfig(c ...string) (err error) {
	// проверяем если количество входных аргументов совпадает с количествоп полей в структуре
	// если совпадает то присваиваем значения иначе возвращаем ошибку
	if len(c) == 5 {
		m.nameDataBase = c[0]
		m.hostName = c[1]
		m.portName = c[2]
		m.userName = c[3]
		m.pass = c[4]

	} else {
		return errors.New("Error: SetConfig Mongo, not correct count of arguments ")
	}
	return nil
}

// получить масив строк всех поля для структуры Mongo
func (m *Mongo) GetConfig() (ret interface{}, err error) {
	var resp [5]string
	resp[0] = m.nameDataBase
	resp[1] = m.hostName
	resp[2] = m.portName
	resp[3] = m.userName
	resp[4] = m.pass
	//проверяем если первые 4 поля с данными иначе кидаем ошибку
	for _, element := range resp {
		if element == "" {
			errs := errors.New("Error: GetConfig Mongo, d'nt exist all arguments ")
			return nil, errs
		}
	}
	return resp, nil
}

//CRUD
func (m *Mongo) createTable(t string) (err error) {
	// TODO
	return nil
}
func (m *Mongo) deleteTable(t string) (err error) {
	// TODO
	return nil
}

//InsertOne - вставить елемент element в таблицу table
func (m *Mongo) InsertOne(table string, element interface{}) (err error) {
	// проверяем если есть соединение с базою
	if err_ping := m.ping(); err_ping != nil {
		m.logger.Error(err_ping)
		return err_ping
	}
	// делаем соединение с коллекцией
	collection := m.client.Database(m.nameDataBase).Collection(table)
	//создаем контекст с ожиданием 5 сек
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//вставляем елемент
	_, err = collection.InsertOne(ctx, element)
	//id := res.InsertedID
	if err != nil {
		m.logger.Error(err)
		return err
	}
	return nil
}

//InsertMany - вставить слайс елементоы elements в таблицу table
func (m *Mongo) InsertMany(table string, elements []interface{}) (err error) {
	//проверяем если есть соединение с базою
	if err_ping := m.ping(); err_ping != nil {
		m.logger.Error(err_ping)
		return err_ping
	}
	//связываемся с коллекциею из базв
	collection := m.client.Database(m.nameDataBase).Collection(table)
	//создаем контекст с задержкой в 5 сек
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//устанавливаем опции для результата "без сортировке"
	opts := options.InsertMany().SetOrdered(false)
	// вставляем слайс елементов в коллекцию
	_, err = collection.InsertMany(ctx, elements, opts)
	if err != nil {
		m.logger.Error(err)
		return err
	}
	return nil
}

// DeleteOne - удаляет 1 елемент (первый попавшийся) из таблицы table
//который соответствует фильтру filter
func (m *Mongo) DeleteOne(table string, filter interface{}) (err error) {
	//проверяем соединение с базею
	if err_ping := m.ping(); err_ping != nil {
		m.logger.Error(err_ping)
		return err_ping
	}
	//создаем екземпляр коллекции
	collection := m.client.Database(m.nameDataBase).Collection(table)
	//создаем контекст с задержкой в 5 сек
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	// устанавливаем опции для удаления
	opts := options.Delete().SetCollation(&options.Collation{
		Locale:    "en_US",
		Strength:  1,
		CaseLevel: false,
	})
	//удаляем елемент
	_, err = collection.DeleteOne(ctx, filter, opts)
	if err != nil {
		m.logger.Error(err)
		return err
	}
	return nil
}

// DeleteMany - удаляет все елементы  из таблицы table
//которые соответствуют фильтру filter
func (m *Mongo) DeleteMany(table string, element interface{}) (err error) {
	//проверяем если есть соединение с базою
	if err_ping := m.ping(); err_ping != nil {
		m.logger.Error(err_ping)
		return err_ping
	}
	//создаем екземпляр коллекции
	collection := m.client.Database(m.nameDataBase).Collection(table)
	//создаем контекст с задержкою в 5 сек
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	// устанавливаем опции для удаления
	opts := options.Delete().SetCollation(&options.Collation{
		Locale:    "en_US",
		Strength:  1,
		CaseLevel: false,
	})
	//удаляем елементы из базы
	_, err = collection.DeleteMany(ctx, element, opts)
	if err != nil {
		m.logger.Error(err)
		return err
	}
	return nil
}

// Find - поиск всех документов в базе которые соответсвуют фильтру filter
//результата возвращает слайс документов(елементов)
func (m *Mongo) Find(table string, filter interface{}) (result []interface{}, err error) {
	//проверяем если есть подключение к базе
	if err_ping := m.ping(); err_ping != nil {
		m.logger.Error(err_ping)
		return nil, err_ping
	}
	//создаем екземпляр коллекции
	collection := m.client.Database(m.nameDataBase).Collection(table)
	//создаем конетекст с задержкой 5 сек
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//делаем поиск в коллекции по фильтру
	cursor, err := collection.Find(ctx, filter, nil)
	if err != nil {
		m.logger.Error(err)
		return nil, err
	}
	// get a list of all returned documents and print them out
	// see the mongo.Cursor documentation for more examples of using cursors
	var results []interface{} //bson.M
	if err = cursor.All(ctx, &results); err != nil {
		m.logger.Error(err)
		return nil, err
	}

	return results, nil
}

// FindOne - поиск документа в базе который соответсвуют фильтру filter
//результата возвращает  документа(елемента)
func (m *Mongo) FindOne(table string, element interface{}) (result interface{}, err error) {
	//проверяем если есть связь с базою
	if err_ping := m.ping(); err_ping != nil {
		m.logger.Error(err_ping)
		return nil, err_ping
	}
	//создаем екземпляр колекции
	collection := m.client.Database(m.nameDataBase).Collection(table)
	// создаем контекст с задержклй 5 сек
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//указываем опции для результата поиска
	opts := options.FindOne().SetSort(bson.D{})
	//делаем поиск и готовим результат
	err = collection.FindOne(ctx, element, opts).Decode(&result)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			m.logger.Error(err)
			return nil, nil
		}
		//log.Fatal(err)
		m.logger.Error(err)
		return nil, err
	}
	return &result, nil
}

// FindOneAndDelete - поиск документа в базе который соответсвуют фильтру filter
//и удаляем его
func (m *Mongo) FindOneAndDelete(table string, element interface{}) (err error) {
	//проверяем если есть соединение с базорю
	if err_ping := m.ping(); err_ping != nil {
		m.logger.Error(err_ping)
		return err_ping
	}
	//создаем екземпляр коллекции
	collection := m.client.Database(m.nameDataBase).Collection(table)
	//создаем контекст с задержкой 5 сек
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//делаем поиск и удаляем
	var result interface{}
	err = collection.FindOneAndDelete(ctx, element, nil).Decode(&result)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			m.logger.Error(err)
			return err
		}
		log.Fatal(err)
		m.logger.Error(err)
		return err
	}
	return nil
}

// FindOneAndDelete - поиск документа в базе который соответсвуют фильтру filter
//и замещаем его
func (m *Mongo) FindOneAndReplace(table string, filter interface{}, element interface{}) (err error) {
	//проверяем связь с базою
	if err_ping := m.ping(); err_ping != nil {
		m.logger.Error(err_ping)
		return err_ping
	}
	//делаем екземпляр коллекции
	collection := m.client.Database(m.nameDataBase).Collection(table)
	//создаем контекст с задержкой 5 сек
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//делаем сам поиск и замещение
	var result interface{}
	err = collection.FindOneAndReplace(ctx, filter, element, nil).Decode(&result)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			m.logger.Error(err)
			return err
		}
		//log.Fatal(err)
		m.logger.Error(err)
		return err
	}
	return nil
}

// FindOneAndDelete - поиск документа в базе который соответсвуют фильтру filter
//и обновляем его
func (m *Mongo) FindOneAndUpdate(table string, filter interface{}, element interface{}) (err error) {
	//проверяем связь с базою
	if err_ping := m.ping(); err_ping != nil {
		m.logger.Error(err_ping)
		return err_ping
	}
	//создаем екземпляр коллекции
	collection := m.client.Database(m.nameDataBase).Collection(table)
	//создаем контекс с задержкой 5 сек
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//делаем сам поиск и обновление
	var result interface{}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	update := bson.D{{"$set", element}}
	err = collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			m.logger.Error(err)
			return err
		}
		//log.Fatal(err)
		m.logger.Error(err)
		return err
	}
	return nil
}

// UpdateOne - обновляем первый документа в базе который соответсвуют фильтру filter
func (m *Mongo) UpdateOne(table string, filter interface{}, element interface{}) (err error) {
	//проверяем связь с базою
	if err_ping := m.ping(); err_ping != nil {
		m.logger.Error(err_ping)
		return err_ping
	}
	//создаем екземпляр коллекции
	collection := m.client.Database(m.nameDataBase).Collection(table)
	//создаем контекс с задержкой 5 сек
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	// устанавливаем поции обновления
	opts := options.Update().SetUpsert(true)
	update := bson.D{{"$set", element}}
	//делаем само обновление
	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		m.logger.Error(err)
		return err
	}
	if result.MatchedCount != 0 {
		m.logger.Info("matched and replaced an existing document")

	}
	if result.UpsertedCount != 0 {
		m.logger.Info("inserted a new document with ID \n", result.UpsertedID)
	}
	return nil
}

// UpdateMany - обновляем все документы в базе который соответсвуют фильтру filter
func (m *Mongo) UpdateMany(table string, filter interface{}, element interface{}) (err error) {
	//проверяем связь с базою
	if err_ping := m.ping(); err_ping != nil {
		m.logger.Error(err_ping)
		return err_ping
	}
	//создаем екземпляр коллекции
	collection := m.client.Database(m.nameDataBase).Collection(table)
	//создаем контекс с задержкой 5 сек
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//устанавливаем опции обновления
	update := bson.D{{"$set", element}}
	//проводим само обновление
	result, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		m.logger.Error(err)
		return err
	}

	if result.MatchedCount != 0 {
		m.logger.Info("matched and replaced an existing document")
		return
	}
	if result.UpsertedCount != 0 {
		m.logger.Info("inserted a new document with ID \n", result.UpsertedID)
	}
	return nil
}

//CountDocuments - определяем количество все документы в базе который
//соответсвуют фильтру filter возвращаем число
func (m *Mongo) CountDocuments(table string, filter interface{}) (count int64, err error) {
	//проверяем если есть связь с базою
	if err_ping := m.ping(); err_ping != nil {
		m.logger.Error(err_ping)
		return 0, err_ping
	}
	//создаем екземпляр коллекции
	collection := m.client.Database(m.nameDataBase).Collection(table)
	//создаем контекст с задержкою 5 сек
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//устанавливаем опциизапроса
	opts := options.Count().SetMaxTime(2 * time.Second)
	//получаем колличество совпадений
	count, err = collection.CountDocuments(ctx, filter, opts)
	if err != nil {
		log.Fatal(err)
		m.logger.Error(err)
		return 0, err
	}
	return count, nil
}
