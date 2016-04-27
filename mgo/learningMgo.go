package main

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"io/ioutil"
	// "log"
	"github.com/MerinEREN/OlduOlacak/collections"
	"net/http"
	"regexp"
	"time"
)

const (
	// mongoDBLocal = "mongodb://localhost"
	mongoDBLocal = "localhost"
)

type Page struct {
	Title string
	Body  []byte
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

	// mgo
	session, err := mgo.Dial(mongoDBLocal)
	if err != nil {
		fmt.Println(session, err)
		panic(err)
	}
	defer session.Close()
	// mgo DropDatabase
	// if err = session.DB("OO").DropDatabase(); err != nil {
	// fmt.Println("Error occurred when deleting OO db !!!")
	// }
	c := session.DB("OO").C("accounts")
	m := bson.M{}
	// var results Accounts
	var results2 Accounts

	// mgo Find and mgo All and mgo One and mgo Explain
	// mgo FindId(bson.M{"_id": id}) UpdateId, UpsertId...
	// err = c.Find(bson.M{"name": "Merin"}).All(&results)
	// if err != nil {
	// return err
	// }
	err = c.Find(Account{Name: "Merin", LastName: "EREN"}).All(&results2)
	if err != nil {
		// return err
	}
	err = c.Find(Account{Name: "Merin", LastName: "EREN"}).Explain(m)
	if err != nil {
		// return err
	}
	fmt.Println(results2)
	fmt.Println("\n")
	fmt.Printf("Explain: %#v\n", m)
	if err != nil {
		fmt.Printf("Explain: %#v\n", m)
	}

	// for _, v := range results {
	// 	fmt.Printf("Name: %v     Last Name: %v\n",
	// 		v.Name, v.LastName)
	// }
	for i, v := range results2 {
		fmt.Printf("Doc-%v:     Name: %v     Last Name: %v\n",
			i+1, v.Name, v.LastName)
	}

	// mgo Count
	var accountCount int
	accountCount, err = c.Find(Account{Name: "Merin", LastName: "EREN"}).
		Count()
	if err != nil {
		// return err
	}
	fmt.Println(accountCount)

	// mgo Insert
	err = c.Insert(Account{Name: "Merin", LastName: "EREN"})
	if err != nil {
		// return err
	}

	// mgo Update and UpdateAll. UpdateAll also returns info *ChangeInfo
	// like all other ...All's.
	type M map[string]interface{}
	change := M{"$set": Account{Name: "Merin", LastName: "EREN"}}
	var updateAllInfo *mgo.ChangeInfo
	updateAllInfo, err = c.UpdateAll(nil, change)
	if err != nil {
		// return err
		fmt.Println("Can't find anyone with LastName: EREN !!!")
	}
	fmt.Println(*updateAllInfo)

	//mgo Upsert
	upsertChange := M{"$set": Account{Name: "Metin", LastName: "EREN"}}
	var upsertInfo *mgo.ChangeInfo
	upsertInfo, err = c.Upsert(Account{Name: "Metin"}, upsertChange)
	if err != nil {
		fmt.Println("Upsert error !!!")
	}
	fmt.Println(*upsertInfo)

	// mgo Iter
	iter := c.Find(Account{Name: "Merin", LastName: "EREN"}).Iter()
	var account Account
	for iter.Next(&account) {
		fmt.Printf("Account: %v\n", account)
	}
	// return iter.Err() // YOU HAVE TO DO THIS CONTROL AFTER EVERY
	// ITERATION TO BE SURE ABOUT EVERY STEP !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

	// mgo Remove and RemoveAll
	info, err := c.RemoveAll(Account{Name: "Merin"})
	if err != nil {
		fmt.Println("Error when Deleting all accounts !!!")
	}
	fmt.Println(info)
	fmt.Println(&info)
	fmt.Println(*info)

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
