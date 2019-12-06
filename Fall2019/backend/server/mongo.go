package main

import (
	"context"
	"github.com/aut-ce/Web101/models"
	"go.mongodb.org/mongo-driver/bson"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MongoUrl     = "mongodb://localhost:27017"
	DatabaseName = "kilid"
	// collections
	Houses = "houses"
	Mags   = "mags"
)

func GetClient() *mongo.Client {
	clientOptions := options.Client().ApplyURI(MongoUrl)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func GetHouse(client *mongo.Client, id string) *models.House {
	opts := options.Find().SetProjection(bson.D{{"pic.image", 0}})
	cur, err := client.Database(DatabaseName).Collection(Houses).Find(context.TODO(), bson.M{"id": id}, opts)
	if err != nil {
		log.Fatal("Error on Finding the document", id, err)
	}
	for cur.Next(context.TODO()) {
		var house models.House
		err = cur.Decode(&house)
		if err != nil {
			log.Fatal("Error on Decoding the document", err)
		}
		return &house
	}
	return nil
}

func GetOccasion(client *mongo.Client, limit int, skip int) *models.Occasion {
	var occasion models.Occasion

	opts := options.Find().
		SetSort(bson.D{{"created_at", -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetProjection(
			bson.D{
				{"parkings", 0},
				{"_id", 0},
				{"location.lat", 0},
				{"location.long", 0},
				{"breadcrumb", 0},
				{"pic.images", 0},
				{"estate.phone", 0},
			})
	cur, err := client.Database(DatabaseName).Collection(Houses).Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		log.Fatal("Error on Finding the document ", err)
	}

	for cur.Next(context.TODO()) {
		var house models.House
		err = cur.Decode(&house)
		if err != nil {
			log.Fatal("Error on Decoding the document", err)
		}
		occasion.Items = append(occasion.Items, house)
	}
	return &occasion
}

func GetAllMagazine(client *mongo.Client) *models.MagazineResponse {
	var mags models.MagazineResponse

	cur, err := client.Database(DatabaseName).Collection(Mags).Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal("Error on Finding the documents", err)
	}
	for cur.Next(context.TODO()) {
		var mag models.Magazine
		err = cur.Decode(&mag)
		if err != nil {
			log.Fatal("Error on Decoding the document", err)
		}
		mags.Items = append(mags.Items, mag)
	}
	return &mags
}
