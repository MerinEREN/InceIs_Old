package content

import (
	"github.com/MerinEREN/InceIs/account"
	"github.com/MerinEREN/InceIs/cookie"
	usr "github.com/MerinEREN/InceIs/user"
	"io/ioutil"
	"net/http"
)

// II Language and page Sturcts
// Languages colection
/* type Languages []Language

type Language struct {
	Id    string `json:"id"` // EN, TR ...
	pages []page `json:"pages"`
}

type pages []page

type page struct {
	Title string `json:"title"`
	Body  Body `json:"body"`
	// Templates maybe
}

type Body struct {
	Header Header `json:"header"`
	// Others ...
	Footer Footer `json:"footer"`
}

type Header struct {
	// Should be created their own types in the future !!!!!!!!!!!!!!!!!!!!
	SearchPlaceHolder []byte `json:"searchPlaceHolder"`
	MenuButtonText []byte `json:"menuButtonText"`
}

type Footer struct {
	// Should be created their own types in the future !!!!!!!!!!!!!!!!!!!!
	SearchPlaceHolder []byte `json:"searchPlaceHolder"`
	MenuButtonText []byte `json:"menuButtonText"`
} */

/* func (p *page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
} */

// Delete this dummy struct
type Page struct {
	Title       string           `json:"title"`
	User        *usr.User        `json:"user"`
	Body        []byte           `json:"body"`
	Account     *account.Account `json:"account"`
	Form        form             `json:"form"`
	ProfilePic  string           `json:"profile_pic"`
	RedirectURL string           `json:"redirect_url"`
}

type form struct {
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

// GET CONTENT FROM PAGES COLLECTION AND COOKIES !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// RETURN VALUE COULD BE CHANGE AS (interface{}, error) AT THE FUTURE !!!!!!!!!!!!!!!!!!!!!
func Get(r *http.Request, title string) (*Page, error) {
	filename := title + ".html"
	//USE CURRENT WORKING DIRECTORY IN PATH !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// BECAUSE ioutil.ReadFile USES CALLAR PACKAGE'S DIRECTORY AS CURRENT WORKING
	// DIRECTORY !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	body, err := ioutil.ReadFile("../page/templates/" + filename)
	if err != nil {
		return nil, err
	}
	p := &Page{
		Body:  body,
		Title: title,
	}
	if title == "logIn" {
		cd, errCookie := cookie.GetData(r, title)
		if errCookie != http.ErrNoCookie {
			pp := *cd
			p.ProfilePic = pp.Photo
		}
	}
	return p, nil
}
