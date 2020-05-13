// Package ragger defines server work and functions to configure server.
package logistic

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"github.com/IgorRybak2055/logistic-service/internal/models"
)

// JwtAuthentication is middleware for checks JWT token in requests.
var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenHeader = r.Header.Get("Authorization")

		if tokenHeader == "" {

			err := newError(http.StatusForbidden, errors.New("missing auth token"))
			RespondError(w, err)
			return
		}

		var tk = &models.Token{}

		var token, err = jwt.ParseWithClaims(tokenHeader, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_PASSWORD")), nil
		})
		if err != nil {
			log.Println("malformed authentication token")
			nErr := newError(http.StatusForbidden, errors.New("malformed authentication token"))
			RespondError(w, nErr)
			return
		}

		if !token.Valid {
			log.Println("token is not valid")
			nErr := newError(http.StatusForbidden, errors.New("token is not valid"))
			RespondError(w, nErr)
			return
		}

		var ctx = context.WithValue(r.Context(), "user", tk.UserID)
		ctx = context.WithValue(ctx, "company", tk.CompanyID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
