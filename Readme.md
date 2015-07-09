## Kolkata

A golang library for more personable apps.

### Warning

*This is v.01 -- it has no tests, performance is probably awful and it's not proven safe in production anywhere.*  But maybe it'll scratch an itch?

### WTF

`Kolkata` is an opinionated user/session auth API that is designed to make a variety of user sign up scenarios easy.

### Quick Start?  Use `Calcutta`!

`Calcutta` implements a basic sign in form for `Kolkata`, so if you want to get started immediately, check out [likestripes/calcutta](https://www.github.com/likestripes/calcutta)

Otherwise, read on to set it up...

### Install / Import

`go get -u github.com/likestripes/kolkata`

```go
import (
	"github.com/likestripes/kolkata"
)
```

### Dependency on `Pacific`

`Kolkata` uses [likestripes/pacific](https://www.github.com/likestripes/pacific) as an opinionated ORM.  `Pacific` currently supports AppEngine and Postgres.

Google AppEngine: `goapp serve` works out of the box (they include the buildtag for you)

Postgres: `go run -tags 'postgres' main.go` -- details in the [pacific/Readme](https://github.com/likestripes/pacific/blob/master/readme.md).


### Overview

`Kolkata` is designed to export a `Person` struct that can be mixed into whatever `WildAndCrazyUser` model your app requires:

```go
type Person struct {
	PersonId    int64
	PersonIdStr string
	Timestamp   time.Time
	Secret      string
	Anon        bool   `datastore:"-" sql:"-" json:"-"`
	Scope       *Scope `datastore:"-" sql:"-" json:"-"`
}
```
The `Person` has _n_ `SignIn`s (a token + password used for authentication) and attaches itself to the current session from your app via `person, err := kolkata.Current(w, r)`


#### TODO
- [ ] logging
- [ ] documentation!
- [ ] tests!
- [ ] benchmarking

- [x] Contributors welcome!
