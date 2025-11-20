package test

import (
	"ganxue-server/initialize"
	"ganxue-server/utils/token"
	"testing"
)

// TestTokenGeneration 测试token生成
func TestTokenGeneration(t *testing.T) {
	initialize.InitAll()

	// 测试短token生成
	userID := uint(123)
	shortToken, err := token.GenerateShortToken(userID)
	if err != nil {
		t.Errorf("生成短token失败: %v", err)
	}
	if shortToken == "" {
		t.Error("短token为空")
	}

	// 测试长token生成
	longToken, err := token.GenerateLongToken(userID)
	if err != nil {
		t.Errorf("生成长token失败: %v", err)
	}
	if longToken == "" {
		t.Error("长token为空")
	}

	t.Logf("生成的短token: %s", shortToken)
	t.Logf("生成的长token: %s", longToken)
}

// TestShortTokenParsing 测试短token解析
func TestShortTokenParsing(t *testing.T) {
	initialize.InitAll()

	// 生成短token
	userID := uint(456)
	shortToken, err := token.GenerateShortToken(userID)
	if err != nil {
		t.Fatalf("生成短token失败: %v", err)
	}

	// 解析短token
	parsedUserID, err := token.ParseShortToken(shortToken)
	if err != nil {
		t.Fatalf("解析短token失败: %v", err)
	}

	if parsedUserID != userID {
		t.Errorf("用户ID不匹配: 期望 %d, 得到 %d", userID, parsedUserID)
	}
}

// TestLongTokenParsing 测试长token解析
func TestLongTokenParsing(t *testing.T) {
	initialize.InitAll()

	// 生成长token
	userID := uint(789)
	longToken, err := token.GenerateLongToken(userID)
	if err != nil {
		t.Fatalf("生成长token失败: %v", err)
	}

	// 解析长token
	parsedUserID, err := token.ParseToken(longToken)
	if err != nil {
		t.Fatalf("解析长token失败: %v", err)
	}

	if parsedUserID != userID {
		t.Errorf("用户ID不匹配: 期望 %d, 得到 %d", userID, parsedUserID)
	}
}

// TestTokenRefresh 测试短token刷新
func TestTokenRefresh(t *testing.T) {
	initialize.InitAll()

	// 生成短token
	userID := uint(321)
	oldShortToken, err := token.GenerateShortToken(userID)
	if err != nil {
		t.Fatalf("生成短token失败: %v", err)
	}

	// 刷新短token
	newShortToken, err := token.RefreshShortToken(oldShortToken, userID)
	if err != nil {
		t.Fatalf("刷新短token失败: %v", err)
	}

	if newShortToken == "" {
		t.Error("新短token为空")
	}

	if newShortToken == oldShortToken {
		t.Error("新短token与旧短token相同，应该不同")
	}

	// 测试旧token应该无效
	_, err = token.ParseShortToken(oldShortToken)
	if err == nil {
		t.Error("旧短token应该无效，但解析成功了")
	}

	// 测试新token应该有效
	parsedUserID, err := token.ParseShortToken(newShortToken)
	if err != nil {
		t.Errorf("新短token解析失败: %v", err)
	}

	if parsedUserID != userID {
		t.Errorf("新短token用户ID不匹配: 期望 %d, 得到 %d", userID, parsedUserID)
	}
}