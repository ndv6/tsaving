package tokens

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ndv6/tsaving/constants"

	"github.com/go-chi/jwtauth"
	"github.com/ndv6/tsaving/helpers"
)

type JWT struct {
	*jwtauth.JWTAuth
}

func New(secret []byte) *JWT {
	return &JWT{jwtauth.New("HS256", []byte("secret"), nil)}
}

func (j *JWT) Encode(token Token) string {
	_, tokenString, _ := j.JWTAuth.Encode(&token)
	return tokenString
}

func (j *JWT) Decode(tokenString string) (token Token, err error) {
	jwtToken, err := j.JWTAuth.Decode(tokenString)
	fmt.Println(jwtToken.Claims)
	if err != nil {
		return
	}
	return
}

func (j *JWT) GetToken(r *http.Request) Token {
	token, ok := r.Context().Value("token").(Token)
	if !ok {
		return Token{}
	}
	return token
}

func (j *JWT) GetTokenAdmin(r *http.Request) TokenAdmin {
	token, ok := r.Context().Value("token").(TokenAdmin)
	if !ok {
		return TokenAdmin{}
	}
	return token
}

func (j *JWT) AuthAdminMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(constants.ContentType, constants.Json)
		jwtTokenAdmin, err := jwtauth.VerifyRequest(j.JWTAuth, r, TokenFromHeader)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, "Error Verifying Token")
			return
		}

		var claimsAdmin TokenAdmin
		b, err := json.Marshal(jwtTokenAdmin.Claims) //Encode Token
		if err != nil {
			helpers.HTTPError(w, http.StatusUnauthorized, "Invalid Token")
			return
		}

		err = json.Unmarshal(b, &claimsAdmin)
		if err != nil {
			helpers.HTTPError(w, http.StatusUnauthorized, "Unable to Parse Token")
			return
		}

		if len(claimsAdmin.Username) < 1 {
			helpers.HTTPError(w, http.StatusUnauthorized, "Invalid User")
			return
		}

		err = claimsAdmin.Valid()
		if err != nil {
			helpers.HTTPError(w, http.StatusUnauthorized, "Token Not Valid")
			return
		}

		ctx := context.WithValue(r.Context(), "token", claimsAdmin)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (j *JWT) AuthMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(constants.ContentType, constants.Json)
		jwtToken, err := jwtauth.VerifyRequest(j.JWTAuth, r, TokenFromHeader)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, "Error Verifying Token")
			return
		}

		var claims Token
		b, err := json.Marshal(jwtToken.Claims) //Encode Token
		if err != nil {
			helpers.HTTPError(w, http.StatusUnauthorized, "Invalid Token")
			return
		}

		err = json.Unmarshal(b, &claims)
		if err != nil {
			helpers.HTTPError(w, http.StatusUnauthorized, "Unable to Parse Token")
			return
		}

		if len(claims.AccountNum) < 1 {
			helpers.HTTPError(w, http.StatusUnauthorized, "Invalid User")
			return
		}

		err = claims.Valid()
		if err != nil {
			helpers.HTTPError(w, http.StatusUnauthorized, "Token Expired")
			return
		}

		ctx := context.WithValue(r.Context(), "token", claims)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (j *JWT) ValidateAccount(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(constants.ContentType, constants.Json)
		jwtToken, err := jwtauth.VerifyRequest(j.JWTAuth, r, TokenFromHeader)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, "Error Verifying Token")
			return
		}

		var claims Token
		b, err := json.Marshal(jwtToken.Claims)
		if err != nil {
			helpers.HTTPError(w, http.StatusUnauthorized, "Invalid Token")
			return
		}

		err = json.Unmarshal(b, &claims)
		if err != nil {
			helpers.HTTPError(w, http.StatusUnauthorized, "Unable to Parse Token")
			return
		}

		if claims.AccountExpiration.Before(time.Now()) {
			helpers.HTTPError(w, http.StatusBadRequest, "Card expired, please renew it")
			return
		}
		handler.ServeHTTP(w, r)
		return
	})
}

func TokenFromHeader(r *http.Request) string {
	return r.Header.Get("Authorization")
}
