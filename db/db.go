package db

import (
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
)

const (
	DB_HOST string = "localhost"
	DB_NAME string = "biddinginfo"
)

type Item struct {
	Title     string
	Category  string
	Region    string
	Industry  string
	Date      string
	AgentName string
	UrlDetail string
}

func withCollection(collName string, f func(c *mgo.Collection) error) error {
	session, err := mgo.Dial(DB_HOST)
	if err != nil {
		return err
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true) // Optional. Switch the session to a monotonic behavior.
	c := session.DB(DB_NAME).C(collName)
	return f(c)
}

func SaveAll(collName string, docs []*Item) error {
	return withCollection(collName, func(c *mgo.Collection) error {
		for _, doc := range docs {
			err := c.Insert(doc)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
