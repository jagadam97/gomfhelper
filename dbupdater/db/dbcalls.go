package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jagadam97/dbupdater/mfdata"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MyMFS struct {
	GrowLink   string      `json:"growLink"`
	MFAPICode  string      `json:"mfapicode"`
	MFCode     string      `json:"mfcode"`
	MFName     string      `json:"mfname"`
}

type MfNavHistory struct {
	MFAPICode  string      `bson:"mfapicode"`
	DATE       string      `bson:"date"`
	NAV        float64     `bson:"nav"`
}

const uri = "mongodb://localhost:27017/"

func GetWatchedMfs() []string {
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

	var results []string
	for cur.Next(context.TODO()) {
		var elem MyMFS
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem.MFAPICode)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

    return results
}

func UpdateNavHistory(apicode string) {
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

	mfdata := mfdata.GetMFData(apicode)


	coll := client.Database("db").Collection("mfnavhistory")

	for _, data := range mfdata {

		filter := bson.D{
		{"date", data.Date},
		{"mfapicode", apicode},
		}

		var existing MfNavHistory
		err := coll.FindOne(context.TODO(), filter).Decode(&existing)
		

		if err == mongo.ErrNoDocuments {
			document := bson.D{
				{"mfapicode", apicode},
				{"date", data.Date},
				{"nav", data.Nav},
			}

			_, err = coll.InsertOne(context.TODO(), document)
			if err != nil {
				log.Fatalf("Failed to update or insert NAV history for %s on %s: %v", apicode, data.Date, err)
			}
		}
	}

	fmt.Println("NAV history updated successfully")
}


func DeleteNavData() {
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


	coll := client.Database("db").Collection("mfnavhistory")
	result, err := coll.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Number of documents deleted: %d\n", result.DeletedCount)

}

func UpdateNavHistoryWatched() {
	DeleteNavData()
	for _, mutualFund := range GetWatchedMfs(){
		UpdateNavHistory(mutualFund)
	}
}

func UpdateLatestNav() {
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

	watchedMFs := GetWatchedMfs()
	coll := client.Database("db").Collection("mfnavhistory")
	for _, MfapiCode := range watchedMFs {
		date, nav := mfdata.GetLatestNav(MfapiCode)
		filter := bson.D{
			{"date", date},
			{"mfapicode", MfapiCode},
			}
	
			var existing MfNavHistory
			err := coll.FindOne(context.TODO(), filter).Decode(&existing)
			if err == mongo.ErrNoDocuments {
				document := bson.D{
					{"mfapicode", MfapiCode},
					{"date", date},
					{"nav", nav},
				}
	
				_, err = coll.InsertOne(context.TODO(), document)
				if err != nil {
					log.Fatalf("Failed to update or insert NAV history for %s on %s: %v", MfapiCode, date, err)
				}
			}

	}
}