package auth

import (
	"context"
	"crypto/subtle"
	"net/http"
)

type ClaimsKey string

const (
	ClaimsKeyUserID   ClaimsKey = "user_id"
	ClaimsKeyVerified ClaimsKey = "verified"
)

func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			writeResponse(w, Response{
				Message: "Authentication failed",
			}, http.StatusUnauthorized)
			return
		}

		creds, credUserOk := db[user]
		if !credUserOk || subtle.ConstantTimeCompare([]byte(pass), []byte(creds.Password)) != 1 {
			writeResponse(w, Response{
				Message: "Authentication failed",
			}, http.StatusUnauthorized)
			return
		} else {
			ctx := r.Context()
			ctx = context.WithValue(ctx, ClaimsKeyUserID, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
