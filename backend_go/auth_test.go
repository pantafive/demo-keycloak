package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRoleBaseAuthenticator(t *testing.T) {
	t.Run("TestRoleBaseAuthenticator", func(t *testing.T) {

		// from Keycloak http://localhost:8080/auth/realms/myrealm/
		PublicKeyRS256String := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqufrtL3fmYLRO3MzWfFtTzFnmsby6R9kGvhYAb0cYamYrjCpFIJrU1cF5bGHYBtHp72RRvs40FcoUV+WhRV+29H17uyPSobgUDfinYFDaxvMz/9hfeTEP4wjqxL2qNwvko5DOJrRBvB1O8CuHE7worpQxFJMfT/c7BfvwUaIhgKoqS03LQ5Pu5USKtzm67bwNT2BM5GtHveuPPt35X8ThEqvK+24eg3xfm7s09M/J2vII9fi7e1JEuPtt4LCggjstobpekIFDepZwCG3joFJecDZqgOfyGxZuMERL86aYfuO35H/jCobOp0Xur2r3daQJShYOeC80fzJN/0V0Pd62wIDAQAB"

		// from auth response
		accessToken := "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJUMlEwS2huWlo2ckYtWWJfWDh1UE8yS0FRNWVheWVRc2xhVi1qdVVHdkpjIn0.eyJleHAiOjE2NDkzMzg5NzQsImlhdCI6MTY0OTMzODY3NCwianRpIjoiNTI5YjI2NWUtNjMyNy00MWRkLTg3YTgtNTlmM2VjZTUwOGU2IiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL2F1dGgvcmVhbG1zL215cmVhbG0iLCJzdWIiOiI1NDQwYmViMy1lNzM4LTQyNjUtOTEzZC0yODhkODM0ZTY3MzQiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJiYWNrZW5kIiwiYWNyIjoiMSIsInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJzZXJ2aWNlIl19LCJzY29wZSI6InByb2ZpbGUgZW1haWwiLCJjbGllbnRJZCI6ImJhY2tlbmQiLCJjbGllbnRIb3N0IjoiMTkyLjE2OC45Ni4xIiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJzZXJ2aWNlLWFjY291bnQtYmFja2VuZCIsImNsaWVudEFkZHJlc3MiOiIxOTIuMTY4Ljk2LjEifQ.PReoSW__n5UO_WQ8dHvMH8W1o4Lt06C8VGLXGl7n-q5LZnVoRKImOXwAe4HRecvOPoL4Qn-AUrQqejuBTxjGD40Z6UDhn52t3OgIOyzhVW7Yp36_j5zPmBawN9wLHhkAxNnt9VHlNd4bFs-sP9sbplUBUC1VVeXFlYufuDF55tQB-iHnHYqevjQ_hbeB7RUa09QdFyzKUBSTVm-WT4eNjuv6bhPGNNEj2h2hAEoAZsUsDGNayAIMKKjNFM0r3MxpAo4D-1ZXx-ei4YYZIYddwnWM_bRLt4nufb99466vBfWOa7Lh6UjzDTXxYNc3YVbOSQ_7tNxafIQxgac9H4Ic1w" //nolint:gosec

		publicKeyBlock, _ := pem.Decode([]byte("-----BEGIN PUBLIC KEY-----\n" + PublicKeyRS256String + "\n-----END PUBLIC KEY-----"))

		publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
		require.NoError(t, err)

		tokenAuth = jwtauth.New("RS256", nil, publicKey)

		decode, err := tokenAuth.Decode(accessToken)
		dump, _ := json.MarshalIndent(decode, "", "  ")
		fmt.Printf("%+v", string(dump))

		require.NoError(t, err)
		require.NotNil(t, decode)
	})
}
