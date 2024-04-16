package users

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email string
	Mdp   string
}

var store DatabaseUser

const ERR_LOGIN = "échec du login"
const ERR_NOMAIL = "l'email ne peut pas être vide"
const ERR_INVMAIL = "email invalide"
const ERR_NOMDP = "le mot de passe ne peut pas être vide"
const ERR_BADMDP = "mot de passe incorrect"
const ERR_DIFMDP = "mots de passe différents"
const ERR_SHORTMDP = "mot de passe trop court (6 caractères minimum)"
const ERR_HASHMDP = "erreur d'encodage du mot de passe"
const ERR_DBNIL = "error db nil"

func Init(db DatabaseUser) error {
	if db == nil {
		return fmt.Errorf(ERR_DBNIL)
	}
	store = db
	return nil
}

func NewUser(email string, mdp string) (u User, err error) {
	if email == "" {
		return u, fmt.Errorf("%v", ERR_NOMAIL)
	}
	if mdp == "" {
		return u, fmt.Errorf("%v", ERR_NOMDP)
	} else if len(mdp) < 6 {
		return u, fmt.Errorf("%v", ERR_SHORTMDP)
	}
	if !strings.Contains(email, "@") {
		return u, fmt.Errorf("%v", ERR_INVMAIL)
	}
	u.Email = email
	u.Mdp = mdp
	return u, nil
}

func Signin(email string, mdp string, confirmmdp string) (u User, err error) {
	if mdp != confirmmdp {
		return u, fmt.Errorf("%v", ERR_DIFMDP)
	}
	u, err = NewUser(email, mdp)
	if err != nil {
		return u, err
	}
	u.Mdp, err = HashPassword(mdp)
	if err != nil {
		return u, fmt.Errorf("%v : %v", ERR_HASHMDP, err)
	}
	err = store.AddUserDb(u)
	if err != nil {
		return u, err
	}
	return u, nil
}

func Login(email string, mdp string) (u User, err error) {
	u, err = NewUser(email, mdp)
	if err != nil {
		return u, err
	}
	user, err := store.GetUser(u)
	if err != nil {
		return u, fmt.Errorf("%v : %v", ERR_LOGIN, err)
	}
	if !checkPasswordHash(u.Mdp, user.Mdp) {
		return u, fmt.Errorf(ERR_BADMDP)
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
