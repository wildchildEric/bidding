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

type Page struct {
	PageNum      int
	CountPerPage int
	TotalCount   int
	Items        interface{}
}

func GetPage(pageNum, countPerPage int) (Page, error) {
	var items []*Item
	var total int
	err := withCollection(new(Item).CollectionName(), func(c *mgo.Collection) error {
		var err error
		total, err = c.Find(nil).Count()
		if err != nil {
			return err
		}
		return c.Find(nil).Skip(pageNum * countPerPage).Limit(countPerPage).All(&items)
	})
	return Page{PageNum: pageNum, CountPerPage: countPerPage, TotalCount: total, Items: items}, err
}
