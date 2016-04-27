package main

import (
	// "errors"
	"fmt"
	valid "github.com/asaskevich/govalidator"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	// "strings"
	"time"
)

const (
	// mongoDBLocal = "mongodb:// localhost"
	mongoDBLocal = "localhost"
)

var (
	templates = template.Must(template.ParseGlob("../templates/*.gohtml"))
	index     = renderTemplate("index")
	logIn     = renderTemplate("logIn")
	signUp    = renderTemplate("signUp")

	// CHANGE THE REGEXP BELOW !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	validPath = regexp.MustCompile("^/(|logIn|signUp)$")
)

// OO Language and Page Sturcts
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

type Page struct {
	Title string `bson:"title"`
	Body  []byte `bson:"body"`
}

// USE THIS TO GET PAGE CONTENTS FOR EVERY LANG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func loadPage(title string) (*Page, error) {
	filename := title + ".gohtml"
	body, err := ioutil.ReadFile("../templates/" + filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func main() {
	session, err := mgo.Dial(mongoDBLocal)
	if err != nil {
		fmt.Println(session, err)
		panic(err)
	}
	defer session.Close()
	c := session.DB("OO").C("accounts")
	// m := bson.M{}

	// Us this as a http.HandlerFunc
	signUpTemp(c)

	mux := http.NewServeMux()
	mux.HandleFunc("/", makeHandler(indexHandler))
	mux.HandleFunc("/logIn", makeHandler(logInHandler))
	mux.HandleFunc("/signUp", makeHandler(signUpHandler))
	/* if mux.PostForm("/logIn", data); err != nil {
		http.Err(w, "Internal server error while login",
			http.StatusBadRequest)
	} */
	log.Printf("About to listen on 10443. " +
		"Go to https://192.168.1.100:10443/ " +
		"or https://localhost:10443/")
	// Redirecting to a port or a domain etc.
	go http.ListenAndServe(":8080",
		http.RedirectHandler("https://192.168.1.100:10443", 301))
	err = http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", mux)
	// ListenAndServeTLS always returns a non-nil error !!!
	log.Fatal(err)
}

// OO DB Structs
// Accounts collection
type Accounts []Account

type Account struct {
	Name          int       `bson:"name,omitempty"`
	Type          string    `bson:"type,omitempty"`
	CurrentStatus string    `bson:"current_status,omitempty"`
	AccountStatus string    `bson:"account_status,omitempty"`
	About         string    `bson:"about,omitempty"`
	Tags          Tags      `bson:"tags,omitempty"`
	Ranks         Ranks     `bson:"ranks,omitempty"`
	Card          Card      `bson:"card,omitempty" valid:"creditcard"`
	Users         Users     `bson:"users,omitempty"`
	Registered    time.Time `bson:"registered,omitempty"`
	LastModified  time.Time `bson:"last_modified,omitempty"`
}

type Company struct {
	Name    string  `bson:"name,omitempty"`
	Address Address `bson:"address,omitempty"`
}

type Address struct {
	Description string      `bson:"description,omitempty"`
	Borough     string      `bson:"borough,omitempty"`
	City        string      `bson:"city,omitempty"`
	Country     string      `bson:"country,omitempty"`
	Postcode    string      `bson:"postcode,omitempty"`
	Geolocation Geolocation `bson:"geolocation,omitempty"`
}

type Geolocation struct {
	Lat  string `bson:"lat,omitempty"`  // type could be differnt !!!
	Long string `bson:"Long,omitempty"` // type could be differnt !!!
}

type Tags []Tag

type Tag struct {
	Type string `bson:"type,omitempty"`
}

type Ranks []Rank

type Rank struct {
	Type string `bson:"type,omitempty"`
}

type Card struct {
	Creditcards Creditcards `bson:"creditcards,omitempty"`
	Debitcards  Debitcards  `bson:"debitcards,omitempty"`
}

type Creditcards []Creditcard

type Creditcard struct {
	HolderName string `bson:"holder_name,omitempty"`
	No         string `bson:"no,omitempty"`
	ExpMonth   string `bson:"exp_month,omitempty"`
	ExpYear    string `bson:"exp_year,omitempty"`
	CVV        string `bson:"cvv,omitempty"`
}

type Debitcards []Debitcard

type Debitcard struct {
	HolderName string `bson:"holder_name,omitempty"`
	No         string `bson:"no,omitempty"`
	ExpMonth   string `bson:"exp_month,omitempty"`
	ExpYear    string `bson:"exp_year,omitempty"`
	CVV        string `bson:"cvv,omitempty"`
}

type Users []User

type User struct {
	Email       string `bson:"email,omitempty"`
	Password    string `bson:"password,omitempty"`
	PicturePath string `bson:"picture_path,omitempty"`
	Name        Name   `bson:"name,omitempty"`
	Phone       string `bson:"phone,omitempty"` // Should be struct in
	// the future !!!
	Status       string       `bson:"status,omitempty"`
	Type         string       `bson:"type,omitempty"`
	BirthDate    time.Time    `bson:"birth_date,omitempty"`
	Registered   time.Time    `bson:"registered,omitempty"`
	LastModified time.Time    `bson:"last_modified,omitempty"`
	IsActive     bool         `bson:"is_active,omitempty"`
	ServicePacks ServicePacks `bson:"service_packs",omitempty"`
	// 	PurchasedServices PurchasedServices `bson:"purchasedServices,
	// 	omitempty"`
}

type Name struct {
	First string `bson:"first,omitempty"`
	Last  string `bson:"last,omitempty"`
}

type ServicePacks []ServicePack

type ServicePack struct {
	Id             string            `bson:"id,omitempty"`
	Type           string            `bson:"type,omitempty"`
	Description    string            `bson:"description,omitempty"`
	Duration       string            `bson:"duration,omitempty"`
	Price          Price             `bson:"price,omitempty"`
	Extras         ServicePackExtras `bson:"extras,omitempty"`
	Photos         Photos            `bson:"photos,omitempty"`
	Videos         Videos            `bson:"videos,omitempty"`
	Tags           Tags              `bson:"tags,omitempty"`
	Created        time.Time         `bson:"created,omitempty"`
	LastModified   time.Time         `bson:"last_modified,omitempty"`
	Status         string            `bson:"status,omitempty"`
	Evaluation     Evaluation        `bson:"evaluation,omitempty"`
	CustomerReview string            `bson:"customer_review,omitempty"`
}

type Price struct {
	Amount   float64 `bson:amount,omitempty"`
	Currency string  `bson:currency,omitempty"`
}

type ServicePackExtras []ServicePackOption

type ServicePackOption struct {
	Id          string `bson:"id,omitempty"`
	Description string `bson:"description,omitempty"`
	Duration    string `bson:"duration,omitempty"`
	Price       Price  `bson:"price,omitempty"`
	Photos      Photos `bson:"photos,omitempty"`
	Videos      Videos `bson:"videos,omitempty"`
}

type Photos []Photo

type Photo struct {
	Id           string    `bson:"id,omitempty"`
	Path         string    `bson:"path,omitempty"`
	Title        string    `bson:"title,omitempty"`
	Description  string    `bson:"description,omitempty"`
	Uploaded     time.Time `bson:"uploaded,omitempty"`
	LastModified time.Time `bson:"last_modified,omitempty"`
	Status       string    `bson:"status,omitempty"`
}

type Videos []Video

type Video struct {
	Id           string    `bson:"id,omitempty"`
	Path         string    `bson:"path,omitempty"`
	Title        string    `bson:"title,omitempty"`
	Description  string    `bson:"description,omitempty"`
	Uploaded     time.Time `bson:"uploaded,omitempty"`
	LastModified time.Time `bson:"last_modified,omitempty"`
	Status       string    `bson:"status,omitempty"`
}

type Evaluation struct {
	Technical     byte `bson:"technical,omitempty"`
	Communication byte `bson:"communication,omitempty"`
	Time          byte `bson:"time,omitempty"`
}

type doc interface {
	// Use this for all structs
	// Update()
	// Upsert()
	// Delete()
}

// Make a makeAjaxHandler And Wrap http.HandlerFunc like signUp()
// Sign Up Handler
func signUpTemp(c *mgo.Collection) {
	// Get parameters from HTTP Request !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	email := "ivieren@gmail.com"
	if !valid.IsEmail(email) {
		// Send response !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		fmt.Println("Invalid email")
		return
	}
	password := "0018@8100"
	// Cahange this control and allow special characters !!!!!!!!!!!!!!!!!!
	if !valid.IsAlphanumeric(password) {
		// Send response !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		fmt.Println("Invalid password")
		return
	}
	var accountCount int
	accountCount, err := c.Find(bson.M{"users.email": email}).Count()
	if err != nil {
		// Send response !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		fmt.Printf("Insert account error, because can't verify email. Error: %v\n",
		err)
		return
	}
	if accountCount != 0 {
		// Send response !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		fmt.Printf("Email exists !!!, Count: %v\n", accountCount)
		return
	}
	users := Users{
		User{
			Email:        email,
			Password:     password,
			Status:       "online",
			Type:         "admin",
			IsActive:     true,
			Registered:   time.Now(),
			LastModified: time.Now(),
		},
	}
	// Get the account count and use as name !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	acc := Account{
		Name:          1,
		CurrentStatus: "available",
		AccountStatus: "online",
		Users:         users,
		Registered:    time.Now(),
		LastModified:  time.Now(),
	}
	err = c.Insert(acc)
	if err != nil {
		// Send response !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		fmt.Println("Insert account error")
		return
	}
}

