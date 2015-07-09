package kolkata

import (
	"github.com/likestripes/pacific"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Session struct {
	PersonId  string // TODO: this is not intuitive
	RandomInt int64  `json:",string"`
	SessionId []byte
}

func (scope *Scope) session() (person_id int64, err error) {
	session_cookie, err := scope.Request.Cookie("__Session")
	personid_cookie, err := scope.Request.Cookie("__PersonId")
	if err != nil {
		return 0, err
	}
	var session Session

	query := pacific.Query{
		Context:   scope.Context,
		Kind:      "Session",
		KeyString: session_cookie.Value,
	}

	if err = query.Get(&session); err != nil {
		session_cookie.MaxAge = -1
		http.SetCookie(scope.Writer, session_cookie)

		personid_cookie.MaxAge = -1
		http.SetCookie(scope.Writer, personid_cookie)
		return 0, err
	} else if person_id_str := personid_cookie.Value; person_id_str == session.PersonId {
		scope.SessionCookie(person_id_str, *session_cookie, *personid_cookie)
		person_id, _ := strconv.ParseInt(person_id_str, 10, 64)
		return person_id, nil
	}

	return 0, err
}

func (scope *Scope) NewSession(person_id ...int64) {
	var session_cookie http.Cookie
	var personid_cookie http.Cookie

	person_id_str := scope.Person.PersonIdStr
	if len(person_id) > 0 {
		person_id_str = strconv.FormatInt(person_id[0], 10)
	}

	scope.SessionCookie(person_id_str, session_cookie, personid_cookie)
}

func (scope *Scope) SessionCookie(person_id_str string, session_cookie http.Cookie, personid_cookie http.Cookie) {

	expire := time.Now().AddDate(0, 1, 0)
	random := rand.Int63()

	session := Session{
		PersonId:  person_id_str,
		RandomInt: random,
		SessionId: []byte(randomCharacters(32)),
	}

	session_cookie.MaxAge = -1
	http.SetCookie(scope.Writer, &session_cookie)

	personid_cookie.MaxAge = -1
	http.SetCookie(scope.Writer, &personid_cookie)

	session_cookie = http.Cookie{
		Name:     "__Session",
		Value:    string(session.SessionId),
		Expires:  expire,
		Path:     "/",
		MaxAge:   86400 * 30,
		Secure:   false,
		HttpOnly: false,
	}
	http.SetCookie(scope.Writer, &session_cookie)

	personid_cookie = http.Cookie{
		Name:     "__PersonId",
		Value:    session.PersonId,
		Path:     "/",
		Expires:  expire,
		MaxAge:   86400 * 30,
		Secure:   false,
		HttpOnly: false,
	}
	http.SetCookie(scope.Writer, &personid_cookie)

	query := pacific.Query{
		Context:   scope.Context,
		Kind:      "Session",
		KeyString: string(session.SessionId),
	}
	query.Put(&session)
}

func (scope *Scope) anonCookie() int64 {
	whom_cookie, err := scope.Request.Cookie("__Whom")

	if err != nil || whom_cookie.Value == "0" {
		whom := rand.Int63()

		session_cookie := http.Cookie{
			Name:     "__Whom",
			Value:    strconv.FormatInt(whom, 10),
			Path:     "/",
			Expires:  time.Now().AddDate(0, 0, 1),
			Secure:   false,
			HttpOnly: false,
		}
		http.SetCookie(scope.Writer, &session_cookie)

		return whom
	}
	person_id, _ := strconv.ParseInt(whom_cookie.Value, 10, 64)
	return person_id
}

func (scope *Scope) ClearSession(w http.ResponseWriter, r *http.Request) {
	expire := time.Now()

	session_cookie := http.Cookie{
		Name:     "__Session",
		Value:    "signed out",
		Expires:  expire,
		Path:     "/",
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: false,
	}
	http.SetCookie(scope.Writer, &session_cookie)

	personid_cookie := http.Cookie{
		Name:     "__PersonId",
		Value:    "signed out",
		Expires:  expire,
		Path:     "/",
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: false,
	}
	http.SetCookie(scope.Writer, &personid_cookie)

	anon_cookie := http.Cookie{
		Name:     "__Whom",
		Value:    "signed out",
		Path:     "/",
		Expires:  expire,
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: false,
	}
	http.SetCookie(scope.Writer, &anon_cookie)
}
