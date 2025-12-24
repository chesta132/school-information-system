package config

import "os"

var (
	DB_USER     = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_HOST     = os.Getenv("DB_HOST")
	DB_PORT     = os.Getenv("DB_PORT")
	DB_NAME     = os.Getenv("DB_NAME")
	GO_ENV      = os.Getenv("GO_ENV")
	SERVER_PORT = os.Getenv("PORT")

	REFRESH_SECRET = os.Getenv("REFRESH_SECRET")
	ACCESS_SECRET  = os.Getenv("ACCESS_SECRET")

	INITIATE_ADMIN_KEY = os.Getenv("INITIATE_ADMIN_KEY")
)
