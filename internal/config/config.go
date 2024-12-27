package config

import "os"

type Config string

var (
	dbUser       = os.Getenv("DB_USER")
	dbPass       = os.Getenv("DB_PASS")
	dbHost       = os.Getenv("DB_HOST")
	dbPort       = os.Getenv("DB_PORT")
	dbName       = os.Getenv("DB_NAME")
	dbLog        = os.Getenv("DB_LOG")
	secretKey    = os.Getenv("SECRET_KEY")
	MailPassword = os.Getenv("MAIL_PASS")
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

func GetDBLog() string {
	return dbLog
}

func GetSecretKey() string {
	return secretKey
}

func GetMailPass() string {
	return MailPassword
}

func GetMailUsername() string {
	return "mmdhossein.haghdadi@gmail.com"
}
