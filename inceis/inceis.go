package inceis

import (
	// "fmt"
	"github.com/MerinEREN/InceIs/account"
	"github.com/MerinEREN/InceIs/cookie"
	"github.com/MerinEREN/InceIs/page/content"
	usr "github.com/MerinEREN/InceIs/user"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
	// "io/ioutil"
	// "github.com/MerinEREN/InceIs/page/temp"
	"html/template"
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
	http.HandleFunc("/signUp", makeHandler(signUpHandler))
	http.HandleFunc("/logIn", makeHandler(logInHandler))
	http.HandleFunc("/accounts", makeHandler(accountsHandler))
	http.HandleFunc("/accounts/", makeHandler(accountHandler))
	http.HandleFunc("/logOut", makeHandler(logOutHandler))
	/* if http.PostForm("/logIn", data); err != nil {
		http.Err(w, "Internal server error while login",
			http.StatusBadRequest)
	} */
	fs := http.FileServer(http.Dir("../assets"))
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
	/* p, err := content.Get(r, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} */
	var err error
	p := &content.Page{}
	p.Title = s
	var acc account.Account
	accPointer := &acc
	var pBody string
	ctx := appengine.NewContext(r)
	u1 := user.Current(ctx)
	if u1 == nil {
		var url string
		url, err = user.LoginURL(ctx, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		p.RedirectURL = url
		pBody = `<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<!--Import Google Icon Font-->
      			<link href="http://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">
      			<!--Import materialize.css-->
      			<link type="text/css" rel="stylesheet" href="css/materialize.min.css"  media="screen,projection"/>
			<!--Let browser know website is optimized for mobile-->
      			<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
			<title>InceIs</title>
		</head>
		<body>
			<h1>{{.Title}}</h1>
			<a href="{{.RedirectURL}}">Login With Google Account</a>
		</body>
		</html>`
		// http.Redirect(w, r, url, http.StatusFound)
	} else {
		// var acc *account.Account
		var errAc error
		var u2 *usr.User
		var uKey *datastore.Key
		u2, uKey, err = usr.Exist(ctx, u1.Email)
		switch err {
		case datastore.Done:
			accPointer, u2, errAc = account.Create(r)
			if errAc != nil {
				log.Printf("Error while creating account: %v\n", errAc)
				// ALSO LOG THIS WHITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!
				http.Error(w, errAc.Error(),
					http.StatusInternalServerError)
				return
			}
		case usr.ExistingEmail:
			// GET USERS ACCOUNT DATA AND USE THEM INSTEAD OF USERS !!!!!!!!!!!
			aKey := uKey.Parent()
			// log.Println(uKey, aKey, accPointer)
			errAc = datastore.Get(ctx, aKey, accPointer)
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
		errCookie := cookie.Set(w, r, "session", u2.UUID)
		if errCookie != nil {
			log.Printf("Error while creating session cookie: %v\n", errCookie)
		}
		p.Account = accPointer
		p.User = u2
		pBody = `<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<!--Import Google Icon Font-->
      			<link href="http://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">
      			<!--Import materialize.css-->
      			<link type="text/css" rel="stylesheet" href="css/materialize.min.css"  media="screen,projection"/>
			<!--Let browser know website is optimized for mobile-->
      			<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
			<title>InceIs</title>
		</head>
		<body>
			<h1>{{.Title}}</h1>
			<h2>Hello {{.Account.Name}} =)</h2>
			<h3>AND</h3>
			<h2>Hello {{.User.Email}} =)</h2>
			<h3>This is your account data below</h3>
			<i>{{.Account}}</i>
			<h3>AND</h3>
			<h3>This is your user data below</h3>
			<i>{{.User}}</i>
			<a href="/logOut">Logout</a>
		</body>
		</html>`
	}
	temp := template.Must(template.New("fdsfdfdf").Parse(pBody))
	err = temp.Execute(w, p)
	if err != nil {
		log.Print(err)
	}
	// template.RenderIndex(w, p)
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

func signUpHandler(w http.ResponseWriter, r *http.Request, s string) {
	// w.Header().Set("Content-Type", "text/plain")
	// w.Write([]byte("This is Sign Up page " + s + "\n"))
	/* p, err := content.Get(r, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} */
	/* if r.Method == "POST" {
		acc, UUID, err := account.Create(r)
		if err != nil {
			log.Printf("Error while creating account: %v\n", err)
			// ALSO LOG THIS WHITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cookie.Set(w, r, "session", UUID)
		http.Redirect(w, r, "/accounts/"+acc.Name, 302)
	} */
	p := &content.Page{}
	p.Title = s
	pBody := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="windows-1254">
<!--Import Google Icon Font-->
      	<link href="http://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">
      	<!--Import materialize.css-->
      	<link type="text/css" rel="stylesheet" href="css/materialize.min.css"  media="screen,projection"/>

      	<!--Let browser know website is optimized for mobile-->
      	<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
		<title>InceIs</title>
	</head>
	<body>
		<h1>Sign Up</h1>
<div class="row">
    <form class="col s12">
      <div class="row">
        <div class="input-field col s6">
          <input placeholder="Placeholder" id="first_name" type="text" class="validate">
          <label for="first_name">First Name</label>
        </div>
        <div class="input-field col s6">
          <input id="last_name" type="text" class="validate">
          <label for="last_name">Last Name</label>
        </div>
      </div>
      <div class="row">
        <div class="input-field col s12">
          <input disabled value="I am not editable" id="disabled" type="text" class="validate">
          <label for="disabled">Disabled</label>
        </div>
      </div>
      <div class="row">
      <div class="row">
        <div class="input-field col s12">
          <input id="email" type="email" name="email" class="validate">
          <label for="email">Email</label>
        </div>
        <div class="input-field col s12">
          <input id="password" type="password" name="password" class="validate">
          <label for="password">Sifre</label>
        </div>
      </div>
      </div>
    </form>
  </div>
		<form method="post">
			<input name="email" type="text">
			<input name="password" type="password">
			<input type="submit">
		</form>
<!--Import jQuery before materialize.js-->
      	<script type="text/javascript" src="https://code.jquery.com/jquery-2.1.1.min.js"></script>
      	<script type="text/javascript" src="js/lib/materialize.min.js"></script><a href="/logIn">Login</a>
	</body>
	</html>`
	temp := template.Must(template.New("fdsfdfdf").Parse(pBody))
	err := temp.Execute(w, p)
	if err != nil {
		log.Print(err)
	}
	// template.RenderSignUp(w, p)
}

func logInHandler(w http.ResponseWriter, r *http.Request, s string) {
	/* p, err := content.Get(r, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} */
	/* if r.Method == "POST" {
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
					cookie.Set(w, r, "session", u.UUID)
					// BE SURE THIS RETURNS ONLY OUTSIDE OF FOR THEN
					// ACTIVATE IT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
					// return
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
	} */
	/* p := &content.Page{}
		p.Title = s
		pBody := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>InceIs</title>
	</head>
	<body>
		<h1>{{.Title}}</h1>
		<div align="center">
			<img src="{{.ProfilePic}}" style="width: 13%">
			<form method="post">
				<ul style="list-style-type: none">
					<li><input name="email" type="text" value="{{.Form.Email}}"></li>
					<li><input name="password" type="password"></li>
					<li><input type="submit"></li>
					</li>
			</form>
			<a href="/signUp">Sign Up</a>
		</div>
	</body>
	</html>`
		temp := template.Must(template.New("fdsfdfdf").Parse(pBody))
		err := temp.Execute(w, p)
		if err != nil {
			log.Print(err)
		} */
	// template.RenderLogIn(w, p)
}

func accountsHandler(w http.ResponseWriter, r *http.Request, s string) {
	/* p, err := content.Get(r, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} */
	p := &content.Page{}
	p.Title = s
	pBody := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>InceIs</title>
</head>
<body>
	<i>{{.Title}} gonna be listed here =)</i>
	<a href="/logOut">Log Out</a>
</body>
</html>`
	temp := template.Must(template.New("fdsfdfdf").Parse(pBody))
	err := temp.Execute(w, p)
	if err != nil {
		log.Print(err)
	}
	// template.RenderAccounts(w, p)
}

func accountHandler(w http.ResponseWriter, r *http.Request, s string) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, url, http.StatusFound)
		// w.Header.Set("Location", url)
		// w.WriteHeader(http.StatusFound)
	} else {
		// USE USER DATA TO GET THE REST FROM DATASTORE, SET COOKIE AND MEMCACHE !!
	}
	/* p, err := content.Get(r, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} */
	p := &content.Page{}
	p.Title = s
	pBody := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>InceIs</title>
</head>
<body>
	<i>This is {{.Title}} page =)</i>	
	<a href="/logOut">Log Out</a>
</body>
</html>`
	temp := template.Must(template.New("fdsfdfdf").Parse(pBody))
	err := temp.Execute(w, p)
	if err != nil {
		log.Print(err)
	}
	// template.RenderAccount(w, p)
}

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
