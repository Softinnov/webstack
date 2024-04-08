package web

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"webstack/metier"
	"webstack/models"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	jwt.StandardClaims
	UserEmail string `json:"useremail"`
}

const SECRET_KEY = "codesecret123"

var now = time.Now()
var token_exp = now.Add(time.Hour * 12) // Délai d'expiration du token : 12h

// var invalidatedTokens = make(map[string]bool)

func getActiveCookieTkn(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("cookie")
	if err != nil {
		if err == http.ErrNoCookie {
			err = fmt.Errorf("aucun token des les cookies : %v", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return ""
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return ""
	}

	tokenStr := cookie.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			err = fmt.Errorf("erreur signature invalide : %v", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return ""
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if !tkn.Valid {
		err = fmt.Errorf("erreur token invalide : %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return ""
	}
	return tokenStr
}

// func invalidateToken(token string) {
// 	invalidatedTokens[token] = true
// }

// func validateToken(token string) bool {
// 	_, invalidated := invalidatedTokens[token]
// 	return !invalidated
// }

// func checkExpiredTokens() error {
// 	for tokenStr := range invalidatedTokens {
// 		if !validateToken(tokenStr) {
// 			tknclaims, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
// 				return []byte(SECRET_KEY), nil
// 			})
// 			if err != nil {
// 				return fmt.Errorf("%v : %v", err, http.StatusInternalServerError)
// 			}

// 			claims, ok := tknclaims.Claims.(*Claims)
// 			if !ok || !tknclaims.Valid {
// 				return fmt.Errorf("%v : %v", err, http.StatusInternalServerError)
// 			}
// 			fmt.Println(claims.ExpiresAt, now.Unix())
// 			if claims.ExpiresAt <= now.Unix() {
// 				invalidateToken(tokenStr)
// 			}
// 		}
// 	}
// 	return nil
// }

func jsonwebToken(u models.User) string {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: token_exp.Unix(),
		},
		UserEmail: u.Email,
	})
	token, err := t.SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Fatalln("error signedstring :", err)
	}
	return token
}

func getUserEmail(tokenStr string) string {
	claims := &Claims{}
	jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	return claims.UserEmail
}

func HandleSignin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmpassword := r.FormValue("confirmpassword")

	user, err := metier.AddUser(email, password, confirmpassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token := jsonwebToken(user)
	http.SetCookie(w,
		&http.Cookie{
			Name:    "cookie",
			Value:   token,
			Expires: token_exp,
		})
	// fmt.Printf("Utilisateur enregistré : %v", getUserEmail(token))
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := metier.Login(email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token := jsonwebToken(user)
	http.SetCookie(w,
		&http.Cookie{
			Name:     "cookie",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour),
			SameSite: http.SameSiteLaxMode,
		})
	http.SetCookie(w,
		&http.Cookie{
			Name:     "cookie",
			Value:    token,
			Expires:  token_exp,
			SameSite: http.SameSiteLaxMode,
		})
	// fmt.Println("handle login ok", invalidatedTokens)
	// fmt.Println(r.Cookie("cookie"))
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	// invalidateToken(getActiveCookieTkn(w, r))
	http.SetCookie(w, &http.Cookie{
		Name:    "cookie",
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
		Path:    "/",
	})
	// fmt.Println(invalidatedTokens)
}
