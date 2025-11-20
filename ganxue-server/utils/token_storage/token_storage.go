package token_storage

import (
	"context"
	"fmt"
	"ganxue-server/global"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	// ShortTokenPrefix 短token在Redis中的键前缀
	ShortTokenPrefix = "short_token"
	// ShortTokenExpiry 短token过期时间
	ShortTokenExpiry = 15 * time.Minute
)

// StoreShortToken 存储短token到Redis
func StoreShortToken(shortToken string, userID uint, customExpiry ...time.Duration) error {
	key := fmt.Sprintf("%s:%s", ShortTokenPrefix, shortToken)

	expiry := ShortTokenExpiry
	if len(customExpiry) > 0 {
		expiry = customExpiry[0]
	}

	err := global.RDB.Set(context.Background(), key, userID, expiry).Err()
	if err != nil {
		return fmt.Errorf("存储短token到Redis失败: %w", err)
	}

	return nil
}

// GetUserIDByShortToken 根据短token获取用户ID
func GetUserIDByShortToken(shortToken string) (uint, error) {
	key := fmt.Sprintf("%s:%s", ShortTokenPrefix, shortToken)

	userIDStr, err := global.RDB.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, fmt.Errorf("短token不存在或已过期")
		}
		return 0, fmt.Errorf("从Redis获取短token失败: %w", err)
	}

	var userID uint
	_, err = fmt.Sscanf(userIDStr, "%d", &userID)
	if err != nil {
		return 0, fmt.Errorf("解析用户ID失败: %w", err)
	}

	return userID, nil
}

// DeleteShortToken 删除短token
func DeleteShortToken(shortToken string) error {
	key := fmt.Sprintf("%s:%s", ShortTokenPrefix, shortToken)

	err := global.RDB.Del(context.Background(), key).Err()
	if err != nil {
		return fmt.Errorf("删除短token失败: %w", err)
	}

	return nil
}

// RefreshShortToken 刷新短token（删除旧的，生成新的）
func RefreshShortToken(oldShortToken, newShortToken string, userID uint, customExpiry ...time.Duration) error {
	// 删除旧token
	if err := DeleteShortToken(oldShortToken); err != nil {
		return err
	}

	// 存储新token
	return StoreShortToken(newShortToken, userID, customExpiry...)
}

// CheckShortTokenExists 检查短token是否存在
func CheckShortTokenExists(shortToken string) (bool, error) {
	key := fmt.Sprintf("%s:%s", ShortTokenPrefix, shortToken)

	result, err := global.RDB.Exists(context.Background(), key).Result()
	if err != nil {
		return false, fmt.Errorf("检查短token存在性失败: %w", err)
	}

	return result > 0, nil
}

// ExtendShortTokenExpiry 延长短token过期时间
func ExtendShortTokenExpiry(shortToken string, additionalTime time.Duration) error {
	key := fmt.Sprintf("%s:%s", ShortTokenPrefix, shortToken)

	err := global.RDB.Expire(context.Background(), key, additionalTime).Err()
	if err != nil {
		return fmt.Errorf("延长短token过期时间失败: %w", err)
	}

	return nil
}

// GetUserShortTokens 获取用户的所有短token（用于用户登出时清理）
func GetUserShortTokens(userID uint) ([]string, error) {
	pattern := fmt.Sprintf("%s:*", ShortTokenPrefix)

	keys, err := global.RDB.Keys(context.Background(), pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("获取短token键列表失败: %w", err)
	}

	var userTokens []string
	for _, key := range keys {
		userIDStr, err := global.RDB.Get(context.Background(), key).Result()
		if err != nil {
			continue // 跳过获取失败的键
		}

		var tokenUserID uint
		_, err = fmt.Sscanf(userIDStr, "%d", &tokenUserID)
		if err != nil {
			continue
		}

		if tokenUserID == userID {
			// 提取token部分（去掉前缀）
			if len(key) > len(ShortTokenPrefix)+1 {
				token := key[len(ShortTokenPrefix)+1:]
				userTokens = append(userTokens, token)
			}
		}
	}

	return userTokens, nil
}

// DeleteAllUserShortTokens 删除用户的所有短token
func DeleteAllUserShortTokens(userID uint) error {
	tokens, err := GetUserShortTokens(userID)
	if err != nil {
		return err
	}

	for _, token := range tokens {
		if err := DeleteShortToken(token); err != nil {
			// 记录错误但继续删除其他token
			continue
		}
	}

	return nil
}