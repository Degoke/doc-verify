package verification

import "github.com/gin-gonic/gin"

func VerificationRouter(r *gin.RouterGroup) {
	r.POST("/", CreateVerificationHandler)
	r.POST("/:id/verify", VerifyDetailsHandler)
	r.GET("/", GetVerificationsHandler)
}

func UploadRouter(r *gin.RouterGroup) {
	r.PATCH("/:id/selfie", UploadSelfieHandler)
	r.PATCH("/:id/document", UploadDocumentHandler)
}