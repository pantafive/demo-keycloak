package main

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwt"
	"net/http"
)

func RoleBaseAuthenticator(roles ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())

			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if token == nil || jwt.Validate(token) != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			realmAccess, ok := token.PrivateClaims()["realm_access"].(map[string]interface{})
			if !ok {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			data, ok := realmAccess["roles"].([]interface{})
			if !ok {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			var userRoles []string
			for _, v := range data {
				userRoles = append(userRoles, v.(string))
			}

			requireRolesCounter := len(roles)
			for _, role := range roles {
				for _, userRole := range userRoles {
					if role == userRole {
						requireRolesCounter = requireRolesCounter - 1
						break
					}
				}
			}

			if requireRolesCounter != 0 {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

}
