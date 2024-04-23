package users

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	email string
	mdp   string
}

var store DatabaseUser

const ErrLogin = "échec du login"
const ErrNoMail = "l'email ne peut pas être vide"
const ErrInvMail = "email invalide"
const ErrNoMdp = "le mot de passe ne peut pas être vide"
const ErrBadMdp = "mot de passe incorrect"
const ErrDifMdp = "mots de passe différents"
const ErrShortMdp = "mot de passe trop court (6 caractères minimum)"
const ErrHashMdp = "erreur d'encodage du mot de passe"
const ErrDBNil = "erreur db nil"
const MinPasswordLen = 6

func Init(db DatabaseUser) error {
	if db == nil {
		return fmt.Errorf(ErrDBNil)
	}

	store = db

	return nil
}

func GetMdp(u User) string {
	return u.mdp
}

func GetEmail(u User) string {
	return u.email
}

func SetMdp(mdp string) (u User) {
	u.mdp = mdp
	return u
}

func SetEmail(mail string) (u User) {
	u.email = mail
	return u
}

func NewUser(email, mdp string) (u User, err error) {
	if email == "" {
		return u, fmt.Errorf("%v", ErrNoMail)
	}

	if mdp == "" {
		return u, fmt.Errorf("%v", ErrNoMdp)
	} else if len(mdp) < MinPasswordLen {
		return u, fmt.Errorf("%v", ErrShortMdp)
	}

	if !strings.Contains(email, "@") {
		return u, fmt.Errorf("%v", ErrInvMail)
	}

	u.email = email
	u.mdp = mdp

	return u, nil
}

func Signin(email, mdp, confirmmdp string) (u User, err error) {
	if mdp != confirmmdp {
		return u, fmt.Errorf("%v", ErrDifMdp)
	}

	u, err = NewUser(email, mdp)
	if err != nil {
		return u, err
	}

	u.mdp, err = HashPassword(mdp)
	if err != nil {
		return u, fmt.Errorf("%v : %v", ErrHashMdp, err)
	}

	err = store.AddUserDb(u)
	if err != nil {
		return u, err
	}

	return u, nil
}

func Login(email, mdp string) (u User, err error) {
	u, err = NewUser(email, mdp)
	if err != nil {
		return u, err
	}

	user, err := store.GetUser(u)
	if err != nil {
		return u, fmt.Errorf("%v : %v", ErrLogin, err.Error())
	}

	if !checkPasswordHash(u.mdp, user.mdp) {
		return u, fmt.Errorf(ErrBadMdp)
	}

	return u, nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
