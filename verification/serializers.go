package verification

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VerificationSerializer struct {
	Verification VerificationModel
}

type VerificationsSerializer struct {
	Verifications []VerificationModel
}

type VerificationResponse struct {
	ID primitive.ObjectID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Selfie string `json:"selfie"`
	Values map[string]string `json:"values"`
	Document string `json:"document"`
	ExtractedText string `json:"extractedText"`
	VerificationStatus Status `json:"verificationStatus"`
	VerificationResponse map[string]bool `json:"verificationResponse"`
	UserId primitive.ObjectID `json:"userId"`
}

func (v *VerificationSerializer) Response() VerificationResponse {
	return VerificationResponse{
		ID: v.Verification.ID,
		CreatedAt: v.Verification.CreatedAt,
		UpdatedAt: v.Verification.UpdatedAt,
		Selfie: v.Verification.Selfie,
		Values: v.Verification.Values,
		Document: v.Verification.Document,
		ExtractedText: v.Verification.ExtractedText,
		VerificationStatus: v.Verification.VerificationStatus,
		VerificationResponse: v.Verification.VerificationResponse,
		UserId: v.Verification.UserId,
	}
}

func (v *VerificationsSerializer) Response() []VerificationResponse {
	var verifications []VerificationResponse
	for _, verification := range v.Verifications {
		verifications = append(verifications, VerificationResponse{
			ID: verification.ID,
			CreatedAt: verification.CreatedAt,
			UpdatedAt: verification.UpdatedAt,
			Selfie: verification.Selfie,
			Values: verification.Values,
			Document: verification.Document,
			ExtractedText: verification.ExtractedText,
			VerificationStatus: verification.VerificationStatus,
			VerificationResponse: verification.VerificationResponse,
			UserId: verification.UserId,
		})
	}
	return verifications
}