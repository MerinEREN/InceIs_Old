package content

import (
	"io/ioutil"
)

// II Language and Page Sturcts
// Languages colection
/* type Languages []Language

type Language struct {
	Id    string `bson:"id"` // EN, TR ...
	Pages []Page `bson:"pages"`
}

type Pages []Page

type Page struct {
	Title string `bson:"title"`
	Body  Body `bson:"body"`
	// Templates maybe
}

type Body struct {
	Header Header `bson:"header"`
	// Others ...
	Footer Footer `bson:"footer"`
}

type Header struct {
	// Should be created their own types in the future !!!!!!!!!!!!!!!!!!!!
	SearchPlaceHolder []byte `bson:"searchPlaceHolder"`
	MenuButtonText []byte `bson:"menuButtonText"`
}

type Footer struct {
	// Should be created their own types in the future !!!!!!!!!!!!!!!!!!!!
	SearchPlaceHolder []byte `bson:"searchPlaceHolder"`
	MenuButtonText []byte `bson:"menuButtonText"`
} */

/* func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
} */

// Delete this dummy struct
type Page struct {
	Title string `bson:"title"`
	Body  []byte `bson:"body"`
}

// GET CONTENT FROM PAGES COLLECTION !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// RETURN VALUE COULD BE CHANGE AS (interface{}, error) AT THE FUTURE !!!!!!!!!!!!!!!!!!!!!
func Get(title string) (*Page, error) {
	filename := title + ".html"
	//USE CURRENT WORKING DIRECTORY IN PATH !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// BECAUSE ioutil.ReadFile USES CALLAR PACKAGE'S DIRECTORY AS CURRENT WORKING
	// DIRECTORY !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	body, err := ioutil.ReadFile("../page/templates/" + filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
