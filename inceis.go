/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package inceis

import (
	"encoding/json"
	"fmt"
	"github.com/MerinEREN/iiPackages/account"
	"github.com/MerinEREN/iiPackages/cookie"
	"github.com/MerinEREN/iiPackages/page/content"
	usr "github.com/MerinEREN/iiPackages/user"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/user"
	// "io/ioutil"
	"github.com/MerinEREN/iiPackages/page/template"
	// "html/template"
	"log"
	// "mime/multipart"
	"net/http"
	"regexp"
	"time"
)

var _ memcache.Item // For debugging, delete when done.

var (
	// CHANGE THE REGEXP BELOW !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	validPath = regexp.MustCompile("^/|[/A-Za-z0-9]$")
)

// type LoginURLs map[string]string

func init() {
	// http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/",
		http.TimeoutHandler(http.HandlerFunc(makeHandlerFunc(indexHandler)),
			130*time.Millisecond,
			"This is http.TimeoutHandler(handler, time.Duration, message) "+
				"message bitch =)"))
	http.HandleFunc("/roles/", makeHandlerFunc(rolesHandler))
	http.HandleFunc("/userSettings/", makeHandlerFunc(userSettingsHandler))
	http.HandleFunc("/accountSettings/", makeHandlerFunc(accountSettingsHandler))
	// http.HandleFunc("/signUp", makeHandlerFunc(signUpHandler))
	// http.HandleFunc("/logIn", makeHandlerFunc(logInHandler))
	// http.HandleFunc("/accounts", makeHandlerFunc(accountsHandler))
	// http.HandleFunc("/accounts/", makeHandlerFunc(accountHandler))
	http.HandleFunc("/logOut/", makeHandlerFunc(logOutHandler))
	/* if http.PostForm("/logIn", data); err != nil {
		http.Err(w, "Internal server error while login",
			http.StatusBadRequest)
	} */
	fs := http.FileServer(http.Dir("../iiClient/src"))
	// http.Handle("/css/", fs)
	http.Handle("/img/", fs)
	http.Handle("/js/", fs)
	/* log.Printf("About to listen on 10443. " +
	"Go to https://192.168.1.100:10443/ " +
	"or https://localhost:10443/") */
	// Redirecting to a port or a domain etc.
	// go http.ListenAndServe(":8080",
	// http.RedirectHandler("https://192.168.1.100:10443", 301))
	// err := http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil)
	// ListenAndServe and ListenAndServeTLS always returns a non-nil error !!!
	// log.Fatal(err)
}

