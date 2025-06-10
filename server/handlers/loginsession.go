package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strings"
	"templateapp/database"
	"templateapp/logger"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); nil != err {
		return err
	}
	if req.Email == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	account, err := database.FindAccountByHandleNameAndPassword(req.Email, req.Password)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.SendStatus(fiber.StatusNotFound)
	} else if nil != err {
		logger.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	logger.Info(fmt.Sprintf("account: %v", account))

	//todo 만든 account로
	s := c.Locals("session").(*session.Session)
	user := database.FindByName(req.Email)
	switch {
	default:
		return c.SendStatus(fiber.StatusNotFound)
	case nil != user && user.ID != 0:

		// JWT 토큰 생성
		expireTime := time.Now().Add(time.Hour * 72)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":   user.ID,
			"exp":       expireTime.Unix(),
			"iat":       time.Now().Unix(),
			"user_name": user.Name,
		})

		// Sign the token with a secret key
		secretKey := []byte("your_secret_key")
		tokenString, err := token.SignedString(secretKey)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		split := strings.Split(tokenString, ".")
		// 세션에 사용자 정보 저장
		s.Set("user", *user)

		// token string을 . 기준으로 나누어서 서명은 http only 쿠키로, 페이로드는 일반 쿠키로 저장
		c.Cookie(&fiber.Cookie{
			Name:     "jwt_server", // 쿠키 이름
			Value:    split[2],
			Expires:  expireTime,
			HTTPOnly: true,
		})

		c.Cookie(&fiber.Cookie{
			Name:    "jwt_client", // 쿠키 이름
			Value:   split[0] + "." + split[1],
			Expires: time.Now().Add(72 * time.Hour),
		})

		refreshExpire := time.Now().Add(7 * 24 * time.Hour)
		// jwt refresh token
		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":   user.ID,
			"exp":       refreshExpire.Unix(),
			"iat":       time.Now().Unix(),
			"user_name": user.Name,
		})
		refreshTokenString, err := refreshToken.SignedString(secretKey)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		c.Cookie(&fiber.Cookie{
			Name:     "refresh", // 쿠키 이름
			Value:    refreshTokenString,
			Expires:  refreshExpire, // 7일 동안 유지
			HTTPOnly: true,          // JavaScript에서 접근 금지
			Secure:   false,         // HTTPS 환경에서는 true로 설정
		})
		return c.JSON(user)
	}
}

// Logout 로그아웃
func Logout(c *fiber.Ctx) error {
	// 쿠키 삭제
	c.Cookie(&fiber.Cookie{
		Name:     "auto_login",
		Expires:  time.Now().Add(-time.Hour), // 과거로 설정하여 삭제
		HTTPOnly: true,
		Secure:   false,
	})

	c.Locals("session").(*session.Session).Destroy()

	return c.JSON(fiber.Map{"message": "Logged out successfully!" + c.Params("sessionID")})
}

func InitLoginSession(router fiber.Router) {
	rt := router.Group("/login-sessions")
	rt.Post("", Login)
	rt.Get("", func(c *fiber.Ctx) error {
		return c.JSON(
			map[string]interface{}{
				"session": c.Locals("session").(*session.Session).ID(),
			})
	})
	rt.Get(":sessionID", func(c *fiber.Ctx) error {
		return c.SendString("Session ID: " + c.Params("sessionID") + c.Locals("session").(*session.Session).ID())
	})
	rt.Delete(":sessionID", Logout)
}
