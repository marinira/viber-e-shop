package db

// interface for DataBase

type DataBaseInterface interface {
	// func initialisation Connect to data base
	InitConnect() (err error)
	Disconnect() (err error)
	SetNameDataBase(n string)
	GetNameDataBase() string
	SetConfig(c ...string) (err error)
	GetConfig() (ret interface{}, err error)
	InsertOne(table string, element interface{}) (err error)
	InsertMany(table string, element []interface{}) (err error)
	DeleteOne(table string, element interface{}) (err error)
	DeleteMany(table string, element interface{}) (err error)
	UpdateOne(table string, filter interface{}, element interface{}) (err error)
	UpdateMany(table string, filter interface{}, element interface{}) (err error)
	CountDocuments(table string, filter interface{}) (count int64, err error)
	Find(table string, element interface{}) (result []interface{}, err error)
	FindOne(table string, element interface{}) (result interface{}, err error)
	FindOneAndDelete(table string, element interface{}) (err error)
	FindOneAndReplace(table string, filter interface{}, element interface{}) (err error)
	FindOneAndUpdate(table string, filter interface{}, element interface{}) (err error)
}
