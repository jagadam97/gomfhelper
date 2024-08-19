package apis

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mymfs struct {
	Mfapicode string `json:"mfapicode"`
	Mfcode    string `json:"mfcode"`
	Mfname    string `json:"mfname"`
	Growlink  string `json:"growLink"`
}

const uri = "mongodb://localhost:27017/"

func GetMutualFunds(c *gin.Context) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	
	coll := client.Database("db").Collection("mymfs")
	findOptions := options.Find()
	cur, err := coll.Find(context.TODO(), bson.D{{}}, findOptions)

	if err != nil {
		log.Fatal(err)
	}
	
	defer cur.Close(context.TODO())

	var results []mymfs
	for cur.Next(context.TODO()) {
		var elem mymfs
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem)
	}
	fmt.Println(cur)

    c.IndentedJSON(http.StatusOK, results)
}

func AddMutualFund(c *gin.Context) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	
	coll := client.Database("db").Collection("mymfs")
	findOptions := options.Find()
	cur, err := coll.Find(context.TODO(), bson.D{{}}, findOptions)

	if err != nil {
		log.Fatal(err)
	}
	
	defer cur.Close(context.TODO())

	var results []mymfs
	for cur.Next(context.TODO()) {
		var elem mymfs
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem)
	}
	fmt.Println(cur)
    c.IndentedJSON(http.StatusOK, results)
}
    