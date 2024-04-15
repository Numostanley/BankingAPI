package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type GetEnv struct {
	PortString       string
	PostgresHost     string
	PostgresUser     string
	PostgresDB       string
	PostgresPassword string
	PostgresPort     string
	AllowedHosts     string
}

func (env *GetEnv) LoadEnv() {
	err := godotenv.Load("../env/dev/.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	env.getPortString()
	env.getPostgresHost()
	env.getPostgresUser()
	env.getPostgresDB()
	env.getPostgresPassword()
	env.getPostgresPort()
	env.getAllowedHosts()
}

func (env *GetEnv) getPortString() {
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT not found in the environment")
	}
	env.PortString = portString
}

func (env *GetEnv) getPostgresHost() {
	postgresHost := os.Getenv("PG_HOST")
	if postgresHost == "" {
		log.Fatal("PG_HOST not found in the environment")
	}
	env.PostgresHost = postgresHost
}

func (env *GetEnv) getPostgresUser() {
	postgresUser := os.Getenv("PG_USER")
	if postgresUser == "" {
		log.Fatal("PG_USER not found in the environment")
	}
	env.PostgresUser = postgresUser
}

func (env *GetEnv) getPostgresDB() {
	postgresDB := os.Getenv("PG_DATABASE")
	if postgresDB == "" {
		log.Fatal("PG_DATABASE not found in the environment")
	}
	env.PostgresDB = postgresDB
}

func (env *GetEnv) getPostgresPassword() {
	postgresPassword := os.Getenv("PG_PASSWORD")
	if postgresPassword == "" {
		log.Fatal("PG_PASSWORD not found in the environment")
	}
	env.PostgresPassword = postgresPassword
}

func (env *GetEnv) getPostgresPort() {
	postgresPort := os.Getenv("PG_PORT")
	if postgresPort == "" {
		log.Fatal("PG_PORT not found in the environment")
	}
	env.PostgresPort = postgresPort
}

func (env *GetEnv) getAllowedHosts() {
	allowedHosts := os.Getenv("ALLOWED_HOSTS")
	if allowedHosts == "" {
		log.Fatal("ALLOWED_HOSTS not found in the environment")
	}
	env.AllowedHosts = allowedHosts
}
