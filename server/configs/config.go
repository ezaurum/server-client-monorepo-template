package configs

import (
	"fmt"
	"github.com/spf13/viper"
	"templateapp/logger"
)

type Config struct {
	*viper.Viper
	Database ConfigGetter
	Redis    ConfigGetter
}

type ConfigGetter interface {
	Get(key string) interface{}
	GetString(key string) string
	GetInt(key string) int
}

func New() *Config {
	v := Config{
		Viper: viper.New(),
	}

	v.AutomaticEnv()
	v.SetConfigName(".env")
	v.SetConfigType("dotenv")
	v.AddConfigPath(".")

	v.setDefaults()
	v.setDefaultLoginConfig()
	v.SetDefault("ROOT_PASSWORD", "test")

	err := v.ReadInConfig()
	if err != nil {
		logger.Warn("Failed to read config file", err)
	}

	dbConfigPath := v.Get("DB_CONFIG_PATH")
	if dbConfigPath != nil {
		// If DB_CONFIG_PATH is set, use it as the config file path
		// json format
		dbConfig := viper.New()
		dbConfig.SetConfigFile(dbConfigPath.(string))
		dbConfigError := dbConfig.ReadInConfig()
		if dbConfigError != nil {
			logger.Warn("Failed to read database config file", dbConfigError)
		}
		v.Database = dbConfig
	} else {
		v.Database = &viper.Viper{}
	}
	redisConfigPath := v.Get("REDIS_CONFIG_PATH")
	if redisConfigPath != nil {
		// If REDIS_CONFIG_PATH is set, use it as the config file path
		// json format
		redisConfig := viper.New()
		redisConfig.SetConfigFile(redisConfigPath.(string))
		redisConfigError := redisConfig.ReadInConfig()
		if redisConfigError != nil {
			logger.Warn("Failed to read redis config file", redisConfigError)
		}
		v.Redis = redisConfig
	} else {
		v.Redis = &viper.Viper{}
	}

	return &v
}

func (c *Config) setDefaults() {
	c.SetDefault("PORT", 3000)
	c.SetDefault("PREFORK", false)
}

func (c *Config) Port() int {
	return c.GetInt("PORT")
}

func (c *Config) ListenString() string {
	return fmt.Sprintf(":%d", c.Port())
}