func indexHandler(w http.ResponseWriter, r *http.Request, s string) {
	// The "/" pattern matches everything, so we need to check
	// that we're at the root here.
	if r.URL.Path == "/favicon.ico" {
		return
	} else if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	ctx := appengine.NewContext(r)
	p, err := content.Get(ctx, s)
	if p == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err != nil {
		log.Printf("Error while getting %s page content. Error: %v\n", s,
			err)
	}
	u1 := user.Current(ctx)
	currentOAuth, _ := user.CurrentOAuth(ctx, "")
	isAdmin := user.IsAdmin(ctx)
	oAuthConsumerKey, err := user.OAuthConsumerKey(ctx)
	log.Println(u1, currentOAuth, isAdmin, oAuthConsumerKey, err)
	// Login or get data needed
	if u1 == nil {
		switch r.Method {
		case "POST":
			gURL, err := user.LoginURL(ctx, r.URL.String())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			loginURLs := make(map[string]string)
			loginURLs["Google"] = gURL
			loginURLs["LinkedIn"] = gURL
			loginURLs["Twitter"] = gURL
			loginURLs["Facebook"] = gURL
			// http.Error(w, "Google login URL =)",
			//http.StatusInternalServerError)
			// http.NotFound(w, r)
			// http.Redirect(w, r, "/MerinEREN", http.StatusFound)
			// log.Println(r.BasicAuth())
			// log.Println(r.RemoteAddr)
			// log.Println(r.URL)
			// log.Println(r.TLS)
			// log.Println(r.Close)
			// log.Println(r.RequestURI)
			// log.Println(http.ProxyFromEnvironment(r))
			// log.Println(r.Referer())
			// log.Println(r.Cookies())
			// log.Printf("Request method is %v\n", r.Method)
			// _ = r.Write(w)
			// w.Flush()
			// log.Println(w.Header())
			// w.Header().Set("Content-Type", "application/json")
			// contentLength, err := w.Write([]byte(fmt.Sprintf("%v", p.D)))
			// log.Println(contentLength, err)
			/* t := &http.Transport{}
			t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
			c := &http.Client{Transport: t}
			res, err := c.Get("file:///etc/passwd")
			log.Println(res, err) */
			// When u get request with Body
			// defer r.Body.Close()
			// to read the request Body
			// b, err := ioutil.ReadAll(r.Body)
			// if err != nil {
			// http.Error(w, "Error reading body", http.StatusBadRequest)
			// return
			// }
			// To respond to request without any data
			// w.WriteHeader(StatusOK)
			b, err := json.Marshal(loginURLs)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write(b)
		default:
			template.RenderIndex(w, p)
		}
	} else {
		switch r.Method {
		case "POST":
			acc := new(account.Account)
			u2, uKey, err := usr.Exist(ctx, u1.Email)
			switch err {
			case datastore.Done:
				acc, u2, err = account.Create(r)
				if err != nil {
					log.Printf("Error while creating account: %v\n",
						err)
					// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
					http.Error(w, err.Error(),
						http.StatusInternalServerError)
					return
				}
			case usr.ExistingEmail:
				aKey := uKey.Parent()
				err = datastore.Get(ctx, aKey, acc)
				if err != nil {
					log.Printf("Error while getting user's account"+
						"data: %v\n", err)
					// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
					http.Error(w, err.Error(),
						http.StatusInternalServerError)
					return
				}
			case usr.FindUserError:
				log.Printf("Error while login user: %v\n", err)
				// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err = cookie.Set(w, r, "session", u2.UUID); err != nil {
				// CHECK FOR DISABLED COOKIE CLIENTS
				if _, err = r.Cookie(s); err == http.ErrNoCookie {
					p.D.URLUUID = "?uuid=" + u2.UUID
					// ALSO SET URL PATH WITH UUID !!!!!!!!!!!!!!!!!!!!
				}
				log.Printf("Error while creating session cookie: %v\n",
					err)
			}
			p.D.Account = acc
			p.D.User = u2
			p.D.LogoutURL, err = user.LogoutURL(ctx, "/")
			if err != nil {
				log.Println(err)
				// CHANGE http.Status.....Error !!!!!!!!!!!!!!!!!!!!!!!!!!!
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// log.Printf("Selected language by user is %s",
			// r.FormValue("lang"))
			b, err := json.Marshal(p.D)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(b)
		default:
			template.RenderIndex(w, p)
		}
	}
	/* temp := template.Must(template.New("fdsfdfdf").Parse(pBody))
	err = temp.Execute(w, p)
	if err != nil {
		log.Print(err)
	} */
	// THE IF CONTROL BELOW IS IMPORTANT
	// WHEN PAGE LOADS THERE IS NO FILE SELECTED AND THIS CAUSE A PROBLEM FOR
	/* if r.Method == "POST" {
		var f multipart.File
		key := "uploadedFile"
		f, _, err := r.FormFile(key)
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
	} */
	// w.Write([]byte(p.D))
}

func rolesHandler(w http.ResponseWriter, r *http.Request, s string) {
	ctx := appengine.NewContext(r)
	u1 := user.Current(ctx)
	log.Println(ctx, u1)
	if u1 == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		u2, uKey, err := usr.Exist(ctx, u1.Email)
		if err == usr.FindUserError {
			log.Printf("Error while login user: %v\n", err)
			// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if u2.Type == "inHouse" && u2.Status != "frozen" && (u2.IsAdmin() || u2.IsContentEditor()) {
			acc := new(account.Account)
			aKey := uKey.Parent()
			err = datastore.Get(ctx, aKey, acc)
			if err != nil {
				log.Printf("Error while getting user's account data: %v\n",
					err)
				// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!
				http.Error(w, err.Error(),
					http.StatusInternalServerError)
				return
			}
			p := new(content.Page)
			p, err = content.Get(ctx, s)
			if p == nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err != nil {
				log.Printf("Error while getting page content. Error: %v\n", err)
			}
			p.D.Account = acc
			p.D.User = u2
			// template.RenderRoles(w, p)
			// keyValue := *role
			// log.Println(keyValue)
			// log.Println(role.StringID())
			// log.Println(role.IntID())
			// log.Println(role.Parent())
			// log.Println(role.AppID())
			// log.Println(role.Kind())
			// log.Println(role.Namespace())
		} else {
			log.Printf("Unauthorized user %s trying to see "+
				"roles page !!!", u2.Email)
			fmt.Fprintf(w, "Permission denied !!!")
			return
		}
	}
}

func userSettingsHandler(w http.ResponseWriter, r *http.Request, s string) {
	ctx := appengine.NewContext(r)
	u1 := user.Current(ctx)
	if u1 == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		u2, uKey, err := usr.Exist(ctx, u1.Email)
		if err == usr.FindUserError {
			log.Printf("Error while login user: %v\n", err)
			// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if u2.Status == "frozen" {
			log.Printf("Unauthorized user %s trying to see "+
				"user settings page !!!", u2.Email)
			fmt.Fprintf(w, "Permission denied !!!")
			return
		}
		acc := new(account.Account)
		aKey := uKey.Parent()
		err = datastore.Get(ctx, aKey, acc)
		if err != nil {
			log.Printf("Error while getting user's account data: %v\n",
				err)
			// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}
		p := new(content.Page)
		p, err = content.Get(ctx, s)
		if p == nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err != nil {
			log.Printf("Error while getting page content. Error: %v\n", err)
		}
		p.D.Account = acc
		p.D.User = u2
		// template.RenderUserSettings(w, p)
	}
}

