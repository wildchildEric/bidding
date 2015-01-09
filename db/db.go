package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const (
	DB_HOST string = "localhost"
	DB_NAME string = "biddinginfo"
)

type Item struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
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
	// for _, it := range items {
	// 	log.Printf("%+v", it)
	// }
	return Page{CurrentPage: pageNum + 1, CountPerPage: countPerPage, TotalCount: total, Items: items}, err
}

func GetItems(ids []string) ([]*Item, error) {
	var items []*Item
	err := withCollection(new(Item).CollectionName(), func(c *mgo.Collection) error {
		obj_ids := make([]bson.ObjectId, 0, len(ids))
		for _, id := range ids {
			obj_ids = append(obj_ids, bson.ObjectIdHex(id))
		}
		return c.Find(bson.M{"_id": bson.M{"$in": obj_ids}}).All(&items)
	})
	return items, err
}

func Try() {
	var items []*Item
	err := withCollection(new(Item).CollectionName(), func(c *mgo.Collection) error {
		ids := []bson.ObjectId{bson.ObjectIdHex("54af9a318587fbdade997b1c"), bson.ObjectIdHex("54af9a328587fbdade998629")}
		return c.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&items)
	})
	if err != nil {
		log.Println(err)
	}
	for _, it := range items {
		log.Printf("%+v", it)
	}
}
