package token

import (
	mError "ganxue-server/utils/error"
	"ganxue-server/utils/token_storage"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"

	"github.com/google/uuid"
)

// GenerateShortToken 生成短token (UUID存储在Redis中)
func GenerateShortToken(userID uint) (string, *mError.Error) {
	// 生成UUID作为短token
	shortToken := uuid.New().String()

	// 使用token存储工具存储到Redis
	err := token_storage.StoreShortToken(shortToken, userID)
	if err != nil {
		return "", mError.New(mError.GenerateTokenError, err, "存储短token到Redis失败")
	}

	return shortToken, nil
}

// GenerateLongToken 生成长token
func GenerateLongToken(userID uint) (string, *mError.Error) {
	// 生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	// 生成token字符串
	tokenString, err := token.SignedString([]byte(os.ExpandEnv("$JWT_SECRET")))
	if err != nil {
		return "", mError.New(mError.GenerateTokenError, err, "生成token失败")
	}
	return tokenString, nil
}

// ParseShortToken 解析短token (从Redis中获取用户ID)
func ParseShortToken(shortToken string) (uint, *mError.Error) {
	// 使用token存储工具从Redis中获取用户ID
	userID, err := token_storage.GetUserIDByShortToken(shortToken)
	if err != nil {
		return 0, mError.New(mError.ParseTokenError, err, err.Error())
	}

	return userID, nil
}

// ParseToken 解析长token (JWT)
func ParseToken(tokenString string) (uint, *mError.Error) {
	// 解析token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.ExpandEnv("$JWT_SECRET")), nil
	})
	if err != nil {
		return 0, mError.New(mError.ParseTokenError, err, "解析token失败")
	}

	// 获取token中的claims信息
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, mError.New(mError.ParseTokenError, err, "解析token失败")
	}
	return uint(claims["userID"].(float64)), nil
}

// RefreshShortToken 刷新短token (删除旧的，生成新的)
func RefreshShortToken(oldShortToken string, userID uint) (string, *mError.Error) {
	// 生成新的短token
	newShortToken, err := GenerateShortToken(userID)
	if err != nil {
		return "", err
	}

	// 删除旧的短token
	_ = token_storage.DeleteShortToken(oldShortToken)
	// 忽略删除错误，因为新的token已经生成成功

	return newShortToken, nil
}

// LogoutUser 用户登出，删除所有短token
func LogoutUser(userID uint) *mError.Error {
	err := token_storage.DeleteAllUserShortTokens(userID)
	if err != nil {
		return mError.New(mError.RedisError, err, "删除用户短token失败")
	}
	return nil
}
