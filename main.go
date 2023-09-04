package main

/* ================
Project 	: aboutme, Aug2023
Link		: https://github.com/kneerunjun/aboutme
Author		: kneerunjun@gmail.com
Copyright	: Eensymachines
Desc		: Website hosted on docker to publish a public face, host blog,
and all about niranjan awati as a profile. Also stores the latest resume. This can be one stop that HR/ recruiters can download details from
================ */
import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/kneerunjun/aboutme/data"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	FVerbose, FLogF, FSeed bool
	logFile                string
	// IMP: setting this to true would mean all the recent changes to the dtabase are lost and overriden with seed data from within the code
)

/*
================
- from the env gets a configuration flags
- sets the global variables
- logging configuration
  - direction of logs
  - verbosity of logs
  - file configuration of logs if it needs to be ilogged to files and not to console

================
*/
func init() {
	if val := os.Getenv("LOG_VERBOSITY"); val == "y" {
		FVerbose = true
	}
	if val := os.Getenv("FILE_LOG"); val == "y" {
		FLogF = true
	}
	if val := os.Getenv("DB_SEED"); val == "y" {
		FSeed = true
	}

	// Setting up log configuration for the api
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
		ForceColors:   true,
	})
	log.SetReportCaller(false)
	// By default the log output is stdout and the level is info
	log.SetOutput(os.Stdout)     // FLogF will set it main, but dfault is stdout
	log.SetLevel(log.DebugLevel) // default level info debug but FVerbose will set it main
	logFile = os.Getenv("LOGF")
	log.WithFields(log.Fields{
		"seed": FSeed,
	}).Debug("now chcking for the seed variable")
}

/*================
MakeInsertDBConn:  Will insert DB connection object
incase of failed connection will revert with 502 gateway failed
================*/

func MakeInsertDBConn(collName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		coll, err := data.NewDbConn(&data.MongoConfig{DBName: data.DB_NAME, CollName: collName})
		if err != nil {
			log.Error("failed to connect to database")
			c.AbortWithStatus(http.StatusBadGateway)
			return
		}
		c.Set("conn", coll)
	}
}

// CloseDBconn : closes the database connection on the way out
func CloseDBconn(c *gin.Context) {
	val, ok := c.Get("conn")
	if !ok {
		return
	}
	coll, _ := val.(*mgo.Collection)
	if coll == nil {
		return
	}
	coll.Database.Session.Close()
}

/*
================
Middleware that helps insert the db connection to the chain of handlers.
this is deprecated since we have MakeInsertDBConn which can give customizable collection object from the name
================
*/
func InsertDBConn(c *gin.Context) {
	resumeColl, err := data.NewDbConn(&data.MongoConfig{DBName: data.DB_NAME, CollName: data.COLL_NAME})
	if err != nil {
		log.Error("failed to connect to database")
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}
	c.Set("conn", resumeColl)
}
func ServeIndexHtml(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"Title": "About me"})
}

