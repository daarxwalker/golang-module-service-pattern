package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"

	"example/core/helper/enviromentHelper"
	"example/core/helper/pathHelper"
)

type configApp struct {
	Name               string
	Port               int
	TwoFactorImageSize int
	ReCaptchaSecret    string
	SpacesName         string
	PDFHostname        string
	CronSecret         string
}

type configCache struct {
	Addr     string
	TestAddr string
	Password string
	DB       int
}

type configDatabase struct {
	Addr            string
	User            string
	Password        string
	Dbname          string
	CertificatePath string
}

type configBitrix struct {
	Webhook string
}

type Config interface {
	GetApp() configApp
	GetBitrix() configBitrix
	GetCache() configCache
	GetDatabase() configDatabase
}

type config struct {
	app      configApp
	bitrix   configBitrix
	cache    configCache
	database configDatabase
}

func New() Config {
	instance := &config{}

	readConfig()
	instance.parseAppConfig()
	instance.parseBitrixConfig()
	instance.parseCacheConfig()
	instance.parseDatabaseConfig()

	return instance
}

func readConfig() {
	viper.SetConfigName("example.config")
	viper.SetConfigType("yml")

	rootPath, err := pathHelper.GetRoot()
	if err != nil {
		log.Fatalln("fatal error reading config file: \n", err)
	}

	viper.AddConfigPath(fmt.Sprintf("%s", rootPath))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("fatal error reading config file: \n", err)
	}
}

func (c *config) parseAppConfig() {
	c.app = configApp{
		Name: viper.GetString("app.name"),
		Port: viper.GetInt("app.port"),
	}
}

func (c *config) parseBitrixConfig() {
	c.bitrix = configBitrix{
		Webhook: viper.GetString("bitrix.webhook"),
	}
}

func (c *config) parseCacheConfig() {
	c.cache = configCache{
		Addr:     viper.GetString("cache.address"),
		TestAddr: viper.GetString("cache.testAddress"),
		Password: viper.GetString("cache.password"),
		DB:       viper.GetInt("cache.db"),
	}
}

func (c *config) parseDatabaseConfig() {
	dbType := enviromentHelper.GetEnviroment()
	c.database = configDatabase{
		Addr:            viper.GetString(fmt.Sprintf("database.%s.address", dbType)),
		User:            viper.GetString(fmt.Sprintf("database.%s.user", dbType)),
		Password:        viper.GetString(fmt.Sprintf("database.%s.password", dbType)),
		Dbname:          viper.GetString(fmt.Sprintf("database.%s.dbname", dbType)),
		CertificatePath: viper.GetString(fmt.Sprintf("database.%s.certificatePath", dbType)),
	}
}

func (c config) GetApp() configApp {
	return c.app
}

func (c config) GetBitrix() configBitrix {
	return c.bitrix
}

func (c config) GetCache() configCache {
	return c.cache
}

func (c config) GetDatabase() configDatabase {
	return c.database
}
