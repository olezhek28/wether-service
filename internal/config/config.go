package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// Config содержит все конфигурационные параметры сервиса.
type Config struct {
	Env  string `yaml:"env" env-default:"local"`
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
	DB   DBConfig
}

// DBConfig определяет параметры подключения к базе данных.
type DBConfig struct {
	DBPort   string `env:"DB_PORT"`     // Порт БД
	SSLMode  string `env:"SSL_MODE"`    // Режим SSL-соединения
	Username string `env:"DB_USER"`     // Имя пользователя БД
	Password string `env:"DB_PASSWORD"` // Пароль БД (загружается из переменной окружения)
	DBName   string `env:"DB_NAME"`     // Имя базы данных (может быть переопределено через ENV)
	DBHost   string `env:"DB_HOST"`     // Адрес хоста БД
}

// MustLoad загружает конфигурацию из файла или завершает работу при ошибке.
// Функция ищет путь к конфигурационному файлу через флаги командной строки
// или переменные окружения. Если путь не указан - вызывает панику.
func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("config file path is empty")
	}

	return MustLoadByPath(path)
}

// MustLoadByPath загружает конфигурацию из конкретного файла.
// Проверяет существование файла, парсит конфигурацию и возвращает объект Config.
// Вызывает панику при любой ошибке загрузки или парсинга.
func MustLoadByPath(path string) *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Env file does not exist: ", err.Error())
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found: " + path)
	}

	var config Config

	err := cleanenv.ReadConfig(path, &config)
	if err != nil {
		panic("failed to read config: " + err.Error())
	}
	return &config

}

// fetchConfigPath извлекает путь к конфигурационному файлу из:
// 1. Флагов командной строки (--config_path)
// 2. Переменной окружения CONFIG_PATH
// Если ни один источник не указан - возвращает пустую строку.
func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config_path", "", "Path to config path")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}
	return path
}
