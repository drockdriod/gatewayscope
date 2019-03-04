# Gateway Scope
Permission based API for managing external user base groups, built in Go lang

## Resources:
- github.com/gin-gonic/gin
- golang.org/x/crypto/bcrypt
- github.com/dgrijalva/jwt-go
- github.com/appleboy/gin-jwt

## Important Notes:

### Gin JWT 
This project uses the library gin-gwt (github.com/appleboy/gin-jwt). I had to modify the library to include a alternative 
MiddlewareFunc() function called MiddlewareFunc() in `auth_jwt.go`:
```golang
// MiddlewareFunc makes GinJWTMiddleware implement the Middleware interface with gin.Context as an argument for 
// flexibility 
func (mw *GinJWTMiddleware) MiddlewareFuncAlt(c *gin.Context) {
	mw.middlewareImpl(c)
}
```
This newly created function allows me to decode the JWT for the purpose of grabbing the clientId from the token. 
The clientId would be used to instantiate the gin-jwt middleware with the correct RSA public/private keys.

### Crypto
There was some issues with generating dynamic RSA private/ public PEMs. Generating a public PEM based on the private predicesor resulted in the public key to read as "invalid" upon validation with the JWT token. I used the https://golang.org/pkg/crypto/x509 library, specifically the MarshalPKCS1PublicKey ¶ function, to perform the public PEM export.

It turns out that function was not returning the result I needed:
https://github.com/golang/go/issues/29141

The ParsePKIXPublicKey function was what was needed to export the corresponding format to the private PEM file.
