package verification

import (
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/Degoke/doc-verify/common"
	"github.com/Degoke/doc-verify/user"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/textract"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateVerificationHandler(c *gin.Context) {
	var newVerification CreateVerificationValidator
	var verification VerificationModel

	userId := c.MustGet("userId");
	if userId == primitive.NilObjectID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	filter := bson.D{{Key: "_id", Value: userId}}
	user, err := user.FindOneUser(&filter)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	if err := c.ShouldBindJSON(&newVerification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	verification.ID = primitive.NewObjectID()
	verification.CreatedAt = time.Now()
	verification.UpdatedAt = time.Now()
	verification.VerificationStatus = PENDING
	verification.Values = newVerification.Values
	verification.UserId = user.ID

	verification.Save()

	c.JSON(http.StatusCreated, verification)
}

func GetVerificationsHandler(c *gin.Context) {

	userId := c.MustGet("userId");
	if userId == primitive.NilObjectID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	filter := bson.D{{Key: "_id", Value: userId}}
	user, err := user.FindOneUser(&filter)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	filter = bson.D{{Key: "userId", Value: user.ID}}

	verifications, err := FindAllVerifications(&filter)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	serializer := VerificationsSerializer{verifications}

	c.JSON(http.StatusOK, serializer.Response())
}

func UploadSelfieHandler(c *gin.Context) {
	var err error
	paramId := c.Param("id")

	verificationId, _ := primitive.ObjectIDFromHex(paramId)
	filter := bson.D{{Key: "_id", Value: verificationId}}
	_, err = FindOneVerification(&filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "verification not found"})
		return
	}

	userId := c.MustGet("userId");
	if userId == primitive.NilObjectID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	filter = bson.D{{Key: "_id", Value: userId}}
	_, err = user.FindOneUser(&filter)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	file, err := c.FormFile("selfie")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if file.Size > 1 << 20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large"})
		return
	}

	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File must be an image"})
		return
	}

	uploadedFile, err := file.Open()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
    }
    defer uploadedFile.Close()

    fileName := uuid.New().String() + filepath.Ext(file.Filename)

	session := common.GetSession()
	s3Bucket := common.GetENV("S3_BUCKET")

	uploader := s3manager.NewUploader(session)
	uploadOutput, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3Bucket),
		Key: aws.String("selfies/" + fileName),
		Body: uploadedFile,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	filter = bson.D{{Key: "_id", Value: verificationId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "selfie", Value: uploadOutput.Location}}}}

	err = UpdateVerification(&filter, &update)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	
	c.JSON(http.StatusOK, gin.H{"success": true})

}

func UploadDocumentHandler(c *gin.Context) {
	var err error
	paramId := c.Param("id")

	verificationId, _ := primitive.ObjectIDFromHex(paramId)
	filter := bson.D{{Key: "_id", Value: verificationId}}
	_, err = FindOneVerification(&filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "verification not found"})
		return
	}

	userId := c.MustGet("userId");
	if userId == primitive.NilObjectID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	filter = bson.D{{Key: "_id", Value: userId}}
	_, err = user.FindOneUser(&filter)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	file, err := c.FormFile("document")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if file.Size > 1 << 20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large"})
		return
	}

	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File must be an image"})
		return
	}

	uploadedFile, err := file.Open()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
    }
    defer uploadedFile.Close()

    fileName := uuid.New().String() + filepath.Ext(file.Filename)

	session := common.GetSession()
	s3Bucket := common.GetENV("S3_BUCKET")

	uploader := s3manager.NewUploader(session)
	uploadOutput, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3Bucket),
		Key: aws.String("documents/" + fileName),
		Body: uploadedFile,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	filter = bson.D{{Key: "_id", Value: verificationId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "document", Value: uploadOutput.Location}}}}

	err = UpdateVerification(&filter, &update)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func VerifyDetailsHandler(c *gin.Context) {
	var err error
	paramId := c.Param("id")
	verificationId, _ := primitive.ObjectIDFromHex(paramId)

	userId := c.MustGet("userId");
	if userId == primitive.NilObjectID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	filter := bson.D{{Key: "_id", Value: userId}}
	_, err = user.FindOneUser(&filter)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	filter = bson.D{{Key: "_id", Value: verificationId}}
	verification, err := FindOneVerification(&filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "verification not found"})
		return
	}

	if verification.Selfie == "" || verification.Document == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please upload document and selfie to continue"})
		return
	}

	documentUrl, err := url.Parse(verification.Document)
	if err != nil {
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "verificationStatus", Value: FAILED}}}}
		err := UpdateVerification(&filter, &update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	documentPath := strings.TrimLeft(documentUrl.Path, "/")

	session := common.GetSession()
	s3Bucket := common.GetENV("S3_BUCKET")

	textractClient := textract.New(session)
	textractInput := &textract.DetectDocumentTextInput{
		Document: &textract.Document{
			S3Object: &textract.S3Object{
				Bucket: aws.String(s3Bucket),
				Name: aws.String(documentPath),
			},
		},
	}

	textractOutput, err := textractClient.DetectDocumentText(textractInput)

	if err != nil {
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "verificationStatus", Value: FAILED}}}}
		err := UpdateVerification(&filter, &update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, details := checkTextMatchesDetails(textractOutput, verification)

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "verificationResponse", Value: details}}}}
		err = UpdateVerification(&filter, &update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	isFaceMatch := false

	selfieUrl, err := url.Parse(verification.Selfie)
	if err != nil {
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "verificationStatus", Value: FAILED}}}}
		err := UpdateVerification(&filter, &update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	selfiePath := strings.TrimLeft(selfieUrl.Path, "/")

	rekognitionClient := rekognition.New(session)
	rekognitionInput := &rekognition.CompareFacesInput{
		SourceImage: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: aws.String(s3Bucket),
				Name: aws.String(selfiePath),
			},
		},
		TargetImage: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: aws.String(s3Bucket),
				Name: aws.String(documentPath),
			},
		},
	}

	rekognitionOutput, err := rekognitionClient.CompareFaces(rekognitionInput)

	if err != nil {
		update = bson.D{{Key: "$set", Value: bson.D{{Key: "verificationStatus", Value: FAILED}}}}
		err := UpdateVerification(&filter, &update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}


	if len(rekognitionOutput.FaceMatches) == 0 {
		isFaceMatch = false

	} else {
		isFaceMatch = true
	}

	details["isFaceMatch"] = isFaceMatch

	update = bson.D{{Key: "$set", Value: bson.D{{Key: "verificationResponse", Value: details}, {Key: "verificationStatus", Value: COMPLETED}}}}
		err = UpdateVerification(&filter, &update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "data": details})
	
}

func checkTextMatchesDetails(textractOutput *textract.DetectDocumentTextOutput, verification VerificationModel) (bool, map[string]bool) {
	details := make(map[string]bool)
	results := make(map[string]bool)

	for _, value := range verification.Values {
		details[common.NormalizeString(value)] = false
	}
	
	// valuesSlice := strings.Split(verification.Values, ":")

	// for i := 0; i < len(valuesSlice)-1; i++ {
	// 	details[normalizeString(valuesSlice[i])] = false
	// }

	for _, item := range textractOutput.Blocks {
		if *item.BlockType == "LINE" {
			line := common.NormalizeString(*item.Text)

			for detail := range details {
				if strings.Contains(line, detail) {
					details[detail] = true
				}
			}
		}
	}
	
	for key, value := range verification.Values {
		results[key] = details[value]
	}

	// for _, found := range details {
	// 	if !found {
	// 		return false, details
	// 	}
	// }

	return true, results
}