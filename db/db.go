package db

import (
    "github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	bsonPrimitive "github.com/mongodb/mongo-go-driver/bson/primitive"
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

func FindOne(collection string, filter interface{}) (interface{}, error) {
	var v map[string]interface{}

	err := client.Database(dbDetails["dbName"]).Collection(collection).FindOne(dbContext, filter).Decode(&v)

	return v, err
}

func GetItems(collection string, filter interface{}) ([]bson.M) {

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

func InsertObj(collection string, jsonBody interface{}) (*mongo.InsertOneResult, error, bsonPrimitive.ObjectID){
	objectId := bsonPrimitive.NewObjectID()

	var bsonObj map[string]interface{}

	var body1, err = bson.Marshal(jsonBody)

	body2, err := bson.Marshal(bson.M{"_id": objectId })
	
	bson.Unmarshal(body1, &bsonObj)
	bson.Unmarshal(body2, &bsonObj)

	res, err := client.Database(dbDetails["dbName"]).Collection(collection).InsertOne(dbContext, bsonObj)

	return res, err, objectId
}

func UpdateObj(collection string,filter interface{}, jsonBody interface{}) (*mongo.UpdateResult, error){
	return client.Database(dbDetails["dbName"]).Collection(collection).UpdateOne(dbContext, filter, jsonBody)
}

// func Client() (*mongo.Database, context.Context, error){

// 	getConnectionDetails()

// 	if(client == nil){
// 		return nil, nil, fmt.Errorf("Client is not connected")
// 	}

// 	return client.Database(dbDetails["dbName"]), dbContext, nil
// }