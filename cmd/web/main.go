package main

import (
	"fmt"
	"os"
	"strconv"

	autorepair "github.com/falagansoftware/auto-repair/internal"
	"github.com/falagansoftware/auto-repair/internal/http"
	"github.com/falagansoftware/auto-repair/internal/postgres"
	"github.com/falagansoftware/auto-repair/pkg/translator"
	"github.com/go-playground/validator"
)

func main() {
	m := NewMain()
	m.Run()
}

type Main struct {
	DB         *postgres.DB
	HTTPServer *http.Server
	// Epoxose common services to test
	I18n      *translator.Translator
	Validator *validator.Validate
	// Expose domain services to test
	UserService autorepair.UserService
}

type Config struct {
	DB         DBConfig
	HTTPServer HTTPConfig
}

type DBConfig struct {
	Dsn string
}

type HTTPConfig struct {
	Address string
	Port    int
}

func NewMain() *Main {
	config := NewConfig()
	return &Main{
		DB: postgres.NewDB(config.DB.Dsn),
		HTTPServer: http.NewServer(
			config.HTTPServer.Address,
			config.HTTPServer.Port,
			http.WithTimeout(1000),
		),
	}
}

func (m *Main) Run() {
	err := m.DB.Open()
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	m.I18n = translator.NewTranslator(translator.WithDefaultLang("es-es"))
	m.Validator = validator.New()
	m.UserService = postgres.NewUserService(m.DB)
	// Re attach for testing porpoises
	m.HTTPServer.I18n = m.I18n
	m.HTTPServer.Validator = m.Validator
	m.HTTPServer.UserService = m.UserService
	err = m.HTTPServer.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	} else {
		fmt.Printf("running: url=%q dsn=%q", m.HTTPServer.URL(), m.DB.Dsn)

	}
}

func NewConfig() Config {
	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		port = 6060
	}
	return Config{
		DB: DBConfig{
			Dsn: os.Getenv("DATABASE_URL"),
		},
		HTTPServer: HTTPConfig{
			Address: os.Getenv("SERVER_ADDRESS"),
			Port:    port,
		},
	}
}
