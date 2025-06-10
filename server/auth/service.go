package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

// SessionStorage 세션 저장용 인터페이스 정의
type SessionStorage interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

// Service Auth 구조체
type Service struct {
	SessionStorage     SessionStorage
	SignedCookieName   string
	UnsignedCookieName string
	JwtSecret          []byte
}

// NewAuthService 생성자
func NewAuthService(rdb SessionStorage, signedCookie, unsignedCookie string, secret []byte) *Service {
	return &Service{
		SessionStorage:     rdb,
		SignedCookieName:   signedCookie,
		UnsignedCookieName: unsignedCookie,
		JwtSecret:          secret,
	}
}

// JWT 생성
func (s *Service) createJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.JwtSecret)
}

// JWT 문자열을 `header.payload.signature`로 분리
func (s *Service) splitJWT(token string) (string, string) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", ""
	}
	return parts[0] + "." + parts[1], parts[2] // header.payload와 signature
}

// JWT를 재조합
func (s *Service) combineJWT(headerPayload, signature string) string {
	return headerPayload + "." + signature
}

// Login 로그인 처리 함수
func (s *Service) Login(ctx context.Context, userID uint) (string, error) {
	tokenString, err := s.createJWT(userID)
	if err != nil {
		return "", err
	}
	headerPayload, signature := s.splitJWT(tokenString)
	// 세션 스토리지에 저장
	err = s.SessionStorage.Set(ctx, s.sessionKey(userID), tokenString, 24*time.Hour)
	return s.combineJWT(headerPayload, signature), err
}

// 세션 키 생성
func (s *Service) sessionKey(userID uint) string {
	return "session:" + fmt.Sprint(userID)
}
