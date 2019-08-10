package mongo

import (
	"gopkg.in/mgo.v2"
)

const (
	MongodbUrl = "127.0.0.1:27017"
	MongodbName = "test"
	MongodbCollection = "ctest"
)

func MgoLink() (s *mgo.Session, err error){
	//创建连接
	session, err := mgo.Dial(MongodbUrl)
	if err != nil {
		return session, err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	return session, nil
}

func MgoInsert(table []TableItem) (err error){
	session, err := MgoLink()
	if err != nil {
		return err
	}

	c := session.DB(MongodbName).C(MongodbCollection)
	if err := c.Insert(table); err != nil {
		return err
	}
	return nil
}