package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"sync"
	"templateapp/configs"
	"templateapp/models"
)

var (
	db *gorm.DB
)

func DB() *gorm.DB {
	return db
}

type Connection struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// Database 인터페이스는 DB 객체를 추상화하여 테스트 시 MockDB를 넣을 수 있도록 한다.
type Database interface {
	Create(value interface{}) *gorm.DB
	Find(dest interface{}, conds ...interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
}

func (c Connection) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		c.Host, c.User, c.Password, c.Database, c.Port,
	)
}

const TablePrefix = "templateapp_"

var _once sync.Once

func ConnectWithConfig(c configs.ConfigGetter) *gorm.DB {
	return Connect(Connection{
		Host:     c.GetString("HOST"),
		Port:     c.GetInt("PORT"),
		User:     c.GetString("USER"),
		Password: c.GetString("PASSWORD"),
		Database: c.GetString("DATABASE"),
	})
}

// Connect with database
func Connect(c Connection) *gorm.DB {
	_once.Do(func() {
		// 복수형 사용하지 않음
		_db, err := gorm.Open(postgres.Open(c.DSN()), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   TablePrefix,
				SingularTable: true,
			},
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			panic("failed to connect database")
		}
		err = _db.AutoMigrate(&models.User{}, &models.Post{}, &models.GatheringUser{}, &models.Gathering{}, &models.GatheringSession{}, &models.GatheringSessionAttendanceSetting{}, &models.AttendanceLog{}, &models.CheckInLog{},
			&models.CheckOutLog{},
			&models.AttendanceAuditLog{},
			&models.AttendanceAggregate{},
			&models.LoginPasswordGrant{},
			&models.LoginHandleNameGrant{},
		)
		if err != nil {
			panic("failed to migrate database")
		}
		db = _db
		log.Println("Connected with Database")
	})
	return db
}

func Insert(user *models.User) {
	tx := db.Create(user)
	if tx.Error != nil {
		panic(tx.Error)
	}
}

func Get() []*models.User {
	var users []*models.User
	db.Find(&users)
	return users
}

func FindByName(name string) *models.User {
	var user models.User
	row := db.Where("name = ?", name).First(&user)
	if row.Error != nil {
		return nil
	}
	return &user
}
