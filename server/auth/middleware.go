package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Middleware struct {
	Auth           *Service
	CookieDuration time.Duration
}

func NewMiddleware(authService *Service) fiber.Handler {
	m := &Middleware{Auth: authService}
	m.CookieDuration = 24 * time.Hour
	return m.handle
}

func (m *Middleware) Login(c *fiber.Ctx) error {
	type req struct {
		UserID uint `json:"user_id"`
	}
	var data req
	if err := c.BodyParser(&data); err != nil {
		return fiber.ErrBadRequest
	}
	_, err := m.Auth.Login(c.Context(), data.UserID)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	c.Locals(
		"userID", data.UserID)

	return c.JSON(fiber.Map{"message": "logged in"})
}

func (m *Middleware) handle(c *fiber.Ctx) error {
	signed := c.Cookies(m.Auth.SignedCookieName)
	unsigned := c.Cookies(m.Auth.UnsignedCookieName)
	if signed == "" || unsigned == "" {
		return fiber.ErrUnauthorized
	}
	// 헤더.payload와 signature 조합
	tokenStr := m.Auth.combineJWT(unsigned, signed)
	// 검증
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return m.Auth.JwtSecret, nil
	})
	if err != nil || !token.Valid {
		return fiber.ErrUnauthorized
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["userID"].(float64))
	// Redis 세션 검증
	key := m.Auth.sessionKey(userID)
	stored, err := m.Auth.SessionStorage.Get(c.Context(), key)
	if err != nil || stored != tokenStr {
		return fiber.ErrUnauthorized
	}
	c.Locals("userID", userID)
	nextErr := c.Next()
	if nextErr != nil {
		return nextErr
	}

	// 모든 미들웨어 또는 핸들러 후에 쿠키 재설정
	expires := time.Now().Add(m.CookieDuration)
	c.Cookie(&fiber.Cookie{
		Name:     m.Auth.SignedCookieName,
		Value:    signed,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Expires:  expires,
	})
	c.Cookie(&fiber.Cookie{
		Name:     m.Auth.UnsignedCookieName,
		Value:    unsigned,
		Secure:   false,
		SameSite: "Lax",
		Expires:  expires,
	})
	return nil
}
