package utils

import (
	"errors"
	"time"

	"mysite/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, username, role string) (string, error) {
	cfg := config.AppConfig.JWT
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.ExpireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ========== 管理员专用Token函数 ==========

// GenerateAdminToken 生成管理员Token（使用独立的secret）
func GenerateAdminToken(userID uint, username, role string) (string, error) {
	cfg := config.AppConfig.JWT
	
	// 使用管理员专用secret，如果未配置则使用默认secret+后缀
	adminSecret := cfg.AdminSecret
	if adminSecret == "" {
		adminSecret = cfg.Secret + "-admin-secret"
	}

	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.ExpireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "admin-service", // 标识为管理服务签发
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(adminSecret))
}

// ParseAdminToken 解析管理员Token（使用独立的secret）
func ParseAdminToken(tokenString string) (*Claims, error) {
	cfg := config.AppConfig.JWT
	
	// 使用管理员专用secret
	adminSecret := cfg.AdminSecret
	if adminSecret == "" {
		adminSecret = cfg.Secret + "-admin-secret"
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(adminSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// 额外验证：必须是管理服务签发的token
		if claims.Issuer != "admin-service" {
			return nil, errors.New("invalid admin token: wrong issuer")
		}
		return claims, nil
	}

	return nil, errors.New("invalid admin token")
}

