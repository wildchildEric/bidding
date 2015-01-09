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
	CurrentPage  int
	CountPerPage int
	TotalCount   int
	Items        interface{}
}

//pageNum is form 1
func GetPage(pageNum, countPerPage int) (Page, error) {
	if pageNum <= 0 {
		pageNum = 0
	} else {
		pageNum -= 1
	}
	var items []*Item
	var total int
	err := withCollection(new(Item).CollectionName(), func(c *mgo.Collection) error {
		var err error
		total, err = c.Count()
		if err != nil {
			return err
		}
		return c.Find(nil).Skip(pageNum * countPerPage).Limit(countPerPage).All(&items)
	})
	return Page{CurrentPage: pageNum + 1, CountPerPage: countPerPage, TotalCount: total, Items: items}, err
}