func renderTemplate(title string) func(w http.ResponseWriter, p *Page) {
	return func(w http.ResponseWriter, p *Page) {
		err := templates.ExecuteTemplate(w, title+".gohtml", p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			// maybe log.Fatal(err) or http.Redirect(w, r, "/logIn/", 
			// http.StatusFound)
		}
	}
}

// MAKE "index" AND "logIn" PAGE AS ONE PAGE !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func indexHandler(w http.ResponseWriter, r *http.Request, s string) {
	// w.Header().Set("Content-Type", "text/plain")
	// w.Write([]byte("This is Main page " + s + "\n"))
	fmt.Printf("index s is %s", s)
	s = "index"
	p, err := loadPage(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	index(w, p)
}

func logInHandler(w http.ResponseWriter, r *http.Request, s string) {
	// w.Header().Set("Content-Type", "text/plain")
	// w.Write([]byte("This is Log In page " + s + "\n"))
	fmt.Printf("logIn s is %s", s)
	p, err := loadPage(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logIn(w, p)
}

func signUpHandler(w http.ResponseWriter, r *http.Request, s string) {
	// w.Header().Set("Content-Type", "text/plain")
	// w.Write([]byte("This is Sign Up page " + s + "\n"))
	fmt.Printf("signUp s is %s", s)
	p, err := loadPage(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	signUp(w, p)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request,
	string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			log.Printf("Invalid Path: %s\n", r.URL.Path)
			// Writing "404 Not Found" error to the HTTP connection
			http.NotFound(w, r)
			return
		}
		for _, val := range m {
			fmt.Println(val)
		}
		switch len(m) {
		case 2:
			fn(w, r, m[1])
		case 3:
			fn(w, r, m[2])
		default:
			fn(w, r, "index")
		}
	}
	/* urlSlice := strings.Split(r.URL.Path, "/")
	for _, v := range urlSlice {
		fmt.Println(v)
	}
	switch len(urlSlice) {
	case 2:
		fn(w, r, urlSlice[1])
	case 3:
		fn(w, r, urlSlice[2])
	default:
		fn(w, r, "index")
	} */
}
