package kolkata

import (
	"github.com/likestripes/pacific"
	"math/rand"
	"net/http"
)

type Scope struct {
	Writer    http.ResponseWriter
	Request   *http.Request
	Context   pacific.Context
	Bootstrap string
	Host      string
	Person    Person
	NewUser   bool
	Min       bool
}

func Current(w http.ResponseWriter, r *http.Request) (person Person, err error) {

	scope := CreateScope(w, r)
	person, err = scope.Session()

	if r.Method != "GET" && r.URL.Path != "/user/create" && r.URL.Path != "/user/auth" && r.URL.Path != "/user/sign_out" && person.Anon { //TODO: this is ugly, kolkata shouldn't have to know about calcutta
		if person.PersonId == 0 {
			person.PersonId = rand.Int63()
		}
		person.Save()
	}

	return person, err
}

func CreateScope(w http.ResponseWriter, r *http.Request) (scope Scope) {

	c := pacific.NewContext(r)

	host := r.URL.Scheme + "://" + r.Host
	if r.URL.Scheme == "" {
		host = "http://" + r.Host
	}

	scope = Scope{
		Writer:  w,
		Request: r,
		Context: c,
		Host:    host,
	}

	return
}
