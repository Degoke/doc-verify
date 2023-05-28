package verification

import "github.com/gin-gonic/gin"

func VerificationRouter(r *gin.RouterGroup) {
	r.POST("/", CreateVerificationHandler)
	r.POST("/:id/verify", VerifyDetailsHandler)
}

func VerificationsRouter(r *gin.RouterGroup) {
	r.GET("/", GetVerificationsHandler)
}

func UploadRouter(r *gin.RouterGroup) {
	r.PATCH("/:id/selfie", UploadSelfieHandler)
	r.PATCH("/:id/document", UploadDocumentHandler)
}