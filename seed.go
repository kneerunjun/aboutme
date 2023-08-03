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

type ProfileContact struct {
	FBLink     string `json:"fblink" bson:"fblink"`
	GmailLink  string `json:"gmaillink" bson:"gmaillink"`
	LinkedLink string `json:"linkedlink" bson:"linkedlink"`
	GitLink    string `json:"gitlink" bson:"gitlink"`
	Phone      string `json:"phone" bson:"phone"`
	Email      string `json:"email" bson:"email"`
	Address    string `json:"address" bson:"address"`
}

type EducQual struct {
	Start       string `json:"start" bson:"start"`
	End         string `json:"end" bson:"end"`
	Degree      string `json:"degree" bson:"degree"`
	ShortDegree string `json:"sdegree" bson:"sdegree"`
	GovnBody    string `json:"govnbody" bson:"govnbody"`
	Desc        string `json:"desc" bson:"desc"`
	ShortDesc   string `json:"sdesc" bson:"sdesc"`
}
type Skill struct {
	Title string `json:"title" bson:"title"`
	Level int8   `json:"level" bson:"level"`
	Desc  string `json:"desc" bson:"desc"`
	Span  string `json:"span" bson:"span"`
}
type Accolade struct {
	Title string `json:"title" bson:"title"`
	Year  string `json:"year" bson:"year"`
}

type Workexp struct {
	ImgSrc   string `json:"imgsrc" bson:"imgsrc"`
	ImgHt    int    `json:"imght" bson:"imght"`
	ImgWd    int    `json:"imgwd" bson:"imgwd"`
	Desig    string `json:"desig" bson:"desig"`
	Employer string `json:"employer" bson:"employer"`
	Span     string `json:"span" bson:"span"`
}

type Resume struct {
	ID          string         `json:"id" bson:"id"`
	Title       string         `json:"title" bson:"-"`
	FullName    string         `json:"fullname" bson:"fullname"`
	Photo       ProfilePhoto   `json:"photo" bson:"photo"`
	ShortDesc   string         `json:"shortdesc" bson:"shortdesc"`
	ShortDescSm string         `json:"shortdescsm" bson:"shortdescsm"`
	Contact     ProfileContact `json:"contact" bson:"contact,inline"`
	Education   EducQual       `json:"educ" bson:"educ"`
	TopSkills   []Skill        `json:"skills" bson:"skills"`
	Accolades   []Accolade     `json:"accolades" bson:"accolades"`
	Experience  []Workexp      `json:"experience" bson:"experience"`
}

// Data model of a blog, the model that is stored in the database
// When the blog is requested, the data, content is retrieved from the database to be displayed on as html
type Blog struct {
	Id string `json:"id" bson:"id"`
	// cover and title of the blog
	Cover struct {
		Img   string   `json:"img" bson:"img"`
		Title string   `json:"title" bson:"title"`
		Tags  []string `json:"tags" bson:"tags"`
	} `json:"cover" bson:"cover"`
	Preface    string   `json:"preface" bson:"preface"`
	Intro      string   `json:"intro" bson:"intro"`
	Body       string   `json:"body" bson:"body"`
	Conclusion string   `json:"conclusion" bson:"conclusion"`
	References []string `json:"references" bson:"references"`
	Mob        struct {
		Preface    string   `json:"preface" bson:"preface"`
		Intro      string   `json:"intro" bson:"intro"`
		Body       string   `json:"body" bson:"body"`
		Conclusion string   `json:"conclusion" bson:"conclusion"`
		References []string `json:"references" bson:"references"`
	} `json:"mob" bson:"mob"`
}

type DBConfig interface {
	DbName() string
	CollOrTable() string
}

type MongoConfig struct {
	dbName   string
	collName string
}

func (mcfg *MongoConfig) DbName() string {
	return mcfg.dbName
}

func (mcfg *MongoConfig) CollOrTable() string {
	return mcfg.collName
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

// AddResume : adds a new resume to the database
func AddResume(r *Resume) error {
	coll, err := NewDbConn(&MongoConfig{dbName: DB_NAME, collName: COLL_NAME})
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
		ID:       "niranjanawati",
		FullName: "Niranjan Awati",
		Photo: ProfilePhoto{
			Location: "/images/meb_w.jpg",
			HeightPx: 140,
			WidthPx:  140,
		},
		ShortDesc: `Seasoned Go Lang developer with solid 18 years of total experience. An avid IoT junkie, building prototype
		solutions atop single board computers for their sensing capabilities & cloud connectivity. He is adept at
		developing containerized REST API for the web & concurrent applications on IoT devices using Go Lang. He has
		also, in his past contributed extensively to learning functions of his organization.`,
		ShortDescSm: `Seasoned Go Lang developer with solid 18 years of total experience. An avid IoT junkie,He is adept at
		developing containerized REST API for the web & concurrent applications on IoT devices using Go Lang.`,
		Contact: ProfileContact{
			FBLink:     "https://www.facebook.com/kneerunjun/",
			GmailLink:  "mailto:kneerunjun@gmail.com?subject=Reference to your online profile",
			LinkedLink: "https://www.linkedin.com/in/niranjan-awati-a2395856/",
			GitLink:    "https://github.com/kneerunjun",
			Phone:      "+91 8390302623",
			Email:      "kneerunjun@gmail.com",
			Address:    "Sangria, Megapolis Hinjewadi Phase-III, Pune 411057",
		},
		Education: EducQual{
			Start:       "2000",
			End:         "2004",
			Degree:      "Bachelor of Engineering, Mechanical",
			ShortDegree: "B.E. Mechanical",
			GovnBody:    "University of Pune",
			Desc:        "Pursued a 4y bachelor's degree at Maharashtra Institute Technology,Pune. Internal combustion engines as the elective subject in the final year & graduate trainee stint at TATA motors in the year 2004.",
			ShortDesc:   "Pursued a 4y bachelor's degree at Maharashtra Institute Technology,Pune.",
		},
		TopSkills: []Skill{
			{Desc: "Building REST API over HTTP, programming IoT u-controllers using TinyGo", Title: "GoLang", Level: 85, Span: "2017-today"},
			{Desc: "Deep exposure to docker, docker-componse in building portable/scalable apps.", Title: "Docker", Level: 70, Span: "2018-today"},
			{Desc: "Can build single page, responsive apps from ground up.", Title: "AngularJs", Level: 60, Span: "2016-2021"},
			{Desc: "Can build single page, responsive apps from ground up.", Title: "Python", Level: 60, Span: "2016-2021"},
		},
		Accolades: []Accolade{
			{Title: "M.V.P, Infosys", Year: "2007"},
			{Title: "Pride, Boeing", Year: "2008"},
			{Title: "Pride, Boeing", Year: "2009"},
		},
		Experience: []Workexp{
			{ImgSrc: "/images/infy_logo.png", ImgHt: 40, ImgWd: 60, Desig: "Pr. Consultant", Employer: "Infosys Ltd.", Span: "2005-2022"},
			{ImgSrc: "/images/dheeti.jpeg", ImgHt: 45, ImgWd: 40, Desig: "Sr. Developer", Employer: "Dheeti Technologies", Span: "2022-2022"},
			{ImgSrc: "/images/ncs_logo.png", ImgHt: 45, ImgWd: 50, Desig: "Sr. Programmer", Employer: "NCS Technologies", Span: "2022-2023"},
			{ImgSrc: "/images/persistent_logo.png", ImgHt: 45, ImgWd: 50, Desig: "Sr. Architect", Employer: "Persistent", Span: "2023-today"},
		},
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
