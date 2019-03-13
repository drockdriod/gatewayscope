package jwt

import (
    "time"
    "crypto/rsa"
    "github.com/drockdriod/gatewayscope/db"
    commonModels "github.com/drockdriod/gatewayscope/models/common"
    "github.com/drockdriod/gatewayscope/models"
    "github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
    ginJWT "github.com/appleboy/gin-jwt"
    "github.com/mongodb/mongo-go-driver/bson"
    "github.com/drockdriod/gatewayscope/utils/crypto/bcrypt"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "fmt"
)

func getClientJwtConfig(clientId string) (map[string]interface{}) {
    
    privateRsaPath := fmt.Sprintf("./keys/%s/private.pem", clientId)
    publicRsaPath := fmt.Sprintf("./keys/%s/public.pem", clientId)


    mp := mapJwtConfig(privateRsaPath, publicRsaPath)
    
    return mp
}

func getServerJwtConfig() (map[string]interface{}) {
    privateRsaPath := os.Getenv("PRIVATE_RSA_PATH")
    publicRsaPath := os.Getenv("PUBLIC_RSA_PATH")

    mp := mapJwtConfig(privateRsaPath, publicRsaPath)

    return mp
}

func getExp() time.Duration {
    return time.Hour * 2400
}

func mapJwtConfig(privateRsaPath string, publicRsaPath string) (map[string]interface{}) {

    absPath, _ := filepath.Abs(privateRsaPath)
    absPathPublic, _ := filepath.Abs(publicRsaPath)
    log.Println(absPathPublic)

    var dat, err = ioutil.ReadFile(absPath)
    if err != nil { 
        log.Fatal("error read")
        log.Println(err) 
    }

    mp := make(map[string]interface{})

    mp["key"] = dat
    mp["alg"] = "RS512"
    mp["iss"] = "GatewayScope"
    mp["exp"] = time.Now().Add(getExp()).Unix()

    signedKey, err := jwt.ParseRSAPrivateKeyFromPEM(dat)

    if(err != nil){
        log.Println(err.Error()) 
        log.Fatal("error: rsa private key")
    }

    mp["signedKey"] = signedKey
    mp["privateKeyPath"] = absPath
    mp["publicKeyPath"] = absPathPublic
    
    return mp
}

func InitClientAuthMiddleware(clientId string) (*ginJWT.GinJWTMiddleware, error){
    jwtConfig := getClientJwtConfig(clientId)

    authMiddleware, err := initAuthMiddleware(jwtConfig)

    return authMiddleware, err
}


func initAuthMiddleware(jwtConfig map[string]interface{}) (*ginJWT.GinJWTMiddleware, error){
    authMiddleware, err := ginJWT.New(&ginJWT.GinJWTMiddleware{
        Realm: "Shadow Realm",
        Key: jwtConfig["key"].([]byte),
        PrivKeyFile: jwtConfig["privateKeyPath"].(string),
        PubKeyFile: jwtConfig["publicKeyPath"].(string),
        SigningAlgorithm: jwtConfig["alg"].(string),
        IdentityHandler: func(c *gin.Context) interface{} {

            claims := ginJWT.ExtractClaims(c)

            log.Println("IdentityHandler")
            log.Println(claims)
            return claims
        },
        Authenticator: func(c *gin.Context) (interface{}, error) {
            // Add database call here to check if account holder can be authenticated
            log.Println("HERE")

            var jsonBody commonModels.Login
            var user models.User
            err := c.BindJSON(&jsonBody)

            if err != nil {
                return nil, ginJWT.ErrFailedAuthentication
            }

            items := db.FindOne("users", bson.M{
                "account.email": jsonBody.Email,   
            })

            body1, err := bson.Marshal(items)

            if err != nil {
                log.Println(err.Error())
            }

            bson.Unmarshal(body1, &user)

            allowAccess := bcrypt.ComparePasswords(user.Account.HashPassword, []byte(jsonBody.Password))

            if allowAccess {


                return gin.H{
                    "token": "tokenString",
                }, nil
            } else {
                return nil, ginJWT.ErrFailedAuthentication
            }
            // var loginVals login
            // if err := c.ShouldBind(&loginVals); err != nil {
            //     return "", ginJWT.ErrMissingLoginValues
            // }
            // userID := loginVals.Username
            // password := loginVals.Password

            // if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
            //     return &User{
            //         UserName:  userID,
            //         LastName:  "Bo-Yi",
            //         FirstName: "Wu",
            //     }, nil
            // }

            return nil, ginJWT.ErrFailedAuthentication
        },
        Authorizator: func(data interface{}, c *gin.Context) bool {
            // TODO: Add database call here to get user permissions
            log.Println("Authorizator")
            log.Println(data)

            return true
        },
        Unauthorized: func(c *gin.Context, code int, message string) {
            c.JSON(code, gin.H{
                "code":    code,
                "message": "Unauthorized access",
            })
        },
    })

    if(err != nil){
        log.Fatal(err.Error())
    }

    return authMiddleware, err
}


func GenerateToken(aud string, sub ...string) (string, error) {
    log.Println("sub")
    log.Println(sub)

    jwtConfig := getClientJwtConfig(aud)
    log.Println(jwtConfig)

    // Create the token
    token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
    	"iss": jwtConfig["iss"],
    	"aud": aud,
        "sub": sub[0],
        "exp": jwtConfig["exp"].(int64),
    })

    // Sign and get the complete encoded token as a string
    tokenString, err := token.SignedString(jwtConfig["signedKey"].(*rsa.PrivateKey))
    return tokenString, err
}



