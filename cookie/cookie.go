package cookie

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// CHANGE THIS DUMMY COOKIE STRUCT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
type sessionData struct {
	Name  string
	Value string
}
type logInData struct {
	Photo string
}

// ADDING UUID, NECESSARY USER DATA, COUNT AND HASH TO THE COOKIE AND CHECK HASH CODE
// CREATE logIn COOKIE AND SET PROFILE PIC !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func Create(w http.ResponseWriter, r *http.Request, s, uuid string) {
	// COOKIE IS A PART OF THE HEADER, SO U SHOULD SET THE COOKIE BEFORE EXECUTING A
	// TEMPLATE OR WRITING SOMETHING TO THE BODY !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	c, err := r.Cookie(s)
	if err == http.ErrNoCookie {
		c = newCookie(s, uuid)
	}
	if isUserDataChanged(c) {
		// DELETING CORRUPTED COOKIE AND CREATING NEW ONE !!!!!!!!!!!!!!!!!!!!!!!!!
		// c.MaxAge = -1
		// http.SetCookie(w, c)
		c = newCookie(s, uuid)
	}
	// CREATING A COOKIE IS NOT ENOUGH, YOU HAVE TO SET THE COOKIE TO USE IT !!!!!!!!!!
	http.SetCookie(w, c)
}

func newCookie(s, uuid string) (c *http.Cookie) {
	c = &http.Cookie{
		Name: s,
		// U CAN USE UUID AS VALUE !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		Value: uuid,
		// NOT GOOD PRACTICE
		// ADDING USER DATA TO A COOKIE
		// WITH NO WAY OF KNOWING WHETER OR NOT THEY MIGHT HAVE ALTERED
		// THAT DATA !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		// HMAC WOULD ALLOW US TO DETERMINE WHETHER OR NOT THE DATA IN THE
		// COOKIE WAS ALTERED !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		// HOWEVER, BEST TO STORE USER DATA ON THE SERVER AND KEEP
		// BACKUPS !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		// Value: "emil = merin@inceis.net" + "JSON data" + "whatever",
		// IF SECURE IS TRUE THIS COOKIE ONLY SEND WITH HTTP2 !!!!!!!!!!!!!!!!!!!!!
		Secure: true,
		// HttpOnly: true MEANS JAVASCRIPT CAN NOT ACCESS THE COOKIE !!!!!!!!!!!!!!
		HttpOnly: false,
	}
	setValue(c)
	return
}

func setValue(c *http.Cookie) {
	// Setting different kind of struct for different cookies
	var cd interface{}
	switch c.Name {
	case "session":
		cd = sessionData{
			Name:  "cookiedataName",
			Value: "cookiedataValue",
		}
	case "logIn":
		cd = logInData{
			Photo: "userPhoto",
		}
	}
	bs, err := json.Marshal(cd)
	if err != nil {
		log.Printf("%s cookie marshaling error. %v\n", c.Name, err)
	}
	log.Printf("Marshalled cookie data is %s\n", string(bs))
	c.Value += "|" + base64.StdEncoding.EncodeToString(bs)
	code := getCode(c.Value)
	c.Value += "|" + code
	fmt.Printf("Cookie value for "+c.Name+" is: %s\n", c.Value)
}

// Checking data with "hmac"
func getCode(s string) string {
	h := hmac.New(sha256.New, []byte("someKey"))
	io.WriteString(h, s)
	return fmt.Sprintf("%v", h.Sum(nil))
}

func isUserDataChanged(c *http.Cookie) bool {
	cvSlice := strings.Split(c.Value, "|")
	uuidData := cvSlice[0] + "|" + cvSlice[1]
	returnedCode := getCode(uuidData)
	if returnedCode != cvSlice[2] {
		log.Printf("%s cookie value is corrupted. Cookie HMAC is %s, "+
			"genereted HMAC is %s", c.Name, cvSlice[2], returnedCode)
		decodedBase64, err := base64.StdEncoding.DecodeString(cvSlice[1])
		if err != nil {
			log.Printf("Error while decoding %s cookie data. Error "+
				"is %v\n", c.Name, err)
		}
		var returnedCookieData sessionData
		err = json.Unmarshal(decodedBase64, &returnedCookieData)
		if err != nil {
			log.Printf("%s cookie unmarshaling error. %v\n", c.Name, err)
		}
		log.Printf("Returned cookie data is %v", returnedCookieData)
		// DID NOT CHECKED DELETING AND CREATING NEW COOKIE YET, SO CHECK THEM !!!!

		// cookie.MaxAge = -1        // Decleration of deleting the cookie
		return true
	}
	return false
}
