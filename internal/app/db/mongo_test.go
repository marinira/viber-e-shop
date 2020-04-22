package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMongoSetConfig(t *testing.T) {
	var m_db DataBaseInterface
	m_db = new(Mongo)
	err := m_db.SetConfig()
	assert.Error(t, err, "1) set config error arguments ")
	err = m_db.SetConfig("a0", "a1", "a2", "a3", "a4")
	assert.Equal(t, err, nil, "2) set config correct arguments ")
	var ret interface{}
	ret, err = m_db.GetConfig()
	assert.Equal(t, err, nil, "3,1) get config not error ")
	assert.Equal(t, ret, [5]string{"a0", "a1", "a2", "a3", "a4"}, "3,2) get config correct arguments ")
	ret1 := m_db.GetNameDataBase()
	assert.Equal(t, ret1, "a0", "3,3) get config getNameDataBase ")
}

func TestMongo_InitConnect(t *testing.T) {
	//t.Run("InitConnect ", func(t *testing.T) {
	var m_db DataBaseInterface
	m_db = NewMongo()
	err := m_db.InitConnect()
	assert.NoError(t, err, "expected valid deployment, got nil")
	//})
}

func TestMongo_Disconnect(t *testing.T) {
	//t.Run("Disconnect Correct ", func(t *testing.T) {
	var m_db DataBaseInterface
	m_db = NewMongo()
	err := m_db.InitConnect()
	err = m_db.Disconnect()
	assert.NoError(t, err, "expected valid deployment, got nil")
	//})
	//t.Run("Disconnect Not Correct", func(t *testing.T) {
	//var m_db DataBaseInterface
	m_db = NewMongo()
	err = m_db.Disconnect()
	assert.Error(t, err, "expected valid deployment, got nil")
	//})
}

func TestMongo_InsertOne(t *testing.T) {
	//t.Run("InsertOne Correct ", func(t *testing.T) {
	var m_db DataBaseInterface
	m_db = NewMongo()
	err := m_db.InitConnect()
	type Trainer struct {
		Name string
		Age  int
		City string
	}
	el := Trainer{Name: "n1", Age: 5, City: "C1"}
	err = m_db.InsertOne("testcollection", el)
	assert.NoError(t, err, "expected valid deployment, got nil")
	err = m_db.Disconnect()
	//})
}

func TestMongo_InsertMany(t *testing.T) {
	//t.Run("InsertMany Correct ", func(t *testing.T) {
	var m_db DataBaseInterface
	m_db = NewMongo()
	err := m_db.InitConnect()
	type Trainer struct {
		Name string
		Age  int
		City string
	}
	var al_arr = make([]interface{}, 5)

	for i := 0; i < 5; i++ {
		el := Trainer{Name: "n" + string(i), Age: i, City: "C" + string(i)}
		al_arr[i] = el
	}

	err = m_db.InsertMany("testcollection", al_arr)
	assert.NoError(t, err, "expected valid deployment, got nil")
	err = m_db.Disconnect()
	//})
}

func TestMongo_DeleteOne(t *testing.T) {
	//t.Run("InsertOne Correct ", func(t *testing.T) {
	var m_db DataBaseInterface
	m_db = NewMongo()
	err := m_db.InitConnect()
	type Trainer struct {
		Name string
		Age  int
		City string
	}
	el := Trainer{Name: "n1", Age: 5, City: "C1"}
	err = m_db.DeleteOne("testcollection", el)
	assert.NoError(t, err, "expected valid deployment, got nil")
	err = m_db.Disconnect()
	//})
}

func TestMongo_DeleteMany(t *testing.T) {
	//t.Run("InsertMany Correct ", func(t *testing.T) {
	var m_db DataBaseInterface
	m_db = NewMongo()
	err := m_db.InitConnect()
	type Trainer struct {
		Name string
		Age  int
		City string
	}
	type Filter struct {
		Age int
	}
	var al_arr interface{}

	el := Filter{Age: 3}
	al_arr = el
	err = m_db.DeleteMany("testcollection", al_arr)
	assert.NoError(t, err, "expected valid deployment, got nil")
	err = m_db.Disconnect()
	//})
}

func TestMongo_UpdateOne(t *testing.T) {
	//t.Run("InsertOne Correct ", func(t *testing.T) {
	var m_db DataBaseInterface
	m_db = NewMongo()
	err := m_db.InitConnect()
	type Trainer struct {
		Name string
		Age  int
		City string
	}
	el := Trainer{Name: "n1", Age: 5, City: "C1"}
	type Filter struct {
		Age int
	}
	var fl_arr interface{}

	fl := Filter{Age: 0}
	fl_arr = fl
	err = m_db.UpdateOne("testcollection", fl_arr, el)
	assert.NoError(t, err, "expected valid deployment, got nil")
	err = m_db.Disconnect()
	//})
}

