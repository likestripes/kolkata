package kolkata

import (
	"github.com/likestripes/pacific"
	"strconv"
	"time"
)

type Person struct {
	PersonId    int64
	PersonIdStr string
	Timestamp   time.Time
	Secret      string
	Anon        bool   `datastore:"-" sql:"-" json:"-"`
	Scope       *Scope `datastore:"-" sql:"-" json:"-"`
}

func (scope *Scope) Session() (person Person, err error) {

	person_id, err := scope.session()

	if person_id != 0 {
		person = scope.get(person_id)
	} else {
		person.PersonId = scope.anonCookie()
		person.Anon = true
	}
	person.PersonIdStr = strconv.FormatInt(person.PersonId, 10)
	person.Scope = scope
	scope.Person = person

	return person, err
}

func (scope *Scope) query(property string, value string) int64 {

	var results []Person
	property = property + "= "

	query := pacific.Query{
		Context: scope.Context,
		Kind:    "Person",
		Limit:   1,
		Filters: map[string]interface{}{
			property: value,
		},
	}

	query.GetAll(&results)

	if len(results) == 0 {
		return 0
	}
	return results[0].PersonId

}

func (scope *Scope) get(person_id int64) (person Person) {

	query := pacific.Query{
		Context: scope.Context,
		Kind:    "Person",
		KeyInt:  person_id,
	}
	query.Get(&person)

	return person
}

func (person *Person) Save(sign_ins ...map[string]string) int64 {

	scope := person.Scope
	existing := scope.get(person.PersonId)

	if existing.PersonId != person.PersonId {

		person.Timestamp = time.Now()
		person.Secret = randomCharacters(128)
		person.PersonIdStr = strconv.FormatInt(person.PersonId, 10)

		query := pacific.Query{
			Context: scope.Context,
			Kind:    "Person",
			KeyInt:  person.PersonId,
		}
		err := query.Put(person)
		if err != nil {
			scope.Context.Errorf(err.Error())
		}

	}

	if len(sign_ins) > 0 {
		for token, password := range sign_ins[0] {
			if token != "" {
				signin := SignIn{
					Token:    token,
					PersonId: person.PersonId,
					Unsafe:   password,
					Context:  person.Scope.Context,
				}
				signin.Save()
			}
		}
	}

	return person.PersonId
}
