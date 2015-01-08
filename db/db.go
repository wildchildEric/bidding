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

func (p *Item) CollectionName() string {
	return "chinabiddings"
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

func SaveAll(docs []*Item) error {
	return withCollection(new(Item).CollectionName(), func(c *mgo.Collection) error {
		for _, doc := range docs {
			err := c.Insert(doc)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func GetAll(page int) ([]*Item, error) {
	perPage := 100
	var items []*Item
	err := withCollection(new(Item).CollectionName(), func(c *mgo.Collection) error {
		return c.Find(nil).Skip(page * perPage).Limit(perPage).All(&items)
	})
	return items, err
}
