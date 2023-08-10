package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	mgo "gopkg.in/mgo.v2"
)

func TestMgoConnection(t *testing.T) {
	session, err := mgo.Dial("mongodb://localhost:47017")
	assert.Nil(t, err, "failed to dial session")

	session.SetMode(mgo.Monotonic, true)
	coll := session.DB(DB_NAME).C(COLL_NAME)
	assert.NotNil(t, coll, "nil collection, unexpected")

}

func TestSlugify(t *testing.T) {
	samples := []string{
		"Sample blog",
		"saMPle",
		"Sample Test Blog",
	}
	for _, v := range samples {
		result, err := SlugifyTitle(v)
		assert.Nil(t, err, "Unexpected error when sugifying")
		t.Log(result)
	}
}
