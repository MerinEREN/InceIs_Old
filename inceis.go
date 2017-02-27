/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package inceis

import (
	// "encoding/json"
	// "fmt"
	"strings"
	// "github.com/MerinEREN/iiPackages/account"
	"github.com/MerinEREN/iiPackages/apis/account"
	"github.com/MerinEREN/iiPackages/apis/accountSettings"
	"github.com/MerinEREN/iiPackages/apis/index"
	"github.com/MerinEREN/iiPackages/apis/logout"
	"github.com/MerinEREN/iiPackages/apis/roles"
	"github.com/MerinEREN/iiPackages/apis/userSettings"
	// "github.com/MerinEREN/iiPackages/cookie"
	"github.com/MerinEREN/iiPackages/page/content"
	"github.com/MerinEREN/iiPackages/page/template"
	"golang.org/x/net/context"
	// usr "github.com/MerinEREN/iiPackages/user"
	"google.golang.org/appengine"
	// "google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/user"
	// "io/ioutil"
	// "html/template"
	"log"
	// "mime/multipart"
	"net/http"
	// "regexp"
	"time"
)

var _ memcache.Item // For debugging, delete when done.

var (
// CHANGE THE REGEXP BELOW !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// validPath = regexp.MustCompile("^/|[/A-Za-z0-9]$")
)

// type LoginURLs map[string]string

func init() {
	// http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/",
		http.TimeoutHandler(http.HandlerFunc(makeHandlerFunc(index.Handler)),
			1000*time.Millisecond,
			"This is http.TimeoutHandler(handler, time.Duration, message) "+
				"message bitch =)"))
	// http.HandleFunc("/", makeHandlerFunc(index.Handler))
	http.HandleFunc("/roles/", makeHandlerFunc(roles.Handler))
	http.HandleFunc("/userSettings/", makeHandlerFunc(userSettings.Handler))
	http.HandleFunc("/accountSettings/", makeHandlerFunc(accountSettings.Handler))
	// http.HandleFunc("/signUp", makeHandlerFunc(signUpHandler))
	// http.HandleFunc("/logIn", makeHandlerFunc(logInHandler))
	// http.HandleFunc("/accounts", makeHandlerFunc(accountsHandler))
	http.HandleFunc("/accounts/", makeHandlerFunc(account.Handler))
	http.HandleFunc("/logout/", makeHandlerFunc(logout.Handler))
	/* if http.PostForm("/logIn", data); err != nil {
		http.Err(w, "Internal server error while login",
			http.StatusBadRequest)
	} */
	fs := http.FileServer(http.Dir("../iiClient/public"))
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

type handlerFuncWithContextAndUser func(context.Context, http.ResponseWriter,
	*http.Request, *user.User)

func makeHandlerFunc(fn handlerFuncWithContextAndUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		// Authenticate the client
		// ug is google user
		ug := user.Current(ctx)
		if ug == nil && r.URL.Path != "/" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		/* m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			log.Printf("Invalid Path: %s\n", r.URL.Path)
			http.NotFound(w, r)
			return
		} */
		/* for _, val := range m {
			fmt.Println(val)
		}*/
		// CHANGE CONTENT AND TEMPLATE THINGS !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		if strings.Contains(r.Header.Get("Accept"), "text/html") {
			pc, err := content.Get(ctx, "index")
			if err != nil {
				log.Printf("Error while getting page content. Error: %v\n",
					err)
			}
			log.Println("Getting template !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			template.RenderIndex(w, pc)
		} else if strings.Contains(r.Header.Get("Accept"), "text/plain") {
			log.Println("Getting data !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			fn(ctx, w, r, ug)
		}
	}
}

// HEADER ALWAYS SHOULD BE SET BEFORE ANYTHING WRITE A PAGE BODY !!!!!!!!!!!!!!!!!!!!!!!!!!
// w.Header().Set("Content-Type", "text/html"; charset=utf-8")
//fmt.Fprintln(w, things...) // Writes to the body
