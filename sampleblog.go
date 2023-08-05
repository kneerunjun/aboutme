package main

/* ==================
a sample blog that can be pushed to the database while development / testing
this blog hasnt got much production significance
=====================*/

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

var sampleBlog Blog = Blog{
	Id: bson.NewObjectId().Hex(),
	Cover: BlogCover{
		Img:   BlogImg{Src: "/images/yelopencil.jpg", Ht: 200},
		Title: "Sample blog",
		Tags:  []string{"sample", "telegram", "test"},
	},
	Preface:    `Bacon ipsum dolor amet alcatra cow andouille bresaola fatback. Tongue kielbasa pancetta, flank capicola turducken burgdoggen tail cupim. Ground round short ribs chislic andouille rump pork loin shankle brisket filet mignon shank t-bone jerky leberkas. Meatball cupim bacon beef ribs. Capicola kielbasa picanha leberkas, cow meatloaf filet mignon turkey. Strip steak shankle andouille doner jerky.`,
	Intro:      `Boudin beef picanha, short ribs flank ribeye capicola chicken. Andouille sausage meatball bresaola ribeye corned beef tenderloin. Chuck porchetta fatback cupim chislic, landjaeger short loin tongue frankfurter alcatra ham. Tongue beef bacon jerky meatloaf tenderloin picanha jowl ham hock.`,
	Body:       `Chislic prosciutto buffalo t-bone turducken, jowl hamburger frankfurter pork chop. Turkey pork loin shoulder, picanha pork belly beef ribs chislic leberkas. Jerky beef spare ribs tenderloin, prosciutto kielbasa ground round shankle tongue alcatra kevin. Pastrami biltong alcatra jowl sausage kielbasa. Turducken chicken short loin, rump tongue alcatra swine salami ham hock landjaeger tail jowl short ribs pig. Fatback leberkas hamburger meatloaf tri-tip shankle shank short ribs pork belly sirloin chislic andouille. Shoulder sausage cupim frankfurter venison drumstick, alcatra chislic shank.`,
	Conclusion: `Turkey cupim pork loin prosciutto, tail pig porchetta beef chuck meatloaf spare ribs kevin filet mignon jerky. Sirloin chicken tongue, ham hock bacon beef ribs frankfurter. Buffalo doner chuck, short ribs jowl ribeye shoulder shank salami cupim chicken jerky. Sirloin pastrami salami, bacon beef ribs brisket beef short ribs buffalo.`,
	Mob: MobBlog{
		Preface:    `Bacon ipsum dolor amet alcatra cow andouille bresaola fatback. Tongue kielbasa pancetta, flank capicola turducken burgdoggen tail cupim.`,
		Intro:      `Boudin beef picanha, short ribs flank ribeye capicola chicken. Andouille sausage meatball bresaola ribeye corned beef tenderloin.`,
		Body:       `erky beef spare ribs tenderloin, prosciutto kielbasa ground round shankle tongue alcatra kevin. Pastrami biltong alcatra jowl sausage kielbasa. Turducken chicken short loin, rump tongue alcatra swine salami ham hock landjaeger tail jowl short ribs pig. Fatback leberkas hamburger meatloaf tri-tip shankle shank short ribs pork belly sirloin chislic andouille.`,
		Conclusion: `Turkey cupim pork loin prosciutto, tail pig porchetta beef chuck meatloaf spare ribs kevin filet mignon jerky. Sirloin chicken tongue, ham hock bacon beef ribs frankfurter.`,
	},
	References: []string{
		"https://baconipsum.com/?paras=5&type=all-meat&start-with-lorem=1",
	},
}

func SampleBlog() error {
	if err := AddBlog(&sampleBlog); err != nil {
		return fmt.Errorf("failed to add new sample blog to database %s", err)
	}
	return nil
}
