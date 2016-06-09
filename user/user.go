package user

import (
	"errors"
	// "fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
	// "github.com/MerinEREN/InceIs/cookie"
	valid "github.com/asaskevich/govalidator"
	"github.com/nu7hatch/gouuid"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	// "io"
	// "log"
	"net/http"
	"time"
)

var (
	EmailNotExist   = errors.New("Email Not Exist")
	ExistingEmail   = errors.New("Existing Email")
	InvalidEmail    = errors.New("Invalid Email")
	InvalidPassword = errors.New("Invalid Password")
	PutUserError    = errors.New("Error while putting user into the datastore.")
	FindUserError   = errors.New("Error while checking email existincy.")
)

type Users []User

type User struct {
	UUID       string `datastore: "" json:"uuid"`
	Email      string `datastore: "" json:"email"`
	Password   string `datastore: "" json:"password"`
	ProfilePic string `datastore: "" json:"profile_pic"`
	Name       Name   `datastore: "" json:"name"`
	Phone      string `datastore: "" json:"phone"` // Should be struct in
	// the future !!!
	Status       string    `datastore: "" json:"status"`
	Type         string    `datastore: "" json:"type"`
	BirthDate    time.Time `datastore: "" json:"birth_date"`
	Registered   time.Time `datastore: "" json:"registered"`
	LastModified time.Time `datastore: "" json:"last_modified"`
	IsActive     bool      `datastore: "" json:"is_active"`
	// CAN'T USE [][] in DATASTORE, CHACK A BETTER WAY FOR THIS
	ServicePacks string `datastore: "-" json:"service_packs""`
	// 	PurchasedServices PurchasedServices `datastore: "-"
	// json:"purchasede_srvices"`
}

type Name struct {
	First string `datastore: "" json:"first"`
	Last  string `datastore: "" json:"last"`
}

type ServicePacks []ServicePack

type ServicePack struct {
	Id             string            `datastore: "" json:"id"`
	Type           string            `datastore: "" json:"type"`
	Title          string            `datastore: "" json:"title"`
	Description    string            `datastore: "" json:"description"`
	Duration       string            `datastore: "" json:"duration"`
	Price          Price             `datastore: "" json:"price"`
	Extras         ServicePackExtras `datastore: "" json:"extras"`
	Photos         Photos            `datastore: "" json:"photos"`
	Videos         Videos            `datastore: "" json:"videos"`
	Tags           Tags              `datastore: "" json:"tags"`
	Created        time.Time         `datastore: "" json:"created"`
	LastModified   time.Time         `datastore: "" json:"last_modified"`
	Status         string            `datastore: "" json:"status"`
	Evaluation     Evaluation        `datastore: "" json:"evaluation"`
	CustomerReview string            `datastore: "" json:"customer_review"`
}

type Price struct {
	Amount   float64 `datastore: "" json:amount"`
	Currency string  `datastore: "" json:currency"`
}

type ServicePackExtras []ServicePackOption

type ServicePackOption struct {
	Id          string `datastore: "" json:"id"`
	Description string `datastore: "" json:"description"`
	Duration    string `datastore: "" json:"duration"`
	Price       Price  `datastore: "" json:"price"`
	Photos      Photos `datastore: "" json:"photos"`
	Videos      Videos `datastore: "" json:"videos"`
}

type Photos []Photo

type Photo struct {
	Id           string    `datastore: "" json:"id"`
	Path         string    `datastore: "" json:"path"`
	Title        string    `datastore: "" json:"title"`
	Description  string    `datastore: "" json:"description"`
	Uploaded     time.Time `datastore: "" json:"uploaded"`
	LastModified time.Time `datastore: "" json:"last_modified"`
	Status       string    `datastore: "" json:"status"`
}

type Videos []Video

type Video struct {
	Id           string    `datastore: "" json:"id"`
	Path         string    `datastore: "" json:"path"`
	Title        string    `datastore: "" json:"title"`
	Description  string    `datastore: "" json:"description"`
	Uploaded     time.Time `datastore: "" json:"uploaded"`
	LastModified time.Time `datastore: "" json:"last_modified"`
	Status       string    `datastore: "" json:"status"`
}

type Tags []Tag

type Tag struct {
	Value string `datastore: "" json:"value"`
}

type Evaluation struct {
	Technical     int
	Timing        int
	Communication int
}

type Doc interface {
	// Use this for all structs
	// Update()
	// Upsert()
	// Delete()
}

func Add(r *http.Request, parentKey *datastore.Key) (u *User, err error) {
	ctx := appengine.NewContext(r)
	var email string
	if r.Method == "POST" {
		k := "email"
		email = r.PostFormValue(k)
		if !valid.IsEmail(email) {
			err = InvalidEmail
			return
		}
	} else {
		u1 := user.Current(ctx)
		email = u1.Email
	}
	/* k = "password"
	password := r.PostFormValue(k)
	// Cahange this control and allow special characters !!!!!!!!!!!!!!!!!!
	if !valid.IsAlphanumeric(password) {
		err = InvalidPassword
		return
	} */
	u, _, err = Exist(ctx, email)
	if err == datastore.Done {
		u4, errUUID := uuid.NewV4()
		if errUUID != nil {
			err = errUUID
			return
		}
		UUID := u4.String()
		u = &User{
			UUID:  UUID,
			Email: email,
			// Password:     GetHmac(password),
			Status:       "online",
			Type:         "admin",
			IsActive:     true,
			Registered:   time.Now(),
			LastModified: time.Now(),
		}
		key := datastore.NewKey(ctx, "Users", u.UUID, 0, parentKey)
		_, err = datastore.Put(ctx, key, u)
		if err != nil {
			return
		}
	}
	return
}

func Exist(ctx context.Context, email string) (u *User, key *datastore.Key, err error) {
	// q := datastore.NewQuery("Users").Project("UUID", "Email")
	q := datastore.NewQuery("Users")
	for t := q.Run(ctx); ; {
		u = &User{}
		key, err = t.Next(u)
		// log.Printf("USERRRRRRRRRRRRRR: %v\n", u)
		if err == datastore.Done {
			return
		}
		if err != nil {
			err = FindUserError
			return
		}
		if email == u.Email {
			err = ExistingEmail
			return
		}
	}
}
