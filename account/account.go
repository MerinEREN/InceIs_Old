package account

import (
	"errors"
	"fmt"
	"github.com/MerinEREN/InceIs/cookie"
	valid "github.com/asaskevich/govalidator"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	EmailNotExist   = errors.New("Invalid Email")
	ExistingEmail   = errors.New("Existing Email")
	InvalidPassword = errors.New("Invalid Password")
)

// II DB Structs
// Accounts collection
type Accounts []account

type account struct {
	Name          string    `bson:"name,omitempty"`
	Type          string    `bson:"type,omitempty"`
	CurrentStatus string    `bson:"current_status,omitempty"`
	AccountStatus string    `bson:"account_status,omitempty"`
	About         string    `bson:"about,omitempty"`
	Tags          Tags      `bson:"tags,omitempty"`
	Ranks         Ranks     `bson:"ranks,omitempty"`
	Card          card      `bson:"card,omitempty" valid:"creditcard"`
	Users         Users     `bson:"users,omitempty"`
	Registered    time.Time `bson:"registered,omitempty"`
	LastModified  time.Time `bson:"last_modified,omitempty"`
}

type company struct {
	Name    string  `bson:"name,omitempty"`
	Address address `bson:"address,omitempty"`
}

type address struct {
	Description string      `bson:"description,omitempty"`
	Borough     string      `bson:"borough,omitempty"`
	City        string      `bson:"city,omitempty"`
	Country     string      `bson:"country,omitempty"`
	Postcode    string      `bson:"postcode,omitempty"`
	Geolocation geolocation `bson:"geolocation,omitempty"`
}

type geolocation struct {
	Lat  string `bson:"lat,omitempty"`  // type could be differnt !!!
	Long string `bson:"Long,omitempty"` // type could be differnt !!!
}

type Tags []tag

type tag struct {
	Type string `bson:"type,omitempty"`
}

type Ranks []rank

type rank struct {
	Type string `bson:"type,omitempty"`
}

type card struct {
	Creditcards Creditcards `bson:"creditcards,omitempty"`
	Debitcards  Debitcards  `bson:"debitcards,omitempty"`
}

type Creditcards []creditcard

type creditcard struct {
	HolderName string `bson:"holder_name,omitempty"`
	No         string `bson:"no,omitempty"`
	ExpMonth   string `bson:"exp_month,omitempty"`
	ExpYear    string `bson:"exp_year,omitempty"`
	CVV        string `bson:"cvv,omitempty"`
}

type Debitcards []debitcard

type debitcard struct {
	HolderName string `bson:"holder_name,omitempty"`
	No         string `bson:"no,omitempty"`
	ExpMonth   string `bson:"exp_month,omitempty"`
	ExpYear    string `bson:"exp_year,omitempty"`
	CVV        string `bson:"cvv,omitempty"`
}

type Users []user

type user struct {
	UUID        string `bson:"uuid"`
	Email       string `bson:"email,omitempty"`
	Password    string `bson:"password,omitempty"`
	PicturePath string `bson:"picture_path,omitempty"`
	Name        name   `bson:"name,omitempty"`
	Phone       string `bson:"phone,omitempty"` // Should be struct in
	// the future !!!
	Status       string       `bson:"status,omitempty"`
	Type         string       `bson:"type,omitempty"`
	BirthDate    time.Time    `bson:"birth_date,omitempty"`
	Registered   time.Time    `bson:"registered,omitempty"`
	LastModified time.Time    `bson:"last_modified,omitempty"`
	IsActive     bool         `bson:"is_active,omitempty"`
	ServicePacks ServicePacks `bson:"service_packs",omitempty"`
	// 	PurchasedServices PurchasedServices `bson:"purchasedServices,
	// 	omitempty"`
}

type name struct {
	First string `bson:"first,omitempty"`
	Last  string `bson:"last,omitempty"`
}

type ServicePacks []servicePack

type servicePack struct {
	Id             string            `bson:"id,omitempty"`
	Type           string            `bson:"type,omitempty"`
	Description    string            `bson:"description,omitempty"`
	Duration       string            `bson:"duration,omitempty"`
	Price          price             `bson:"price,omitempty"`
	Extras         ServicePackExtras `bson:"extras,omitempty"`
	Photos         Photos            `bson:"photos,omitempty"`
	Videos         Videos            `bson:"videos,omitempty"`
	Tags           Tags              `bson:"tags,omitempty"`
	Created        time.Time         `bson:"created,omitempty"`
	LastModified   time.Time         `bson:"last_modified,omitempty"`
	Status         string            `bson:"status,omitempty"`
	Evaluation     evaluation        `bson:"evaluation,omitempty"`
	CustomerReview string            `bson:"customer_review,omitempty"`
}

