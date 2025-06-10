package templateapp

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"templateapp/controllers"
	"templateapp/models"
	"templateapp/services"
)

//todo 여기서 계정을 만들고, 계정에 로그인 정보를 넣어서 로그인을 시도한다.

func SetupMockDB() (*gorm.DB, error) {
	// 메모리에 존재하는 sqlite DB를 setup 합니다.
	db, err := gorm.Open(sqlite.Open(":memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 테스트 환경의 테이블 자동 생성
	db.AutoMigrate(&models.User{})

	return db, nil
}

func TestRegisterAndLogin(t *testing.T) {
	db, err := SetupMockDB()
	assert.NoError(t, err)

	authService := services.UserService{
		DB: db,
	}

	email := "test@example.com"
	password := "testpassword"
	name := "John Doe"

	// 회원가입 테스트
	user, err := authService.Register(email, password, name)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)

	// 중복 이메일 가입 테스트
	_, err = authService.Register(email, password, name)
	assert.Error(t, err)

	// 로그인 성공 테스트
	loginUser, err := authService.Login(email, password)
	assert.NoError(t, err)
	assert.Equal(t, email, loginUser.Email)

	// 잘못된 비밀번호 로그인 테스트
	_, err = authService.Login(email, "wrongpassword")
	assert.Error(t, err)

	// 등록되지 않은 이메일 로그인 테스트
	_, err = authService.Login("wrong@example.com", password)
	assert.Error(t, err)
}

func SetupAppWithMock() *fiber.App {
	mockDB := NewMockDBUser()
	authService := services.UserService{DB: mockDB}
	authController := controllers.AuthController{Service: &authService}

	app := fiber.New()
	app.Post("/api/login", authController.Login)

	// 사전에 회원가입 처리 (로그인 테스트를 위해)
	authService.Register("test@example.com", "password123", "John Doe")

	return app
}

func TestLoginAPI(t *testing.T) {
	app := SetupAppWithMock()

	payload := `{"email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)

	// 상태코드 체크
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// 응답 바디 체크
	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "Login successful", body["message"])

	// 쿠키 체크
	cookies := resp.Cookies()
	var jwtCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "jwt_token" {
			jwtCookie = c
		}
	}
	assert.NotNil(t, jwtCookie)
	assert.Equal(t, "dummy-jwt-token", jwtCookie.Value)

	// 헤더 체크
	userIdHeader := resp.Header.Get("X-User-Id")
	assert.NotEmpty(t, userIdHeader)
	assert.Equal(t, "1", userIdHeader)
}

func TestLoginAPIFailure(t *testing.T) {
	app := SetupAppWithMock()

	// 비밀번호 오류 테스트
	payload := `{"email":"test@example.com","password":"wrongpassword"}`
	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)

	// 상태코드 체크 (Unauthorized 예상)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	// 응답 바디 체크
	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Contains(t, body["message"].(string), "invalid credentials")

	// 쿠키 체크 (쿠키 설정 없어야 함)
	cookies := resp.Cookies()
	var jwtCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "jwt_token" {
			jwtCookie = c
		}
	}
	assert.Nil(t, jwtCookie)

	// 헤더 체크 (헤더 설정 없어야 함)
	userIdHeader := resp.Header.Get("X-User-Id")
	assert.Empty(t, userIdHeader)
}
