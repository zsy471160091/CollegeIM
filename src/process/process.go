package process

import (
	"github.com/donnie4w/go-logger/logger"
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
)

var g_DB_URL string

func InitDB(db_url string) error {
	session, err := mgo.Dial(db_url)
	if err != nil {
		logger.Fatal(err.Error())
		return err
	}
	defer session.Close()
	g_DB_URL = db_url
	return nil
}
