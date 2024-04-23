package web

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"webstack/metier/users"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	jwt.StandardClaims
	UserEmail string `json:"useremail"`
}

const CookieName = "cookie"
const SecretKey = "codesecret123"
const ErrNoCookie = "connexion expirée"
const ErrInvToken = "token invalide"
const MinSubTokenStr = 2

var tokenExp time.Time // Délai d'expiration du token : 1h

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

func HandleSignin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmpassword := r.FormValue("confirmpassword")

	user, err := users.Signin(email, password, confirmpassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := jsonwebToken(user)
	http.SetCookie(w,
		&http.Cookie{
			Name:     CookieName,
			Value:    token,
			Expires:  tokenExp,
			SameSite: http.SameSiteStrictMode,
		})
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := users.Login(email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := jsonwebToken(user)
	http.SetCookie(w,
		&http.Cookie{
			Name:     CookieName,
			Value:    token,
			Expires:  tokenExp,
			SameSite: http.SameSiteStrictMode,
		})
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	tokenStr, err := getTokenString(r, CookieName)
	if err != nil {
		err = fmt.Errorf(ErrInvToken)
		http.Error(w, err.Error(), http.StatusForbidden)

		return
	}

	invalidateToken(tokenStr)
	http.SetCookie(w, &http.Cookie{
		Name:    CookieName,
		Value:   "",
		MaxAge:  -1,
		Expires: time.Now().Add(-time.Hour),
		Path:    "/",
	})
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
						Name:    CookieName,
						Value:   "",
						MaxAge:  -1,
						Expires: time.Now().Add(-1 * time.Hour),
						Path:    "/",
					})

				err = fmt.Errorf(ErrNoCookie)
				http.Error(w, err.Error(), http.StatusForbidden)
			} else {
				err = fmt.Errorf("err getToken")
				http.Error(w, err.Error(), http.StatusForbidden)
			}

			return
		}

		if !validateToken(tokenStr) {
			err = fmt.Errorf(ErrInvToken)
			http.Error(w, err.Error(), http.StatusForbidden)

			return
		}

		if ok, errAuth := isAuthenticated(tokenStr, info.PrivateKey); !ok || errAuth != nil {
			http.Error(w, errAuth.Error(), http.StatusUnauthorized)
			return
		}

		m := make(map[string]any)

		err = parseTokenString(tokenStr, m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(w, r)
	}
}

func parseTokenString(tokenStr string, v map[string]any) (err error) {
	encodedStrings := strings.Split(tokenStr, ".")
	if len(encodedStrings) < MinSubTokenStr {
		err = http.ErrNoCookie
		return err
	}

	b, err := base64.RawURLEncoding.DecodeString(encodedStrings[1])
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &v)
}

func getTokenString(r *http.Request, cookieName string) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func isAuthenticated(tokenStr, privateKey string) (bool, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
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

func jsonwebToken(u users.User) string {
	tokenExp = time.Now().Add(time.Hour * 1)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExp.Unix(),
		},
		UserEmail: users.GetEmail(u),
	})

	token, err := t.SignedString([]byte(SecretKey))
	if err != nil {
		fmt.Print(err)
	}

	return token
}

func getUserEmail(tokenStr string) string {
	claims := &Claims{}
	jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		return []byte(SecretKey), nil
	})

	return claims.UserEmail
}
