package token

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth"
	helper "github.com/jocelyntjahyadi/tsaving/helpers"
)

type JWT struct {
	*jwtauth.JWTAuth
}

func New(secret []byte) *JWT {
	//secret is a secret keyy, bebas mau apa aja. mustinya yang rahasia.
	return &JWT{jwtauth.New("HS256", []byte("secret"), nil)}
}

func (j *JWT) Encode(token Token) string {
	_, tokenString, _ := j.JWTAuth.Encode(&token)
	return tokenString
}

//untuk ngecek konteks yang disimpen
func (j *JWT) GetToken(r *http.Request) Token {
	token, ok := r.Context().Value("token").(Token) // r.Context() -> variable per request
	//ngambil konteks dari request, dan cari value yang keynya token yang didefine dibawah
	// (Token) -> convert jadi objek token
	//

	if !ok {
		return Token{}
	}
	return token
}

func (j *JWT) AuthMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		jwtToken, err := jwtauth.VerifyRequest(j.JWTAuth, r, TokenFromHeader)

		if err != nil {
			msg := fmt.Sprintf("error verfying token : %v", err)
			helper.HTTPError(w, http.StatusBadRequest, msg)
			return
		}

		var claims Token
		b, err := json.Marshal(jwtToken.Claims) // tipe datanya beda jadi di marshal, terus di unmarshal.
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "invalid token %#v", jwtToken.Raw)
			return
		}

		err = json.Unmarshal(b, &claims)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unable to parse token", jwtToken.Raw)
			return
		}

		err = claims.Valid() // kalau expired
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, err.Error()) //lempar status error.
			return
		}

		// r.Context() = {request_id: 00000}
		//

		ctx := context.WithValue(r.Context(), "token", claims)
		// Context value mau menandakan objek baru dalam konteks yang sudah ada
		// kalau sudah ada dia bikin token baru ?
		// r context akan nge return token sama tapi konteks yang berbeda
		// konteks yang sebelumnya dengan tambahan token claims.

		//kalau engga expired, bisa dapetin usernya di claims.user.

		// fmt.Fprintf(w, "Login as %v\n", claims.User)
		handler.ServeHTTP(w, r.WithContext(ctx))

	})
}

func TokenFromHeader(r *http.Request) string {
	return r.Header.Get("Authorization")
}
