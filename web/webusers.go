package web

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"webstack/metier"
	"webstack/models"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	jwt.StandardClaims
	UserEmail string `json:"useremail"`
}

const COOKIE_NAME = "cookie"
const SECRET_KEY = "codesecret123"

var now = time.Now()
var token_exp = now.Add(time.Minute * 1) // Délai d'expiration du token : 12h

type TokenInfo struct {
	CookieName string
	PrivateKey string
	Auth
}

type Auth struct {
	Name       string
	IsRequired bool
}

func WrapAuth(handler http.Handler, info TokenInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !info.Auth.IsRequired {
			handler.ServeHTTP(w, r)
			return
		}

		tokenString, err := GetTokenString(r, info.CookieName)
		if err != nil {
			if err == http.ErrNoCookie {
				http.SetCookie(w,
					&http.Cookie{
						Name:    COOKIE_NAME,
						Value:   "",
						MaxAge:  -1,
						Expires: time.Now().Add(-1 * time.Hour),
						Path:    "/",
					})
				token_exp = now.Add(time.Minute * 1)
				// fmt.Println("Reset date expiration du token : ", token_exp)
				err = fmt.Errorf("aucun utilisateur connecté, reconnectez vous")
				http.Error(w, err.Error(), http.StatusForbidden)
			}
			// http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if ok, err := isAuthenticated(tokenString, info.PrivateKey); !(ok && err == nil) {
			err = fmt.Errorf("isAuthenticated : %v", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		m := make(map[string]interface{})

		err = ParseTokenString(tokenString, &m)
		if err != nil {
			err = fmt.Errorf("ParseToken : %v", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		handler.ServeHTTP(w, r)
	}
}

func ParseTokenString(tokenString string, v *map[string]interface{}) (err error) {
	encodedStrings := strings.Split(tokenString, ".")
	if len(encodedStrings) < 2 {
		err = http.ErrNoCookie
		return
	}
	b, err := base64.RawURLEncoding.DecodeString(encodedStrings[1])
	if err != nil {
		return
	}
	return json.Unmarshal(b, v)
}

func GetTokenString(r *http.Request, cookieName string) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func isAuthenticated(tokenString string, privateKey string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(privateKey), nil
	})

	if !(token.Valid && err == nil) {
		return false, err
	}

	return true, nil
}

// var invalidatedTokens = make(map[string]bool)

// func getActiveCookieTkn(w http.ResponseWriter, r *http.Request) (string, error) {
// 	cookie, err := r.Cookie(COOKIE_NAME)
// 	if err != nil {
// 		if err == http.ErrNoCookie {
// 			err = fmt.Errorf("aucun token des les cookies : %v", err)
// 			http.Error(w, err.Error(), http.StatusUnauthorized)
// 			return "", err
// 		}
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return "", err
// 	}

// 	tokenStr := cookie.Value

// 	claims := &Claims{}

// 	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
// 		return []byte(SECRET_KEY), nil
// 	})

// 	if err != nil {
// 		if err == jwt.ErrSignatureInvalid {
// 			err = fmt.Errorf("erreur signature invalide : %v", err)
// 			http.Error(w, err.Error(), http.StatusUnauthorized)
// 			return "", err
// 		}
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 	}

// 	if !tkn.Valid {
// 		err = fmt.Errorf("erreur token invalide : %v", err)
// 		http.Error(w, err.Error(), http.StatusUnauthorized)
// 		return "", err
// 	}
// 	return tokenStr, nil
// }

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
			Name:     COOKIE_NAME,
			Value:    token,
			Expires:  token_exp,
			SameSite: http.SameSiteStrictMode,
		})
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
			Name:     COOKIE_NAME,
			Value:    token,
			Expires:  token_exp,
			SameSite: http.SameSiteStrictMode,
		})
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    COOKIE_NAME,
		Value:   "",
		MaxAge:  -1,
		Expires: time.Now().Add(-time.Hour),
		Path:    "/",
	})
}
