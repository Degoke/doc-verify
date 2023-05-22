package main

import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

var db *gorm.DB
var sess *session.Session

type Verification struct {
	gorm.Model
	Identifier string
	Selfie string
	Details map[string]string
	Documennt string
}

type Client struct {
	gorm.Model
	Email string
	Firstname string
	Lastname string
	Password string
	PrivateKey string
	Verifications []Verification
}

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}

	db.AutoMigrate(&Client{})
	db.AutoMigrate(&Verification{})

	sess, err = session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
	})

	if err != nil {
		panic("Failed to create AWS session")
	}

	r := gin.Default()

	r.POST("/verification", createVerificationHandler)

	r.MaxMultipartMemory = 8 << 20
	r.POST("/selfie/:verificationId", uploadSelfieHandler)

	r.Run(":8080")
}

func createVerificationHandler(c *gin.Context) {
	var verification Verification

	if err := c.ShouldBindJSON(&verification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Save(&verification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, verification)
}

func uploadSelfieHandler(c *gin.Context) {
	verificationId := c.Param("verificationId");

	log.Println(verificationId)

	file, _ := c.FormFile("file")

	log.Println(file)
	

}