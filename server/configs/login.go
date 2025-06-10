package configs

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis/v3"
	"github.com/golang-jwt/jwt/v5"
	"runtime"
	"templateapp/logger"
	"templateapp/models"
	"time"
)

const SessionExpirationKey = "SESSION_EXPIRATION"
const LoginSessionTokenSecretKey = "LOGIN_SESSION_TOKEN_SECRET_KEY"
const LoginSessionContextKey = "LOGIN_SESSION_CONTEXT_KEY"
const LoginSessionUserKey = "user"

type LoginSessionManager struct {
	Store      *session.Store
	Expiration time.Duration
	SecretKey  []byte
}

func (c *Config) setDefaultLoginConfig() {
	c.SetDefault(SessionExpirationKey, 30*time.Minute)
}

func (c *Config) SetupLogin(app *fiber.App) *LoginSessionManager {
	sessionExpiration := time.Duration(c.GetInt(SessionExpirationKey))
	secretKey := []byte(c.GetString(LoginSessionTokenSecretKey))

	// Initialize redis config from .env
	storage := redis.New(redis.Config{
		Host:      c.Redis.GetString("HOST"),
		Port:      c.Redis.GetInt("PORT"),
		Username:  "",
		Password:  "",
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	})

	sessionStore := session.New(session.Config{
		Storage:    storage,
		Expiration: sessionExpiration,
	})

	app.Use(func(c *fiber.Ctx) error {
		// 쿠키에 저장된 세션 정보를 가져옴
		jwtServer := c.Cookies("jwt_server")
		jwtClient := c.Cookies("jwt_client")
		//refresh := c.Cookies("refresh")
		if jwtServer != "" && jwtClient != "" {
			jwtTokenString := jwtClient + "." + jwtServer
			// 토큰 검증
			claims := jwt.MapClaims{}
			_, err := jwt.ParseWithClaims(jwtTokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return secretKey, nil
			})
			if err != nil {
				logger.Error("jwt parse error:", err)
				//return c.SendStatus(fiber.StatusUnauthorized)
			}
		}

		currentLoginSession, err := sessionStore.Get(c)
		if err != nil {
			return err
		}
		defer func(get *session.Session) {
			if err := get.Save(); err != nil {
				logger.Error("session save error:", err)
			}
		}(currentLoginSession)
		if currentLoginSession.Fresh() {
			currentLoginSession.Set(LoginSessionUserKey, models.User{})
		} else {
			// 세션 만료 시간 연장
			currentLoginSession.SetExpiry(sessionExpiration)
		}

		c.Locals(LoginSessionContextKey, currentLoginSession)
		return c.Next()
	})

	manager := LoginSessionManager{
		Store:      sessionStore,
		Expiration: sessionExpiration,
		SecretKey:  secretKey,
	}
	return &manager

}
