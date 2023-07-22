package main

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	SERVER_IP = "mongodb://mongostore:27017"
	DB_NAME   = "aboutme"
	COLL_NAME = "resume"
)

var (
	COLL_NAMES = []string{"resume"} // enlisting all the collection names
)

type ProfilePhoto struct {
	Location string // web location of the photo that appears on the front page splash
	HeightPx int    // height in pixels for the photo
	WidthPx  int    // width in pixels for the photo
}

type Resume struct {
	FullName  string       `json:"fullname" bson:"fullname"`
	Photo     ProfilePhoto `json:"photo" bson:"photo"`
	ShortDesc string       // short description of the profile to summary what the candidate is
}

type DBConfig interface {
}

type MongoConfig struct {
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
	return session.DB(DB_NAME).C(COLL_NAME), nil
}

// AddResume : adds a new resume to the database
func AddResume(r *Resume) error {
	coll, err := NewDbConn(&MongoConfig{})
	if err != nil {
		return fmt.Errorf("failed to connect to database :%s", err)
	}
	if coll == nil {
		return fmt.Errorf("invaild/nil collection, cannot add resume")
	}
	return coll.Insert(r)
}

// NiranjanAwati : seeds niranjan's reume to the database
func NiranjanAwati() error {
	res := &Resume{
		FullName: "Niranjan Awati",
		Photo: ProfilePhoto{
			Location: "/images/meb_w.jpg",
			HeightPx: 200,
			WidthPx:  200,
		},
		ShortDesc: `Seasoned Go Lang developer with solid 18 years of total experience. An avid IoT junkie, building prototype
		solutions atop single board computers for their sensing capabilities & cloud connectivity. He is adept at
		developing containerized REST API for the web & concurrent applications on IoT devices using Go Lang. He has
		also, in his past contributed extensively to learning functions of his organization.`,
	}
	if err := AddResume(res); err != nil {
		return fmt.Errorf("failed to add resume to the database: %s", err)
	}
	return nil
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
