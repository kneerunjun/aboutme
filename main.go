package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	FVerbose, FLogF, FSeed bool
	logFile                string
	// IMP: setting this to true would mean all the recent changes to the dtabase are lost and overriden with seed data from within the code
)

func init() {
	/* ======================
	-verbose=true would mean log.Debug can work
	-verbose=false would mean log.Debug will be hidden
	-flog=true: all the log output shall be onto a file
	-flog=false: all the log output shall be on stdout
	- We are setting the default log level to be Info level
	======================= */
	// flag.BoolVar(&FVerbose, "verbose", false, "Level of logging messages are set here")
	// flag.BoolVar(&FLogF, "flog", false, "Direction in which the log should output")
	// flag.BoolVar(&FSeed, "dbseed", false, "flag to force seed the database, use it at your own risk")

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

// ServeIndexHtml : will dispatch the index.html page
func ServeIndexHtml(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
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
	r.Static("/views", fmt.Sprintf("%s/views/", dirStatic))
	r.LoadHTMLGlob(fmt.Sprintf("%s/pages/*", dirStatic))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"app": "aboutme",
		})
	})
	r.GET("/", ServeIndexHtml)
	log.Fatal(r.Run(":8080"))
}
