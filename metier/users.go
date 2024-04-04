package metier

import (
	"fmt"
	"webstack/models"

	"golang.org/x/crypto/bcrypt"
)

var storeUser DatabaseUser

func InitUser(db DatabaseUser) error {
	if db == nil {
		return fmt.Errorf("db est nil")
	}
	storeUser = db
	return nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func AddUser(email string, mdp string, confirmmdp string) (models.User, error) {
	user.Email = email
	if mdp != confirmmdp {
		return models.User{}, fmt.Errorf("mots de passe différents, réessayez")
	}

	user.Mdp, err = hashPassword(mdp)
	if err != nil {
		return models.User{}, fmt.Errorf("hashPassword error : %v", err)
	}
	err = storeUser.AddUserDb(user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func Login(email string, mdp string) (models.User, error) {
	user.Email = email
	user.Mdp = mdp
	u, err := storeUser.GetUser(user)
	if err != nil {
		return models.User{}, fmt.Errorf("erreur de login : %v", err)
	}
	if !checkPasswordHash(user.Mdp, u.Mdp) {
		return models.User{}, fmt.Errorf("mot de passe incorrect")
	}
	return user, nil
}
