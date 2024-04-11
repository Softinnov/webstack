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

const ERR_NOTAUTH = "aucun utilisateur connecté"
const ERR_INVTOKEN = "token invalide"

var token_exp = time.Now().Add(time.Hour * 12) // Délai d'expiration du token : 12h

var invalidatedTokens = make(map[string]bool)

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

		tokenStr, err := getTokenString(r, info.CookieName)
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
				token_exp = time.Now().Add(time.Hour * 12)
				err = fmt.Errorf(ERR_NOTAUTH)
				http.Error(w, err.Error(), http.StatusForbidden)
			}
			return
		}

		if !validateToken(tokenStr) {
			err = fmt.Errorf(ERR_INVTOKEN)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if ok, err := isAuthenticated(tokenStr, info.PrivateKey); !(ok && err == nil) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		m := make(map[string]interface{})

		err = parseTokenString(tokenStr, &m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		handler.ServeHTTP(w, r)
	}
}

func parseTokenString(tokenStr string, v *map[string]interface{}) (err error) {
	encodedStrings := strings.Split(tokenStr, ".")
	if len(encodedStrings) < 2 {
		err = http.ErrNoCookie
		return
	}
	b, err := base64.RawURLEncoding.DecodeString(encodedStrings[1])
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func getTokenString(r *http.Request, cookieName string) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func isAuthenticated(tokenStr string, privateKey string) (bool, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
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

func invalidateToken(tokenStr string) {
	invalidatedTokens[tokenStr] = true
}

func validateToken(tokenStr string) bool {
	_, invalidated := invalidatedTokens[tokenStr]
	return !invalidated
}

func jsonwebToken(u models.User) string {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: token_exp.Unix(),
		},
		UserEmail: u.Email,
	})
	token, err := t.SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Fatalln(err)
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
	tokenStr, err := getTokenString(r, COOKIE_NAME)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	invalidateToken(tokenStr)
	http.SetCookie(w, &http.Cookie{
		Name:    COOKIE_NAME,
		Value:   "",
		MaxAge:  -1,
		Expires: time.Now().Add(-time.Hour),
		Path:    "/",
	})
}
