package auth

import (
	"fmt"
	"github.com/spf13/cast"
	"med/internal/configs"
	"med/internal/schemas"
	"med/pkg/logger"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const HeaderKeyAuthorization = "Authorization"

const ContextKeyAuthorization = "ContextKeyAuthorization"
const ContextKeyUserID = "userID"

type Claims struct {
	UserID string `json:"userID"`
	jwt.StandardClaims
}

func GenerateJWTToken(arg *schemas.TokenPayload, log logger.ILogger) (*schemas.TokenResponse, error) {

	accessNewJWt := jwt.New(jwt.SigningMethodHS256)

	accessClaims := accessNewJWt.Claims.(jwt.MapClaims)
	accessClaims["id"] = arg.ID
	accessClaims["role"] = arg.Role
	accessClaims["iat"] = time.Now().Unix()
	accessClaims["exp"] = time.Now().Add(configs.AccessExpireTime).Unix()

	accessToken, err := accessNewJWt.SignedString(configs.SignKey)
	if err != nil {
		log.Error("this error is signedString-~~~~~~~~~>ERROR", logger.Error(err))
		return nil, err
	}
	refreshNewJWT := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshNewJWT.Claims.(jwt.MapClaims)
	refreshClaims["id"] = arg.ID
	refreshClaims["role"] = arg.Role
	refreshClaims["iat"] = time.Now().Unix()
	refreshClaims["exp"] = time.Now().Add(configs.RefreshExpireTime).Unix()

	refreshToken, err := refreshNewJWT.SignedString(configs.SignKey)
	if err != nil {
		log.Error("this error is signed  to string that sign key -~~~~~~~~~ERROR", logger.Error(err))
		return nil, err
	}

	return &schemas.TokenResponse{
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		AccessExpiredTime:  cast.ToFloat64(accessClaims["exp"].(int64) - accessClaims["iat"].(int64)),
		RefreshExpiresTime: cast.ToFloat64(refreshClaims["exp"].(int64) - refreshClaims["iat"].(int64)),
		Success:            true,
	}, nil

}

func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	_, _, err := new(jwt.Parser).ParseUnverified(tokenString, claims)
	if err != nil {
		return nil, err
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, fmt.Errorf("token is expired")
	}

	return claims, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must start with Bearer"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ValidateJWT(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		ctx.Set(ContextKeyUserID, claims.UserID)
		ctx.Next()
	}
}
