package users

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email string
	Mdp   string
}

var store DatabaseUser
var user User
var err error

const ERR_LOGIN = "échec du login"
const ERR_NOMAIL = "l'email ne peut pas être vide"
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

func Signin(email string, mdp string, confirmmdp string) (User, error) {
	if email == "" {
		return User{}, fmt.Errorf("%v", ERR_NOMAIL)
	}
	if mdp == "" {
		return User{}, fmt.Errorf("%v", ERR_NOMDP)
	}
	user.Email = email
	if mdp != confirmmdp {
		return User{}, fmt.Errorf("%v", ERR_DIFMDP)
	} else if len(mdp) < 6 {
		return User{}, fmt.Errorf("%v", ERR_SHORTMDP)
	}

	user.Mdp, err = hashPassword(mdp)
	if err != nil {
		return User{}, fmt.Errorf("%v : %v", ERR_HASHMDP, err)
	}
	err = store.AddUserDb(user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func Login(email string, mdp string) (User, error) {
	if email == "" {
		return User{}, fmt.Errorf("%v", ERR_NOMAIL)
	}
	if mdp == "" {
		return User{}, fmt.Errorf("%v", ERR_NOMDP)
	}
	user.Email = email
	user.Mdp = mdp
	u, err := store.GetUser(user)
	if err != nil {
		return User{}, fmt.Errorf("%v : %v", ERR_LOGIN, err)
	}
	if !checkPasswordHash(user.Mdp, u.Mdp) {
		return User{}, fmt.Errorf(ERR_BADMDP)
	}
	return user, nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
