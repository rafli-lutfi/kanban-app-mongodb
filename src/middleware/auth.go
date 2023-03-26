package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rafli-lutfi/kanban-app-mongodb/src/models"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var headerType = r.Header.Get("Content-Type")
		c, err := r.Cookie("token")

		if err != nil {
			if headerType == "application/json" {
				models.ResponeWithError(w, http.StatusUnauthorized, err.Error())
				return
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
		}

		token, err := jwt.Parse(c.Value, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		// Token is either expired or not active yet
		if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			models.ResponeWithError(w, http.StatusUnauthorized, "Token is either expired or not active yet")
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)
			exp := claims["exp"].(float64)

			if time.Now().Unix() > int64(exp) {
				http.SetCookie(w, &http.Cookie{
					Name:     "token",
					Value:    "",
					MaxAge:   -1,
					Domain:   "localhost",
					HttpOnly: false,
					Secure:   true,
				})
				models.ResponeWithError(w, http.StatusUnauthorized, "token is expired")
			}
			ctx := context.WithValue(r.Context(), "id", id)

			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			models.ResponeWithError(w, http.StatusUnauthorized, "Couldn't handle this token")
		}
	})
}
