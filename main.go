package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	FVerbose, FLogF, FSeed bool
	logFile                string
	// IMP: setting this to true would mean all the recent changes to the dtabase are lost and overriden with seed data from within the code
)

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

func InsertDBConn(c *gin.Context) {
	resumeColl, err := NewDbConn(&MongoConfig{dbName: DB_NAME, collName: COLL_NAME})
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
	result := Resume{}
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
		FlushDB()
		if err := NiranjanAwati(); err != nil {
			log.WithFields(log.Fields{"err": err}).Error("Error seeding the database")
		}
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
	r.GET("/myprofile/:userid", InsertDBConn, renderMyProfile)
	// r.GET("/views/:name", InsertDBConn, ServeView)
	log.Fatal(r.Run(":8080"))
}