// renderBlogList : this sends out a list of all the blogs chrnologically
// from there on the link to the actual blog
// this also caters to list of blogs that are filtered on the search bar
func renderBlogList(c *gin.Context) {
	val, ok := c.Get("conn")
	if !ok {
		log.Error("Cannot server Html, no connection to database")
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}
	searchPhrase := c.Query("search") // this is when user is trying to search for specific blogs with title
	pttrnMtch := (regexp.MustCompile(`^[\w\d\s]*$`)).Match([]byte(searchPhrase))
	if !pttrnMtch {
		log.WithFields(log.Fields{
			"phrase": searchPhrase,
		}).Warn("Suspicious search phrase")
		c.HTML(http.StatusBadRequest, "400.html", data.ErrPayload{
			Code:   http.StatusBadRequest,
			Msg:    "Invalid search phrase,Search phrases are simple alphanumeric for searching the blogs by the title. Check the search phrase and try all over again",
			Status: "Bad Request",
			GoBack: "/blogs/",
		})
		return
	}
	flt := bson.M{}
	if searchPhrase != "" {
		// https://stackoverflow.com/questions/10610131/checking-if-a-field-contains-a-string
		flt = bson.M{"title": bson.M{"$regex": searchPhrase, "$options": "i"}}
		log.WithFields(log.Fields{
			"phrase": searchPhrase,
		}).Debug("we have a search phrase")
	}

	coll, ok := val.(*mgo.Collection) // blogs collection
	if !ok {
		log.Error("invalid object for mgo.Collection, check and try again")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	result := data.BlogListResult{List: []data.Blog{}, ClearSearch: false} // result list of all the blogs
	coll.Find(flt).All(&result.List)                                       // getting the list of all the blogs
	// Settign the clear flag
	if searchPhrase != "" {
		result.ClearSearch = true
	}

	log.WithFields(log.Fields{
		"count_blogs": len(result.List),
	}).Debug("requested for the list of all the blogs")
	c.HTML(http.StatusOK, "blog-list.html", result)
}

// renderBlog : handler when the client requests a blog of the single id
// since we ARENT designing a single page application, here each blog would have its individual page.
// each page for the blog would have some elements that are part of the template and majority of the content that is specific to the blog
// cover image, title, references, can be tempalated but the body of the blog remains specific
func renderBlog(c *gin.Context) {
	blogid, ok := c.Params.Get("blogid")
	if !ok {
		log.Error("invalid or empty blog ID")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	log.WithFields(log.Fields{
		"blog": blogid,
	}).Debug("rendering blog")
	val, ok := c.Get("conn")
	if !ok {
		log.Error("Cannot server Html, no connection to database")
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}
	coll, ok := val.(*mgo.Collection)
	if !ok {
		log.Error("invalid object for mgo.Collection, check and try again")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	result := data.Blog{}
	err := coll.Find(bson.M{"slug": blogid}).One(&result)
	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			log.WithFields(log.Fields{
				"id": blogid,
			}).Error("failed to get blog of slug")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		} else {
			// case when the query has failed - this could be of failed gateway , but would be reported as InternalError
			log.WithFields(log.Fields{
				"err": err,
			}).Error("query to get profile failed")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
	log.WithFields(log.Fields{
		"title": result.Title,
	}).Debug("found blog in database")
	c.HTML(http.StatusOK, fmt.Sprintf("%s.html", result.Slug), result)
}

// renderMyProfile : will dispatch the index.html page
func renderMyProfile(c *gin.Context) {
	userid, ok := c.Params.Get("userid")
	if !ok {
		log.Error("invalid or empty userid for /myprofile")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	val, ok := c.Get("conn")
	if !ok {
		log.Error("Cannot server Html, no connection to database")
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}
	coll, ok := val.(*mgo.Collection)
	if !ok {
		log.Error("invalid object for mgo.Collection, check and try again")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	result := data.Resume{}
	err := coll.Find(bson.M{"id": userid}).One(&result)
	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			log.WithFields(log.Fields{
				"id": userid,
			}).Error("failed to get profile of userid")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		} else {
			// case when the query has failed - this could be of failed gateway , but would be reported as InternalError
			log.WithFields(log.Fields{
				"err": err,
			}).Error("query to get profile failed")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
	result.Title = "About me"
	log.WithFields(log.Fields{
		"title":    result.Title,
		"fullname": result.FullName,
	}).Debug("spitting out the result")
	c.HTML(http.StatusOK, "index.html", result)
}

func main() {
	flag.Parse() // command line flags are parsed
	log.WithFields(log.Fields{
		"verbose": FVerbose,
		"flog":    FLogF,
		"seed":    FSeed,
	}).Info("Log configuration..")
	if FVerbose {
		log.SetLevel(log.DebugLevel)
	}
	if FLogF {
		lf, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("Failed to connect to log file, kindly check the privileges")
		} else {
			log.Infof("Check log file for entries @ %s", logFile)
			log.SetOutput(lf)
		}
	}
	if FSeed {
		log.Warn("Seed flag set to true, flushing the data. This will be replaced with seed data..")
		data.FlushDB()
		if err := data.NiranjanAwati(); err != nil {
			log.WithFields(log.Fields{"err": err}).Error("Error seeding the database")
		}
		data.SeedBlogs()
	}
	// Loading all environment variables
	dirStatic := os.Getenv("DIR_STATIC")
	log.WithFields(log.Fields{
		"static_dir": dirStatic,
	}).Debug("echoing static directory")
	log.Info("Starting server..")
	defer log.Warn("Server now shutting down..")

	// Seed the database only if the seed flag is on
	// TODO:  incase the seed flag is set the database details have to be dropped

	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	r.Static("/images", fmt.Sprintf("%s/images/", dirStatic))
	r.Static("/js", fmt.Sprintf("%s/js/", dirStatic))
	r.LoadHTMLGlob(fmt.Sprintf("%s/templates/**/*", dirStatic))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"app": "aboutme",
		})
	})
	r.GET("/myprofile/:userid", InsertDBConn, renderMyProfile, CloseDBconn)
	r.GET("/blogs/", MakeInsertDBConn("blogs"), renderBlogList, CloseDBconn)
	r.GET("/blogs/:blogid", MakeInsertDBConn(data.BLOGS_COLL), renderBlog, CloseDBconn)
	// r.GET("/views/:name", InsertDBConn, ServeView)
	log.Fatal(r.Run(":8080"))
}
