package template

import (
	"github.com/MerinEREN/InceIs/page/content"
	"html/template"
	"log"
	"net/http"
)

var (
	templates      = template.Must(template.ParseGlob("../page/templates/*.html"))
	RenderIndex    = renderTemplate("index")
	RenderSignUp   = renderTemplate("signUp")
	RenderLogIn    = renderTemplate("logIn")
	RenderAccounts = renderTemplate("accounts")
	RenderAccount  = renderTemplate("account")
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
