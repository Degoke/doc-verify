package user

import (
	"net/http"

	"github.com/Degoke/doc-verify/common"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func stripBearerPrefixFromToken(token string) (string, error) {
	if len(token) > 7 && token[:7] == "Bearer " {
		return token[7:], nil
	}
	return token, nil
}

var AuthorizationHeaderExtractor = &request.PostExtractionFilter{
	Extractor: request.HeaderExtractor{"Authorization"},
	Filter:    stripBearerPrefixFromToken,
}

func UpdateContextUserModel(c *gin.Context, userId primitive.ObjectID) {
	c.Set("userId", userId)
}

func AuthMiddleware(auto401 bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		UpdateContextUserModel(c, primitive.NilObjectID)
		token, err := request.ParseFromRequest(c.Request, AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
			secret := common.GetENV("JWT_SECRET")
			return []byte(secret), nil
		})

		if err != nil {
			if auto401 {
				c.AbortWithError(http.StatusUnauthorized, err)
			}
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			stringId := claims["id"].(string)
			userId, err := primitive.ObjectIDFromHex(stringId)
			if err != nil {
				if auto401 {
					c.AbortWithError(http.StatusUnauthorized, err)
				}
				return
			}

			UpdateContextUserModel(c, userId)
		} else {
			if auto401 {
				c.AbortWithError(http.StatusUnauthorized, err)
			}
			return
		}

	}
}