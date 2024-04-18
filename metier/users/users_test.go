package users

import (
	"fmt"
	"strings"
	"testing"
)

type fakeDb struct {
	users []User
}

func (f *fakeDb) AddUserDb(u User) error {
	for _, user := range f.users {
		if user.email == u.email {
			return fmt.Errorf("email déjà utilisé")
		}
	}
	f.users = append(f.users, u)
	return nil
}
func (f *fakeDb) GetUser(u User) (User, error) {
	for _, user := range f.users {
		if user.email == u.email {
			return user, nil
		}
	}
	return User{}, fmt.Errorf("error")
}

func setupFakeDb() fakeDb {
	db := fakeDb{}

	mdp1, _ := HashPassword("25mai1995")
	mdp2, _ := HashPassword("sortla8.6")

	user1 := User{email: "mail20@mail.com", mdp: mdp1}
	user2 := User{email: "clement@caramail.com", mdp: mdp2}

	db.AddUserDb(user1)
	db.AddUserDb(user2)

	return db
}

func TestLogin(t *testing.T) {
	db := setupFakeDb()
	Init(&db)

	var tests = []struct {
		name, entryEmail, entryPassword, want string
	}{
		{"Cas normal", "mail20@mail.com", "25mai1995", "mail20@mail.com"},
		{"Email vide", "", "12azerty", ERR_NOMAIL},
		{"Mot de passe incorrect", "mail20@mail.com", "azerty", ERR_BADMDP},
		{"Email invalide", "ma@mail.com", "25mai1995", ERR_LOGIN},
		{"Mot de passe vide", "clement@caramail.com", "", ERR_NOMDP},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Login(tt.entryEmail, tt.entryPassword)
			if (!strings.Contains(got.email, tt.want)) && !strings.Contains(err.Error(), tt.want) {
				t.Errorf("expected response to contain '%s', but got '%s'", tt.want, err.Error())
			}
		})
	}
}

func TestSignin(t *testing.T) {
	db := setupFakeDb()
	Init(&db)

	var tests = []struct {
		name, entryEmail, entryPassword, entryConfirm, want string
	}{
		{"Cas normal", "mail2018@mail.com", "29mai1995", "29mai1995", "mail2018@mail.com"},
		{"Mots de passes différents", "mail2019@mail.com", "29mai1995", "2mai1995", ERR_DIFMDP},
		{"Email vide", "", "12azerty", "12azerty", ERR_NOMAIL},
		{"Email invalide", "mail2018mailcom", "29mai1995", "29mai1995", ERR_INVMAIL},
		{"Mot de passe trop court", "mail@mail.com", "azey", "azey", ERR_SHORTMDP},
		{"Email déjà utilisé", "mail20@mail.com", "2mai1995", "2mai1995", "email déjà utilisé"},
		{"Mot de passe vide", "clem@caramail.com", "", "", ERR_NOMDP},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Signin(tt.entryEmail, tt.entryPassword, tt.entryConfirm)
			if (!strings.Contains(got.email, tt.want)) && !strings.Contains(err.Error(), tt.want) {
				t.Errorf("expected response to contain '%s', but got '%s'", tt.want, err.Error())
			}
		})
	}
}
