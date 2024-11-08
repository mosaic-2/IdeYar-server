package config

import "os"

type Config string

var (
	dbUser = os.Getenv("DB_USER")
	dbPass = os.Getenv("DB_PASS")
	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
	dbName = os.Getenv("DB_NAME")
)

func GetDBUser() string {
	return dbUser
}

func GetDBPass() string {
	return dbPass
}

func GetDBHost() string {
	return dbHost
}

func GetDBPort() string {
	return dbPort
}

func GetDBName() string {
	return dbName
}
