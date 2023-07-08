package farm

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Geo struct {
	Code    string
	Country string
}

type ProfileName struct {
	Prefix string
	Geo    Geo
	Suffix string
	Number int
}

type Farmer struct {
	Name string
}

type Status struct {
	Name        string
	Description string
	IsFarmable  bool
	Color       string
}

type Tag struct {
	Name        string
	Description string
	Color       string
}

type StrategyDescription struct {
	Name        string
	Description string
}

type LoginCredentials struct {
	Login string
	Pass  string
}

func (l *LoginCredentials) String() string {
	return fmt.Sprintf("%v:%v", l.Login, l.Pass)
}

func NewLoginCredentials(login, pass string) (*LoginCredentials, error) {

	if !isValidEmail(login) {
		return nil, errors.New("provided string is not valid email")
	}

	return &LoginCredentials{
		Login: login,
		Pass:  pass,
	}, nil
}

func isValidEmail(login string) bool {

	login = strings.TrimSpace(login)

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	return emailRegex.MatchString(login)
}

type Email struct {
	Login    LoginCredentials
	Provider string
	Reserve  []Email
}

type Person struct {
	FirstName  string
	SecondName string
	// [ ] Documents: {FileID, FileName}
}

func (p *Person) String() string {
	return fmt.Sprintf("%v:%v", p.FirstName, p.SecondName)
}

type Profile struct {
	Name        ProfileName
	Responsible Farmer

	Strategy StrategyDescription
	Status   Status
	Tags     []Tag

	Emails []Email

	Person Person
}
