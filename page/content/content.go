package content

import (
	"github.com/MerinEREN/InceIs/account"
	// "github.com/MerinEREN/InceIs/cookie"
	usr "github.com/MerinEREN/InceIs/user"
	// "io/ioutil"
	//"net/http"
)

// II Language and page Sturcts
// Languages colection
// STORE ALL OF THOSE IN TO THE DATASTORE !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
type Languages []Language

type Language struct {
	Id    string `datastore:"" json:"id"` // EN, TR ...
	Pages []Page `datastore:"" json:"pages"`
}

type Pages []Page

type Page struct {
	C C `datastore:"" json:"content"`
	D D `datastore:"" json:"data"`
}

type C struct {
	Title string `datastore:"" json:"title"`
	Body  Body   `datastore:"" json:"body"`
}

type Body struct {
	Header Header `datastore:"" json:"header"`
	// Others ...
	Footer Footer `datastore:"" json:"footer"`
}

type Header struct {
	// Should be created their own types in the future !!!!!!!!!!!!!!!!!!!!
	SearchPlaceHolder []byte `datastore:"" json:"searchPlaceHolder"`
	MenuButtonText    []byte `datastore:"" json:"menuButtonText"`
}

type Footer struct {
	// Should be created their own types in the future !!!!!!!!!!!!!!!!!!!!
	SearchPlaceHolder []byte `datastore:"" json:"searchPlaceHolder"`
	MenuButtonText    []byte `datastore:"" json:"menuButtonText"`
}

type D struct {
	User     *usr.User        `datastore:"-" json:"user"`
	Account  *account.Account `datastore:"-" json:"account"`
	LoginURL string           `datastore:"-" json:"login_url"`
	URLUUID  string           `datastore:"-" json:"url_uuid"`
}

/* func (p *page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
} */

// GET PAGE CONTENTS FROM DATASTORE WITH USERS SELECTED LANGUAGE !!!!!!!!!!!!!!!!!!!!!!!!!!
func Get(title string) (*Page, error) {
	// filename := title + ".html"
	//USE CURRENT WORKING DIRECTORY IN PATH !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// BECAUSE ioutil.ReadFile USES CALLAR PACKAGE'S DIRECTORY AS CURRENT WORKING
	// DIRECTORY !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	/* body, err := ioutil.ReadFile("../page/templates/" + filename)
	if err != nil {
		return nil, err
	} */
	content := C{
		Title: title,
	}
	p := &Page{}
	p.C = content
	return p, nil
}