type price struct {
	Amount   float64 `bson:amount,omitempty"`
	Currency string  `bson:currency,omitempty"`
}

type ServicePackExtras []servicePackOption

type servicePackOption struct {
	Id          string `bson:"id,omitempty"`
	Description string `bson:"description,omitempty"`
	Duration    string `bson:"duration,omitempty"`
	Price       price  `bson:"price,omitempty"`
	Photos      Photos `bson:"photos,omitempty"`
	Videos      Videos `bson:"videos,omitempty"`
}

type Photos []photo

type photo struct {
	Id           string    `bson:"id,omitempty"`
	Path         string    `bson:"path,omitempty"`
	Title        string    `bson:"title,omitempty"`
	Description  string    `bson:"description,omitempty"`
	Uploaded     time.Time `bson:"uploaded,omitempty"`
	LastModified time.Time `bson:"last_modified,omitempty"`
	Status       string    `bson:"status,omitempty"`
}

type Videos []video

type video struct {
	Id           string    `bson:"id,omitempty"`
	Path         string    `bson:"path,omitempty"`
	Title        string    `bson:"title,omitempty"`
	Description  string    `bson:"description,omitempty"`
	Uploaded     time.Time `bson:"uploaded,omitempty"`
	LastModified time.Time `bson:"last_modified,omitempty"`
	Status       string    `bson:"status,omitempty"`
}

type evaluation struct {
	Technical     int
	Timing        int
	Communication int
}

type doc interface {
	// Use this for all structs
	// Update()
	// Upsert()
	// Delete()
}

func Create(w http.ResponseWriter, r *http.Request, c *mgo.Collection,
	email, password string) {
	if !valid.IsEmail(email) {
		// Inform client
		fmt.Fprintln(w, "Invalid email")
		return
	}
	// Cahange this control and allow special characters !!!!!!!!!!!!!!!!!!
	if !valid.IsAlphanumeric(password) {
		// Inform client
		fmt.Fprintln(w, "Invalid password")
		return
	}
	acc, err := VerifyUser(c, email, password)
	switch err {
	case EmailNotExist:
		u4, errUUID := uuid.NewV4()
		if errUUID != nil {
			// status code could be wrong
			http.Error(w, errUUID.Error(), http.StatusNotImplemented)
			log.Fatalf("Can't create UUID when signUp, error: %v\n",
				errUUID)
		}
		users := Users{
			user{
				UUID:         u4.String(),
				Email:        email,
				Password:     password,
				Status:       "online",
				Type:         "admin",
				IsActive:     true,
				Registered:   time.Now(),
				LastModified: time.Now(),
			},
		}
		accCount, errCount := c.Find(bson.M{}).Count()
		accCount++
		if errCount != nil {
			// status code could be wrong
			http.Error(w, errCount.Error(), http.StatusNotImplemented)
			log.Fatalln(errCount)
		}
		acc = &account{
			Name:          "Account_" + strconv.Itoa(accCount),
			CurrentStatus: "available",
			AccountStatus: "online",
			Users:         users,
			Registered:    time.Now(),
			LastModified:  time.Now(),
		}
		errInsert := c.Insert(acc)
		if errInsert != nil {
			// status code could be wrong
			http.Error(w, errInsert.Error(), http.StatusNotImplemented)
			log.Fatalln(errInsert)
		}
		cookie.Create(w, r, "session", u4.String())
		http.Redirect(w, r, "/accounts/"+acc.Name, 302)
	case ExistingEmail:
		fmt.Fprintln(w, err)
	case InvalidPassword:
		fmt.Fprintln(w, err)
	default:
		// status code could be wrong
		http.Error(w, err.Error(), http.StatusNotImplemented)
		log.Fatalln(err)
	}
}

func VerifyUser(c *mgo.Collection, e, p string) (result *account, err error) {
	// CHECK THIS QUERY, IT IS WRONG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	err = c.Find(bson.M{"users.email": e}).One(&result)
	if err != nil {
		if err.Error() == "not found" {
			err = EmailNotExist
		} else {
			return
		}
	} else {
		err = ExistingEmail
		for _, user := range result.Users {
			if user.Password != p {
				err = InvalidPassword
			}
		}
	}
	return
}
