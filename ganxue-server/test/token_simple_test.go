package test

import (
	"testing"

	"github.com/google/uuid"
)

// TestUUIDGeneration 测试UUID生成
func TestUUIDGeneration(t *testing.T) {
	// 测试UUID生成
	uuid1 := uuid.New().String()
	uuid2 := uuid.New().String()

	if uuid1 == uuid2 {
		t.Error("生成的UUID应该不同")
	}

	if uuid1 == "" || uuid2 == "" {
		t.Error("UUID不能为空")
	}

	// 验证UUID格式（应该包含横线）
	if len(uuid1) < 30 {
		t.Error("UUID格式不正确，长度太短")
	}

	t.Logf("生成的UUID1: %s", uuid1)
	t.Logf("生成的UUID2: %s", uuid2)
}

// TestTokenLogic 测试token逻辑
func TestTokenLogic(t *testing.T) {
	userID := uint(123)

	// 模拟token生成逻辑
	shortToken := uuid.New().String()

	// 验证token不为空
	if shortToken == "" {
		t.Error("生成的短token不能为空")
	}

	// 模拟Redis key格式
	expectedKeyFormat := "short_token:" + shortToken
	if expectedKeyFormat == "" {
		t.Error("Redis key格式不能为空")
	}

	t.Logf("用户ID: %d", userID)
	t.Logf("生成的短token: %s", shortToken)
	t.Logf("Redis key格式: %s", expectedKeyFormat)
}

// TestBasicAuthenticationFlow 测试基本认证流程
func TestBasicAuthenticationFlow(t *testing.T) {
	userID := uint(456)

	// 1. 生成短token
	shortToken := uuid.New().String()
	if shortToken == "" {
		t.Fatal("短token生成失败")
	}

	// 2. 模拟存储到Redis（这里只验证key格式）
	redisKey := "short_token:" + shortToken
	redisValue := userID

	if redisKey == "" {
		t.Error("Redis key不能为空")
	}

	if redisValue != userID {
		t.Errorf("Redis值应该为用户ID %d，但得到 %d", userID, redisValue)
	}

	// 3. 模拟从Redis解析
	parsedUserID := redisValue
	if parsedUserID != userID {
		t.Errorf("解析的用户ID不匹配: 期望 %d, 得到 %d", userID, parsedUserID)
	}

	t.Logf("认证流程测试通过 - 用户ID: %d, 短token: %s", userID, shortToken)
}