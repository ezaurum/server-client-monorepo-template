package main

import (
	"context"
	"errors"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	ml "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/template/html/v2"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"templateapp/auth"
	"templateapp/configs"
	"templateapp/database"
	"templateapp/handlers"
	"templateapp/logger"
	"templateapp/models"
	"time"
)

func main() {
	config := configs.New()

	// Connected with database
	db := database.ConnectWithConfig(config.Database)

	// Create fiber app
	app := fiber.New(fiber.Config{
		Prefork: config.GetBool("PREFORK"), // go run app.go -prod
		Views:   html.New("./static/private", ".html"),
	})
	rootUser := models.User{
		Name: "root",
	}
	if tx := db.Find(&rootUser, 1); nil != tx.Error {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			//todo rootUser.Password = utils.HashPassword(config.GetString("ROOT_PASSWORD"))
			db.Create(&rootUser)
		} else {
			logger.Error("%v", tx.Error)
		}
	}

	// Middleware
	app.Use(recover.New(
		recover.Config{
			EnableStackTrace: true,
		}))
	app.Use(requestid.New())
	app.Use(ml.New(
		ml.Config{
			Format: "${pid} ${latency} ${locals:requestid} ${status} - ${method} ${path}\n",
		}))

	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return strings.Contains(origin, "localhost")
		},
		AllowCredentials: true,
	}))

	config.SetupLogin(app)

	const HeaderName = "X-CSRF-Token"
	const CookieName = "_Host-csrf_"
	app.Use(csrf.New(
		csrf.Config{
			KeyLookup:    "header:" + HeaderName,
			CookieName:   CookieName,
			Expiration:   1 * time.Hour,
			KeyGenerator: utils.UUIDv4,
			Extractor: func(c *fiber.Ctx) (string, error) {
				if s, e := csrf.CsrfFromQuery("_csrf")(c); nil == e {
					return s, e
				}
				if s, e := csrf.CsrfFromParam("_csrf")(c); nil == e {
					return s, e
				}
				if s, e := csrf.CsrfFromHeader(HeaderName)(c); nil == e {
					return s, e
				}
				return csrf.CsrfFromForm("_csrf")(c)
			},
			//todo 로그인 세션과 별로도로 해야할지? 고민이 필요
			//todo Session:           sessionStore,
			//todo csrf가 세션에 저장될 필요가 있을지?
			//todo SessionKey:        "fiber.csrf.token",
			HandlerContextKey: "fiber.csrf.handler",
			ContextKey:        "fiber.csrf.token_string",
		}))

	_ = auth.CasbinMiddleware(db)
	// Create a /api/v1 endpoint
	// auto login
	v1 := app.Group("/api/v1", func(c *fiber.Ctx) error {
		username := c.Cookies("auto_login")
		if username != "" {
			user := database.FindByName(username)
			if user != nil {
				s := c.Locals("session").(*session.Session)
				s.Set("user", *user)
			}
		}
		return c.Next()
	})

	// Bind handlers
	v1.Get("/users", handlers.UserList)
	v1.Post("/users", func(c *fiber.Ctx) error {
		s := c.Locals("session")
		if s == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		ses := s.(*session.Session)
		u := ses.Get("user")
		if u == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		user := u.(models.User)
		if user.Name != "test" {
			return c.SendStatus(fiber.StatusForbidden)
		}
		logger.Debug("user: %v", user)
		return c.Next()
	},
		handlers.UserCreate)
	v1.Delete("/users", handlers.UserDelete)
	v1.Get("/users/:id", handlers.UserGet)

	handlers.InitLoginSession(v1)
	v1.Post("/login", handlers.Login)
	v1.Post("/logout", handlers.Logout)

	// Setup static files
	app.Static("/", "./static/public")

	app.Get("/login", func(c *fiber.Ctx) error {
		// 템플릿 렌더링 시 토큰 전달
		return c.Render("login", fiber.Map{
			"csrfToken": c.Locals("fiber.csrf.token_string"),
		})
	})
	app.Get("/users", func(c *fiber.Ctx) error {
		// 템플릿 렌더링 시 토큰 전달
		return c.Render("users", fiber.Map{
			"csrfToken": c.Locals("fiber.csrf.token_string"),
		})
	})

	// websocket
	app.Get("/ws/:id", websocket.New(handlers.WebSocket))

	// Handle not founds
	app.Use(handlers.NotFound)

	// Listen from a different goroutine
	go func() {
		logger.Error("%v", app.Listen(config.ListenString()))
	}()

	// Wait for interrupt signal to gracefully shutdown the app

	shutdownChannel := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(shutdownChannel, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	_ = <-shutdownChannel // This blocks the main thread until an interrupt is received
	logger.Info("Gracefully shutting down...")

	// 컨텍스트 생성 및 타임아웃 설정
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("Error shutting down: %v", err)
	}

	logger.Info("Fiber was successful shutdown.")

	logger.Info("Running cleanup tasks...")

	// Your cleanup tasks go here
	// db.Close()
	// redisConn.Close()
	// 여기
	logger.Info("All app was successful shutdown.")
}
