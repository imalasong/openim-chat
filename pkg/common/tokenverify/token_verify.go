package tokenverify

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/chat/pkg/common/constant"
	utils "github.com/OpenIMSDK/open_utils"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const (
	TokenUser  = constant.NormalUser
	TokenAdmin = constant.AdminUser
)

type claims struct {
	UserID   string
	UserType int32
	jwt.RegisteredClaims
}

func buildClaims(userID string, userType int32, ttlDay int64) claims {
	now := time.Now()
	before := now.Add(-time.Minute * 5)
	return claims{
		UserID:   userID,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(ttlDay*24) * time.Hour)), //Expiration time
			IssuedAt:  jwt.NewNumericDate(now),                                           //Issuing time
			NotBefore: jwt.NewNumericDate(before),                                        //Begin Effective time
		}}
}

func CreateToken(UserID string, userType int32, ttlDay int64) (string, error) {
	if !(userType == TokenUser || userType == TokenAdmin) {
		return "", errs.ErrTokenUnknown.Wrap("token type unknown")
	}
	claims := buildClaims(UserID, userType, ttlDay)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(*config.Config.Secret))
	if err != nil {
		return "", utils.Wrap(err, "")
	}
	return tokenString, nil
}

func secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(*config.Config.Secret), nil
	}
}

func getToken(t string) (string, int32, error) {
	token, err := jwt.ParseWithClaims(t, &claims{}, secret())
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return "", 0, errs.ErrTokenMalformed.Wrap()
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return "", 0, errs.ErrTokenExpired.Wrap()
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return "", 0, errs.ErrTokenNotValidYet.Wrap()
			} else {
				return "", 0, errs.ErrTokenUnknown.Wrap()
			}
		} else {
			return "", 0, errs.ErrTokenNotValidYet.Wrap()
		}
	} else {
		if claims, ok := token.Claims.(*claims); ok && token.Valid {
			return claims.UserID, claims.UserType, nil
		}
		return "", 0, errs.ErrTokenNotValidYet.Wrap()
	}
}

func GetToken(token string) (string, int32, error) {
	userID, userType, err := getToken(token)
	if err != nil {
		return "", 0, err
	}
	if !(userType == TokenUser || userType == TokenAdmin) {
		return "", 0, errs.ErrTokenUnknown.Wrap("token type unknown")
	}
	return userID, userType, nil
}

func GetAdminToken(token string) (string, error) {
	userID, userType, err := getToken(token)
	if err != nil {
		return "", err
	}
	if userType != TokenAdmin {
		return "", errs.ErrTokenInvalid.Wrap("token type error")
	}
	return userID, nil
}

func GetUserToken(token string) (string, error) {
	userID, userType, err := getToken(token)
	if err != nil {
		return "", err
	}
	if userType != TokenUser {
		return "", errs.ErrTokenInvalid.Wrap("token type error")
	}
	return userID, nil
}