func TestMongo_UpdateMany(t *testing.T) {
	//t.Run("InsertOne Correct ", func(t *testing.T) {
	var m_db DataBaseInterface
	m_db = NewMongo()
	err := m_db.InitConnect()
	type Trainer struct {
		Name string
		Age  int
		City string
	}
	el := Trainer{Name: "n1", Age: 5, City: "C1"}
	type Filter struct {
		Age int
	}
	var fl_arr interface{}

	fl := Filter{Age: 0}
	fl_arr = fl
	err = m_db.UpdateMany("testcollection", fl_arr, el)
	assert.NoError(t, err, "expected valid deployment, got nil")
	err = m_db.Disconnect()
	//})
}

func TestMongo_Find(t *testing.T) {
	//t.Run("InsertMany Correct ", func(t *testing.T) {
	var m_db DataBaseInterface
	m_db = NewMongo()
	err := m_db.InitConnect()
	type MFilter struct {
		Age int
	}
	var al_arr interface{}
	al_arr = MFilter{Age: 5}

	_, err = m_db.Find("testcollection", al_arr)
	//logrus.Info(arr)
	assert.NoError(t, err, "expected valid deployment, got nil")
	err = m_db.Disconnect()
	//})
}

func TestMongo_FindOne(t *testing.T) {
	//t.Run("InsertMany Correct ", func(t *testing.T) {
	var m_db DataBaseInterface
	m_db = NewMongo()
	err := m_db.InitConnect()
	type MFilter struct {
		Age int
	}
	var al_arr interface{}
	al_arr = MFilter{Age: 5}

	_, err = m_db.FindOne("testcollection", al_arr)
	//logrus.Info(arr)
	assert.NoError(t, err, "expected valid deployment, got nil")
	err = m_db.Disconnect()
	//})
}
func TestMongo_FindOneAndDelete(t *testing.T) {
	//t.Run("InsertMany Correct ", func(t *testing.T) {
	var m_db DataBaseInterface
	m_db = NewMongo()
	err := m_db.InitConnect()
	type MFilter struct {
		Age int
	}
	var al_arr interface{}
	al_arr = MFilter{Age: 1}

	err = m_db.FindOneAndDelete("testcollection", al_arr)
	assert.NoError(t, err, "expected valid deployment, got nil")
	err = m_db.Disconnect()
	//})
}

func TestMongo_FindOneAndReplace(t *testing.T) {
	//t.Run("InsertMany Correct ", func(t *testing.T) {
	var m_db DataBaseInterface
	m_db = NewMongo()
	err := m_db.InitConnect()

	type Trainer struct {
		Name string
		Age  int
		City string
	}
	el := Trainer{Name: "n1", Age: 50, City: "C1"}
	type MFilter struct {
		Age int
	}
	var al_arr interface{}
	al_arr = MFilter{Age: 5}

	err = m_db.FindOneAndReplace("testcollection", al_arr, el)
	assert.NoError(t, err, "expected valid deployment, got nil")
	err = m_db.Disconnect()
	//})
}

func TestMongo_FindOneAndUpdate(t *testing.T) {
	//t.Run("InsertMany Correct ", func(t *testing.T) {
	var m_db DataBaseInterface
	m_db = NewMongo()
	err := m_db.InitConnect()

	type Trainer struct {
		Name string
		Age  int
		City string
	}
	el := Trainer{Name: "n1", Age: 51, City: "C1"}
	type MFilter struct {
		Age int
	}
	var al_arr interface{}
	al_arr = MFilter{Age: 50}

	err = m_db.FindOneAndUpdate("testcollection", al_arr, el)
	assert.NoError(t, err, "expected valid deployment, got nil")
	err = m_db.Disconnect()
	//})
}

func TestMongo_CountDocuments(t *testing.T) {
	//t.Run("InsertMany Correct ", func(t *testing.T) {
	var m_db DataBaseInterface
	m_db = NewMongo()
	err := m_db.InitConnect()
	type MFilter struct {
		Age int
	}
	var al_arr interface{}
	al_arr = MFilter{Age: 51}

	_, err = m_db.CountDocuments("testcollection", al_arr)
	//logrus.Info(arr)
	assert.NoError(t, err, "expected valid deployment, got nil")
	err = m_db.Disconnect()
	//})
}
