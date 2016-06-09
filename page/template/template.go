package template

import (
	"github.com/MerinEREN/InceIs/page/content"
	"html/template"
	"log"
	"net/http"
)

var (
	templates          = template.Must(template.ParseGlob("/templates/*.html"))
	RenderIndex        = renderTemplate("index")
	RenderUserSettings = renderTemplate("userSettings")
	// RenderAccounts     = renderTemplate("accounts")
	// RenderSignUp   = renderTemplate("signUp")
	// RenderLogIn    = renderTemplate("logIn")
	// RenderAccount  = renderTemplate("account")
)

func renderTemplate(title string) func(w http.ResponseWriter, p *content.Page) {
	return func(w http.ResponseWriter, p *content.Page) {
		err := templates.ExecuteTemplate(w, title+".html", p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatalln(err)
		}
	}
}
