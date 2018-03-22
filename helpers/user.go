package helpers

import (
	"net/http"
	"strconv"

	c "github.com/astaxie/beego/context"
	"github.com/gorilla/securecookie"
)

var s = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

//KeepLogin login
func KeepLogin(w *c.Response, username string, roleID int) (ok bool, err string) {
	value := map[string]string{
		"username": username,
		"role":     strconv.Itoa(roleID),
		"req-only": "1",
	}
	if encoded, errs := s.Encode("fixman", value); errs != nil {
		ok = false
		err = errs.Error()
	} else {
		cookie := http.Cookie{
			Name:     "fixman",
			Value:    encoded,
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(w.ResponseWriter, &cookie)
		ok = true
		err = ""
	}
	return ok, err
}

//LogOut login
func LogOut(w *c.Response) {

	cookie := http.Cookie{
		Name:     "fixman",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}
	http.SetCookie(w.ResponseWriter, &cookie)
}

//GetUser get user from cookie
func GetUser(r *http.Request) string {
	if cookie, err := r.Cookie("fixman"); err == nil {
		value := make(map[string]string)
		if err = s.Decode("fixman", cookie.Value, &value); err == nil {
			return value["username"]
		}
	}
	return ""
}

//GetRole get role from cookie
func GetRole(r *http.Request) string {

	if cookie, err := r.Cookie("fixman"); err == nil {
		value := make(map[string]string)
		if err = s.Decode("fixman", cookie.Value, &value); err == nil {
			return value["role"]
		}
	}
	return ""
}
