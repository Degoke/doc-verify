package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUserHandler(c *gin.Context) {
	var user UserModel
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	emailFilter := bson.D{{Key: "email", Value: user.Email}}
	_, err := FindOneUser(&emailFilter)

	if err != nil {
		if err != mongo.ErrNoDocuments {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
			return
		}
	}

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Emasil already in use"})
			return
	}

	user.SetPassword(user.Password)
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	user.Save()

	serializer := UserSerializer{user}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": serializer.Response()})
}

func LoginUserHandler(c *gin.Context) {
	var login LoginValidator
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	filter := bson.D{{Key: "email", Value: login.Email}}
	user, err := FindOneUser(&filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "user does not exist"})
			return
		} else {
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
				return
		}
	}

	err = user.CheckPassword(login.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	serializer := LoginSerializer{user}
	
	c.JSON(http.StatusOK, gin.H{"success": true, "data": serializer.Response()})

}

func GetUsersHandler(c *gin.Context) {
	allFilter := bson.D{}

	users, err := FindAllUsers(&allFilter)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	serializer := UsersSerializer{users}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": serializer.Response()})
}

func GetUserHandler(c *gin.Context) {
	userId := c.MustGet("userId");
	if userId == primitive.NilObjectID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	filter := bson.D{{Key: "_id", Value: userId}}
	user, err := FindOneUser(&filter)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	serializer := UserSerializer{user}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": serializer.Response()})
}
