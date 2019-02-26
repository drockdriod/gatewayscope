package jwt

import (
    "time"
	"github.com/dgrijalva/jwt-go"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
)

var signedKey []byte
var rsaPassword string

func getJwtConfig() (map[string]interface{}) {
    privateRsaPath := os.Getenv("PRIVATE_RSA_PATH")
    rsaPassword = os.Getenv("PRIVATE_RSA_PASSWORD")

    absPath, _ := filepath.Abs(privateRsaPath)
    log.Println(absPath)

    var dat, err = ioutil.ReadFile(absPath)
    if err != nil { 
        log.Fatal("error read")
        log.Println(err) 
    }

    signedKey = dat

    mp := make(map[string]interface{})

    mp["secret"] = dat
    mp["iss"] = "GatewayScope"
    
    return mp
}

func getExp() time.Duration {
    return time.Hour * 2400
}

func GenerateToken(aud string) (string, error) {
    log.Println(aud)

    jwtConfig := getJwtConfig()
    log.Println(jwtConfig)

    // Create the token
    token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
    	"iss": jwtConfig["iss"],
    	"aud": aud,
        "exp": time.Now().Add(getExp()).Unix(),
    })


    signKey, err := jwt.ParseRSAPrivateKeyFromPEMWithPassword(signedKey, rsaPassword)

    // Sign and get the complete encoded token as a string
    tokenString, err := token.SignedString(signKey)
    return tokenString, err
}