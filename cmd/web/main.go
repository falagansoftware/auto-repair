package main

import (
	"fmt"
	"os"
	"strconv"

	gohtmx "github.com/falagansoftware/auto-repair/internal"
	"github.com/falagansoftware/auto-repair/internal/http"
	"github.com/falagansoftware/auto-repair/internal/postgres"
)

func main() {
	m := NewMain()
	m.Run()
}

type Main struct {
	DB         *postgres.DB
	HTTPServer *http.Server
	// Expose Services to e2e test
	UserService gohtmx.UserService
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
		DB:         postgres.NewDB(config.DB.Dsn),
		HTTPServer: http.NewServer(config.HTTPServer.Address, config.HTTPServer.Port, http.WithTimeout(1000)),
	}
}

func (m *Main) Run() {
	err := m.DB.Open()
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	m.UserService = postgres.NewUserService(m.DB)
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
