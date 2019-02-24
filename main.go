package main

import (
	"github.com/gin-gonic/gin"
    "github.com/drockdriod/gatewayscope/db"
	// "net/http"
	"fmt"
	"os"
    // "encoding/json"
    "context"
	_ "github.com/joho/godotenv/autoload"
    "log"
    authRouter "github.com/drockdriod/gatewayscope/routes/auth"
    // "./models"
)

type element map[string]interface{}

func main() {
	ctx := context.Background()

	v := os.Getenv("VERSION_NO")


	client, err := db.Connect(ctx)

	if err != nil { 
		log.Fatal("error")
		log.Fatal(err) 
	}

	log.Println(client)

	r := gin.Default()

	root := r.Group(fmt.Sprintf("/v%s",v))
	{
		authRouter.GetGroup(root) 
		
	}


	r.Run() // listen and serve on 0.0.0.0:8080
}