package mango

import (
	"context"
	"testing"

	bf "github.com/samtech09/bsonquery"
	"go.mongodb.org/mongo-driver/bson"
)

type testuser struct {
	ID   int    `bson:"_id"`
	Name string `bson:"name"`
	Age  int    `bson:"age"`
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

	u := testuser{1, "test1", 21}
	err := ses.SaveSingle(c, u)
	if err != nil {
		t.Errorf("Error saving doc: %s", err.Error())
		t.FailNow()
	}

	usr1 := testuser{2, "test2", 45}
	usr2 := testuser{3, "test3", 30}
	err = ses.InsertBulk(c, usr1, usr2)
	if err != nil {
		t.Errorf("Error insert bulk: %s", err.Error())
		t.FailNow()
	}

	// insert using slice
	usr := []testuser{}
	usr = append(usr, testuser{4, "test4", 48})
	usr = append(usr, testuser{5, "test5", 32})
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

	_, err = c.Find(context.Background(), bson.M{})
	if err != nil {
		t.Errorf("Error finding docs: %s", err.Error())
		t.FailNow()
	}

	// test filter
	filter := bf.Builder().
		And(bf.C().EQ("name", "test4"), bf.C().GT("age", 30)).
		Build()
	cur, err := c.Find(context.Background(), filter)
	if err != nil {
		t.Errorf("Error finding docs with filter: %s", err.Error())
		t.FailNow()
	}

	count := 0
	for cur.Next(context.TODO()) {
		count++
	}
	if err := cur.Err(); err != nil {
		t.Errorf("Error interating curser: %s", err.Error())
	}

	if count != 1 {
		t.Errorf("Filter failed. Expected %d,  Got: %d", 1, count)
	}

}
