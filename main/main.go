package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MerinEREN/InceIs/page/content"
	valid "github.com/asaskevich/govalidator"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// mongoDBLocal = "mongodb:// localhost"
	mongoDBLocal = "localhost"
)

var (
	templates = template.Must(template.ParseGlob("../page/templates/*.html"))
	// CHANGE THE REGEXP BELOW !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	validPath       = regexp.MustCompile("^/|[/A-Za-z0-9]$")
	renderIndex     = renderTemplate("index")
	renderSignUp    = renderTemplate("signUp")
	renderLogIn     = renderTemplate("logIn")
	renderAccounts  = renderTemplate("accounts")
	renderAccount   = renderTemplate("account")
	EmailNotExist   = errors.New("Invalid Email")
	ExistingEmail   = errors.New("Existing Email")
	InvalidPassword = errors.New("Invalid Password")
)

func main() {
	session, err := mgo.Dial(mongoDBLocal)
	if err != nil {
		fmt.Println(session, err)
		panic(err)
	}
	defer session.Close()
	c := session.DB("II").C("accounts")
	// m := bson.M{}
	mux := http.NewServeMux()
	mux.Handle("/favicon.ico", http.NotFoundHandler())
	mux.HandleFunc("/", makeHandler(c, indexHandler))
	mux.HandleFunc("/signUp", makeHandler(c, signUpHandler))
	mux.HandleFunc("/logIn", makeHandler(c, logInHandler))
	mux.HandleFunc("/accounts", makeHandler(c, accountsHandler))
	mux.HandleFunc("/accounts/", makeHandler(c, accountHandler))
	/* if http.PostForm("/logIn", data); err != nil {
		http.Err(w, "Internal server error while login",
			http.StatusBadRequest)
	} */
	fs := http.FileServer(http.Dir("../public"))
	mux.Handle("/css/", fs)
	mux.Handle("/img/", fs)
	mux.Handle("/js/", fs)
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

// II DB Structs
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
	UUID        string `bson:"uuid"`
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

func indexHandler(w http.ResponseWriter, r *http.Request, c *mgo.Collection, s string) {
	// HANDLE FOR /favicon.ico REQUEST !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	/* if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	} */
	p, err := content.Get(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderIndex(w, p)
	// THE IF CONTROL BELOW IS IMPORTANT
	// WHEN PAGE LOADS THERE IS NO FILE SELECTED AND THIS CAUSE A PROBLEM FOR
	// r.FormFile(key)
	if r.Method == "POST" {
		var f multipart.File
		key := "uploadedFile"
		f, _, err = r.FormFile(key)
		if err != nil {
			fmt.Println("File input is empty.")
			return
		}
		defer f.Close()
		var bs []byte
		bs, err = ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "File: %s\n Error: %v\n", string(bs), err)
	}
}

func signUpHandler(w http.ResponseWriter, r *http.Request, c *mgo.Collection, s string) {
	// w.Header().Set("Content-Type", "text/plain")
	// w.Write([]byte("This is Sign Up page " + s + "\n"))
	p, err := content.Get(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderSignUp(w, p)
	if r.Method == "POST" {
		key := "email"
		email := r.PostFormValue(key)
		key = "password"
		password := r.PostFormValue(key)
		fmt.Fprintf(w, "Sign Up Email: %s\n Sign Up Password: %s\n", email,
			password)
		createAccount(w, r, c, email, password)
	}
}

func logInHandler(w http.ResponseWriter, r *http.Request, c *mgo.Collection, s string) {
	p, err := content.Get(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderLogIn(w, p)
	if r.Method == "POST" {
		key := "email"
		email := r.PostFormValue(key)
		key = "password"
		password := r.PostFormValue(key)
		fmt.Fprintf(w, "Login Email: %s\n Login Password: %s\n", email, password)
		acc, err := verifyUser(c, email, password)
		switch err {
		case EmailNotExist:
			fmt.Fprintln(w, err)
			return
		case ExistingEmail:
			// GET ONLY NECESSARY DATAS VIA PROJECTION !!!!!!!!!!!!!!!!!!!!!!!!
			for _, user := range acc.Users {
				if user.Email == email {
					// ALLWAYS CREATE COOKIE BEFORE EXECUTING TEMPLATE
					defer createCookie(w, r, "session", user.UUID)
					// BE SURE THIS RETURNS ONLY OUTSIDE OF FOR THEN
					// ACTIVATE IT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
					// return
				}
			}
			http.Redirect(w, r, "https://192.168.1.100:1043/accounts/"+
				strconv.Itoa(acc.Name), 302)
		case InvalidPassword:
			fmt.Fprintln(w, err)
			return
		default:
			// Inform client
			log.Println(err)
			// status code could be wrong
			http.Redirect(w, r, "https://192.168.1.100:10443", 307)
		}
	}

}

func accountsHandler(w http.ResponseWriter, r *http.Request, c *mgo.Collection, s string) {
	p, err := content.Get(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderAccounts(w, p)
}

func accountHandler(w http.ResponseWriter, r *http.Request, c *mgo.Collection, s string) {
	p, err := content.Get(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderAccount(w, p)
}

// MAKE ITS OWN PACKAGE
func renderTemplate(title string) func(w http.ResponseWriter, p *content.Page) {
	return func(w http.ResponseWriter, p *content.Page) {
		err := templates.ExecuteTemplate(w, title+".html", p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			// maybe log.Fatal(err) or http.Redirect(w, r, "/logIn/",
			// http.StatusFound)
		}
	}
}

func createAccount(w http.ResponseWriter, r *http.Request, c *mgo.Collection,
	email, password string) {
	if !valid.IsEmail(email) {
		// Inform client
		fmt.Fprintln(w, "Invalid email")
		return
	}
	// Cahange this control and allow special characters !!!!!!!!!!!!!!!!!!
	if !valid.IsAlphanumeric(password) {
		// Inform client
		fmt.Fprintln(w, "Invalid password")
		return
	}
	acc, err := verifyUser(c, email, password)
	switch err {
	case EmailNotExist:
		u4, errUUID := uuid.NewV4()
		if errUUID != nil {
			log.Printf("Can't create UUID when signUp, error: %v\n",
				errUUID)
			return
		}
		users := Users{
			User{
				UUID:         u4.String(),
				Email:        email,
				Password:     password,
				Status:       "online",
				Type:         "admin",
				IsActive:     true,
				Registered:   time.Now(),
				LastModified: time.Now(),
			},
		}
		// Get the account count and use as name !!!!!!!!!!!!!!!!!!!!!!!!!!
		acc = &Account{
			Name:          1,
			CurrentStatus: "available",
			AccountStatus: "online",
			Users:         users,
			Registered:    time.Now(),
			LastModified:  time.Now(),
		}
		errInsert := c.Insert(acc)
		if errInsert != nil {
			// Inform client
			log.Println(errInsert)
			http.Redirect(w, r, "/", 304) // status code could be wrong
		}
		// ALLWAYS CREATE COOKIE BEFORE EXECUTING TEMPLATE
		createCookie(w, r, "session", u4.String())
		http.Redirect(w, r, "https://192.168.1.100:10443/accounts/"+
			strconv.Itoa(acc.Name), 302)
	case ExistingEmail:
		fmt.Fprintln(w, err)
	case InvalidPassword:
		fmt.Fprintln(w, ExistingEmail)
	default:
		// Inform client
		log.Println(err)
		// status code could be wrong
		http.Redirect(w, r, "https://192.168.1.100:10443", 307)
	}
}

// USE PROJECTION AND GET USER ONLY, NOT ACCOUNT
func verifyUser(c *mgo.Collection, e, p string) (result *Account, err error) {
	err = c.Find(bson.M{"users.email": e}).One(&result)
	if err != nil {
		if err.Error() == "not found" {
			err = EmailNotExist
		} else {
			return
		}
	} else {
		err = ExistingEmail
		for _, user := range result.Users {
			if user.Password != p {
				err = InvalidPassword
			}
		}
	}
	return
}

// ADDING UUID, NECESSARY USER DATA, COUNT AND HASH TO THE COOKIE AND CHECK HASH CODE
func createCookie(w http.ResponseWriter, r *http.Request, s string, uuid string) {
	// REMOVE THIS DUMMY COOKIE STRUCT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	type cookieData struct {
		Name  string
		Value string
	}
	cookie, err := r.Cookie(s)
	if err == http.ErrNoCookie {
		cookie = &http.Cookie{
			Name: s,
			// U CAN USE UUID AS VALUE !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			Value: uuid,
			// NOT GOOD PRACTICE
			// ADDING USER DATA TO A COOKIE
			// WITH NO WAY OF KNOWING WHETER OR NOT THEY MIGHT HAVE ALTERED
			// THAT DATA !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			// HMAC WOULD ALLOW US TO DETERMINE WHETHER OR NOT THE DATA IN THE
			// COOKIE WAS ALTERED !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			// HOWEVER, BEST TO STORE USER DATA ON THE SERVER AND KEEP
			// BACKUPS !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			// Value: "emil = merin@inceis.net" + "JSON data" + "whatever",
			// IF SECURE IS TRUE THIS COOKIE ONLY SEND WITH HTTP2 !!!!!!!!!!!!!
			Secure: true,
			// HttpOnly: true MEANS JAVASCRIPT CAN NOT ACCESS THE COOKIE !!!!!!
			HttpOnly: false,
		}
		cd := cookieData{
			Name:  "cookiedataName",
			Value: "cookiedataValue",
		}
		var bs []byte
		bs, err = json.Marshal(cd)
		log.Printf("Marshalled cookie data is %s\n", string(bs))
		if err != nil {
			log.Printf("%s cookie marshaling error. %v\n", cookie.Name, err)
		}
		cookie.Value += "|" + base64.StdEncoding.EncodeToString(bs) + "|0"
		code := getCode(cookie.Value)
		cookie.Value += "|" + code
		fmt.Printf("Cookie value for "+cookie.Name+" is: %s\n", cookie.Value)
	}
	xs := strings.Split(cookie.Value, "|")
	uuidDataCount := xs[0] + "|" + xs[1] + "|" + xs[2]
	returnedCode := getCode(uuidDataCount)
	if returnedCode != xs[3] {
		log.Printf("%s cookie value of uuid %v is corrupted. Cookie HMAC is %s, "+
			"genereted HMAC is %s", cookie.Name, uuid, xs[3], returnedCode)
		var decodedBase64 []byte
		decodedBase64, err = base64.StdEncoding.DecodeString(xs[1])
		if err != nil {
			log.Printf("Error while decoding %s cookie data. Error "+
				"is %v\n", cookie.Name, err)
		}
		var returnedCookieData cookieData
		err = json.Unmarshal(decodedBase64, &returnedCookieData)
		if err != nil {
			log.Printf("%s cookie unmarshaling error. %v\n", cookie.Name, err)
		}
		log.Printf("Returned cookie data is %v", returnedCookieData)
		// DID NOT CHECKED DELETING AND CREATING NEW COOKIE YET, SO CHECK THEM !!!!
		// DELETING CORRUPTED COOKIE AND CREATING NEW ONE !!!!!!!!!!!!!!!!!!!!!!!!!
		cookie.MaxAge = -1 // Deleting cookie
		createCookie(w, r, s, uuid)
		return
	}
	count, err := strconv.Atoi(xs[2])
	if err != nil {
		fmt.Println("Can't convert " + s + " cookie Atoi!")
	}
	count++
	cookie.Value = xs[0] + "|" + xs[1] + "|" + strconv.Itoa(count) + "|" + xs[3]
	fmt.Printf(cookie.Name+" page vizitor count: %s\n", strconv.Itoa(count))
	// CREATING A COOKIE IS NOT ENOUGH, YOU HAVE TO SET THE COOKIE TO USE IT !!!!!!!!!!
	http.SetCookie(w, cookie)
}

// Checking data with "hmac"
func getCode(s string) string {
	h := hmac.New(sha256.New, []byte("someKey"))
	io.WriteString(h, s)
	return fmt.Sprintf("%v", h.Sum(nil))
}

func makeHandler(c *mgo.Collection, fn func(http.ResponseWriter, *http.Request,
	*mgo.Collection, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			log.Printf("Invalid Path: %s\n", r.URL.Path)
			// Writing "404 Not Found" error to the HTTP connection
			http.NotFound(w, r)
			return
		}
		/* for _, val := range m {
			fmt.Println(val)
		}
		switch len(m) {
		case 2:
			createCookie(w, r, m[1])
			fn(w, r, m[1])
		case 3:
			createCookie(w, r, m[2])
			fn(w, r, m[2])
		default:
			createCookie(w, r, "index")
			fn(w, r, "index")
		} */
		switch r.URL.Path {
		case "/":
			// createCookie(w, r, "index", nil)
			fn(w, r, c, "index")
		case "/signUp":
			// createCookie(w, r, "signUp", nil)
			fn(w, r, c, "signUp")
		case "/logIn":
			fn(w, r, c, "logIn")
		case "/accounts":
			// createCookie(w, r, "accounts", nil)
			fn(w, r, c, "accounts")
		default:
			fn(w, r, c, "account")
		}
	}
}

// HEADER ALWAYS SHOULD BE SET BEFORE ANYTHING WRITE A PAGE !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// w.Header().Set("Content-Type", "text/html"; charset=utf-8")
//fmt.Fprintln(w, things...)
