package tokens

import(
	"net/http" 
	"github.com/go-chi/jwtauth"
	"encoding/json"
	"github.com/ndv6/tsaving/helpers"
	"context"
)

type JWT struct {
	*jwtauth.JWTAuth
}

func New(secret []byte) *JWT{
	return &JWT{jwtauth.New("HS256", []byte("secret"), nil)}
}

func (j *JWT) Endcode(token Token) string{
	_, tokenString, _ := j.JWTAuth.Encode(&token)
	return tokenString
}

func (j *JWT) GetToken(r *http.Request) Token {
	token, ok := r.Context().Value("token").(Token)
	if !ok {
		return Token{}
	}
	return token
}

func (j *JWT) AuthMiddleware(handler http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		jwtToken, err := jwtauth.VerifyRequest(j.JWTAuth, r, TokenFromHeader)
		if err != nil{
			helpers.HTTPError(w, http.StatusBadRequest, "Error Verifying Token")
			return
		}

		var claims Token
		b, err := json.Marshal(jwtToken.Claims) //Encode Token
		if err != nil{
			helpers.HTTPError(w, http.StatusBadRequest, "Invalid Token")
			return
		}

		err = json.Unmarshal(b, &claims)
		if err != nil{
			helpers.HTTPError(w, http.StatusBadRequest, "Unable to Parse Token")
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

func TokenFromHeader(r *http.Request) string {
	return r.Header.Get("Authorization")
}