package templateapp

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"strings"
	"templateapp/auth"
	"testing"
)

func TestAuthMiddleware_SetsCookiesAtEnd(t *testing.T) {
	ctx := context.Background()

	// 목 Redis
	rdb := NewMockRedis()

	// 서비스 생성
	authSvc := auth.NewAuthService(rdb, "jwt_signed", "jwt_unsigned", []byte("test_secret"))
	middleware := auth.NewMiddleware(authSvc)

	app := fiber.New()

	// 로그인 엔드포인트
	app.Post("/login", func(c *fiber.Ctx) error {
		type req struct {
			UserID uint `json:"user_id"`
		}
		var data req
		if err := c.BodyParser(&data); err != nil {
			return fiber.ErrBadRequest
		}
		token, err := authSvc.Login(ctx, data.UserID)
		if err != nil {
			return fiber.ErrInternalServerError
		}
		// 응답에서 token은 제외, 쿠키는 미들웨어에서 설정
		return c.JSON(fiber.Map{"token": token})
	})

	// 보호된 그룹 미들웨어 적용
	protected := app.Group("/protected", middleware)
	protected.Get("/", func(c *fiber.Ctx) error {
		userID := c.Locals("userID")
		return c.JSON(fiber.Map{"userID": userID})
	})

	// 1. 로그인 요청
	req := httptest.NewRequest(
		"POST", "/login", strings.NewReader(`{"user_id": 123}`))
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// 헤더에서 Set-Cookie 추출
	cookies := resp.Header.Values("Set-Cookie")
	var signedCookie, unsignedCookie string
	for _, c := range cookies {
		cStr := string(c)
		if strings.HasPrefix(cStr, "jwt_signed") {
			signedCookie = cStr
		} else if strings.HasPrefix(cStr, "jwt_unsigned") {
			unsignedCookie = cStr
		}
	}
	assert.NotEmpty(t, signedCookie)
	assert.NotEmpty(t, unsignedCookie)

	// 다음 요청에서 쿠키 포함시키기
	req2 := httptest.NewRequest("GET", "/protected/", nil)

	resp2, err := app.Test(req2)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp2.StatusCode)

	// 응답 바디에서 userID 검증
	body, err := io.ReadAll(resp2.Body)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	assert.Contains(t, string(body), `"userID":123`)
}
