package middleware

import (
	"encoding/json"
	"errors"
	"github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"

	"net/http"
)

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func CreateAuth0MiddleWare(aud, iss string) *jwtmiddleware.JWTMiddleware{
	var t jwt.Keyfunc = func(token *jwt.Token) (interface{}, error) {
		// Verify 'aud' claim

		//This section has been ignored due to a bug in the libary, and is insecure
		//aud := "http://192.168.0.28:3500"
		//checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
		//if !checkAud {
		//	return token, errors.New("Invalid audience.")
		//}
		// Verify 'iss' claim
		iss := "https://dev-e5s8h580.us.auth0.com/"
		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
		if !checkIss {
			return token, errors.New("Invalid issuer.")
		}

		cert, err := getPemCert(iss,token)
		if err != nil {
			panic(err.Error())
		}
		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	}

	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter:  t,
		SigningMethod: jwt.SigningMethodRS256,
	})

}

func getPemCert(iss string, token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get(iss + ".well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}
