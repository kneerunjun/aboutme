package data

/*================
For each of the blogs there is metadata that is stored onto the database
This has identification and references of the blog but NOT the content
File here describes the datamodel of such, and also the seed data
================*/

import (
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

// BlogListResult : on the page where you have list of all the blogs with search feature
// this has a list of all the blogs + some flags for the UI components
type BlogListResult struct {
	List        []Blog `json:"list"`        // actual result
	ClearSearch bool   `json:"clearsearch"` // if the list is based on search, the clear button can help you get back to listing all the blogs

}

// Blog : the database model to record the blog
// skeleton information of the blog
// blog body to be loaded as html from the page
type Blog struct {
	CoverImg    string              `bson:"coverimg" json:"coverimg"` // cover image appears on the blog page
	Slug        string              `bson:"slug" json:"slug"`         // slug that uniquely represents the blog also appears in the url as text. Is derived from the title see :SlugifyTitle
	Title       string              `bson:"title" json:"title"`
	Summary     string              `bson:"summary" json:"summary"`
	Tags        []string            `bson:"tags" json:"tags"`
	References  []map[string]string `bson:"refs" json:"refs"`
	PubDate     string              `bson:"pubdate" json:"pubdate"` // date of publishing the blog
	PubLoc      string              `bson:"publoc" json:"publoc"`   //location from which the blog was published
	AuthorName  string              `json:"author" bson:"author"`
	AboutAuthor string              `bson:"aboutauthor" json:"aboutauthor"`
	AuthorEmail string              `bson:"authoremail" json:"authoremail"`
}

// SlugifyTitle : takes the title of the blog and makes a abrdiged string that can used as a slug
// this slug can be used to uniquely identify the blog in the database
// also the slug gets used in the url as the url param to uniquely identify the blog yet keep it human readable.
/*
	sampleTitle := "this is a sample title"
	abridged, _:= Slugify(sampleTitle)
	fmt.Println(abridged)
	// "this-is"
*/
func SlugifyTitle(title string) (string, error) {
	// this is the pattern we are looking for in the title
	titlPtrn, err := regexp.Compile(`^[\w\d]+[\s]*[\w\d]*`)
	if err != nil {
		return "", fmt.Errorf("error in compiling pattern expression")
	}
	abridged := titlPtrn.FindString(title)
	if abridged == "" {
		// case when the string is empty and this happens when the title is empty or there isnt any match
		return "", fmt.Errorf("title is empty, or not according to as expected")
	}
	splitSlug := strings.Split(abridged, " ")
	var result string
	for i, v := range splitSlug {
		if i > 0 {
			result = fmt.Sprintf("%s-%s", result, strings.ToLower(v))
		} else {
			result = strings.ToLower(v)
		}
	}
	return result, nil
}

// below is the seed of all the blogs.
// All the blogs that I'd written were on multiple sites, hence getting them under one roof
// https://kneerunjun.wordpress.com/
var (
	blogSeed = []Blog{
		{
			CoverImg: "/images/rpicloseup.jpg",
			Slug:     "",
			Title:    "Reading LM35 with RaspberryPi using amplified RC timer",
			Tags: []string{
				"iot", "rctimer", "raspberrypi", "lm35", "temperature",
			},
			References: []map[string]string{
				{"text": "Raspberry Pi in teaching", "link": "https://www.raspberrypi.org/teach/"},
				{"text": "Maximum amperage thru a Raspberry Pi", "link": "https://raspberrypi.stackexchange.com/questions/9298/what-is-the-maximum-current-the-gpio-pins-can-output"},
				{"text": "RC charging circuit", "link": "https://www.electronics-tutorials.ws/rc/rc_1.html"},
				{"text": "Raspberry Pi computers aboard the International Space Station (ISS).", "link": "https://astro-pi.org/"},
			},
			PubDate:     "25-OCT-2016",
			PubLoc:      "Pune, India",
			AuthorName:  "Niranjan Awati",
			AboutAuthor: "Niranjan is an IoT junkie & GoLang developer",
			AuthorEmail: "kneerunjun@gmail.com",
			Summary:     "When measuring low level analogue voltages would you use a chip or roll up your own. Say you are onto a prototype, what would be your choice?",
		},
		{
			CoverImg: "/images/angularjs.png",
			Slug:     "",
			Title:    "Tabby Tab Angular tab-control in a jiffy!",
			Tags: []string{
				"webdev", "angularjs", "re-usecontrols", "frontend", "javascript",
			},
			References: []map[string]string{
				{"text": "Source code on GitHub", "link": "https://github.com/kneerunjun/tabby-tab"},
			},
			PubDate:     "25-SEP-2015",
			PubLoc:      "Pune, India",
			AuthorName:  "Niranjan Awati",
			AboutAuthor: "Niranjan is an IoT junkie & GoLang developer",
			AuthorEmail: "kneerunjun@gmail.com",
			Summary:     "Tabs aren't a popular choice when it comes to mobile-first design. Just incase you need an angularjs directive to get an array of tabs up and running in no time",
		},
		{
			CoverImg: "/images/angularphoto.jpg",
			Slug:     "",
			Title:    "Angular directive with conditional transclusion & discrete compile",
			Tags: []string{
				"webdev", "angularjs", "re-usecontrols", "frontend", "javascript",
			},
			References: []map[string]string{
				{"text": "Source code on GitHub", "link": "https://github.com/kneerunjun/tabby-tab"},
			},
			PubDate:     "14-FEB-2016",
			PubLoc:      "Pune, India",
			AuthorName:  "Niranjan Awati",
			AboutAuthor: "Niranjan is an IoT junkie & GoLang developer",
			AuthorEmail: "kneerunjun@gmail.com",
			Summary:     "Bending your mind around the 'deep-sea' concepts of Angualrjs. Rarely used in common scenarious but can make your life a tad bit easier when understood.",
		},
		{
			CoverImg: "/images/angularrelativity.png",
			Slug:     "",
			Title:    "That relativity of angular broadcasts",
			Tags: []string{
				"webdev", "angularjs", "re-usecontrols", "frontend", "javascript",
			},
			References: []map[string]string{
				{"text": "GithubGist", "link": "https://gist.github.com/kneerunjun/7d95d3c1db15c1e62352"},
			},
			PubDate:     "07-NOV-2015",
			PubLoc:      "Pune, India",
			AuthorName:  "Niranjan Awati",
			AboutAuthor: "Niranjan is an IoT junkie & GoLang developer",
			AuthorEmail: "kneerunjun@gmail.com",
			Summary:     "Caution when using $broadcast( ) in AngularJS",
		},
		{
			CoverImg: "/images/dockershipping.png",
			Slug:     "",
			Title:    "Testing Django apps live on docker containers",
			Tags: []string{
				"webdev", "docker", "django", "python", "devops",
			},
			References: []map[string]string{
				{"text": "GithubGist", "link": "https://gist.github.com/kneerunjun/7d95d3c1db15c1e62352"},
			},
			PubDate:     "04-MAR-2017",
			PubLoc:      "Pune, India",
			AuthorName:  "Niranjan Awati",
			AboutAuthor: "Niranjan is an IoT junkie & GoLang developer",
			AuthorEmail: "kneerunjun@gmail.com",
			Summary:     "On how to quickstart setting up docker containers for Django Apps. Sounds very basic but a couple of easy pitfalls can waste a lot of effort.",
		},
		{
			CoverImg: "images/helppage.jpg",
			Slug:     "",
			Title:    "Help page ecosystem for your angular SPAs",
			Tags: []string{
				"angularjs", "webdev",
			},
			References: []map[string]string{
				{"text": "Single Page Applications, MDN glossary", "link": "https://developer.mozilla.org/en-US/docs/Glossary/SPA"},
			},
			PubDate:     "04-MAR-2017",
			PubLoc:      "Pune, India",
			AuthorName:  "Niranjan Awati",
			AboutAuthor: "Niranjan is an IoT junkie & GoLang developer",
			AuthorEmail: "kneerunjun@gmail.com",
			Summary:     "Everyone needs help pages, its so frustrating to not find any when required the most. Here is how you can jumpstart",
		},
		{
			CoverImg: "images/purebool.png",
			Slug:     "",
			Title:    "Binding pure boolean values to scope of isolated Angular directives",
			Tags: []string{
				"angularjs", "webdev",
			},
			References: []map[string]string{
				{"text": "Single Page Applications, MDN glossary", "link": "https://developer.mozilla.org/en-US/docs/Glossary/SPA"},
			},
			PubDate:     "28-NOV-2015",
			PubLoc:      "Pune, India",
			AuthorName:  "Niranjan Awati",
			AboutAuthor: "Niranjan is an IoT junkie & GoLang developer",
			AuthorEmail: "kneerunjun@gmail.com",
			Summary:     "",
		},
		{
			CoverImg: "images/angular-3.jpg",
			Slug:     "",
			Title:    "Nomenclature blip in AngularJS providers",
			Tags: []string{
				"angularjs", "webdev", "javascript",
			},
			References: []map[string]string{
				{"text": "Single Page Applications, MDN glossary", "link": "https://developer.mozilla.org/en-US/docs/Glossary/SPA"},
			},
			PubDate:     "10-OCT-2015",
			PubLoc:      "Pune, India",
			AuthorName:  "Niranjan Awati",
			AboutAuthor: "Niranjan is an IoT junkie & GoLang developer",
			AuthorEmail: "kneerunjun@gmail.com",
			Summary:     "AngularJS providers have that tiny dark spot, miss it and it can pester you for hours. Save time and read this before you make your own providers.",
		},
		{
			CoverImg: "images/raspwifi.jpg",
			Slug:     "",
			Title:    "Autoconnect WiFi on Raspbian Stretch",
			Tags: []string{
				"raspberrypi", "wifi", "network",
			},
			References:  []map[string]string{},
			PubDate:     "05-JUN-2018",
			PubLoc:      "Pune, India",
			AuthorName:  "Niranjan Awati",
			AboutAuthor: "Niranjan is an IoT junkie & GoLang developer",
			AuthorEmail: "kneerunjun@gmail.com",
			Summary:     "Unless you are a beginner you'd be running raspbian on headless mode on all your Pis. Here is how you can auto connect WiFi on your device on setup.",
		},
	}
)

// SeedBlogs :  this will take the seed data of the blogs and push to database
// Will also convert the blog title to slugs
func SeedBlogs() {
	for _, b := range blogSeed {
		if err := AddBlog(&b); err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("failed to add seed blog database")
			continue
		}
	}
}
