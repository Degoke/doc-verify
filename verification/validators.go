package verification

type CreateVerificationValidator struct {
	Values map[string]string `json:"values" binding:"required"`
}