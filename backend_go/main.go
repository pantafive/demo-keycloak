package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var tokenAuth *jwtauth.JWTAuth

func init() {

	keycloakBaseUrl := os.Getenv("KEYCLOAK_BASE_URL")
	if keycloakBaseUrl == "" {
		keycloakBaseUrl = "http://localhost:8080"
	}

	var resp *http.Response
	for repeatCount := 10; ; repeatCount-- {
		var err error
		url := fmt.Sprintf("%s/auth/realms/myrealm/", keycloakBaseUrl)
		resp, err = http.Get(url)
		if err != nil && repeatCount == 0 {
			log.Fatal(err)
		}
		if err != nil {
			log.Printf("Warning: %v", err)
			time.Sleep(3 * time.Second) //nolint:gomnd
			continue
		}
		break
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	realm := struct {
		Realm           string `json:"realm"`
		PublicKey       string `json:"public_key"`
		TokenService    string `json:"token-service"`
		AccountService  string `json:"account-service"`
		TokensNotBefore int    `json:"tokens-not-before"`
	}{}
	if err := json.Unmarshal(body, &realm); err != nil {
		log.Fatal(err)
	}

	publicKeyBlock, _ := pem.Decode([]byte("-----BEGIN PUBLIC KEY-----\n" + realm.PublicKey + "\n-----END PUBLIC KEY-----"))

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	tokenAuth = jwtauth.New("RS256", nil, publicKey)
	log.Println("Token auth initialized")
}

func main() {
	addr := ":3333"
	fmt.Printf("Starting server on %v\n", addr)
	http.ListenAndServe(addr, router())
}

func router() http.Handler {
	r := chi.NewRouter()

	// Protected routes
	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))

		// Handle valid / invalid tokens. In this example, we use
		// the provided authenticator middleware, but you can write your
		// own very easily, look at the Authenticator method in jwtauth.go
		// and tweak it, its not scary.
		r.Use(jwtauth.Authenticator)

		r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["user_id"])))
		})
	})

	// RoleBased routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(RoleBaseAuthenticator("service"))

		r.Get("/service", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("welcome service"))
		})
	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("welcome anonymous"))
		})
	})

	return r
}
