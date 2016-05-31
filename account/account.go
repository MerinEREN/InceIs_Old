package account

import (
	// "fmt"
	"google.golang.org/appengine"
	// "github.com/MerinEREN/InceIs/cookie"
	usr "github.com/MerinEREN/InceIs/user"
	// valid "github.com/asaskevich/govalidator"
	// "google.golang.org/appengine/user"
	// "github.com/nu7hatch/gouuid"
	"google.golang.org/appengine/datastore"
	// "log"
	"net/http"
	"strconv"
	"time"
)

type Accounts []Account

type Account struct {
	Name          string    `datastore: "" json:"name"`
	Company       Company   `datastore: "" json:"company"`
	CurrentStatus string    `datastore: "" json:"current_status"`
	AccountStatus string    `datastore: "" json:"account_status"`
	About         string    `datastore: "" json:"about"`
	Tags          Tags      `datastore: "" json:"tags"`
	Ranks         Ranks     `datastore: "" json:"ranks"`
	Card          Card      `datastore: "" json:"card" valid:"creditCard"`
	Users         usr.Users `datastore: "-" json:"users"`
	Registered    time.Time `datastore: "" json:"registered"`
	LastModified  time.Time `datastore: "" json:"last_modified"`
}

type Company struct {
	Name    string  `datastore: "" json:"name"`
	Type    string  `datastore: "" json:"type"`
	Address Address `datastore: "" json:"address"`
}

type Address struct {
	Description string      `datastore: "" json:"description"`
	Borough     string      `datastore: "" json:"borough"`
	City        string      `datastore: "" json:"city"`
	Country     string      `datastore: "" json:"country"`
	Postcode    string      `datastore: "" json:"postcode"`
	Geolocation Geolocation `datastore: "" json:"geolocation"`
}

type Geolocation struct {
	Lat  string `datastore: "" json:"lat"`  // type could be differnt !!!
	Long string `datastore: "" json:"Long"` // type could be differnt !!!
}

type Tags []Tag

type Tag struct {
	Value string `datastore: "" json:"value"`
}

type Ranks []Rank

type Rank struct {
	Value string `datastore: "" json:"value"`
}

type Card struct {
	CreditCards CreditCards `datastore: "" json:"creditCards"`
	DebitCards  DebitCards  `datastore: "" json:"debitCards"`
}

type CreditCards []CreditCard

type CreditCard struct {
	HolderName string `datastore: "" json:"holder_name"`
	No         string `datastore: "" json:"no"`
	ExpMonth   string `datastore: "" json:"exp_month"`
	ExpYear    string `datastore: "" json:"exp_year"`
	CVV        string `datastore: "" json:"cvv"`
}

type DebitCards []DebitCard

type DebitCard struct {
	HolderName string `datastore: "" json:"holder_name"`
	No         string `datastore: "" json:"no"`
	ExpMonth   string `datastore: "" json:"exp_month"`
	ExpYear    string `datastore: "" json:"exp_year"`
	CVV        string `datastore: "" json:"cvv"`
}

type Doc interface {
	// Use this for all structs
	// Update()
	// Upsert()
	// Delete()
}

func Create(r *http.Request) (acc *Account, u *usr.User, err error) {
	/* k = "password"
	password := r.PostFormValue(k)
	if !valid.IsEmail(email) {
		err = usr.InvalidEmail
		return
	} */
	// CAHANGE THIS CONTROL AND ALLOW SPECIAL CHARACTERS !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	/* if !valid.IsAlphanumeric(password) {
		err = usr.InvalidPassword
		return
	} */
	ctx := appengine.NewContext(r)
	q := datastore.NewQuery("Accounts")
	var accCount int
	accCount, err = q.Count(ctx)
	if err != nil {
		return
	}
	accCount++
	acc = &Account{
		Name:          "Account_" + strconv.Itoa(accCount),
		CurrentStatus: "available",
		AccountStatus: "online",
		Registered:    time.Now(),
		LastModified:  time.Now(),
	}
	key := datastore.NewIncompleteKey(ctx, "Accounts", nil)
	// fmt.Println(key)
	var parentKey *datastore.Key
	parentKey, err = datastore.Put(ctx, key, acc)
	// fmt.Println(parentKey)
	if err != nil {
		return
	}
	u, err = usr.Add(r, parentKey)
	if err != nil {
		// DELETE CREATED ACCOUNT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		return
	}
	return
}

/* func AddTags(s ...string) bool {
	return
} */
