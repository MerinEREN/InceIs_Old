package main

import (
	"fmt"
	"github.com/MerinEREN/InceIs/account"
	"github.com/MerinEREN/InceIs/cookie"
	"github.com/MerinEREN/InceIs/page/content"
	"github.com/MerinEREN/InceIs/page/templates"
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"regexp"
)

const (
	// mongoDBLocal = "mongodb:// localhost"
	mongoDBLocal = "localhost"
)

var (
	// CHANGE THE REGEXP BELOW !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	validPath = regexp.MustCompile("^/|[/A-Za-z0-9]$")
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
	template.RenderIndex(w, p)
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
	if r.Method == "POST" {
		key := "email"
		email := r.PostFormValue(key)
		key = "password"
		password := r.PostFormValue(key)
		account.Create(w, r, c, email, password)
	}
	template.RenderSignUp(w, p)
}

func logInHandler(w http.ResponseWriter, r *http.Request, c *mgo.Collection, s string) {
	p, err := content.Get(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Method == "POST" {
		key := "email"
		email := r.PostFormValue(key)
		key = "password"
		password := r.PostFormValue(key)
		// fmt.Fprintf(w, "Login Email: %s\n Login Password: %s\n", email, password)
		acc, err := account.VerifyUser(c, email, password)
		switch err {
		case account.EmailNotExist:
			fmt.Fprintln(w, err)
		case account.ExistingEmail:
			for _, u := range acc.Users {
				if u.Email == email {
					// ALLWAYS CREATE COOKIE BEFORE EXECUTING TEMPLATE
					cookie.Create(w, r, "session", u.UUID)
					// BE SURE THIS RETURNS ONLY OUTSIDE OF FOR THEN
					// ACTIVATE IT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
					// return
				}
			}
			// NEWER EXECUTE TEPLATE OR WRITE ANYTHING TO THE BODY BEFORE
			// REDIRECT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			http.Redirect(w, r, "/accounts/"+acc.Name, 302)
		case account.InvalidPassword:
			fmt.Fprintln(w, err)
		default:
			// Status code could be wrong
			http.Error(w, err.Error(), http.StatusNotImplemented)
			log.Fatalln(err)
		}
	}
	template.RenderLogIn(w, p)
}

func accountsHandler(w http.ResponseWriter, r *http.Request, c *mgo.Collection, s string) {
	p, err := content.Get(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	template.RenderAccounts(w, p)
}

func accountHandler(w http.ResponseWriter, r *http.Request, c *mgo.Collection, s string) {
	p, err := content.Get(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	template.RenderAccount(w, p)
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
			cookie.Create(w, r, m[1])
			fn(w, r, m[1])
		case 3:
			cookie.Create(w, r, m[2])
			fn(w, r, m[2])
		default:
			cookie.Create(w, r, "index")
			fn(w, r, "index")
		} */
		switch r.URL.Path {
		case "/":
			// cookie.Create(w, r, "index", nil)
			fn(w, r, c, "index")
		case "/signUp":
			// cookie.Create(w, r, "signUp", nil)
			fn(w, r, c, "signUp")
		case "/logIn":
			fn(w, r, c, "logIn")
		case "/accounts":
			// cookie.Create(w, r, "accounts", nil)
			fn(w, r, c, "accounts")
		default:
			fn(w, r, c, "account")
		}
	}
}

// HEADER ALWAYS SHOULD BE SET BEFORE ANYTHING WRITE A PAGE BODY !!!!!!!!!!!!!!!!!!!!!!!!!!
// w.Header().Set("Content-Type", "text/html"; charset=utf-8")
//fmt.Fprintln(w, things...) // Writes to the body
