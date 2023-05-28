package common

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DatabaseCollection string

const (
	VERIFICATION DatabaseCollection = "verification"
	USER DatabaseCollection = "user"
)

var DB *mongo.Database
var AwsSession *session.Session

var ctx = context.TODO()

func Init() *mongo.Client {
	mongoUri := GetENV("MONGO_URI")
	mongoDatabase := GetENV("MONGO_DATABASE")
	s3Region := GetENV("S#_REGION")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
        panic(err)
	}

	db := client.Database(mongoDatabase)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(s3Region),
	})

	if err != nil {
		log.Fatal("Failed to create AWS session")
	}

	AwsSession = sess
	DB = db
	return client

}

func GetDB() *mongo.Database {
	return DB
}

func GetSession() *session.Session {
	return AwsSession
}