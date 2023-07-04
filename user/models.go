package user

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/Degoke/doc-verify/common"
)

type UserModel struct {
	ID primitive.ObjectID `bson:"_id"`
	Email string `bson:"email"`
	Password string `bson:"password"`
	Key string `bson:"key"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (u *UserModel) SetPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password should not be empty")
	}

	bytePassword := []byte(password)

	passwordHash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)

	if err != nil {
		log.Println("Model Error: (user:setPassword)", err)
		return errors.New("an error occured")
	}
	u.Password = string(passwordHash)
	return nil
}

func (u *UserModel) CheckPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.Password)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

func (u *UserModel) Save() error {
	db := common.GetDB()
	userCollection := db.Collection("user")
	ctx := context.TODO()

	_, err := userCollection.InsertOne(ctx, &u)

	return err
}

func FindOneUser(filter *primitive.D) (UserModel, error) {
	db := common.GetDB()
	userCollection := db.Collection("user")
	ctx := context.TODO()

	var user UserModel

	err := userCollection.FindOne(ctx, filter).Decode(&user)

	return user, err
}

func FindAllUsers(filter *primitive.D) ([]UserModel, error) {
	db := common.GetDB()
	userCollection := db.Collection("user")
	ctx := context.TODO()

	var users []UserModel

	cursor, err := userCollection.Find(ctx, filter)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return users, nil
		} else {
			return nil, err
		}
	}

	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	for _, user := range users {
		cursor.Decode(&user)
		_, err := json.MarshalIndent(user, "", "    ")
		if err != nil {
			return nil, err
		}
	}

	return users, nil
}