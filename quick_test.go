package main

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)
// TestCompanyName : when requesting for the resume to be sent using email, requestor needs to provide the company name 
// this test helps us fixate the regex used to verify if the company name is valid
// we need to weed out nuissance value 
func TestCompanyName(t *testing.T) {
	companyNames := []string{
		"Schumm-Conroy",
		"Gleason LLC",
		"Connelly, Rempel and Wolf",
		"Collier-O'Conner",
		"FDG Inc, India",
		"Satterfield, Lubowitz & Torphy",
		"Padberg and 2Sons",
	}
	patt := regexp.MustCompile(`^[a-zA-Z'-]+[\s]*[a-zA-Z0-9\s'-,&]*$`)
	for _, val := range companyNames {
		assert.Equal(t, true, patt.MatchString(val), fmt.Sprintf("pattern failed to verify %s", val))
	}
	notOkNames := []string{
		" ",
		"",
		"$%#$%",
		"_",
		"Schumm-Conroy %",
		"12313",
	}
	for _, val := range notOkNames {
		assert.Equal(t, false, patt.MatchString(val), fmt.Sprintf("pattern failed to verify %s", val))
	}
}

func TestEmailvalidation(t *testing.T) {
	emails := []string{
		"niranjan_awati@gmail.com",
		"kneerunjun@gmail.com",
		"niranjan.awati@gmail.com",
		"niranjan1_awati@gmail.com",
		"niranjan-awati@gmail.com",
		"niranjan@gmail.co1",
		"324343_awati@gmail.com",
		"niranjan_324343@gmail.com",
	}
	patt := regexp.MustCompile(`^[a-zA-Z0-9]+[-_.]{0,1}[a-zA-Z0-9]*@[a-zA-Z0-9]+.[a-zA-Z0-9]+$`)
	for _, val := range emails {
		assert.Equal(t, true, patt.MatchString(val), fmt.Sprintf("pattern failed to verify %s", val))
	}
	notokEmails := []string{
		"",
		" ",
		"@gmail.com",
		"-@gmail.com",
		"_@gmail.com",
		"niranjan__awati@gmail.com",
		"niranjan%awati@gmail.com",
		"niranjan@gmail.co.in",
	}
	for _, val := range notokEmails {
		assert.Equal(t, false, patt.MatchString(val), fmt.Sprintf("pattern failed to verify %s", val))
	}

}
