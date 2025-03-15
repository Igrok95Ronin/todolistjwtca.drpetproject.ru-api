package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
	"sync"
)

// Структура конфигурации
type Config struct {
	Port string         `yaml:"port"`
	DB   DatabaseConfig `yaml:"db"`
}

// Подконфигурация для базы данных
type DatabaseConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	SslMode  string `yaml:"sslMode"`
	TimeZone string `yaml:"timeZone"`
}

// Глобальная переменная для хранения конфигурации
var instance *Config
var once sync.Once

// Функция получения конфигурации
func GetConfig() *Config {

	// Загружаем .env
	if err := godotenv.Load(); err != nil {
		fmt.Println("Не удалось загрузить .env файл, переменные окружения могут отсутствовать")
	}

	once.Do(func() {
		instance = &Config{}

		if err := cleanenv.ReadConfig("./config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			fmt.Println(help)
		}

		// Читаем переменные окружения и заменяем значения в конфиге
		overrideWithEnv(instance)
	})
	return instance
}

// overrideWithEnv перезаписывает конфигурацию значениями из переменных окружения (если они заданы)
func overrideWithEnv(cfg *Config) {

	if user := os.Getenv("POSTGRES_USER"); user != "" {
		cfg.DB.User = user
	}
	if password := os.Getenv("POSTGRES_PASSWORD"); password != "" {
		cfg.DB.Password = password
	}
	if dbName := os.Getenv("POSTGRES_DB"); dbName != "" {
		cfg.DB.DBName = dbName
	}

}