func accountSettingsHandler(w http.ResponseWriter, r *http.Request, s string) {
	ctx := appengine.NewContext(r)
	u1 := user.Current(ctx)
	if u1 == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		u2, uKey, err := usr.Exist(ctx, u1.Email)
		if err == usr.FindUserError {
			log.Printf("Error while login user: %v\n", err)
			// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if u2.Status == "frozen" || !u2.IsAdmin() {
			log.Printf("Unauthorized user %s trying to see "+
				"account settings page !!!", u2.Email)
			fmt.Fprintf(w, "Permission denied !!!")
			return
		}
		acc := new(account.Account)
		aKey := uKey.Parent()
		err = datastore.Get(ctx, aKey, acc)
		if err != nil {
			log.Printf("Error while getting user's account data: %v\n",
				err)
			// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}
		p := new(content.Page)
		p, err = content.Get(ctx, s)
		if p == nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err != nil {
			log.Printf("Error while getting page content. Error: %v\n", err)
		}
		p.D.Account = acc
		p.D.User = u2
		// template.RenderAccountSettings(w, p)
	}
}

/* func signUpHandler(w http.ResponseWriter, r *http.Request, s string) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is Sign Up page " + s + "\n"))
	p, err := content.Get(r, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Method == "POST" {
		acc, UUID, err := account.Create(r)
		if err != nil {
			log.Printf("Error while creating account: %v\n", err)
			// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cookie.Set(w, r, "session", UUID)
		http.Redirect(w, r, "/accounts/"+acc.Name, 302)
	}
	template.RenderSignUp(w, p)
} */

/* func logInHandler(w http.ResponseWriter, r *http.Request, s string) {
	p, err := content.Get(r, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Method == "POST" {
		key := "email"
		email := r.PostFormValue(key)
		key = "password"
		password := r.PostFormValue(key)
		acc, err := account.VerifyUser(c, email, password)
		switch err {
		case account.EmailNotExist:
			fmt.Fprintln(w, err)
		case account.ExistingEmail:
			for _, u := range acc.Users {
				if u.Email == email {
					// ALLWAYS CREATE COOKIE BEFORE EXECUTING TEMPLATE
					cookie.Set(w, r, "session", u.UUID)
				}
			}
			// NEWER EXECUTE TEMPLATE OR WRITE ANYTHING TO THE BODY BEFORE
			// REDIRECT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			http.Redirect(w, r, "/accounts/"+acc.Name, http.StatusSeeOther)
		case account.InvalidPassword:
			fmt.Fprintln(w, err)
		default:
			// Status code could be wrong
			http.Error(w, err.Error(), http.StatusNotImplemented)
			log.Fatalln(err)
		}
	}
	template.RenderLogIn(w, p)
} */

/* func accountsHandler(w http.ResponseWriter, r *http.Request, s string) {
	p, err := content.Get(r, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	template.RenderAccounts(w, p)
} */

/* func accountHandler(w http.ResponseWriter, r *http.Request, s string) {
	// w.Header.Set("Location", url)
	// w.WriteHeader(http.StatusFound)
	template.RenderAccount(w, p)
} */

func logOutHandler(w http.ResponseWriter, r *http.Request, s string) {
	ctx := appengine.NewContext(r)
	url, err := user.LogoutURL(ctx, "/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Deleting session cookie
	var cookie *http.Cookie
	cookie, err = r.Cookie(s)
	if err != http.ErrNoCookie {
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
	}
	//  CHANGE NECESSARY DB FIELDS OF USER !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	http.Redirect(w, r, url, http.StatusFound)
}

type handlerFuncWithDomainArg func(http.ResponseWriter, *http.Request, string)

func makeHandlerFunc(fn handlerFuncWithDomainArg) http.HandlerFunc {
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
			cookie.Set(w, r, m[1])
			fn(w, r, m[1])
		case 3:
			cookie.Set(w, r, m[2])
			fn(w, r, m[2])
		default:
			cookie.Set(w, r, "index")
			fn(w, r, "index")
		} */
		switch r.URL.Path {
		case "/":
			// cookie.Set(w, r, "index", nil)
			fn(w, r, "index")
		case "/roles":
			// cookie.Set(w, r, "roles", nil)
			fn(w, r, "roles")
		case "/userSettings":
			fn(w, r, "userSettings")
		case "/accountSettings":
			fn(w, r, "accountSettings")
		case "/logOut":
			fn(w, r, "session")
		case "/accounts":
			// If session cookie not exists, redirect to the index page
			/* if !cookie.IsExists(r, "session") {
				// HTTP StatusCode could be wrong
				http.Redirect(w, r, "/", http.StatusSeeOther)
			} */
			// cookie.Set(w, r, "accounts", nil)
			fn(w, r, "accounts")
		default:
			// If session cookie not exists, redirect to the index page
			/* if !cookie.IsExists(r, "session") {
				// HTTP StatusCode could be wrong
				http.Redirect(w, r, "/", http.StatusSeeOther)
			} */
			fn(w, r, "account")
		}
	}
}

// HEADER ALWAYS SHOULD BE SET BEFORE ANYTHING WRITE A PAGE BODY !!!!!!!!!!!!!!!!!!!!!!!!!!
// w.Header().Set("Content-Type", "text/html"; charset=utf-8")
//fmt.Fprintln(w, things...) // Writes to the body
