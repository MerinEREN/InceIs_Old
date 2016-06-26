package inceis

import (
	// "fmt"
	"github.com/MerinEREN/iiPackages/account"
	"github.com/MerinEREN/iiPackages/cookie"
	"github.com/MerinEREN/iiPackages/page/content"
	usr "github.com/MerinEREN/iiPackages/user"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	// "google.golang.org/appengine/memcache"
	"google.golang.org/appengine/user"
	// "io/ioutil"
	"github.com/MerinEREN/iiPackages/page/template"
	// "html/template"
	"log"
	// "mime/multipart"
	"net/http"
	"regexp"
)

var (
	// CHANGE THE REGEXP BELOW !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	validPath = regexp.MustCompile("^/|[/A-Za-z0-9]$")
)

func init() {
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/", makeHandler(indexHandler))
	http.HandleFunc("/roles", makeHandler(rolesHandler))
	// http.HandleFunc("/signUp", makeHandler(signUpHandler))
	// http.HandleFunc("/logIn", makeHandler(logInHandler))
	// http.HandleFunc("/accounts", makeHandler(accountsHandler))
	// http.HandleFunc("/accounts/", makeHandler(accountHandler))
	http.HandleFunc("/logOut", makeHandler(logOutHandler))
	/* if http.PostForm("/logIn", data); err != nil {
		http.Err(w, "Internal server error while login",
			http.StatusBadRequest)
	} */
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/css/", fs)
	http.Handle("/img/", fs)
	http.Handle("/js/", fs)
	/* log.Printf("About to listen on 10443. " +
	"Go to https://192.168.1.100:10443/ " +
	"or https://localhost:10443/") */
	// Redirecting to a port or a domain etc.
	// go http.ListenAndServe(":8080",
	// http.RedirectHandler("https://192.168.1.100:10443", 301))
	// err := http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil)
	// ListenAndServeTLS always returns a non-nil error !!!
	// log.Fatal(err)
}

func indexHandler(w http.ResponseWriter, r *http.Request, s string) {
	// HANDLE FOR /favicon.ico REQUEST !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	/* if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	} */
	ctx := appengine.NewContext(r)
	u1 := user.Current(ctx)
	p := new(content.Page)
	if u1 == nil {
		url, err := user.LoginURL(ctx, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		p, err = content.Get(ctx, "index")
		if p == nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err != nil {
			log.Printf("Error while getting index page content. Error: %v\n",
				err)
		}
		p.D.LoginURL = url
		template.RenderIndex(w, p)
	} else {
		acc := new(account.Account)
		var errAc error
		u2, uKey, err := usr.Exist(ctx, u1.Email)
		switch err {
		case datastore.Done:
			acc, u2, errAc = account.Create(r)
			if errAc != nil {
				log.Printf("Error while creating account: %v\n", errAc)
				// ALSO LOG THIS WHITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!
				http.Error(w, errAc.Error(),
					http.StatusInternalServerError)
				return
			}
		case usr.ExistingEmail:
			aKey := uKey.Parent()
			// log.Println(uKey, aKey, acc)
			errAc = datastore.Get(ctx, aKey, acc)
			if errAc != nil {
				log.Printf("Error while getting user's account data: %v\n",
					errAc)
				// ALSO LOG THIS WHITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!
				http.Error(w, errAc.Error(),
					http.StatusInternalServerError)
				return
			}
		case usr.FindUserError:
			log.Printf("Error while login user: %v\n", err)
			// ALSO LOG THIS WHITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		p, err = content.Get(ctx, "account")
		if p == nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err != nil {
			log.Printf("Error while getting account page content. Error: %v\n",
				err)
		}
		if err = cookie.Set(w, r, "session", u2.UUID); err != nil {
			// CHECK FOR DISABLED COOKIE CLIENTS
			if _, err = r.Cookie(s); err == http.ErrNoCookie {
				p.D.URLUUID = "?uuid=" + u2.UUID
				// ALSO SET URL PATH WITH UUID !!!!!!!!!!!!!!!!!!!!!!!!!!!!
			}
			log.Printf("Error while creating session cookie: %v\n", err)
		}
		p.D.Account = acc
		p.D.User = u2
		template.RenderAccount(w, p)
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
}

func rolesHandler(w http.ResponseWriter, r *http.Request, s string) {
	ctx := appengine.NewContext(r)
	u1 := user.Current(ctx)
	if u1 == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		u2, uKey, err := usr.Exist(ctx, u1.Email)
		if err == usr.FindUserError {
			log.Printf("Error while login user: %v\n", err)
			// ALSO LOG THIS WHITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// CHECK USER ROLES !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		acc := new(account.Account)
		aKey := uKey.Parent()
		err = datastore.Get(ctx, aKey, acc)
		if err != nil {
			log.Printf("Error while getting user's account data: %v\n",
				err)
			// ALSO LOG THIS WHITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}
		p := new(content.Page)
		p, err = content.Get(ctx, "roles")
		if p == nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err != nil {
			log.Printf("Error while getting page content. Error: %v\n", err)
		}
		p.D.Account = acc
		p.D.User = u2
		template.RenderRoles(w, p)
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
			// ALSO LOG THIS WHITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
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
		case "/signUp":
			// cookie.Set(w, r, "signUp", nil)
			fn(w, r, "signUp")
		case "/logIn":
			fn(w, r, "logIn")
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
