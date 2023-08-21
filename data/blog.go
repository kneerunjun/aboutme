package data

import (
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Blog : the database model to record the blog
// skeleton information of the blog
// blog body to be loaded as html from the page
type Blog struct {
	CoverImg    string              `bson:"coverimg" json:"coverimg"`
	Slug        string              `bson:"slug" json:"slug"`
	Title       string              `bson:"title" json:"title"`
	Tags        []string            `bson:"tags" json:"tags"`
	References  []map[string]string `bson:"refs" json:"refs"`
	PubDate     string              `bson:"pubdate" json:"pubdate"`
	PubLoc      string              `bson:"publoc" json:"publoc"`
	AuthorName  string              `json:"author" bson:"author"`
	AboutAuthor string              `bson:"aboutauthor" json:"aboutauthor"`
	AuthorEmail string              `bson:"authoremail" json:"authoremail"`
}

// SlugifyTitle : takes the title of the blog and makes a abrdiged string that can used as a slug
// this slug can be used to uniquely identify the blog in the database
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
