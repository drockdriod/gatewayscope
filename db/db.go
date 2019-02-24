package db

import (
    "github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"context"
	"os"
	"fmt"
)

var client *mongo.Client
var dbDetails map[string]string = make(map[string]string)
var dbContext context.Context

func Connect(ctx context.Context) (*mongo.Database, error){
	
	dbContext = ctx
	getConnectionDetails()

	var err error

	client, err = mongo.NewClient(dbDetails["dbConn"])

	if err != nil { 
		log.Fatal("error")
		log.Fatal(err) 
	}

	err = client.Connect(ctx)
	
	if err != nil { log.Fatal(err) }

	
	return client.Database(dbDetails["dbName"]), nil
}


func getConnectionDetails() {
	var dbName string = ""
	var dbConn string = ""

	env := os.Getenv("ENV")

	switch env {
	case "local":
		dbName = os.Getenv("LOCAL_DB_NAME")
		dbConn = fmt.Sprintf("mongodb://%s:%s", os.Getenv("LOCAL_DB_SERVER_URL"), os.Getenv("LOCAL_DB_PORT"))
	}

	dbDetails["dbName"] = dbName
	dbDetails["dbConn"] = dbConn

}

func GetItems(collection string, filter bson.D) ([]bson.M) {

	cur, err := client.Database(dbDetails["dbName"]).Collection(collection).Find(dbContext, filter)

	if err != nil { log.Fatal(err) }
	defer cur.Close(dbContext)

	jsonArray := []bson.M{}

	for cur.Next(dbContext) {
	   var result bson.M
	   err := cur.Decode(&result)
	   if err != nil { log.Fatal(err) }

	   log.Println(result)
	   jsonArray = append(jsonArray, result)
	   // do something with result....
	}
	if err := cur.Err(); err != nil {
	  log.Fatal(err)
	}

	return jsonArray
}

func InsertObj(collection string, jsonBody interface{}) (*mongo.InsertOneResult, error){
	log.Println("client")
	log.Println(client)
	return client.Database(dbDetails["dbName"]).Collection(collection).InsertOne(dbContext, jsonBody)
}

func Client() (*mongo.Database, context.Context, error){

	getConnectionDetails()

	if(client == nil){
		return nil, nil, fmt.Errorf("Client is not connected")
	}

	return client.Database(dbDetails["dbName"]), dbContext, nil
}