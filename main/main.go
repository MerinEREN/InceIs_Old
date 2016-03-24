package main

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

const (
	mongoDBLocal = "mongodb://localhost"
)

// OO DB Structs
// Accounts collection docs
type Accounts struct {
	Name          string `bson: "name"`
	LastName      string `bson: "lastName"`
	Type          string `bson: "type"`
	CurrentStatus string `bson: "currentStatus"`
	AccountStatus string `bson: "accountStatus"`
	About         string `bson: "about"`
	Tags          Tags
	Ranks         Ranks
	Card          Card
	// Users Users
}

type Company struct {
	Name    string `bson: "name"`
	Address Address
}

type Address struct {
	Details     string `bson : "details"`
	Borough     string `bson : "borough"`
	City        string `bson : "city"`
	Country     string `bson : "country"`
	Postcode    string `bson : "postcode"`
	Geolocation Geolocation
}

type Geolocation struct {
	Lat  string `bson: "lat"`
	Long string `bson: "Long"`
}

type Tags []Tag

type Tag struct {
	Type string `bson: "type"`
}

type Ranks []Rank

type Rank struct {
	Type string `bson: "type"`
}

type Page struct {
	Title string
	Body  []byte
}

type Card struct {
	CreditCards CreditCards
	DebitCards  DebitCards
}

type CreditCards []CreditCard

type CreditCard struct {
	HolderName string `bson: "holderName"`
	No         string `bson: "no"`
	ExpMonth   string `bson: "expMonth"`
	ExpYear    string `bson: "expYear"`
	CVV        string `bson: "cvv"`
}

type DebitCards []DebitCard

type DebitCard struct {
	HolderName string `bson: "holderName"`
	No         string `bson: "no"`
	ExpMonth   string `bson: "expMonth"`
	ExpYear    string `bson: "expYear"`
	CVV        string `bson: "cvv"`
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func main() {
	p1 := &Page{Title: "Test Page", Body: []byte("This is a sample Page.")}
	if p1.save() != nil {
		fmt.Println("Shit happens!!!!!!")
	}
	p2, _ := loadPage("Test Page")
	fmt.Println(string(p2.Body))

	session, err := mgo.Dial(mongoDBLocal)
	if err != nil {
		fmt.Println(session, err)
	}

	c := session.DB("OO").C("accounts")

	c.Insert(&Accounts{Name: "Merin", LastName: "EREN"})
	var result interface{}
	c.Find(Accounts{Name: "Merin", LastName: "EREN"}).One(&result)
	// acc := c.Find(Accounts{Name: "Merin", LastName: "EREN"})
	// acc := c.Find(bson.M{"name": "Merin"})
	fmt.Println(result)

	/*res, err := http.Get("http://www.google.com/robots.txt")
	if err != nil {
		log.Fatal(err)
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", robots)*/
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":8080", nil)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		// The http.Redirect function adds an HTTP status code of
		//http.StatusFound (302) and a Location header to the HTTP
		// response.
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	/*t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)*/
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body") // returns string
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

var templates = template.Must(template.ParseFiles("view.html", "edit.html"))

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request,
	string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r) // Writing "404 Not Found" error
			// to the HTTP connection.
			fmt.Println(errors.New("Invalid Page Title"))
			return
		}
		for _, val := range m {
			fmt.Println(val)
		}
		fn(w, r, m[2])
	}
}
