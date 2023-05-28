package verification

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Degoke/doc-verify/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Status string
const (
	PENDING Status = "pending"
	COMPLETED Status = "completed"
	FAILED Status = "failed"
)

type VerificationModel struct {
	ID primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	Selfie string `bson:"selfie"`
	Values map[string]string `bson:"values"`
	Document string `bson:"document"`
	ExtractedText string `bson:"extractedText"`
	VerificationStatus Status `bson:"verificationStatus"`
	VerificationResponse map[string]bool `bson:"verificationResponse"`
	UserId primitive.ObjectID `bson:"userId"`
}

func (u *VerificationModel) Save() error {
	db := common.GetDB()
	verificationCollection := db.Collection("verification")
	ctx := context.TODO()

	_, err := verificationCollection.InsertOne(ctx, &u)

	return err
}

func UpdateVerification(filter *primitive.D, update *primitive.D) error {
	db := common.GetDB()
	verificationCollection := db.Collection("verification")
	ctx := context.TODO()

	_, err := verificationCollection.UpdateOne(ctx, filter, update)

	return err
}

func FindOneVerification(filter *primitive.D) (VerificationModel, error) {
	db := common.GetDB()
	verificationCollection := db.Collection("verification")
	ctx := context.TODO()

	var verification VerificationModel

	err := verificationCollection.FindOne(ctx, filter).Decode(&verification)

	return verification, err
}

func FindAllVerifications(filter *primitive.D) ([]VerificationModel, error) {
	db := common.GetDB()
	verificationCollection := db.Collection("verification")
	ctx := context.TODO()

	var verifications []VerificationModel

	cursor, err := verificationCollection.Find(ctx, filter)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return verifications, nil
		} else {
			return nil, err
		}
	}

	if err := cursor.All(ctx, &verifications); err != nil {
		return nil, err
	}

	for _, verification := range verifications {
		cursor.Decode(&verification)
		_, err := json.MarshalIndent(verification, "", "    ")
		if err != nil {
			return nil, err
		}
	}

	return verifications, nil
}