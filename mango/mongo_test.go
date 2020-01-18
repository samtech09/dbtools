package mango

import (
	"context"
	"testing"
)

type testuser struct {
	ID   int    `bson:"_id"`
	Name string `bson:"name"`
}

func (t testuser) GetID() interface{} {
	return t.ID
}
func (t *testuser) ToInterface(list []testuser) []interface{} {
	var islice []interface{} = make([]interface{}, len(list))
	for i, d := range list {
		islice[i] = d
	}
	return islice
}

func TestMongo(t *testing.T) {
	cfg := MongoConfig{}
	cfg.Host = "192.168.60.206"
	cfg.Port = 27017
	cfg.DbName = "testdb"

	ses := InitSession(cfg)
	defer ses.Cleanup()
	c := ses.GetColl("testusers")

	c.Drop(context.Background())

	u := testuser{1, "test1"}
	err := ses.SaveSingle(c, u)
	if err != nil {
		t.Errorf("Error saving doc: %s", err.Error())
		t.FailNow()
	}

	usr1 := testuser{2, "test2"}
	usr2 := testuser{3, "test3"}
	err = ses.InsertBulk(c, usr1, usr2)
	if err != nil {
		t.Errorf("Error insert bulk: %s", err.Error())
		t.FailNow()
	}

	// insert using slice
	usr := []testuser{}
	usr = append(usr, testuser{4, "test4"})
	usr = append(usr, testuser{5, "test5"})
	//var islice []interface{} = make([]interface{}, len(usr))
	islice := usr1.ToInterface(usr)
	for i, d := range usr {
		islice[i] = d
	}
	err = ses.InsertBulk(c, islice...)
	if err != nil {
		t.Errorf("Error insert bulk slice: %s", err.Error())
		t.FailNow()
	}

	err = ses.DeleteSingle(c, 1)
	if err != nil {
		t.Errorf("Error deleting doc: %s", err.Error())
		t.FailNow()
	}
}
