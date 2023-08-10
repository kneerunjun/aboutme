package data

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	SERVER_IP  = "mongodb://mongostore:27017"
	DB_NAME    = "aboutme"
	COLL_NAME  = "resume"
	BLOGS_COLL = "blogs"
)

var (
	COLL_NAMES = []string{COLL_NAME, BLOGS_COLL} // enlisting all the collection names
)

type DBConfig interface {
	DbName() string
	CollOrTable() string
}

type MongoConfig struct {
	DBName   string
	CollName string
}

func (mcfg *MongoConfig) DbName() string {
	return mcfg.DBName
}

func (mcfg *MongoConfig) CollOrTable() string {
	return mcfg.CollName
}

// NewDbConn : helps to instantiate a new connection and ping the db server upon successful connection
// will take in the configuration of the datbase as an interface
// will send back a session, or connection object after having to set the mode of the connection
func NewDbConn(cfg DBConfig) (*mgo.Collection, error) {
	session, err := mgo.Dial(SERVER_IP)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, fmt.Errorf("nil session, cannot cotinue")
	}
	session.SetMode(mgo.Monotonic, true)
	return session.DB(cfg.DbName()).C(cfg.CollOrTable()), nil
}
func AddBlog(r *Blog) error {
	coll, err := NewDbConn(&MongoConfig{DBName: DB_NAME, CollName: BLOGS_COLL})
	if err != nil {
		return fmt.Errorf("failed to connect to database :%s", err)
	}
	if coll == nil {
		return fmt.Errorf("invaild/nil collection, cannot add new blog")
	}

	// form the slug and then check if slug is already assigned
	slug, err := SlugifyTitle(r.Title)
	if err != nil {
		return fmt.Errorf("failed to slugify the title %s", err)
	}
	// Getting to see if the slug is unique
	cnt, err := coll.Find(bson.M{"slug": slug}).Count()
	if err != nil {
		return fmt.Errorf("failed to verify if the slug is unique")
	}
	if cnt > 0 {
		// slug isnt unique
		return fmt.Errorf("%s slug is already assigned to one of the blog, try changing the title", slug)
	}
	log.WithFields(log.Fields{
		"url": fmt.Sprintf("/blogs/%s", slug),
	}).Debug("inserting blog")
	r.Slug = slug
	return coll.Insert(r)
}

// AddResume : adds a new resume to the database
func AddResume(r *Resume) error {
	coll, err := NewDbConn(&MongoConfig{DBName: DB_NAME, CollName: COLL_NAME})
	if err != nil {
		return fmt.Errorf("failed to connect to database :%s", err)
	}
	if coll == nil {
		return fmt.Errorf("invaild/nil collection, cannot add resume")
	}
	return coll.Insert(r)
}

func FlushDB() error {
	session, err := mgo.Dial(SERVER_IP)
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("nil session, cannot cotinue")
	}
	for _, name := range COLL_NAMES {
		session.DB(DB_NAME).C(name).RemoveAll(bson.M{})
	}
	return nil
}
