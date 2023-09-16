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
	"io"
	"net/http"
	"os"
	"reflect"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/kneerunjun/aboutme/data"
	gmail "github.com/kneerunjun/aboutme/mail"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	SMTP_SECRET = "/run/secrets/smtp_secret"
)

var (
	FVerbose, FLogF, FSeed bool
	logFile                string
	// IMP: setting this to true would mean all the recent changes to the dtabase are lost and overriden with seed data from within the code
	domain_name, gmailAppPass string
	sender, resumepath        string // when sending email notifications, this is the email address used
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
	// IMP: when a cookie is set for a domain the same cannot be accessed from another
	// say for example cookie is of the domain eensymachines.in, when testing the cookie would fail the domain would be localhost
	// this hence needs to be set from environment variables
	domain_name = "localhost"
	// gmailAppPass = os.Getenv("GMAIL_APPPASS")
	// GMAIL_APPPASS

	sender = os.Getenv("GMAIL_SENDER")
	resumepath = os.Getenv("RESUME_PATH")
	log.WithFields(log.Fields{
		"sender": sender,
		"resume": resumepath,
	}).Debug("Loaded environment")
	// Load passwords from secrets
	f, err := os.Open(SMTP_SECRET)
	if err != nil {
		panic(fmt.Sprintf("failed to load smtp secret from file, cannot proceed %s", err))
	}
	byt, err := io.ReadAll(f)
	if err != nil {
		panic(fmt.Sprintf("failed to load smtp secret from file, cannot proceed %s", err))
	}
	gmailAppPass = string(byt)
	log.WithFields(log.Fields{
		"secret": len(gmailAppPass),
	}).Debug("Loaded SMTP secret")
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

// Hndl200OKRedirect : whenever an operation typically a post http call runs as expected this handler can show the user the contextual success message
// From other handlers whenever redirect to /success this handler here takes over
// All the other params of the request remain the same
func Hndl200OKRedirect(c *gin.Context) {
	cookie, err := c.Cookie("aboutme-200ok")
	if err != nil {
		// Error reading cookie
		// this should not stop us from loading the page.
		// Afterall the operation preceeding this must have completed corrrectly - that is what matters
		log.WithFields(log.Fields{
			"err": err,
		}).Debug("Error reading the success message cookie")
		cookie = ""
	}
	c.HTML(http.StatusOK, "200OK.html", gin.H{"successMsg": cookie})
}

// RequestProfileOnEmail : will take you to the page where resume of the candidate can be requested over email
func RequestProfileOnEmail(c *gin.Context) {
	pyld := gin.H{"invalid_email": false, "invalid_company": false}
	if c.Request.Method == "GET" {
		// sending the page where the profile can be requested
		c.HTML(http.StatusOK, "req-resume.html", gin.H{
			"emailed": false,
		})
	} else if c.Request.Method == "POST" {
		// TODO: this request has to be idempotent - upon refreshing the success page the same request is sent
		// this happens cause despite the page change, the url still remains the same. Only contents of the page change the browser does NOT navigate away from the page
		// this is when the user needs the server to send the resume on email
		// this is only upon sending the correct information from the form
		// will store the details of the requestor on the database and send the email containing the resume
		formData := map[string]string{}
		if err := c.Bind(&formData); err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("RequestProfileOnEmail: failed to bind form data")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		emailPttrn := regexp.MustCompile(`^[a-zA-Z0-9]+[-_.]{0,1}[a-zA-Z0-9]*@[a-zA-Z0-9]+.[a-zA-Z0-9]+$`)
		cpmnyPttrn := regexp.MustCompile(`^[a-zA-Z'-]+[\s]*[a-zA-Z0-9\s'-,&]*$`)
		// TODO: incase of invalid entries on the email and the company name , the same page gets loaded with input boxes being highlighted for relevant fields
		if !emailPttrn.MatchString(formData["reqemail"]) {
			pyld["invalid_email"] = true
		}
		if !cpmnyPttrn.MatchString(formData["reqcompany"]) {
			pyld["invalid_company"] = true
		}
		// Starting a new instance of a notifier
		notifier, err := gmail.NewMailNotify(gmail.MailConfig{
			Host: "smtp.gmail.com", Port: 587, UName: "awatiniranjan@gmail.com", Passwd: "imbilafrkzilxvwv",
		}, reflect.TypeOf(&gmail.GmailNotify{}))
		if err != nil {
			// Gateway error - this should go into the cookie followed by a redirect request
			c.AbortWithStatus(http.StatusBadGateway)
			return
		}
		body := "Hi,<br>As requested I'm attaching my latest resume alongside.<br>Best regards,<br>Niranjan"
		go func(n gmail.MailNotify) {
			err := n.SendFileAttach(sender, formData["reqemail"], "Resume: Niranjan Awati", body, resumepath)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("Failed to send email")
				err := n.SendErrNotification(sender, formData["reqemail"])
				if err != nil {
					// this can happen when the recipient's email id itself is incorrect
					// ideally speaking we shouldnt be even trying to send the error notification to this address
					log.WithFields(log.Fields{
						"err": err,
					}).Error("Error sending the error notification")
				}
			}

		}(notifier)
		ckeMsg := fmt.Sprintf("Kindly check at %s for a pdf copy of the resume", formData["reqemail"])
		c.SetCookie("aboutme-200ok", ckeMsg, 3600, "/success", domain_name, true, true)
		c.Redirect(http.StatusPermanentRedirect, "/success")
	}
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
	r.GET("/notifications/email/myprofile", RequestProfileOnEmail)  // can send the profile to the requestor via email
	r.POST("/notifications/email/myprofile", RequestProfileOnEmail) // can send the profile to the requestor via email

	r.GET("/blogs/", MakeInsertDBConn("blogs"), renderBlogList, CloseDBconn)
	r.GET("/blogs/:blogid", MakeInsertDBConn(data.BLOGS_COLL), renderBlog, CloseDBconn)

	r.POST("/success", Hndl200OKRedirect)

	// r.GET("/views/:name", InsertDBConn, ServeView)
	log.Fatal(r.Run(":8080"))
}
